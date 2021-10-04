package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	_ "github.com/joho/godotenv/autoload"

	"github.com/gin-gonic/gin"
	"github.com/speedrun-website/leaderboard-backend/database"
	"github.com/speedrun-website/leaderboard-backend/router"
)

func main() {
	port := os.Getenv("BACKEND_PORT")

	if err := database.Init(); err != nil {
		log.Fatal(err)
		panic(err)
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	router.InitRoutes(r)

	go func() {
		if err := http.ListenAndServe(":"+port, r); err != nil {
			log.Fatal(err)
			panic(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server ...")
}
