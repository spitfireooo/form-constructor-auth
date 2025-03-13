package request

type User struct {
	Email    string `json:"email" db:"email" validate:"required,email"`
	Nickname string `json:"nickname" db:"nickname" validate:"required,min=3,max=50"`
	Logo     string `json:"logg" db:"logo"`
	Password string `json:"password" db:"password" validate:"required,min=6"`
}

type UserLogin struct {
	Email    string `json:"email" db:"email" validate:"required,email"`
	Password string `json:"password" db:"password" validate:"required,min=6"`
}
