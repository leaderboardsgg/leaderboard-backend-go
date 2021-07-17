package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/graphql/graphiql"
	"github.com/samsarahq/thunder/graphql/introspection"

	"github.com/speedrun-website/leaderboard-backend/data"
	"github.com/speedrun-website/leaderboard-backend/graphql_server"
	"github.com/speedrun-website/leaderboard-backend/middleware"
)

func main() {
	// Instantiate a server, build a server, and serve the schema on port 3030.
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
	server := &(graphql_server.Server{
		Games: games,
		Users: users,
		Runs:  runs,
	})

	schema := server.Schema()
	introspection.AddIntrospectionToSchema(schema)

	// Define a Mux
	router := mux.NewRouter()
	router.Use(middleware.PrometheusMiddleware)
	router.Use(middleware.NewAuthMiddleware)
	// Expose metrics.
	router.Path("/metrics").Handler(promhttp.Handler())
	// Expose schema and graphiql.
	router.Path("/graphql").Handler(graphql.Handler(schema))
	router.Path("/graphql/http").Handler(graphql.HTTPHandler(schema))
	router.PathPrefix("/graphiql/").Handler(http.StripPrefix("/graphiql/", graphiql.Handler()))

	// Initialize Prometheus bindings
	middleware.RegisterPrometheus()

	if err := http.ListenAndServe(":3030", router); err != nil {
		log.Fatal(err)
	}
}
