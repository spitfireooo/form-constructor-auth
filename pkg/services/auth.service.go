package service

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"github.com/spitfireooo/form-constructor-auth/pkg/database"
	"github.com/spitfireooo/form-constructor-auth/pkg/model/entity"
	"github.com/spitfireooo/form-constructor-auth/pkg/model/request"
	"github.com/spitfireooo/form-constructor-auth/pkg/model/response"
	"github.com/spitfireooo/form-constructor-auth/pkg/utils"
	"log"
	"strconv"
	"time"
)

func SignUp(user *request.User) (response.User, error) {
	passwordHash, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Println("Error in generate password of crypt")
	}

	req := new(response.User)

	query := fmt.Sprintf(`
		INSERT INTO %s (email, nickname, password) 
		VALUES ($1, $2, $3) 
		RETURNING id, email, nickname, create_at, update_at
		`, database.UsersTable,
	)
	err = database.Connect.
		QueryRowx(query, user.Email, user.Nickname, passwordHash).
		Scan(&req.ID, &req.Email, &req.Nickname, &req.CreateAt, &req.UpdateAt)

	return *req, err
}

func SignIn(user *request.UserLogin) (response.UserLogin, error) {
	userExist := new(entity.User)

	fmt.Println(user)
	query := fmt.Sprintf(`
		SELECT * FROM %s WHERE email = $1
	`, database.UsersTable)
	err := database.Connect.Get(userExist, query, user.Email)
	if err != nil {
		return response.UserLogin{}, err
	}

	if !utils.CheckPassword(userExist.Password, user.Password) {
		return response.UserLogin{}, fiber.NewError(fiber.StatusBadRequest, "Invalid login or password")
	}

	accessTokenExp, _ := strconv.Atoi(viper.GetString("jwt.access_exp"))
	accessToken, err := utils.GenerateJWT(int64(userExist.ID), time.Duration(accessTokenExp))
	if err != nil {
		log.Println("Error in generate access-token")
		return response.UserLogin{}, err
	}

	refreshTokenExp, _ := strconv.Atoi(viper.GetString("jwt.refresh_exp"))
	refreshToken, err := utils.GenerateJWT(int64(userExist.ID), time.Duration(refreshTokenExp))
	if err != nil {
		log.Println("Error in generate refresh-token")
		return response.UserLogin{}, err
	}

	res := response.UserLogin{
		User: response.User{
			ID:       userExist.ID,
			Email:    userExist.Email,
			Nickname: userExist.Nickname,
			CreateAt: userExist.CreateAt,
			UpdateAt: userExist.UpdateAt,
		},
		AccessToken:  response.Token{Token: accessToken, Expires: int64(accessTokenExp)},
		RefreshToken: response.Token{Token: refreshToken, Expires: int64(refreshTokenExp)},
	}

	return res, err
}
