package entity

type User struct {
	Email    string `json:"email" db:"email"`
	Nickname string `json:"nickname" db:"nickname"`
	Password string `json:"password" db:"password"`
}
