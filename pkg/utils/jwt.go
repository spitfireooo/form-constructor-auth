package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"log"
	"os"
	"strconv"
	"time"
)

type Tokens struct {
	AccessToken  JwtToken
	RefreshToken JwtToken
}

type JwtToken struct {
	Token   string
	Expires int64
}

func GenerateTokens(userId int64) (Tokens, error) {
	accessTokenExp, _ := strconv.Atoi(viper.GetString("jwt.access_exp"))
	accessToken, err := GenerateJWT(userId, time.Duration(accessTokenExp))
	if err != nil {
		log.Println("Error in generate access-token")
		return Tokens{}, err
	}

	refreshTokenExp, _ := strconv.Atoi(viper.GetString("jwt.refresh_exp"))
	refreshToken, err := GenerateJWT(userId, time.Duration(refreshTokenExp))
	if err != nil {
		log.Println("Error in generate refresh-token")
		return Tokens{}, err
	}

	return Tokens{
		AccessToken:  JwtToken{Token: accessToken, Expires: int64(accessTokenExp)},
		RefreshToken: JwtToken{Token: refreshToken, Expires: int64(refreshTokenExp)},
	}, nil
}

func GenerateJWT(userId int64, duration time.Duration) (string, error) {
	exp := time.Now().Add(time.Minute * duration).Unix()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["user_id"] = userId
	claims["exp"] = exp

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid token")
	}
}
