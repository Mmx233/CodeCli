package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Default Defaults `yaml:"default"`
	Storage Storage  `yaml:"storage"`
	Rules   []Rule   `yaml:"rules"`
}
type Defaults struct {
	GitSite    string `yaml:"gitSite"`
	Username   string `yaml:"username"`
	CmdProgram string `yaml:"cmdProgram"`
	Idea       string `yaml:"idea"`
}
type Storage struct {
	ProjectDir string `yaml:"projectDir"`
}
type Rule struct {
	Idea string   `yaml:"idea"`
	File []string `yaml:"file,flow"`
}
type Store struct{ Path string }

func UserStore() (Store, error) {
	home, err := os.UserHomeDir()
	return Store{Path: filepath.Join(home, ".CodeCli.yaml")}, err
}

func DefaultsFor(home string) Config {
	return Config{
		Default: Defaults{GitSite: "github.com"}, Storage: Storage{ProjectDir: filepath.Join(home, "project")},
		Rules: []Rule{
			{Idea: "webstorm", File: []string{"package.json"}}, {Idea: "goland", File: []string{"go.mod"}},
			{Idea: "rustrover", File: []string{"Cargo.toml"}}, {Idea: "pycharm", File: []string{"pyproject.toml", "requirements.txt"}},
			{Idea: "studio", File: []string{"android/build.gradle", "build.gradle"}}, {Idea: "idea", File: []string{"gradlew"}},
		},
	}
}

func (s Store) read() (Config, *yaml.Node, error) {
	cfg := DefaultsFor(filepath.Dir(s.Path))
	data, err := os.ReadFile(s.Path)
	if errors.Is(err, os.ErrNotExist) {
		return cfg, nil, nil
	}
	if err != nil {
		return cfg, nil, fmt.Errorf("read config: %w", err)
	}
	var doc yaml.Node
	dec := yaml.NewDecoder(strings.NewReader(string(data)))
	if err = dec.Decode(&doc); err != nil {
		return cfg, nil, fmt.Errorf("decode config: %w", err)
	}
	var extra yaml.Node
	if err = dec.Decode(&extra); err != io.EOF {
		return cfg, nil, errors.New("config must contain exactly one YAML document")
	}
	if len(doc.Content) != 1 || doc.Content[0].Kind != yaml.MappingNode {
		return cfg, nil, errors.New("config must be a YAML mapping")
	}
	if err = validate(doc.Content[0]); err != nil {
		return cfg, nil, err
	}
	if err = doc.Decode(&cfg); err != nil {
		return cfg, nil, fmt.Errorf("decode config: %w", err)
	}
	return cfg, &doc, nil
}

func (s Store) Load() (Config, error)  { cfg, _, err := s.read(); return cfg, err }
func (c Config) YAML() ([]byte, error) { return yaml.Marshal(c) }

func (c Config) Value(key string) (string, error) {
	switch key {
	case "default.gitSite":
		return c.Default.GitSite, nil
	case "default.username":
		return c.Default.Username, nil
	case "default.cmdProgram":
		return c.Default.CmdProgram, nil
	case "default.idea":
		return c.Default.Idea, nil
	case "storage.projectDir":
		return c.Storage.ProjectDir, nil
	default:
		return "", fmt.Errorf("unknown scalar key %q; rules must be edited manually", key)
	}
}

func mapping(n *yaml.Node, key string) *yaml.Node {
	for i := 0; i < len(n.Content); i += 2 {
		if n.Content[i].Value == key {
			return n.Content[i+1]
		}
	}
	return nil
}

