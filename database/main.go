package database

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "github.com/joho/godotenv/autoload"
	"github.com/speedrun-website/leaderboard-backend/graph/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type dbConfig struct {
	host     string
	port     string
	user     string
	dbname   string
	password string
}

var config = dbConfig{
	os.Getenv("POSTGRES_HOST"),
	os.Getenv("POSTGRES_PORT"),
	os.Getenv("POSTGRES_USER"),
	os.Getenv("POSTGRES_DB"),
	os.Getenv("POSTGRES_PASSWORD"),
}

var lock sync.Mutex
var singleton *gorm.DB

var dns = fmt.Sprintf(
	"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
	config.host, config.port, config.user, config.dbname, config.password)

func GetDatabase() (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	lock.Lock()
	defer lock.Unlock()

	if singleton != nil {
		return singleton, nil
	}

	if db, err = gorm.Open(postgres.Open(dns), &gorm.Config{}); err != nil {
		return nil, err
	}

	singleton = db
	return db, nil
}

func Init() error {
	var db *gorm.DB
	var err error

	if db, err = gorm.Open(postgres.Open(dns), &gorm.Config{}); err != nil {
		return err
	}

	if err = db.AutoMigrate(&model.User{}); err != nil {
		return err
	}

	if err = close(db); err != nil {
		return err
	}

	return nil
}

func Close() error {
	if singleton == nil {
		return nil
	}

	if err := close(singleton); err != nil {
		return err
	}

	singleton = nil
	return nil
}

func close(db *gorm.DB) error {
	var sqlDB *sql.DB
	var err error

	if sqlDB, err = db.DB(); err != nil {
		return err
	}

	return sqlDB.Close()
}
