package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/urfave/negroni"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/rs/zerolog"

	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/graphql/graphiql"
	"github.com/samsarahq/thunder/graphql/introspection"

	"github.com/speedrun-website/leaderboard-backend/data"
	"github.com/speedrun-website/leaderboard-backend/graphql_server"
	"github.com/speedrun-website/leaderboard-backend/logger"
	"github.com/speedrun-website/leaderboard-backend/middleware"
	"github.com/speedrun-website/leaderboard-backend/middleware/mux_adapter"
)

func main() {
	logger.InitLogger(getLogConfig())

	// Build schema to serve.
	schema := buildGraphQlSchema(logger.Logger)

	env, _ := godotenv.Read()

	// Set up middleware.
	middlewares := []negroni.Handler{
		negroni.HandlerFunc(middleware.PrometheusMiddleware),
		negroni.HandlerFunc(middleware.AuthMiddleware),
	}
	if val, ok := env["REQUEST_LOG"]; ok && val == "true" {
		middlewares = append(middlewares, negroni.HandlerFunc(middleware.LoggingMiddleware))
	}

	// Define our routes.
	router := mux.NewRouter()
	router.Use(mux_adapter.Middleware(middlewares...))
	router.Path("/metrics").Handler(promhttp.Handler())
	router.Path("/graphql").Handler(graphql.Handler(schema))
	router.Path("/graphql/http").Handler(graphql.HTTPHandler(schema))
	router.PathPrefix("/graphiql/").Handler(http.StripPrefix("/graphiql/", graphiql.Handler()))

	// Run the server.
	if err := http.ListenAndServe(":3030", router); err != nil {
		log.Fatal(err)
	}
}

func getLogConfig() logger.LogConfig {
	config := logger.LogConfig{}
	env, err := godotenv.Read()
	if err != nil {
		return config
	}

	config.Console = loadBoolFromEnvMap(env, "CONSOLE_LOG")
	config.Debug = loadBoolFromEnvMap(env, "DEBUG_LOG")

	return config
}

func loadBoolFromEnvMap(env map[string]string, varName string) bool {
	val, ok := env[varName]
	return ok && val == "true"
}

// Instantiate a server and build a schema.
func buildGraphQlSchema(logger zerolog.Logger) *graphql.Schema {
	games := []*data.Game{
		{Title: "Great game"},
		{Title: "Unloved game"},
	}
	users := []*data.User{
		{Name: "Fast runner"},
		{Name: "Slow runner"},
		{Name: "Non-runner"},
	}
	runs := []*data.Run{
		{Runner: users[0], Game: games[0], Time: 15 * time.Millisecond},
		{Runner: users[0], Game: games[0], Time: 17 * time.Millisecond},
		{Runner: users[1], Game: games[0], Time: 3 * time.Hour},
	}
	server := &graphql_server.Server{
		Games: games,
		Users: users,
		Runs:  runs,
	}

	schema := server.Schema()
	introspection.AddIntrospectionToSchema(schema)
	return schema
}
