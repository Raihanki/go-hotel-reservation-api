package controllers

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/Raihanki/go-hotel-reservation-api/helpers"
	"github.com/Raihanki/go-hotel-reservation-api/request"
	"github.com/Raihanki/go-hotel-reservation-api/resources"
	"github.com/Raihanki/go-hotel-reservation-api/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ReservationControllerInterface interface {
	Index(ctx *fiber.Ctx) error
	Show(ctx *fiber.Ctx) error
	Store(ctx *fiber.Ctx) error
}

type ReservationController struct {
	reservationService services.ReservationServiceInterface
	roomNumberService  services.RoomNumberServiceInterface
}

func NewReservationController(reservationService services.ReservationServiceInterface, roomNumberService services.RoomNumberServiceInterface) ReservationControllerInterface {
	return &ReservationController{
		reservationService: reservationService, roomNumberService: roomNumberService,
	}
}

func (controller *ReservationController) Index(ctx *fiber.Ctx) error {
	userResource := ctx.Locals("user").(resources.UserResourceWithRole)
	reservations, errReservation := controller.reservationService.GetUserReservation(int(userResource.ID))
	if errReservation != nil {
		return helpers.ApiResponse(ctx, fiber.StatusInternalServerError, "Failed to get reservations", nil)
	}

	return helpers.ApiResponse(ctx, fiber.StatusOK, "Success get all reservations", reservations)
}

func (controller *ReservationController) Show(ctx *fiber.Ctx) error {
	paramID := ctx.Params("id")
	idReservation, _ := strconv.Atoi(paramID)

	userResource := ctx.Locals("user").(resources.UserResourceWithRole)
	reservation, errReservation := controller.reservationService.GetReservationByID(int(userResource.ID), idReservation)
	if errors.Is(errReservation, gorm.ErrRecordNotFound) {
		return helpers.ApiResponse(ctx, fiber.StatusNotFound, "Reservation not found", nil)
	}
	if errReservation != nil {
		return helpers.ApiResponse(ctx, fiber.StatusInternalServerError, "Failed to get reservation", nil)
	}

	return helpers.ApiResponse(ctx, fiber.StatusOK, "Success get reservation", reservation)
}

func (controller *ReservationController) Store(ctx *fiber.Ctx) error {
	reservationRequest := request.ReservationRequest{}
	errBody := ctx.BodyParser(&reservationRequest)
	if errBody != nil {
		log.Print("Error getting request : ", errBody.Error())
		return helpers.ApiResponse(ctx, fiber.StatusBadRequest, "Error getting request", nil)
	}

	//parse checkin time
	parsedCheckin, errParse := time.Parse("2006-01-02", reservationRequest.CheckIn)
	if errParse != nil {
		return helpers.ApiResponse(ctx, fiber.StatusBadRequest, "Failed to parse checkin date", nil)
	}

	reservationRequest.ParsedCheckin = parsedCheckin

	errValidate := validator.New().Struct(reservationRequest)
	if errValidate != nil {
		var errorResponses []helpers.ErrorValidationResponse
		for _, err := range errValidate.(validator.ValidationErrors) {
			errorResponse := helpers.ErrorValidationResponse{
				FailedField: err.Field(),
				Tag:         err.Tag(),
				Message:     err.Error(),
			}
			errorResponses = append(errorResponses, errorResponse)
		}
		return helpers.ApiResponse(ctx, fiber.StatusUnprocessableEntity, "Failed to validate request", errorResponses)
	}

	userResource := ctx.Locals("user").(resources.UserResourceWithRole)
	RoomNumberWithRoomResource, errRoomNumber := controller.roomNumberService.GetRoomNumberByID(int(reservationRequest.RoomNumberID))
	if errors.Is(errRoomNumber, gorm.ErrRecordNotFound) {
		return helpers.ApiResponse(ctx, fiber.StatusNotFound, "Room number not found", nil)
	}
	if errRoomNumber != nil {
		return helpers.ApiResponse(ctx, fiber.StatusInternalServerError, "Failed to get room number", nil)
	}

	reservation, errReservation := controller.reservationService.CreateReservation(int(userResource.ID), RoomNumberWithRoomResource.Room.Price, int(RoomNumberWithRoomResource.Room.ID), reservationRequest)
	if errReservation != nil {
		return helpers.ApiResponse(ctx, fiber.StatusInternalServerError, "Failed to create reservation", nil)
	}

	return helpers.ApiResponse(ctx, fiber.StatusCreated, "Success create reservation", reservation)
}
