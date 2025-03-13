package response

import "github.com/spitfireooo/form-constructor-auth/pkg/utils"

type User struct {
	ID       uint   `json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	Nickname string `json:"nickname" db:"nickname"`
	Logo     string `json:"logo" db:"logo"`
	CreateAt string `json:"create_at" db:"create_at"`
	UpdateAt string `json:"update_at" db:"update_at"`
}

type UserLogin struct {
	User   User `json:"user"`
	Tokens utils.Tokens
}
