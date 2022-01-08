package user

import (
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/speedrun-website/leaderboard-backend/database"
	"gorm.io/gorm"
)

type gormUserStore struct {
	DB *gorm.DB
}

func (s gormUserStore) GetUserIdentifierById(userId uint) (*UserIdentifier, error) {
	var user UserIdentifier
	err := s.DB.Model(&User{}).First(&user, userId).Error
	if err != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func (s gormUserStore) GetUserPersonalById(userId uint) (*UserPersonal, error) {
	var user UserPersonal
	err := s.DB.Model(&User{}).First(&user, userId).Error
	if err != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func (s gormUserStore) GetUserByEmail(email string) (*User, error) {
	var user User
	err := s.DB.Where(User{
		Email: email,
	}).First(&user).Error
	if err != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func (s gormUserStore) CreateUser(user *User) error {
	err := s.DB.Create(user).Error

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return ErrUserNotUnique
		}
		return UserCreationError{
			Err: pgErr,
		}
	}

	return nil
}

func (s gormUserStore) DeleteUser(userId uint) error {
	if err := s.DB.Delete(&User{}, userId).Error; err != nil {
		return err
	}
	return nil
}

func (s gormUserStore) DumpDeleted() error {
	err := s.DB.Unscoped().Where("deleted_at IS NOT NULL").Delete(&User{}).Error
	if err != nil {
		return err
	}
	return nil
}

// Initializes a GORM user store and sets the exported
// user store for application use.
func InitGormStore(db *gorm.DB) error {
	if db == nil {
		db = database.DB
	}

	if err := db.AutoMigrate(&User{}); err != nil {
		return err
	}

	// Store is defined in users.go
	Store = &gormUserStore{
		DB: db,
	}
	return nil
}