func validate(root *yaml.Node) error {
	if root.Anchor != "" {
		return errors.New("config root YAML anchors are unsupported; edit anchored configuration manually")
	}
	for _, section := range []string{"default", "storage"} {
		n := mapping(root, section)
		if n == nil {
			continue
		}
		if n.Kind != yaml.MappingNode || n.Anchor != "" {
			return fmt.Errorf("%s must be a mapping (YAML aliases/merges are not supported)", section)
		}
		for i := 0; i < len(n.Content); i += 2 {
			key := section + "." + n.Content[i].Value
			if _, err := (Config{}).Value(key); err != nil {
				continue
			}
			if n.Content[i+1].Kind != yaml.ScalarNode || n.Content[i+1].Tag != "!!str" {
				return fmt.Errorf("%s must be a string", key)
			}
			if n.Content[i+1].Anchor != "" {
				return fmt.Errorf("%s uses a YAML anchor; edit anchored values manually", key)
			}
		}
	}
	if rules := mapping(root, "rules"); rules != nil {
		if rules.Kind != yaml.SequenceNode {
			return errors.New("rules must be a list")
		}
		for _, r := range rules.Content {
			if r.Kind != yaml.MappingNode {
				return errors.New("each rule must be a mapping")
			}
			idea, files := mapping(r, "idea"), mapping(r, "file")
			if idea == nil || idea.Tag != "!!str" || idea.Value == "" || files == nil || files.Kind != yaml.SequenceNode {
				return errors.New("each rule needs a string idea and a file list")
			}
			for _, f := range files.Content {
				if f.Tag != "!!str" || f.Value == "" {
					return errors.New("rule files must be nonempty strings")
				}
			}
		}
	}
	// Merges can change known keys behind the scalar setter's back.
	for i := 0; i < len(root.Content); i += 2 {
		if root.Content[i].Tag == "!!merge" {
			return errors.New("config YAML merges are not supported")
		}
	}
	for _, sec := range []string{"default", "storage"} {
		if n := mapping(root, sec); n != nil {
			for i := 0; i < len(n.Content); i += 2 {
				if n.Content[i].Tag == "!!merge" {
					return errors.New("config YAML merges are not supported")
				}
			}
		}
	}
	return nil
}

func (s Store) Set(key string, value *string) error {
	if _, err := (Config{}).Value(key); err != nil {
		return err
	}
	if info, err := os.Lstat(s.Path); err == nil && !info.Mode().IsRegular() {
		return errors.New("config writes require a regular file, not a symlink or directory")
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	cfg, doc, err := s.read()
	if err != nil {
		return err
	}
	if doc == nil {
		var root yaml.Node
		if err = root.Encode(cfg); err != nil {
			return err
		}
		doc = &yaml.Node{Kind: yaml.DocumentNode, Content: []*yaml.Node{&root}}
	}
	parts := strings.Split(key, ".")
	root := doc.Content[0]
	parent := mapping(root, parts[0])
	if parent == nil {
		parent = &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
		root.Content = append(root.Content, &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: parts[0]}, parent)
	}
	for i := 0; i < len(parent.Content); i += 2 {
		if parent.Content[i].Value != parts[1] {
			continue
		}
		if value == nil {
			parent.Content = append(parent.Content[:i], parent.Content[i+2:]...)
		} else {
			n := parent.Content[i+1]
			n.Kind = yaml.ScalarNode
			n.Tag = "!!str"
			n.Value = *value
			n.Content = nil
		}
		return s.save(doc)
	}
	if value != nil {
		parent.Content = append(parent.Content, &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: parts[1]}, &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: *value})
	}
	return s.save(doc)
}

func (s Store) save(doc *yaml.Node) error {
	if err := validate(doc.Content[0]); err != nil {
		return err
	}
	var cfg Config
	if err := doc.Decode(&cfg); err != nil {
		return err
	}
	data, err := yaml.Marshal(doc)
	if err != nil {
		return err
	}
	return atomicWrite(s.Path, data)
}

func atomicWrite(name string, data []byte) error {
	f, err := os.CreateTemp(filepath.Dir(name), ".codecli-config-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if _, err = f.Write(data); err != nil {
		f.Close()
		return err
	}
	if err = f.Sync(); err != nil {
		f.Close()
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	return replaceFile(f.Name(), name)
}

func (c Config) ProjectDir() (string, error) {
	p := c.Storage.ProjectDir
	if strings.HasPrefix(p, "~/") || strings.HasPrefix(p, `~\`) {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		p = filepath.Join(home, p[2:])
	}
	if p == "" {
		return "", errors.New("storage.projectDir is empty; set it with code config set")
	}
	return filepath.Abs(p)
}
