package main

import (
	"log"
	"net/http"
	"time"

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

	// Setup middleware for all requests.
	middlewares := []middleware.ChainableMiddleware{
		middleware.NewAuthMiddleware,
	}

	// Expose schema and graphiql.
	http.Handle("/graphql", middleware.NewChainMiddlewareHandler(middlewares, graphql.Handler(schema)))
	http.Handle("/graphql/http", middleware.NewChainMiddlewareHandler(middlewares, graphql.HTTPHandler(schema)))
	http.Handle("/graphiql/", middleware.NewChainMiddlewareHandler(middlewares, http.StripPrefix("/graphiql/", graphiql.Handler())))
	if err := http.ListenAndServe(":3030", nil); err != nil {
		log.Fatal(err)
	}
}
