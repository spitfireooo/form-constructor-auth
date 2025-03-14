package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spitfireooo/form-constructor-auth/pkg/controllers"
	"github.com/spitfireooo/form-constructor-auth/pkg/middlewares"
)

func Router(r *fiber.App) {
	auth := r.Group("/auth")
	{
		auth.Post("/sign-up", controller.SignUp)
		auth.Post("/sign-in", controller.SignIn)
		auth.Get("/current", middlewares.AuthMiddleware, controller.CurrentUser)
		auth.Get("/refresh", controller.RefreshToken)
		auth.Get("/logout", controller.Logout)
		auth.Get("/:id", middlewares.AuthMiddleware, middlewares.IsAuthorMiddleware, func(ctx *fiber.Ctx) error {
			return ctx.JSON(fiber.Map{
				"message": "OK",
			})
		})
	}
}
