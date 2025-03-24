package request

type User struct {
	Email    string `json:"email" db:"email" validate:"required,email"`
	Phone    string `json:"phone" db:"phone" validate:"min=3"`
	Address  string `json:"address" db:"address" validate:"min=3"`
	Password string `json:"password" db:"password" validate:"min=6,max=30"`
	Nickname string `json:"nickname" db:"nickname" validate:"min=3"`
	Logo     string `json:"logo" db:"logo"`
}

type UserLogin struct {
	Email    string `json:"email" db:"email" validate:"required,email"`
	Password string `json:"password" db:"password" validate:"required,min=6"`
}
