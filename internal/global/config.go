package global

import (
	"errors"
	"github.com/Mmx233/CodeCli/internal/global/models"
	"github.com/Mmx233/CodeCli/pkg/file"
	"github.com/Mmx233/config"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"runtime"
)

var ConfigLoader *config.Config
var Config models.Config

func init() {
	logger := log.WithField("component", "config")
	home, err := homedir.Dir()
	if err != nil {
		logger.Fatalln(err)
	}

	defaultConfig := models.Config{
		Default: models.Default{
			GitSite: "github.com",
		},
		Storage: models.Storage{
			ProjectDir: path.Join(home, "project"),
		},
		Rules: []models.IdeaRule{
			{
				Idea: "webstorm",
				File: []string{"package.json"},
			},
			{
				Idea: "goland",
				File: []string{"go.mod"},
			},
			{
				Idea: "rustrover",
				File: []string{"Cargo.toml"},
			},
			{
				Idea: "pycharm",
				File: []string{"pyproject.toml", "requirements.txt"},
			},
			{
				Idea: "studio",
				File: []string{path.Join("android", "build.gradle"), "build.gradle"},
			},
			{
				Idea: "idea",
				File: []string{"gradlew"},
			},
		},
	}

	switch runtime.GOOS {
	case "windows":
		defaultConfig.Default.CmdProgram = "powershell"
	case "linux":
		defaultConfig.Default.CmdProgram = "gnome-terminal"
	}

	home = file.PreparePath(home)
	ConfigLoader = config.NewConfig(&config.Options{
		Path:    path.Join(home, ".CodeCli.yaml"),
		Config:  &Config,
		Default: &defaultConfig,
	})
	if err = ConfigLoader.Load(); err != nil {
		if errors.Is(err, config.IsNewConfig) {
			logger.Infoln(err.Error())
			os.Exit(0)
		} else {
			logger.Fatalln(err)
		}
	}

	if len(Config.Rules) == 0 {
		logger.Warnln("no match rule found")
	}
	Config.Storage.ProjectDir = file.PreparePath(path.Clean(Config.Storage.ProjectDir))
}
