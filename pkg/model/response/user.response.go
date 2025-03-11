package response

type User struct {
	ID       uint   `json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	Nickname string `json:"nickname" db:"nickname"`
	CreateAt string `json:"create_at" db:"create_at"`
	UpdateAt string `json:"update_at" db:"update_at"`
}
