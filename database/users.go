package database

import (
	"errors"

	"github.com/speedrun-website/leaderboard-backend/model"
)

// The globally exported UserStore that the application will use.
var Users UserStore

// The UserStore interface, which defines ways that the application
// can query for users.
type UserStore interface {
	GetUserIdentifierById(uint64) (*model.UserIdentifier, error)
	GetUserPersonalById(uint64) (*model.UserPersonal, error)
	GetUserByEmail(string) (*model.User, error)
	CreateUser(*model.User) error
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
