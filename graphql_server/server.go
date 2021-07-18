package graphql_server

import (
	"context"
	"regexp"

	"github.com/speedrun-website/leaderboard-backend/data"

	"github.com/samsarahq/go/oops"
	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/graphql/schemabuilder"
)

// Server is our graphql Server.
type Server struct {
	Users []*data.User
	Games []*data.Game
	Runs  []*data.Run
}

// Schema builds the graphql Schema.
func (s *Server) Schema() *graphql.Schema {
	builder := schemabuilder.NewSchema()
	s.registerEchoMutation(builder)
	s.registerGame(builder)
	s.registerUsers(builder)
	s.registerRuns(builder)
	return builder.MustBuild()
}

// registerEchoMutation registers the sample echo mutation type.
func (s *Server) registerEchoMutation(schema *schemabuilder.Schema) {
	obj := schema.Mutation()
	obj.FieldFunc("echo", func(args struct{ Message string }) string {
		return args.Message
	})
}

// registerGame registers the game type.
func (s *Server) registerGame(schema *schemabuilder.Schema) {
	queryObj := schema.Query()
	queryObj.FieldFunc("games", func(ctx context.Context, args struct{ TitleRegex *string }) ([]*data.Game, error) {
		if args.TitleRegex == nil {
			return s.Games, nil
		}

		re, err := regexp.Compile(*args.TitleRegex)
		if err != nil {
			return nil, oops.Wrapf(err, "compiling regex")
		}

		var games []*data.Game
		for _, game := range s.Games {
			if re.MatchString(game.Title) {
				games = append(games, game)
			}
		}
		return games, nil
	}, schemabuilder.Paginated)

	obj := schema.Object("Game", data.Game{})
	obj.Key("title")
	obj.FieldFunc("title", func(ctx context.Context, g *data.Game) string {
		return g.Title
	})
	obj.FieldFunc("runs", func(ctx context.Context, g *data.Game) []*data.Run {
		var runs []*data.Run
		for _, run := range s.Runs {
			if *run.Game == *g {
				runs = append(runs, run)
			}
		}
		return runs
	})
}

// registerGame registers the user type.
func (s *Server) registerUsers(schema *schemabuilder.Schema) {
	queryObj := schema.Query()
	queryObj.FieldFunc("users", func(ctx context.Context, args struct{ NameRegex *string }) ([]*data.User, error) {
		if args.NameRegex == nil {
			return s.Users, nil
		}

		re, err := regexp.Compile(*args.NameRegex)
		if err != nil {
			return nil, oops.Wrapf(err, "compiling regex")
		}

		var users []*data.User
		for _, user := range s.Users {
			if re.MatchString(user.Name) {
				users = append(users, user)
			}
		}
		return users, nil
	})

	obj := schema.Object("User", data.User{})
	obj.FieldFunc("name", func(ctx context.Context, u *data.User) string {
		return u.Name
	})
	obj.FieldFunc("runs", func(ctx context.Context, u *data.User) []*data.Run {
		var runs []*data.Run
		for _, run := range s.Runs {
			if *run.Runner == *u {
				runs = append(runs, run)
			}
		}
		return runs
	})
}

// registerGame registers the run type.
func (s *Server) registerRuns(schema *schemabuilder.Schema) {
	queryObj := schema.Query()
	queryObj.FieldFunc("runs", func(ctx context.Context, args struct{ MaxDurationMs *int64 }) ([]*data.Run, error) {
		if args.MaxDurationMs == nil {
			return s.Runs, nil
		}

		var runs []*data.Run
		for _, run := range s.Runs {
			if run.Time.Milliseconds() < *args.MaxDurationMs {
				runs = append(runs, run)
			}
		}
		return runs, nil
	})

	obj := schema.Object("Run", data.Run{})
	obj.FieldFunc("runner", func(ctx context.Context, r *data.Run) *data.User {
		return r.Runner
	})
	obj.FieldFunc("game", func(ctx context.Context, r *data.Run) *data.Game {
		return r.Game
	})
}
