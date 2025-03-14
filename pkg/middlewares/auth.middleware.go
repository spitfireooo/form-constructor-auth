package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spitfireooo/form-constructor-auth/pkg/utils"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func AuthMiddleware(ctx *fiber.Ctx) error {
	var accessToken string

	authorization := ctx.Get("Authorization")
	if strings.HasPrefix(authorization, "Bearer ") {
		accessToken = strings.TrimPrefix(authorization, "Bearer ")
	} else if ctx.Cookies("access_token") != "" {
		accessToken = ctx.Cookies("access_token")
	}

	if accessToken == "" {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "Missing authorization token",
		})
	}

	claims, err := utils.ValidateToken(strings.Replace(accessToken, "Bearer ", "", 1))
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token",
		})
	}

	ctx.Locals("user_id", int64(claims["user_id"].(float64)))

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

	userId := ctx.Locals("user_id").(int64)

	if id != userId {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "You dont have permission",
		})
	}

	return ctx.Next()
}
