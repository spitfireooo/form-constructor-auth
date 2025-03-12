package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/spitfireooo/form-constructor-auth/pkg/model/request"
	"github.com/spitfireooo/form-constructor-auth/pkg/services"
	"log"
)

func SignUp(ctx *fiber.Ctx) error {
	body := new(request.User)

	if err := ctx.BodyParser(body); err != nil {
		log.Println("Error in request parsing", err)
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Error in request parsing",
		})
	}

	if err := validator.New().Struct(body); err != nil {
		log.Println("Validation errors", err)
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Validation errors",
		})
	}

	if req, err := services.SignUp(body); err != nil {
		log.Println("Error in SignUp service", err)
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Error in SignUp service",
		})
	} else {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"data": req,
		})
	}
}

func SignIn(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "SignIn",
	})
}
