package middleware

import (
	"github.com/Raihanki/go-hotel-reservation-api/helpers"
	"github.com/Raihanki/go-hotel-reservation-api/resources"

	"github.com/gofiber/fiber/v2"
)

func AuthenticatedAsAdmin(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(resources.UserResourceWithRole)
	if !user.IsAdmin {
		return helpers.ApiResponse(ctx, fiber.StatusForbidden, "Access Forbidden", nil)
	}

	return ctx.Next()
}
