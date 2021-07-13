package main

import (
	"net/http"
	"time"

	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/graphql/graphiql"
	"github.com/samsarahq/thunder/graphql/introspection"

	"github.com/speedrun-website/leaderboard-backend/data"
	"github.com/speedrun-website/leaderboard-backend/graphql_server"
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

	// Expose schema and graphiql.
	http.Handle("/graphql", graphql.Handler(schema))
	http.Handle("/graphql/http", graphql.HTTPHandler(schema))
	http.Handle("/graphiql/", http.StripPrefix("/graphiql/", graphiql.Handler()))
	http.ListenAndServe(":3030", nil)
}
