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
	roomNumberController := controllers.NewRoomNumberController(services.NewRoomNumberService(configs.DB), services.NewRoomService(configs.DB))
	roomRoute := app.Group("/rooms")
	roomRoute.Get("/", roomController.Index)
	roomRoute.Get("/:room_id/numbers", roomNumberController.Index)
	roomRoute.Get("/:id", roomController.Show)
	roomRoute.Post("/", middleware.Auth, middleware.AuthenticatedAsAdmin, roomController.Store)
	roomRoute.Put("/:id", middleware.Auth, middleware.AuthenticatedAsAdmin, roomController.Update)
	roomRoute.Delete("/:id", middleware.Auth, middleware.AuthenticatedAsAdmin, roomController.Destroy)

	// Room Number
	roomNumberRoute := app.Group("/room-numbers")
	roomNumberRoute.Post("/", middleware.Auth, middleware.AuthenticatedAsAdmin, roomNumberController.Store)
	roomNumberRoute.Delete("/:id", middleware.Auth, middleware.AuthenticatedAsAdmin, roomNumberController.Destroy)

	// Reservation Routes
	reservationController := controllers.NewReservationController(services.NewReservationService(configs.DB), services.NewRoomNumberService(configs.DB))
	reservationRoute := app.Group("/reservations")
	reservationRoute.Get("/", middleware.Auth, reservationController.Index)
	reservationRoute.Get("/:id", middleware.Auth, reservationController.Show)
	reservationRoute.Post("/", middleware.Auth, reservationController.Store)
}
