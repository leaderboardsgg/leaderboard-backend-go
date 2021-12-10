package database

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
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

func getConfig() dbConfig {
	return dbConfig{
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PASSWORD"),
	}
}

func getTestConfig() dbConfig {
	return dbConfig{
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_TEST_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_TEST_DB"),
		os.Getenv("POSTGRES_PASSWORD"),
	}
}

func getDns(config dbConfig) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		config.host, config.port, config.user, config.dbname, config.password)
}

func Init() error {
	config := getConfig()
	dns := getDns(config)
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = db
	return nil
}

func InitTest() error {
	config := getTestConfig()
	dns := getDns(config)
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = db
	return nil
}
