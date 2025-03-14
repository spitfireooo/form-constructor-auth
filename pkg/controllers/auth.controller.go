package controller

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/spitfireooo/form-constructor-auth/pkg/model/request"
	"github.com/spitfireooo/form-constructor-auth/pkg/services"
	"github.com/spitfireooo/form-constructor-auth/pkg/utils"
	"log"
	"net/http"
	"time"
)

func SignUp(ctx *fiber.Ctx) error {
	body := new(request.User)

	if err := ctx.BodyParser(body); err != nil {
		log.Println("Error in request parsing", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in request parsing",
		})
	}

	if err := validator.New().Struct(body); err != nil {
		log.Println("Validation errors", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation errors",
		})
	}

	if file, err := ctx.FormFile("logo"); err != nil {
		log.Println("Error in file upload", err)
	} else if file.Size > 0 {
		if err = utils.CheckContentType(
			file,
			"image/jpg",
			"image/png",
			"image/gif",
			"image/jpeg",
		); err != nil {
			log.Println("Bad ext of file", err)
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message": "Bad ext of file",
			})
		}

		filename := fmt.Sprintf("./static/uploads/%s_%s", uuid.New(), file.Filename)
		if err = ctx.SaveFile(file, filename); err != nil {
			log.Println("Error in save file", err)
		}

		body.Logo = filename
	}

	if res, err := service.SignUp(body); err != nil {
		log.Println("Error in SignUp service", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in SignUp service",
		})
	} else {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"data": res,
		})
	}
}

func SignIn(ctx *fiber.Ctx) error {
	body := new(request.UserLogin)

	if err := ctx.BodyParser(body); err != nil {
		log.Println("Error in request parsing", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in request parsing",
		})
	}

	if err := validator.New().Struct(body); err != nil {
		log.Println("Validation errors", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation errors",
		})
	}

	if res, err := service.SignIn(body); err != nil {
		log.Println("Error in SignIn service", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in SignIn service",
		})
	} else {
		ctx.Cookie(&fiber.Cookie{
			Name:     "access_token",
			Value:    res.Tokens.AccessToken.Token,
			Expires:  time.Now().Add(time.Minute * time.Duration(res.Tokens.AccessToken.Expires)),
			HTTPOnly: true,
			Secure:   true,
		})

		ctx.Cookie(&fiber.Cookie{
			Name:     "refresh_token",
			Value:    res.Tokens.RefreshToken.Token,
			Expires:  time.Now().Add(time.Minute * time.Duration(res.Tokens.RefreshToken.Expires)),
			HTTPOnly: true,
			Secure:   true,
		})

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"data": res,
		})
	}
}

func CurrentUser(ctx *fiber.Ctx) error {
	userId := ctx.Locals("user_id").(int64)

	if user, err := service.CurrentUser(userId); err != nil {
		log.Println("Error in CurrentUser service", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Error in currentUser service",
		})
	} else {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"data": user,
		})
	}
}

func RefreshToken(ctx *fiber.Ctx) error {
	refreshToken := ctx.Cookies("refresh_token")
	if refreshToken == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Missing refresh token",
		})
	}

	token, err := utils.ValidateToken(refreshToken)
	if err != nil {
		log.Println("Invalid refresh token", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid refresh token",
		})
	}

	user_id := int64(token["user_id"].(float64))

	tokens, err := utils.GenerateTokens(user_id)
	if err != nil {
		log.Println("Error in generate tokens")
		return err
	}

	if _, err := service.GetToken(user_id); err != nil {
		if _, err = service.CreateToken(user_id, tokens.RefreshToken.Token); err != nil {
			log.Println("Error in create token", err)
		}
	} else {
		if _, err = service.UpdateToken(user_id, tokens.RefreshToken.Token); err != nil {
			log.Println("Error in update token", err)
		}
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    tokens.AccessToken.Token,
		Expires:  time.Now().Add(time.Minute * time.Duration(tokens.AccessToken.Expires)),
		HTTPOnly: true,
		Secure:   true,
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken.Token,
		Expires:  time.Now().Add(time.Minute * time.Duration(tokens.RefreshToken.Expires)),
		HTTPOnly: true,
		Secure:   true,
	})

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"tokens": tokens,
	})
}

func Logout(ctx *fiber.Ctx) error {
	refreshToken := ctx.Cookies("refresh_token")
	if refreshToken == "" {
		log.Println("Missing refresh token")
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Missing refresh token",
		})
	}

	token, err := utils.ValidateToken(refreshToken)
	if err != nil {
		log.Println("Invalid refresh token", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid refresh token",
		})
	}

	userId := int64(token["user_id"].(float64))
	if err := service.DeleteToken(userId); err != nil {
		log.Println("Error in delete token", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error in delete token",
		})
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	})

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}
