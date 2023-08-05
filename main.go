package main

import (
	"rabbitmq/conf"
	"rabbitmq/service"
)

func main() {
	conf.SetEnv()
	app := service.NewApp()

	app.Router.POST("/events", app.Emit)
	app.Start()
}
