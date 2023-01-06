package global

import (
	"github.com/Mmx233/CodeCli/global/models"
	"github.com/Mmx233/config"
	"log"
)

var Config models.Config

func init() {
	c := config.NewConfig(&config.Options{
		Config: &Config,
		Default: &models.Config{
			Default: models.Default{
				GitSite: "github.com",
			},
		},
	})
	if e := c.Load(); e != nil {
		log.Fatalln(e)
	}
}
