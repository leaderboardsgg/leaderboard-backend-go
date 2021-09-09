package database

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/joho/godotenv/autoload"
	"speedrun.website/graph/model"
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

func getDatabaseUrl() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		config.host, config.port, config.user, config.dbname, config.password)
}

func GetDatabase() (*gorm.DB, error) {
	db, err := gorm.Open("postgres", getDatabaseUrl())
	return db, err
}

func InitDb() error {
	db, err := GetDatabase()

	if err != nil {
		return err
	}

	if !db.HasTable(&model.User{}) {
		db.CreateTable(&model.User{})
	}

	return nil
}
