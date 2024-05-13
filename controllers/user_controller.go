package controllers

import (
	"github.com/Raihanki/go-hotel-reservation-api/helpers"
	"github.com/Raihanki/go-hotel-reservation-api/request"
	"github.com/Raihanki/go-hotel-reservation-api/resources"
	"github.com/Raihanki/go-hotel-reservation-api/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserControllerInterface interface {
	Register(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	Me(ctx *fiber.Ctx) error
}

type UserController struct {
	userService services.UserServiceInterface
}

func NewUserController(userService services.UserServiceInterface) UserControllerInterface {
	return &UserController{userService: userService}
}

func (controller *UserController) Register(ctx *fiber.Ctx) error {
	registerRequest := request.RegisterRequest{}
	errBody := ctx.BodyParser(&registerRequest)
	if errBody != nil {
		return helpers.ApiResponse(ctx, fiber.StatusBadRequest, "Failed to process request", nil)
	}

	validateError := validator.New().Struct(registerRequest)
	if validateError != nil {
		var validationErrors []helpers.ErrorValidationResponse
		for _, err := range validateError.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, helpers.ErrorValidationResponse{
				FailedField: err.Field(),
				Tag:         err.Tag(),
				Message:     err.Error(),
			})
		}
		return helpers.ApiResponse(ctx, fiber.StatusUnprocessableEntity, "Failed to validate request", validationErrors)
	}

	_, errUser := controller.userService.GetUserByEmail(registerRequest.Email)
	if errUser == nil {
		return helpers.ApiResponse(ctx, fiber.StatusConflict, "Email already registered", nil)
	}

	token, errRegister := controller.userService.Register(registerRequest)
	if errRegister != nil {
		return helpers.ApiResponse(ctx, fiber.StatusInternalServerError, "Failed to register user", nil)
	}

	return helpers.ApiResponse(ctx, fiber.StatusCreated, "Success register user", token)
}

func (controller *UserController) Login(ctx *fiber.Ctx) error {
	loginRequest := request.LoginRequest{}
	errBody := ctx.BodyParser(&loginRequest)
	if errBody != nil {
		return helpers.ApiResponse(ctx, fiber.StatusBadRequest, "Failed to process request", nil)
	}

	errValidate := validator.New().Struct(loginRequest)
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

	token, errLogin := controller.userService.Login(loginRequest)
	if errLogin != nil {
		return helpers.ApiResponse(ctx, fiber.StatusBadRequest, "Wrong email or password", nil)
	}

	return helpers.ApiResponse(ctx, fiber.StatusOK, "Success login", token)
}

func (controller *UserController) Me(ctx *fiber.Ctx) error {
	UserResourceWithRole := ctx.Locals("user").(resources.UserResourceWithRole)

	return helpers.ApiResponse(ctx, fiber.StatusOK, "Success get user", resources.UserResource{
		ID:        UserResourceWithRole.ID,
		Name:      UserResourceWithRole.Name,
		Email:     UserResourceWithRole.Email,
		CreatedAt: UserResourceWithRole.CreatedAt,
	})
}
