package database

import (
	"fmt"

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
type UserNotFoundError struct {
	ID    uint64
	Email string
}

func (e UserNotFoundError) Error() string {
	var errString string
	if e.Email != "" {
		errString = fmt.Sprintf("User with email %s was not found", e.Email)
	} else {
		errString = fmt.Sprintf("User with ID %d was not found", e.ID)
	}
	return errString
}

type UserUniquenessError struct {
	User       model.User
	ErrorField string
}

func (e UserUniquenessError) Error() string {
	return "User creation failed"
}

type UserCreationError struct {
	User model.User
	Err  error
}

func (e UserCreationError) Error() string {
	return "The user creation failed with the following error: " + e.Err.Error()
}
