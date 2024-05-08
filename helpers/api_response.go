package helpers

import "github.com/gofiber/fiber/v2"

type DefaultApiResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrorValidationResponse struct {
	FailedField string      `json:"key"`
	Tag         string      `json:"tag"`
	Message     interface{} `json:"message"`
}

func ApiResponse(ctx *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return ctx.Status(statusCode).JSON(DefaultApiResponse{
		Message: message,
		Data:    data,
	})
}
