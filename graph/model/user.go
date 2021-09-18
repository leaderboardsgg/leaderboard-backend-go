package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username" validate:"required"`
	Email    string `json:"email"  validate:"email"`
	Password string `json:"password" validate:"min=8,max=32,alphanum"`
}
