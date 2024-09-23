package store

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type Store interface {
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
}
