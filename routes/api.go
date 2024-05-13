package routes

import (
	"github.com/Raihanki/go-hotel-reservation-api/configs"
	"github.com/Raihanki/go-hotel-reservation-api/controllers"
	"github.com/Raihanki/go-hotel-reservation-api/middleware"
	"github.com/Raihanki/go-hotel-reservation-api/services"
	"github.com/gofiber/fiber/v2"
)

func Router(route *fiber.App) {
	app := route.Group("/api/v1")

	// User Routes
	UserController := controllers.NewUserController(services.NewUserService(configs.DB))
	userRoute := app.Group("/users")
	userRoute.Post("/register", UserController.Register)
	userRoute.Post("/login", UserController.Login)
	userRoute.Get("/me", middleware.Auth, UserController.Me)

	// Hotel Routes
	hotelController := controllers.NewHotelController(services.NewHotelService(configs.DB))
	hotelRoute := app.Group("/hotels")
	hotelRoute.Get("/", hotelController.Index)
	hotelRoute.Get("/:id", hotelController.Show)
	hotelRoute.Post("/", middleware.Auth, middleware.AuthenticatedAsAdmin, hotelController.Store)
	hotelRoute.Put("/:id", middleware.Auth, middleware.AuthenticatedAsAdmin, hotelController.Update)
	hotelRoute.Delete("/:id", middleware.Auth, middleware.AuthenticatedAsAdmin, hotelController.Destroy)

	// Room Routes
	roomController := controllers.NewRoomController(services.NewRoomService(configs.DB), services.NewHotelService(configs.DB))
	roomRoute := app.Group("/hotels/:hotel_id/rooms")
	roomRoute.Get("/", roomController.Index)
	roomRoute.Get("/:id", roomController.Show)
	roomRoute.Post("/", middleware.Auth, middleware.AuthenticatedAsAdmin, roomController.Store)
	roomRoute.Put("/:id", middleware.Auth, middleware.AuthenticatedAsAdmin, roomController.Update)
	roomRoute.Delete("/:id", middleware.Auth, middleware.AuthenticatedAsAdmin, roomController.Destroy)
}
