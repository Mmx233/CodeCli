package global

import (
	"errors"
	"github.com/Mmx233/CodeCli/internal/global/models"
	"github.com/Mmx233/CodeCli/pkg/file"
	"github.com/Mmx233/config"
	"github.com/mitchellh/go-homedir"
	"log"
	"os"
)

var ConfigLoader *config.Config
var Config models.Config

func init() {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatalln(err)
	}
	ConfigLoader = config.NewConfig(&config.Options{
		Path:   file.JoinPath(home, ".CodeCli.yaml"),
		Config: &Config,
		Default: &models.Config{
			Default: models.Default{
				GitSite:    "github.com",
				CmdProgram: "powershell",
			},
			Storage: models.Storage{
				ProjectDir: file.JoinPath(home, "project"),
			},
			Rules: map[string][]string{
				"webstorm": {
					"package.json",
				},
				"goland": {
					"go.mod",
				},
				"rustrover": {
					"Cargo.toml",
				},
				"pycharm": {
					"pyproject.toml", "requirements.txt",
				},
				"studio": {
					file.JoinPath("android", "build.gradle"), "build.gradle",
				},
				"idea": {
					"gradlew",
				},
			},
		},
	})
	if err = ConfigLoader.Load(); err != nil {
		if errors.Is(err, config.IsNewConfig) {
			log.Println(err.Error())
			os.Exit(0)
		} else {
			log.Fatalln(err)
		}
	}
}
