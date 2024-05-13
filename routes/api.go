package routes

import (
	"github.com/Raihanki/go-hotel-reservation-api/configs"
	"github.com/Raihanki/go-hotel-reservation-api/controllers"
	"github.com/Raihanki/go-hotel-reservation-api/services"
	"github.com/gofiber/fiber/v2"
)

func Router(route *fiber.App) {
	app := route.Group("/api/v1")

	// Hotel Routes
	hotelController := controllers.NewHotelController(services.NewHotelService(configs.DB))
	hotelRoute := app.Group("/hotels")
	hotelRoute.Get("/", hotelController.Index)
	hotelRoute.Get("/:id", hotelController.Show)
	hotelRoute.Post("/", hotelController.Store)
	hotelRoute.Put("/:id", hotelController.Update)
	hotelRoute.Delete("/:id", hotelController.Destroy)
}
