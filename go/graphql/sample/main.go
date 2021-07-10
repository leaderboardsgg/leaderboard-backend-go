package main

import (
	"context"
	"net/http"
	"time"

	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/graphql/graphiql"
	"github.com/samsarahq/thunder/graphql/introspection"
	"github.com/samsarahq/thunder/graphql/schemabuilder"
)

type user struct {
	Name string
}

type run struct {
	Runner *user
	Game   *game
	Time   time.Duration
}

type game struct {
	Title string
}

// server is our graphql server.
type server struct {
	users []*user
	games []*game
	runs  []*run
}

// registerEchoMutation registers the sample echo mutation type.
func (s *server) registerEchoMutation(schema *schemabuilder.Schema) {
	obj := schema.Mutation()
	obj.FieldFunc("echo", func(args struct{ Message string }) string {
		return args.Message
	})
}

// registerGame registers the game type.
func (s *server) registerGame(schema *schemabuilder.Schema) {
	queryObj := schema.Query()
	queryObj.FieldFunc("games", func() []*game {
		return s.games
	})

	obj := schema.Object("Game", game{})
	obj.FieldFunc("title", func(ctx context.Context, g *game) string {
		return g.Title
	})
	obj.FieldFunc("runs", func(ctx context.Context, g *game) []*run {
		var runs []*run
		for _, run := range s.runs {
			if *run.Game == *g {
				runs = append(runs, run)
			}
		}
		return runs
	})
}

// registerGame registers the user type.
func (s *server) registerUsers(schema *schemabuilder.Schema) {
	queryObj := schema.Query()
	queryObj.FieldFunc("users", func() []*user {
		return s.users
	})

	obj := schema.Object("User", user{})
	obj.FieldFunc("name", func(ctx context.Context, u *user) string {
		return u.Name
	})
	obj.FieldFunc("runs", func(ctx context.Context, u *user) []*run {
		var runs []*run
		for _, run := range s.runs {
			if *run.Runner == *u {
				runs = append(runs, run)
			}
		}
		return runs
	})
}

// registerGame registers the run type.
func (s *server) registerRuns(schema *schemabuilder.Schema) {
	queryObj := schema.Query()
	queryObj.FieldFunc("runs", func() []*run {
		return s.runs
	})

	obj := schema.Object("Run", run{})
	obj.FieldFunc("runner", func(ctx context.Context, r *run) *user {
		return r.Runner
	})
	obj.FieldFunc("game", func(ctx context.Context, r *run) *game {
		return r.Game
	})
}

// schema builds the graphql schema.
func (s *server) schema() *graphql.Schema {
	builder := schemabuilder.NewSchema()
	s.registerEchoMutation(builder)
	s.registerGame(builder)
	s.registerUsers(builder)
	s.registerRuns(builder)
	return builder.MustBuild()
}

func main() {
	// Instantiate a server, build a server, and serve the schema on port 3030.
	games := []*game{
		{Title: "Great game"},
		{Title: "Unloved game"},
	}
	users := []*user{
		{Name: "Fast runner"},
		{Name: "Slow runner"},
		{Name: "Non-runner"},
	}
	runs := []*run{
		{Runner: users[0], Game: games[0], Time: 15 * time.Millisecond},
		{Runner: users[0], Game: games[0], Time: 17 * time.Millisecond},
		{Runner: users[1], Game: games[0], Time: 3 * time.Hour},
	}
	server := &server{
		games: games,
		users: users,
		runs:  runs,
	}

	schema := server.schema()
	introspection.AddIntrospectionToSchema(schema)

	// Expose schema and graphiql.
	http.Handle("/graphql", graphql.Handler(schema))
	http.Handle("/graphql/http", graphql.HTTPHandler(schema))
	http.Handle("/graphiql/", http.StripPrefix("/graphiql/", graphiql.Handler()))
	http.ListenAndServe(":3030", nil)
}
