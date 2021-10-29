package oauth

import (
	"errors"

	"github.com/speedrun-website/leaderboard-backend/model"
	"gorm.io/gorm"
)

type gormContainer struct {
	DB *gorm.DB
}

type OauthStore interface {
	GetUserByTwitterID(string) (*model.User, error)
	CreateUser(model.User) (*model.User, error)
}

var Oauth OauthStore

//initGormContainer initializes the oauth store with a gorm container
func InitGormContainer(db *gorm.DB) {
	Oauth = &gormContainer{
		DB: db,
	}
}

//GetUserByTwitterID fetches a user by their twitter ID
func (s gormContainer) GetUserByTwitterID(twitterID string) (*model.User, error) {
	var maybeExistingUser model.User
	result := s.DB.Model(&model.User{}).Where("twitter_id = ?", twitterID).First(&maybeExistingUser)

	// No error means we found a user
	if result.Error == nil {
		return &maybeExistingUser, nil
	}
	isNotFoundError := errors.Is(result.Error, gorm.ErrRecordNotFound)

	// We got a real error
	if !isNotFoundError {
		return nil, result.Error
	}
	return nil, nil
}

//CreateUser creates a user
//
//TODO: Replace with userstore method?
func (s gormContainer) CreateUser(user model.User) (*model.User, error) {
	err := s.DB.Create(&user).Error
	return &user, err
}
