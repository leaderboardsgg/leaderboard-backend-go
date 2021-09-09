package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"speedrun.website/graph/model"
)

type dbConfig struct {
	host     string
	port     int
	user     string
	dbname   string
	password string
}

var config = dbConfig{"localhost", 5432, "postgres", "speedrun.website", "example"}

func getDatabaseUrl() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
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

	// todo remove fake data
	users := []model.User{
		{Username: "CapitaineToinon", Email: "capitainetoinon@example.com"},
		{Username: "Johndoe", Email: "Johndoe@example.com"},
	}

	for i := 0; i < len(users); i++ {
		db.Create(&users[i])
	}

	return nil
}
