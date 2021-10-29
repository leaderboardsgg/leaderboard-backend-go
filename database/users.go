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
	CreateUser(model.User) error
}

// Errors
var ErrUserNotFound = errors.New("the requested user was not found")

type UserUniquenessError struct {
	ErrorField string
}

func (e UserUniquenessError) Error() string {
	return "user creation failed"
}

type UserCreationError struct {
	Err error
}

func (e UserCreationError) Error() string {
	return "the user creation failed with the following error: " + e.Err.Error()
}
