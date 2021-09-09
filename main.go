package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	database "speedrun.website/db"
	router "speedrun.website/router"
)

func main() {
	port := os.Getenv("PORT")
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	if port == "" {
		port = "3000"
	}

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
