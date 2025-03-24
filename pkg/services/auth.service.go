package service

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/spitfireooo/form-constructor-auth/pkg/model/entity"
	"github.com/spitfireooo/form-constructor-auth/pkg/model/request"
	"github.com/spitfireooo/form-constructor-auth/pkg/model/response"
	"github.com/spitfireooo/form-constructor-auth/pkg/utils"
	"github.com/spitfireooo/form-constructor-server-v2/pkg/database"
	"log"
)

func SignUp(user *request.User) (response.User, error) {
	passwordHash, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Println("Error in generate password of crypt")
	}

	res := new(response.User)

	query := fmt.Sprintf(`
		INSERT INTO %s (email, password) VALUES ($1, $2) 
		RETURNING id, email, phone, address, nickname, logo, created_at, updated_at
		`, database.UsersTable,
	)
	err = database.Connect.
		QueryRowx(query, user.Email, passwordHash).
		Scan(&res.ID, &res.Email, &res.Phone, &res.Address, &res.Nickname, &res.Logo, &res.CreatedAt, &res.UpdatedAt)

	return *res, err
}

func SignIn(user *request.UserLogin) (response.UserLogin, error) {
	userExist := new(entity.User)

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

	tokens, err := utils.GenerateTokens(int64(userExist.ID))
	if err != nil {
		log.Println("Error in generate refresh-token")
		return response.UserLogin{}, err
	}

	if _, err := GetToken(int64(userExist.ID)); err != nil {
		if _, err = CreateToken(int64(userExist.ID), tokens.RefreshToken.Token); err != nil {
			log.Println("Error in create token", err)
		}
	} else {
		if _, err = UpdateToken(int64(userExist.ID), tokens.RefreshToken.Token); err != nil {
			log.Println("Error in update token", err)
		}
	}

	res := response.UserLogin{
		User: response.User{
			ID:        userExist.ID,
			Email:     userExist.Email,
			Phone:     userExist.Phone,
			Address:   userExist.Address,
			Nickname:  userExist.Nickname,
			CreatedAt: userExist.CreatedAt,
			UpdatedAt: userExist.UpdatedAt,
		},
		Tokens: utils.Tokens{
			AccessToken:  utils.JwtToken{Token: tokens.AccessToken.Token, Expires: tokens.AccessToken.Expires},
			RefreshToken: utils.JwtToken{Token: tokens.RefreshToken.Token, Expires: tokens.RefreshToken.Expires},
		},
	}

	return res, err
}

func CurrentUser(userId int64) (response.User, error) {
	res := new(response.User)

	query := fmt.Sprintf(`SELECT id, email, nickname, logo, create_at, update_at FROM %s WHERE id = $1`, database.UsersTable)
	err := database.Connect.Get(res, query, userId)

	return *res, err
}
