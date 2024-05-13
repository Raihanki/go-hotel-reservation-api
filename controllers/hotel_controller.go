package controllers

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"strconv"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/Raihanki/go-hotel-reservation-api/helpers"
	"github.com/Raihanki/go-hotel-reservation-api/request"
	"github.com/Raihanki/go-hotel-reservation-api/services"
	"github.com/gofiber/fiber/v2"
)

type HotelControllerInterface interface {
	Index(ctx *fiber.Ctx) error
	Show(ctx *fiber.Ctx) error
	Store(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Destroy(ctx *fiber.Ctx) error
}

type HotelController struct {
	hotelService services.HotelServiceInterface
}

func NewHotelController(hotelService services.HotelServiceInterface) HotelControllerInterface {
	return &HotelController{hotelService: hotelService}
}

func (controller *HotelController) Index(ctx *fiber.Ctx) error {
	hotels, errHotel := controller.hotelService.GetAllHotels()
	if errHotel != nil {
		log.Error("Error get all hotels : ", errHotel.Error())
		return helpers.ApiResponse(ctx, fiber.StatusInternalServerError, "Failed to get hotels", nil)
	}

	return helpers.ApiResponse(ctx, fiber.StatusOK, "Success get all hotels", hotels)
}

func (controller *HotelController) Show(ctx *fiber.Ctx) error {
	hotelIdParams := ctx.Params("id")
	hotelId, _ := strconv.Atoi(hotelIdParams)
	hotel, err := controller.hotelService.GetHotelByID(hotelId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return helpers.ApiResponse(ctx, fiber.StatusNotFound, "Hotel not found", nil)
	}
	if err != nil {
		log.Error("Error get hotel by ID : ", err.Error())
		return helpers.ApiResponse(ctx, fiber.StatusInternalServerError, "Failed to get hotel", nil)
	}

	return helpers.ApiResponse(ctx, fiber.StatusOK, "Success get hotel", hotel)
}

func (controller *HotelController) Store(ctx *fiber.Ctx) error {
	hotelRequest := request.HotelRequest{}
	errBody := ctx.BodyParser(&hotelRequest)
	if errBody != nil {
		return helpers.ApiResponse(ctx, fiber.StatusBadRequest, "Failed to process request", nil)
	}

	errValidate := validator.New().Struct(&hotelRequest)
	if errValidate != nil {
		var errValidatorResponse []helpers.ErrorValidationResponse
		for _, err := range errValidate.(validator.ValidationErrors) {
			errValidatorResponse = append(errValidatorResponse, helpers.ErrorValidationResponse{
				FailedField: err.Field(),
				Tag:         err.Tag(),
				Message:     err.Error(),
			})
		}
		return helpers.ApiResponse(ctx, fiber.StatusUnprocessableEntity, "Failed to validate request", errValidatorResponse)
	}

	hotelResources, errHotel := controller.hotelService.CreateHotel(hotelRequest)
	if errHotel != nil {
		log.Error("Error create hotel : ", errHotel.Error())
		return helpers.ApiResponse(ctx, fiber.StatusInternalServerError, "Failed to create hotel", nil)
	}

	return helpers.ApiResponse(ctx, fiber.StatusCreated, "Success create hotel", hotelResources)
}

func (controller *HotelController) Update(ctx *fiber.Ctx) error {
	hotelIdParams := ctx.Params("id")
	hotelId, _ := strconv.Atoi(hotelIdParams)

	hotelRequest := request.HotelRequest{}
	errBody := ctx.BodyParser(&hotelRequest)
	if errBody != nil {
		return helpers.ApiResponse(ctx, fiber.StatusBadRequest, "Failed to process request", nil)
	}

	errValidate := validator.New().Struct(&hotelRequest)
	if errValidate != nil {
		var errValidatorResponse []helpers.ErrorValidationResponse
		for _, err := range errValidate.(validator.ValidationErrors) {
			errValidatorResponse = append(errValidatorResponse, helpers.ErrorValidationResponse{
				FailedField: err.Field(),
				Tag:         err.Tag(),
				Message:     err.Error(),
			})
		}
		return helpers.ApiResponse(ctx, fiber.StatusUnprocessableEntity, "Failed to validate request", errValidatorResponse)
	}

	// check hotel ID
	hotelResource, errGetHotel := controller.hotelService.GetHotelByID(hotelId)
	if errors.Is(errGetHotel, gorm.ErrRecordNotFound) {
		return helpers.ApiResponse(ctx, fiber.StatusNotFound, "Hotel not found", nil)
	}
	if errGetHotel != nil {
		log.Error("Error get hotel by ID : ", errGetHotel.Error())
		return helpers.ApiResponse(ctx, fiber.StatusInternalServerError, "Failed to get hotel", nil)
	}

	errUpdateHotel := controller.hotelService.UpdateHotelByID(hotelRequest, int(hotelResource.ID))
	if errUpdateHotel != nil {
		log.Error("Error update hotel : ", errUpdateHotel.Error())
		return helpers.ApiResponse(ctx, fiber.StatusInternalServerError, "Failed to update hotel", nil)
	}

	return helpers.ApiResponse(ctx, fiber.StatusOK, "Success update hotel", nil)
}

func (controller *HotelController) Destroy(ctx *fiber.Ctx) error {
	hotelIdParams := ctx.Params("id")
	hotelId, _ := strconv.Atoi(hotelIdParams)

	rowsAffected, errHotel := controller.hotelService.DeleleHotelByID(hotelId)
	if errHotel != nil {
		log.Error("Error delete hotel : ", errHotel.Error())
		return helpers.ApiResponse(ctx, fiber.StatusInternalServerError, "Failed to delete hotel", nil)
	}

	if rowsAffected == 0 {
		return helpers.ApiResponse(ctx, fiber.StatusNotFound, "Hotel not found", nil)
	}

	return helpers.ApiResponse(ctx, fiber.StatusNoContent, "Success delete hotel", nil)
}
