package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	handlers "github.com/speedrun-website/leaderboard-backend/handlers"
	"github.com/speedrun-website/leaderboard-backend/middleware"
)

func InitRoutes(router *gin.Engine) {
	router.Use(cors.Default())
	api := router.Group("/api/v1")

	// the jwt middleware
	var authMiddleware = middleware.GetGinJWTMiddleware()

	// public routes
	api.POST("/register", handlers.RegisterUser)
	api.POST("/login", authMiddleware.LoginHandler)
	api.POST("/logout", authMiddleware.LogoutHandler)
	api.GET("/refresh_token", authMiddleware.RefreshHandler)
	api.GET("/ping", handlers.PingHandler)
	api.GET("/users/:id", handlers.GetUser)

	// auth routes
	api.Use(authMiddleware.MiddlewareFunc())
	{
		api.GET("/me", handlers.Me)
	}
}
