package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"github.com/spitfireooo/form-constructor-auth/pkg/model/entity"
	"github.com/spitfireooo/form-constructor-auth/pkg/model/request"
	"github.com/spitfireooo/form-constructor-auth/pkg/services"
	"github.com/spitfireooo/form-constructor-auth/pkg/utils"
	"github.com/spitfireooo/form-constructor-server-v2/pkg/database"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"
)

// @Summary	SignUp
// @Tags Auth
// @Description sign un
// @ID sign-up
// @Accept json
// @Produce	json
// @Param input	body request.UserLogin true "body info"
// @Success 200 {object} response.User
// @Router /auth/sign-up [post]
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

	userExist := new(entity.User)
	query := fmt.Sprintf(`SELECT * FROM %s WHERE email = $1`, database.UsersTable)
	err := database.Connect.Get(userExist, query, body.Email)
	if err == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "This email is used now",
		})
	}

	if file, err := ctx.FormFile("logo"); err != nil {
		log.Println("File not found", err)
	} else {
		f, err := file.Open()
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error in file opening",
			})
		}
		defer f.Close()

		bodyReq := &bytes.Buffer{}
		writer := multipart.NewWriter(bodyReq)

		part, err := writer.CreateFormFile("logo", filepath.Base(file.Filename))
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error in creating form file for request",
			})
		}

		io.Copy(part, f)
		writer.Close()

		url := fmt.Sprintf(
			"%s:%s/api/v1/user/upload",
			viper.GetString("http.addr"),
			viper.GetString("http.port"),
		)
		req, err := http.NewRequest(fiber.MethodPost, url, bodyReq)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error in creating request",
			})
		}

		req.Header.Set("Content-Type", writer.FormDataContentType())

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error in response sending",
			})
		}
		defer res.Body.Close()

		if resReq, err := io.ReadAll(res.Body); err != nil {
			log.Println("Bad request", err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Bad request",
			})
		} else {
			type Filename struct {
				Filename string `json:"filename"`
			}
			var filename Filename
			err = json.Unmarshal(resReq, &filename)
			if err != nil {
				log.Println("Error in unmarshal json", err)
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Error in unmarshal json",
				})
			}

			body.Logo = &filename.Filename
		}
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

// @Summary	SignIn
// @Tags Auth
// @Description sign in
// @ID sign-in
// @Accept json
// @Produce	json
// @Param input	body request.UserLogin true "body info"
// @Success 200 {object} response.UserLogin
// @Router /auth/sign-in [post]
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

// @Summary	CurrentUser
// @Tags Auth
// @Description current user
// @ID current-user
// @Accept json
// @Produce	json
// @Success 200 {object} response.User
// @Router /auth/current [get]
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

// @Summary	RefreshToken
// @Tags Auth
// @Description refresh token
// @ID refresh-token
// @Accept json
// @Produce	json
// @Success 200 {object} response.UserLogin
// @Router /auth/refresh [get]
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

	// delete this
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

// @Summary	Logout
// @Tags Auth
// @Description logout
// @ID logout
// @Accept json
// @Produce	json
// @Success 200 {string} string
// @Router /auth/logout [get]
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
