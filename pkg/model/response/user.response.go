package response

type User struct {
	ID       uint   `json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	Nickname string `json:"nickname" db:"nickname"`
	CreateAt string `json:"create_at" db:"create_at"`
	UpdateAt string `json:"update_at" db:"update_at"`
}

type UserLogin struct {
	User         User  `json:"user"`
	AccessToken  Token `json:"access_token"`
	RefreshToken Token `json:"refresh_token"`
}

type Token struct {
	Token   string
	Expires int64
}
