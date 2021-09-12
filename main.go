package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/gin-gonic/gin"
	"github.com/speedrun-website/leaderboard-backend/database"
	"github.com/speedrun-website/leaderboard-backend/router"
)

func main() {
	port := os.Getenv("BACKEND_PORT")
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	err := database.InitDb()

	if err != nil {
		log.Println("Unable to init database", err)
		panic(err)
	}

	router.InitRoutes(r)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
