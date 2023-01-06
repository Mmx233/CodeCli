package global

import (
	"github.com/Mmx233/CodeCli/global/models"
	"github.com/Mmx233/CodeCli/pkg/file"
	"github.com/Mmx233/config"
	"github.com/mitchellh/go-homedir"
	"log"
)

var Config models.Config

func init() {
	home, e := homedir.Dir()
	if e != nil {
		log.Fatalln(e)
	}
	c := config.NewConfig(&config.Options{
		Path:   file.JoinPath(home, ".CodeCli.yaml"),
		Config: &Config,
		Default: &models.Config{
			Default: models.Default{
				GitSite: "github.com",
			},
			Storage: models.Storage{
				ProjectDir: file.JoinPath(home, "project"),
			},
		},
	})
	if e = c.Load(); e != nil {
		log.Fatalln(e)
	}
}
