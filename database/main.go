package database

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/speedrun-website/leaderboard-backend/graph/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

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

var dns = fmt.Sprintf(
	"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
	config.host, config.port, config.user, config.dbname, config.password)

func Init() error {
	var err error

	if DB, err = gorm.Open(postgres.Open(dns), &gorm.Config{}); err != nil {
		return err
	}

	if err = DB.AutoMigrate(&model.User{}); err != nil {
		return err
	}

	return nil
}
