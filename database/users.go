package database

import (
	"errors"

	"github.com/speedrun-website/leaderboard-backend/model"
)

var Users UserStore

type UserStore interface {
	GetUserIdentifierById(uint64) (*model.UserIdentifier, error)
	GetUserPersonalById(uint64) (*model.UserPersonal, error)
	GetUserByEmail(string) (*model.User, error)
	CreateUser(model.User) error
}

// Errors
var UserNotFoundError = errors.New("The requested user was not found.")

type UserUniquenessError struct {
	ErrorField string
}

func (e UserUniquenessError) Error() string {
	return "User creation failed"
}

type UserCreationError struct {
	Err error
}

func (e UserCreationError) Error() string {
	return "The user creation failed with the following error: " + e.Err.Error()
}
