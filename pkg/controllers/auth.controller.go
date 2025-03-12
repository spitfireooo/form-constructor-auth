package controller

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"github.com/spitfireooo/form-constructor-auth/pkg/model/request"
	"github.com/spitfireooo/form-constructor-auth/pkg/services"
	"github.com/spitfireooo/form-constructor-auth/pkg/utils"
	"log"
	"strconv"
	"time"
)

func SignUp(ctx *fiber.Ctx) error {
	body := new(request.User)

	// get image from req
	// if image exist - save, else - generate
	// import "github.com/teris-io/ld"
	// id := body.email
	// logo := ld.Generate(id, 100)

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
			Value:    res.AccessToken.Token,
			Expires:  time.Now().Add(time.Minute * time.Duration(res.AccessToken.Expires)),
			HTTPOnly: true,
			Secure:   true,
		})

		ctx.Cookie(&fiber.Cookie{
			Name:     "refresh_token",
			Value:    res.RefreshToken.Token,
			Expires:  time.Now().Add(time.Minute * time.Duration(res.RefreshToken.Expires)),
			HTTPOnly: true,
			Secure:   true,
		})

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"data": res,
		})
	}
}

func CurrentUser(ctx *fiber.Ctx) error {
	//userId := ctx.Cookie()

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": "CurrentUser",
	})
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

	user_id := token["user_id"].(int64)

	accessTokenExp, _ := strconv.Atoi(viper.GetString("jwt.access_token"))
	accessToken, err := utils.GenerateJWT(user_id, time.Duration(accessTokenExp))
	if err != nil {
		log.Println("Error in generate access-token")
		return err
	}

	refreshTokenExp, _ := strconv.Atoi(viper.GetString("jwt.refresh_token"))
	refreshToken, err = utils.GenerateJWT(user_id, time.Duration(refreshTokenExp))
	if err != nil {
		log.Println("Error in generate refresh-token")
		return err
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(time.Minute * time.Duration(accessTokenExp)),
		HTTPOnly: true,
		Secure:   true,
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Minute * time.Duration(refreshTokenExp)),
		HTTPOnly: true,
		Secure:   true,
	})

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func Logout(ctx *fiber.Ctx) error {
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
