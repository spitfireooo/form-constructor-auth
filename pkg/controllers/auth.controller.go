package controllers

import "github.com/gofiber/fiber/v2"

func SignUp(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "SignUp",
	})
}

func SignIn(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "SignIn",
	})
}
