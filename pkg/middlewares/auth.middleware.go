package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spitfireooo/form-constructor-auth/pkg/utils"
)

func AuthMiddleware(ctx *fiber.Ctx) error {
	tokenString := ctx.Cookies("access_token")

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

	//token := ctx.Get("Authorization")
	//if token == "" {
	//	return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
	//		"message": "Missing authorization token",
	//	})
	//}
	//
	//claims, err := utils.ValidateToken(strings.Replace(token, "Bearer ", "", 1))
	//if err != nil {
	//	return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
	//		"message": "Invalid token",
	//	})
	//}
	//
	//ctx.Locals("user_id", claims["user_id"])

	return ctx.Next()
}
