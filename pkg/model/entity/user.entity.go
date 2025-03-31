package entity

type User struct {
	ID        uint    `json:"id" db:"id"`
	Email     string  `json:"email" db:"email"`
	Password  string  `json:"password" db:"password"`
	Phone     *string `json:"phone" db:"phone"`
	Address   *string `json:"address" db:"address"`
	Nickname  *string `json:"nickname" db:"nickname"`
	Logo      *string `json:"logo" db:"logo"`
	CreatedAt string  `json:"created_at" db:"created_at"`
	UpdatedAt string  `json:"updated_at" db:"updated_at"`
}
