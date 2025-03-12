package services

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/spitfireooo/form-constructor-auth/pkg/database"
	"github.com/spitfireooo/form-constructor-auth/pkg/model/request"
	"github.com/spitfireooo/form-constructor-auth/pkg/model/response"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strconv"
)

func SignUp(user *request.User) (response.User, error) {
	salt, _ := strconv.Atoi(viper.GetString("crypt.salt"))
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), salt)
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

func SignIn() {

}
