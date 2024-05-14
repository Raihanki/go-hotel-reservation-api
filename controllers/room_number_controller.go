package controllers

import (
	"errors"
	"strconv"

	"github.com/Raihanki/go-hotel-reservation-api/helpers"
	"github.com/Raihanki/go-hotel-reservation-api/request"
	"github.com/Raihanki/go-hotel-reservation-api/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type RoomNumberControllerInterface interface {
	Index(ctx *fiber.Ctx) error
	Show(ctx *fiber.Ctx) error
	Store(ctx *fiber.Ctx) error
	Destroy(ctx *fiber.Ctx) error
}

type RoomNumberController struct {
	roomNumberService services.RoomNumberServiceInterface
	roomService       services.RoomServiceInterface
}

func NewRoomNumberController(roomNumberService services.RoomNumberServiceInterface, roomService services.RoomServiceInterface) RoomNumberControllerInterface {
	return &RoomNumberController{roomNumberService: roomNumberService, roomService: roomService}
}

func (controller *RoomNumberController) Index(ctx *fiber.Ctx) error {
	paramID := ctx.Params("room_id")
	roomID, _ := strconv.Atoi(paramID)

	roomNumbers, errRoomNumbers := controller.roomNumberService.GetAllRoomNumbers(roomID)
	if errRoomNumbers != nil {
		return helpers.ApiResponse(ctx, fiber.StatusInternalServerError, "Failed to get room numbers", nil)
	}

	return helpers.ApiResponse(ctx, fiber.StatusOK, "Success get all room numbers", roomNumbers)
}

func (controller *RoomNumberController) Show(ctx *fiber.Ctx) error {
	paramID := ctx.Params("id")
	roomNumberID, _ := strconv.Atoi(paramID)

	roomNumber, errRoomNumber := controller.roomNumberService.GetRoomNumberByID(roomNumberID)
	if errors.Is(errRoomNumber, gorm.ErrRecordNotFound) {
		return helpers.ApiResponse(ctx, fiber.StatusNotFound, "Room number not found", nil)
	}
	if errRoomNumber != nil {
		return helpers.ApiResponse(ctx, fiber.StatusInternalServerError, "Failed to get room number", nil)
	}

	return helpers.ApiResponse(ctx, fiber.StatusOK, "Success get room number", roomNumber)
}

func (controller *RoomNumberController) Store(ctx *fiber.Ctx) error {
	roomNumberRequest := request.RoomNumberRequest{}
	errBody := ctx.BodyParser(&roomNumberRequest)
	if errBody != nil {
		return helpers.ApiResponse(ctx, fiber.StatusBadRequest, "Failed to process request", nil)
	}

	errValidate := validator.New().Struct(roomNumberRequest)
	if errValidate != nil {
		var errResponses []helpers.ErrorValidationResponse
		for _, err := range errValidate.(validator.ValidationErrors) {
			errResponse := helpers.ErrorValidationResponse{
				FailedField: err.Field(),
				Tag:         err.Tag(),
				Message:     err.Error(),
			}
			errResponses = append(errResponses, errResponse)
		}
		return helpers.ApiResponse(ctx, fiber.StatusBadRequest, "Failed to validate request", errResponses)
	}

	roomNumber, errRoomNumber := controller.roomNumberService.CreateRoomNumber(roomNumberRequest)
	if errRoomNumber != nil {
		return helpers.ApiResponse(ctx, fiber.StatusInternalServerError, "Failed to create room number", nil)
	}

	return helpers.ApiResponse(ctx, fiber.StatusCreated, "Success create room number", roomNumber)
}

func (controller *RoomNumberController) Destroy(ctx *fiber.Ctx) error {
	paramID := ctx.Params("id")
	roomNumberID, _ := strconv.Atoi(paramID)

	affectedRows, _ := controller.roomNumberService.DeleteRoomNumberByID(roomNumberID)
	if affectedRows == 0 {
		return helpers.ApiResponse(ctx, fiber.StatusNotFound, "Room number not found", nil)
	}

	return helpers.ApiResponse(ctx, fiber.StatusNoContent, "Success delete room number", nil)
}
