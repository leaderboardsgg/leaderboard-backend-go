package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"min=8,max=32,alphanum"`
}

type UserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserRegister struct {
	Username        string `json:"username" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"password_confirm" validate:"eqfield=Password"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
