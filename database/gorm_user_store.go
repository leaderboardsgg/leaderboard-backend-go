package database

import (
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/speedrun-website/leaderboard-backend/model"
	"gorm.io/gorm"
)

// gormUserStore is an implementation of the UserStore interface
// defined in users.go that accesses the database using GORM.
type gormUserStore struct {
	DB *gorm.DB
}

func (s gormUserStore) GetUserIdentifierById(userId uint64) (*model.UserIdentifier, error) {
	var user model.UserIdentifier
	err := s.DB.Model(&model.User{}).First(&user, userId).Error
	if err != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func (s gormUserStore) GetUserPersonalById(userId uint64) (*model.UserPersonal, error) {
	var user model.UserPersonal
	err := s.DB.Model(&model.User{}).First(&user, userId).Error
	if err != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func (s gormUserStore) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := s.DB.Where(model.User{
		Email: email,
	}).First(&user).Error
	if err != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func (s gormUserStore) CreateUser(user *model.User) error {
	err := s.DB.Create(user).Error

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return ErrUserNotUnique
		}
		return UserCreationError{
			Err: err,
		}
	}

	return nil
}

// Initializes a GORM user store and sets the exported
// user store for application use.
func initGormUserStore(db *gorm.DB) error {
	if err := db.AutoMigrate(&model.User{}); err != nil {
		return err
	}

	// Users is defined in database/users.go
	Users = &gormUserStore{
		DB: db,
	}
	return nil
}
