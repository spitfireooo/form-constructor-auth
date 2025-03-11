package entity

type User struct {
	ID       uint   `json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	Nickname string `json:"nickname" db:"nickname"`
	Password string `json:"password" db:"password"`
	CreateAt string `json:"create_at" db:"create_at"`
	UpdateAt string `json:"update_at" db:"update_at"`
}
