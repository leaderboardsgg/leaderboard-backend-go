package database

import (
	"fmt"
	"os"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/joho/godotenv/autoload"
	"github.com/speedrun-website/leaderboard-backend/graph/model"
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

func getDatabaseUrl() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		config.host, config.port, config.user, config.dbname, config.password)
}

func GetDatabase() (*gorm.DB, error) {
	lock.Lock()
	defer lock.Unlock()

	if singleton != nil {
		return singleton, nil
	}

	db, err := gorm.Open("postgres", getDatabaseUrl())

	if err != nil {
		return nil, err
	}

	if !db.HasTable(&model.User{}) {
		db.CreateTable(&model.User{})
	}

	singleton = db
	return db, nil
}

func Close() {
	if singleton != nil {
		singleton.Close()
	}
}
