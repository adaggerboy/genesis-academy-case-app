package main

import (
	"log"
	"os"

	"github.com/adaggerboy/genesis-academy-case-app/config"
	"github.com/adaggerboy/genesis-academy-case-app/pkg/database"
	"github.com/adaggerboy/genesis-academy-case-app/pkg/mailer"
	"github.com/adaggerboy/genesis-academy-case-app/pkg/routes"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"

	_ "github.com/adaggerboy/genesis-academy-case-app/pkg/init"
)

func main() {

	conffile := os.Getenv("CONFIG")
	if conffile == "" {
		conffile = "/etc/genesis-academy/currency-api/config.yaml"
	}

	gconfig, err := config.Load(conffile)
	if err != nil {
		log.Fatalf("%s", err)
	}
	config.GlobalConfig = gconfig

	err = database.InitDatabase(config.GlobalConfig.Database)
	if err != nil {
		log.Fatalf("%s", err)
	}

	c := cron.New()
	c.AddFunc(config.GlobalConfig.CronString, func() {
		err = mailer.GoThroughSubscriptions()
		if err != nil {
			log.Fatalf("Cron error: %s", err)
		}
	})
	c.Start()

	r := gin.New()
	routes.DeployRoutes(r)
	r.Run(config.GlobalConfig.HTTPServer.Endpoints...)
}
