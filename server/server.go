package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"

	"github.com/speedrun-website/leaderboard-backend/server/user"
)

func Init(router *gin.Engine) {
	if err := initData(); err != nil {
		log.Fatalf("Could not initialize data stores: %s", err)
	}

	router.Use(cors.New(cors.Options{
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowedHeaders:   []string{"*"},
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            true,
	}))

	authMiddleware := user.GetAuthMiddlewareHandler()
	api := router.Group("/api/v1")

	user.PublicRoutes(api, authMiddleware)

	api.Use(authMiddleware.MiddlewareFunc())
	{
		user.AuthRoutes(api, authMiddleware)
	}
}

func initData() error {
	if err := user.InitGormStore(nil); err != nil {
		return err
	}
	return nil
}
