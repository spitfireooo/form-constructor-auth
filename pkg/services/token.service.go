package service

import (
	"fmt"
	"github.com/spitfireooo/form-constructor-auth/pkg/model/entity"
	"github.com/spitfireooo/form-constructor-server-v2/pkg/database"
)

func CreateToken(userId int64, token string) (entity.Token, error) {
	res := new(entity.Token)

	query := fmt.Sprintf(`
		INSERT INTO %s (user_id, token)
		VALUES ($1, $2)
		RETURNING *
		`, database.TokensTable,
	)

	err := database.Connect.
		QueryRowx(query, userId, token).
		Scan(&res.ID, &res.UserId, &res.Token)

	return *res, err
}

func GetToken(userId int64) (entity.Token, error) {
	res := new(entity.Token)

	query := fmt.Sprintf(`SELECT * FROM %s WHERE user_id = $1`, database.TokensTable)
	err := database.Connect.Get(res, query, userId)

	return *res, err
}

func UpdateToken(userId int64, token string) (entity.Token, error) {
	res := new(entity.Token)

	query := fmt.Sprintf(
		`UPDATE %s SET 
	  	user_id = $1, token = $2
	  	WHERE user_id = $3
	  	RETURNING *`,
		database.TokensTable,
	)
	err := database.Connect.
		QueryRowx(query, userId, token, userId).
		Scan(&res.ID, &res.UserId, &res.Token)

	return *res, err
}

func DeleteToken(userId int64) error {
	query := fmt.Sprintf(`
		DELETE FROM %s 
	    WHERE user_id = $1
	    `, database.TokensTable,
	)
	_, err := database.Connect.Exec(query, userId)

	return err
}
