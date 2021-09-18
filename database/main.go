package database

import (
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
	lock.Lock()
	defer lock.Unlock()

	if singleton != nil {
		return singleton, nil
	}

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&model.User{})

	singleton = db
	return db, nil
}

func Close() {
	if singleton != nil {
		if db, err := singleton.DB(); err != nil {
			db.Close()
			singleton = nil
		}
	}
}
