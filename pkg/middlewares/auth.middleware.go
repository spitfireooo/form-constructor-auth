package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spitfireooo/form-constructor-auth/pkg/utils"
	"log"
	"strconv"
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

func IsAuthorMiddleware(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Println("Bad params!", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Bad params",
		})
	}

	tokenString := ctx.Cookies("access_token")
	token, err := utils.ValidateToken(tokenString)
	if err != nil {
		log.Println("Invalid refresh token", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid refresh token",
		})
	}

	userId := int64(token["user_id"].(float64))

	if id != userId {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "You dont have permission",
		})
	}

	return ctx.Next()
}
