package middleware

import (
	"strings"

	"github.com/Raihanki/go-hotel-reservation-api/configs"
	"github.com/Raihanki/go-hotel-reservation-api/helpers"
	"github.com/Raihanki/go-hotel-reservation-api/models"
	"github.com/Raihanki/go-hotel-reservation-api/resources"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(ctx *fiber.Ctx) error {
	authorizationHeader := ctx.Get("Authorization")
	if authorizationHeader == "" {
		return helpers.ApiResponse(ctx, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	bearerToken := strings.Split(authorizationHeader, " ")[1]
	if bearerToken == "" {
		return helpers.ApiResponse(ctx, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	validatedToken, errValidatedToken := helpers.ValidateJWT(bearerToken)
	if errValidatedToken != nil {
		return helpers.ApiResponse(ctx, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	userId := validatedToken.Claims.(jwt.MapClaims)["user_id"]

	user := new(models.User)
	errUser := configs.DB.Where("id = ?", userId).Take(user).Error
	if errUser != nil {
		return helpers.ApiResponse(ctx, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	UserResourceWithRole := resources.UserResourceWithRole{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		IsAdmin:   user.IsAdmin,
		CreatedAt: user.CreatedAt,
	}
	ctx.Locals("user", UserResourceWithRole)

	return ctx.Next()
}
