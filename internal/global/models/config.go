package models

type Config struct {
	Default Default    `yaml:"default"`
	Storage Storage    `yaml:"storage"`
	Rules   []IdeaRule `yaml:"rules"`
}

type Default struct {
	GitSite    string `yaml:"gitSite"`
	Username   string `yaml:"username"`
	CmdProgram string `yaml:"cmdProgram"`
	Idea       string `yaml:"idea"`
}

type Storage struct {
	ProjectDir string `yaml:"projectDir"`
}

type IdeaRule struct {
	Idea string   `yaml:"idea"`
	File []string `yaml:"file,flow"`
}
