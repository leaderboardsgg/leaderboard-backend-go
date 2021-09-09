package model

type User struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email"  validate:"email"`
	Password string `json:"password" validate:"min=8,max=32,alphanum"`
}
