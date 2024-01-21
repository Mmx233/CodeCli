package global

import (
	"github.com/Mmx233/CodeCli/internal/global/models"
	"github.com/Mmx233/CodeCli/pkg/file"
	"github.com/Mmx233/config"
	"github.com/mitchellh/go-homedir"
	"log"
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
		},
	})
	if err = ConfigLoader.Load(); err != nil {
		if err == config.IsNewConfig {
			log.Println(err.Error())
		} else {
			log.Fatalln(err)
		}
	}
}
