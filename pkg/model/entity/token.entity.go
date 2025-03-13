package entity

type Token struct {
	ID     uint   `json:"id" db:"id"`
	UserId uint   `json:"user_id" db:"user_id"`
	Token  string `json:"token" db:"token"`
}
