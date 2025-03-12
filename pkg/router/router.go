package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spitfireooo/form-constructor-auth/pkg/controllers"
)

func Router(r *fiber.App) {
	auth := r.Group("/auth")
	{
		auth.Post("/sign-up", controllers.SignUp)
		auth.Post("/sign-in", controllers.SignIn)
	}
}
