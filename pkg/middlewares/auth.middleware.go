package middlewares

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/spitfireooo/form-constructor-auth/pkg/utils"
)

func AuthMiddleware(ctx *fiber.Ctx) error {
	tokenString := ctx.Cookies("access_token")
	fmt.Println("TokenString", tokenString)

	if tokenString == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Missing access token",
		})
	}

	if _, err := utils.ValidateToken(tokenString); err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token",
		})
	}

	return ctx.Next()
}
