package router

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	controllers "speedrun.website/controller"
	"speedrun.website/graph"
	"speedrun.website/graph/generated"
	"speedrun.website/middleware"
)

const defaultPort = ":8080"

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
	api.POST("/login", authMiddleware.LoginHandler)
	api.GET("/ping", controllers.PingHandler)
	api.GET("/users", controllers.UsersHandler)
	api.GET("/refresh_token", authMiddleware.RefreshHandler)

	// auth routes
	api.Use(authMiddleware.MiddlewareFunc())
	{
		api.GET("/users/me", controllers.MeHandler)
	}
}
