package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/graphql/graphiql"
	"github.com/samsarahq/thunder/graphql/introspection"

	"github.com/speedrun-website/leaderboard-backend/data"
	"github.com/speedrun-website/leaderboard-backend/graphql_server"
	"github.com/speedrun-website/leaderboard-backend/middleware"
	"github.com/speedrun-website/leaderboard-backend/middleware/mux_adapter"
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

	// Set up middleware.
	middlewares := []negroni.Handler{
		negroni.HandlerFunc(middleware.PrometheusMiddleware),
		negroni.HandlerFunc(middleware.AuthMiddleware),
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
