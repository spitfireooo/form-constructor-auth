package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/spitfireooo/form-constructor-auth/pkg/controllers"
	"github.com/spitfireooo/form-constructor-auth/pkg/middlewares"
	"time"
)

func Router(r *fiber.App) {
	auth := r.Group("/auth")
	{
		auth.Post("/sign-up", controller.SignUp)
		auth.Post("/sign-in", controller.SignIn)
		auth.Get("/current", middleware.IsAuthorized, controller.CurrentUser)
		auth.Get("/refresh", controller.RefreshToken)
		auth.Get("/logout", controller.Logout)
		auth.Get("/:id", middleware.IsAuthorized, middleware.IsAuthor, func(ctx *fiber.Ctx) error {
			return ctx.JSON(fiber.Map{
				"message": "OK",
			})
		})
	}

	utils := r.Group("/utils", middleware.IsAuthorized, middleware.IsAdmin)
	{
		utils.Get("/metrics", monitor.New(monitor.Config{
			Title:   "Metrics Of Auth Service For Form Constructor",
			Refresh: time.Second,
		}))
	}
}
