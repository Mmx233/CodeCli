package global

import (
	"github.com/Mmx233/CodeCli/global/models"
	"github.com/Mmx233/config"
	"github.com/Mmx233/tool"
	"github.com/mitchellh/go-homedir"
	"log"
	"os"
)

var Config models.Config

func init() {
	c := config.NewConfig(&config.Options{
		Config: &Config,
		Default: &models.Config{
			Default: models.Default{
				GitSite: "github.com",
			},
			Storage: models.Storage{
				ProjectDir: "~/project",
			},
		},
	})
	e := c.Load()
	if e != nil {
		log.Fatalln(e)
	}

	Config.Storage.ProjectDir, e = homedir.Expand(Config.Storage.ProjectDir)
	if e != nil {
		log.Fatalln(e)
	}

	if !tool.File.Exists(Config.Storage.ProjectDir) {
		if e = os.Mkdir(Config.Storage.ProjectDir, 0600); e != nil {
			log.Fatalln(e)
		}
	}
}
