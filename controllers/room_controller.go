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

type RoomControllerInterface interface {
	Index(ctx *fiber.Ctx) error
	Show(ctx *fiber.Ctx) error
	Store(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Destroy(ctx *fiber.Ctx) error
}

type RoomController struct {
	roomService  services.RoomServiceInterface
	hotelService services.HotelServiceInterface
}

func NewRoomController(roomService services.RoomServiceInterface, hotelService services.HotelServiceInterface) RoomControllerInterface {
	return &RoomController{roomService: roomService, hotelService: hotelService}
}

func (controller *RoomController) Index(ctx *fiber.Ctx) error {
	idHotelParam := ctx.Params("hotel_id")
	idHotel, _ := strconv.Atoi(idHotelParam)

	// check id hotel
	_, errHotel := controller.hotelService.GetHotelByID(idHotel)
	if errors.Is(errHotel, gorm.ErrRecordNotFound) {
		return helpers.ApiResponse(ctx, fiber.StatusBadRequest, "Hotel not found", nil)
	}
	if errHotel != nil {
		return helpers.ApiResponse(ctx, fiber.StatusInternalServerError, "Failed to get hotel", nil)
	}

	rooms, errRooms := controller.roomService.GetAllRooms(idHotel)
	if errRooms != nil {
		return helpers.ApiResponse(ctx, fiber.StatusInternalServerError, "Failed to get rooms", nil)
	}

	return helpers.ApiResponse(ctx, fiber.StatusOK, "Success get all rooms", rooms)
}

func (controller *RoomController) Show(ctx *fiber.Ctx) error {
	paramsId := ctx.Params("id")
	idRoom, _ := strconv.Atoi(paramsId)

	idHotelParam := ctx.Params("hotel_id")
	idHotel, _ := strconv.Atoi(idHotelParam)

	// check id hotel
	_, errHotel := controller.hotelService.GetHotelByID(idHotel)
	if errors.Is(errHotel, gorm.ErrRecordNotFound) {
		return helpers.ApiResponse(ctx, fiber.StatusBadRequest, "Hotel not found", nil)
	}
	if errHotel != nil {
		return helpers.ApiResponse(ctx, fiber.StatusInternalServerError, "Failed to get hotel", nil)
	}

	roomResource, errRoom := controller.roomService.GetRoomByID(idHotel, idRoom)
	if errors.Is(errRoom, gorm.ErrRecordNotFound) {
		return helpers.ApiResponse(ctx, fiber.StatusNotFound, "Room not found", nil)
	}
	if errRoom != nil {
		return helpers.ApiResponse(ctx, fiber.StatusInternalServerError, "Failed to get room", nil)
	}

	return helpers.ApiResponse(ctx, fiber.StatusOK, "Success get room", roomResource)
}

func (controller *RoomController) Store(ctx *fiber.Ctx) error {
	roomRequest := request.RoomRequest{}
	errBody := ctx.BodyParser(&roomRequest)
	if errBody != nil {
		return helpers.ApiResponse(ctx, fiber.StatusBadRequest, "Failed to process request", nil)
	}

	errValidate := validator.New().Struct(&roomRequest)
	if errValidate != nil {
		var errResponse []helpers.ErrorValidationResponse
		for _, err := range errValidate.(validator.ValidationErrors) {
			errResponse = append(errResponse, helpers.ErrorValidationResponse{
				FailedField: err.Field(),
				Tag:         err.Tag(),
				Message:     err.Error(),
			})
		}
		return helpers.ApiResponse(ctx, fiber.StatusUnprocessableEntity, "Failed to validate request", errResponse)
	}

	// check hotels
	_, errHotel := controller.hotelService.GetHotelByID(int(roomRequest.HotelID))
	if errors.Is(errHotel, gorm.ErrRecordNotFound) {
		return helpers.ApiResponse(ctx, fiber.StatusBadRequest, "Hotel not found", nil)
	}

	roomResource, errRoom := controller.roomService.CreateRoom(roomRequest)
	if errRoom != nil {
		return helpers.ApiResponse(ctx, fiber.StatusInternalServerError, "Failed to create room", nil)
	}

	return helpers.ApiResponse(ctx, fiber.StatusCreated, "Success create room", roomResource)
}

func (controller *RoomController) Update(ctx *fiber.Ctx) error {
	idParams := ctx.Params("id")
	roomID, _ := strconv.Atoi(idParams)

	idHotelParam := ctx.Params("hotel_id")
	idHotel, _ := strconv.Atoi(idHotelParam)

	// check id hotel
	_, errHotel := controller.hotelService.GetHotelByID(idHotel)
	if errors.Is(errHotel, gorm.ErrRecordNotFound) {
		return helpers.ApiResponse(ctx, fiber.StatusBadRequest, "Hotel not found", nil)
	}
	if errHotel != nil {
		return helpers.ApiResponse(ctx, fiber.StatusInternalServerError, "Failed to get hotel", nil)
	}

	roomRequest := request.RoomRequest{}
	errBody := ctx.BodyParser(&roomRequest)
	if errBody != nil {
		return helpers.ApiResponse(ctx, fiber.StatusBadRequest, "Failed to process request", nil)
	}

	errValidate := validator.New().Struct(&roomRequest)
	if errValidate != nil {
		var errResponse []helpers.ErrorValidationResponse
		for _, err := range errValidate.(validator.ValidationErrors) {
			errResponse = append(errResponse, helpers.ErrorValidationResponse{
				FailedField: err.Field(),
				Tag:         err.Tag(),
				Message:     err.Error(),
			})
		}
		return helpers.ApiResponse(ctx, fiber.StatusUnprocessableEntity, "Failed to validate request", errResponse)
	}

	//check room
	_, errRoom := controller.roomService.GetRoomByID(idHotel, roomID)
	if errors.Is(errRoom, gorm.ErrRecordNotFound) {
		return helpers.ApiResponse(ctx, fiber.StatusBadRequest, "Room not found", nil)
	}
	if errRoom != nil {
		return helpers.ApiResponse(ctx, fiber.StatusInternalServerError, "Failed to get room", nil)
	}

	// check hotels
	_, errCheckHotel := controller.hotelService.GetHotelByID(int(roomRequest.HotelID))
	if errors.Is(errCheckHotel, gorm.ErrRecordNotFound) {
		return helpers.ApiResponse(ctx, fiber.StatusBadRequest, "Hotel not found", nil)
	}
	if errCheckHotel != nil {
		return helpers.ApiResponse(ctx, fiber.StatusInternalServerError, "Failed to get hotel", nil)
	}

	errUpdateRoom := controller.roomService.UpdateRoomByID(roomRequest, roomID)
	if errUpdateRoom != nil {
		return helpers.ApiResponse(ctx, fiber.StatusInternalServerError, "Failed to update room", nil)
	}

	return helpers.ApiResponse(ctx, fiber.StatusOK, "Success update room", nil)
}

func (controller *RoomController) Destroy(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	roomID, _ := strconv.Atoi(idParam)

	idHotelParam := ctx.Params("hotel_id")
	idHotel, _ := strconv.Atoi(idHotelParam)

	// check id hotel
	_, errHotel := controller.hotelService.GetHotelByID(idHotel)
	if errors.Is(errHotel, gorm.ErrRecordNotFound) {
		return helpers.ApiResponse(ctx, fiber.StatusBadRequest, "Hotel not found", nil)
	}
	if errHotel != nil {
		return helpers.ApiResponse(ctx, fiber.StatusInternalServerError, "Failed to get hotel", nil)
	}

	rowsAffected, errRoom := controller.roomService.DeleteRoomByID(idHotel, roomID)
	if rowsAffected == 0 {
		return helpers.ApiResponse(ctx, fiber.StatusBadRequest, "Room not found", nil)
	}
	if errRoom != nil {
		return helpers.ApiResponse(ctx, fiber.StatusInternalServerError, "Failed to delete room", nil)
	}

	return helpers.ApiResponse(ctx, fiber.StatusNoContent, "Success delete room", nil)
}
