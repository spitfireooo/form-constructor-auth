package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spitfireooo/form-constructor-auth/pkg/utils"
	"log"
	"net/http"
	"strings"
)

func IsAuthorized(ctx *fiber.Ctx) error {
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

func IsAuthor(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		log.Println("Bad params!", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Bad params",
		})
	}

	userId := ctx.Locals("user_id").(int64)

	if int64(id) != userId {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "You dont have permission",
		})
	}

	return ctx.Next()
}

//func IsAdmin(ctx *fiber.Ctx) error {
//	userId := ctx.Locals("user_id").(int)
//	if user, err := service.GetOneUser(userId); err != nil {
//		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
//			"message": "You dont have permission",
//		})
//	} else {
//		if user.Permission == "ADMIN" {
//			return ctx.Next()
//		}
//
//		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
//			"message": "You dont have permission",
//		})
//	}
//}

//func HasPermission(ctx *fiber.Ctx, mod string) error {
//	userId := ctx.Locals("user_id").(int)
//	if user, err := service.GetOneUser(userId); err != nil {
//		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
//			"message": "You dont have permission",
//		})
//	} else {
//		if user.Permission == mod {
//			return ctx.Next()
//		}
//
//		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
//			"message": "You dont have permission",
//		})
//	}
//}
