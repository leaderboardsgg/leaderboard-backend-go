package router

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"speedrun.website/graph"
	"speedrun.website/graph/generated"
	handlers "speedrun.website/handlers"
	"speedrun.website/middleware"
	"speedrun.website/validators"
)

// Defining the Graphql handler
func graphqlHandler() gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func InitRoutes(router *gin.Engine) {
	router.GET("/playground", playgroundHandler())
	router.POST("/query", graphqlHandler())
	api := router.Group("/api/v1")

	// the jwt middleware
	var authMiddleware = middleware.GetGinJWTMiddleware()

	// public routes
	api.POST("/register", validators.RegisterValidator(), handlers.RegisterHandler)
	api.POST("/login", validators.LoginValidator(), authMiddleware.LoginHandler)
	api.POST("/logout", authMiddleware.LogoutHandler)
	api.GET("/refresh_token", authMiddleware.RefreshHandler)
	api.GET("/ping", handlers.PingHandler)
	api.GET("/users", handlers.UsersHandler)

	// auth routes
	api.Use(authMiddleware.MiddlewareFunc())
	{
		api.GET("/users/me", handlers.MeHandler)
	}
}
