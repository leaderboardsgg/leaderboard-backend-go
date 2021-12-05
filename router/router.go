package router

import (
	"html/template"
	"net/http"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"

	"github.com/speedrun-website/leaderboard-backend/handlers"

	"github.com/speedrun-website/leaderboard-backend/handlers/oauth"
	"github.com/speedrun-website/leaderboard-backend/middleware"
)

var indexTemplate = `
    <p><a href="/api/v1/oauth/authenticate?provider=twitter">Log in with Twitter</a></p>
`

func InitRoutes(router *gin.Engine) {
	router.Use(cors.New(cors.Options{
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowedHeaders:   []string{"*"},
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            true,
	}))
	api := router.Group("/api/v1")

	// the jwt middleware
	var authMiddleware = middleware.GetGinJWTMiddleware()

	// public routes
	api.POST("/register", handlers.RegisterUser)
	api.POST("/login", authMiddleware.LoginHandler)
	api.POST("/logout", authMiddleware.LogoutHandler)
	api.GET("/refresh_token", authMiddleware.RefreshHandler)
	api.GET("/ping", handlers.Ping)
	api.GET("/users/:id", handlers.GetUser)
	api.GET("/oauth/authenticate/", oauth.OauthLogin)
	api.GET("/oauth/callback/:provider", oauth.OauthCallback)
	api.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
		t, _ := template.New("foo").Parse(indexTemplate)
		t.Execute(c.Writer, indexTemplate)
	})

	// auth routes
	api.Use(authMiddleware.MiddlewareFunc())
	{
		api.GET("/me", handlers.Me)
	}
}
