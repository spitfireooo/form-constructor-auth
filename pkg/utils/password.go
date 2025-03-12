package utils

import (
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

func CheckPassword(passwordHash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	salt, _ := strconv.Atoi(viper.GetString("crypt.salt"))
	hash, err := bcrypt.GenerateFromPassword([]byte(password), salt)
	return string(hash), err
}
