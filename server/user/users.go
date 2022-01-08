package user

import (
	"errors"

	"github.com/speedrun-website/leaderboard-backend/database"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Email    string `gorm:"unique"`
	Password []byte `gorm:"size:60"`
}

type UserIdentifier struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

type UserPersonal struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (u User) AsIdentifier() *UserIdentifier {
	return &UserIdentifier{
		ID:       u.ID,
		Username: u.Username,
	}
}

func (u User) AsPersonal() *UserPersonal {
	return &UserPersonal{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	}
}

// The globally exported UserStore that the application will use.
var Store UserStore

// The UserStore interface, which defines ways that the application
// can query for users.
type UserStore interface {
	database.DataStore

	GetUserIdentifierById(uint) (*UserIdentifier, error)
	GetUserPersonalById(uint) (*UserPersonal, error)
	GetUserByEmail(string) (*User, error)
	CreateUser(*User) error
	DeleteUser(uint) error
}

// Errors
var ErrUserNotFound = errors.New("the requested user was not found")

var ErrUserNotUnique = errors.New("attempted to create a user with duplicate data")

type UserCreationError struct {
	Err error
}

func (e UserCreationError) Error() string {
	return "the user creation failed with the following error: " + e.Err.Error()
}
