package router

import (
	"github.com/gin-gonic/gin"

	handlers "github.com/speedrun-website/leaderboard-backend/handlers"
	"github.com/speedrun-website/leaderboard-backend/middleware"
)

func InitRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")

	// the jwt middleware
	var authMiddleware = middleware.GetGinJWTMiddleware()

	// public routes
	api.POST("/register", handlers.RegisterHandler)
	api.POST("/login", authMiddleware.LoginHandler)
	api.POST("/logout", authMiddleware.LogoutHandler)
	api.GET("/refresh_token", authMiddleware.RefreshHandler)
	api.GET("/ping", handlers.PingHandler)
	api.GET("/users/:id", handlers.UserHandler)

	// auth routes
	api.Use(authMiddleware.MiddlewareFunc())
	{
		api.GET("/me", handlers.MeHandler)
	}
}
