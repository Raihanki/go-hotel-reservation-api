package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/Raihanki/go-hotel-reservation-api/configs"
	"github.com/gofiber/fiber/v2"
)

func main() {
	configs.LoadConfigApp()
	configs.LoadDatabase()

	app := fiber.New(fiber.Config{})

	errListen := app.Listen(":" + configs.ENV.APP_PORT)
	if errListen != nil {
		log.Fatal("Error listen : " + errListen.Error())
	}
}
