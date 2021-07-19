package graphql_server

import (
	"context"
	"regexp"

	"github.com/speedrun-website/leaderboard-backend/data"
	"github.com/speedrun-website/leaderboard-backend/data/sql_driver"

	"github.com/samsarahq/go/oops"
	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/graphql/schemabuilder"
)

// Server is our graphql Server.
type Server struct {
	Users []*data.User
	Runs  []*data.Run

	SqlDriver sql_driver.SqlDriver
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
		allGames, err := s.SqlDriver.GetAllGames(ctx)
		if err != nil {
			return nil, oops.Wrapf(err, "getting games")
		}

		if args.TitleRegex == nil {
			return allGames, nil
		}

		re, err := regexp.Compile(*args.TitleRegex)
		if err != nil {
			return nil, oops.Wrapf(err, "compiling regex")
		}

		var games []*data.Game
		for _, game := range allGames {
			if re.MatchString(game.Title) {
				games = append(games, game)
			}
		}
		return games, nil
	})

	mutationObj := schema.Mutation()
	mutationObj.FieldFunc("add_game", func(ctx context.Context, args struct{ Title string }) error {
		if err := s.SqlDriver.InsertGame(ctx, &data.Game{Title: args.Title}); err != nil {
			return oops.Wrapf(err, "inserting row")
		}
		return nil
	})

	obj := schema.Object("Game", data.Game{})
	obj.FieldFunc("title", func(ctx context.Context, g *data.Game) string {
		return g.Title
	})
	obj.FieldFunc("runs", func(ctx context.Context, g *data.Game) []*data.Run {
		var runs []*data.Run
		for _, run := range s.Runs {
			if run.Game.Title == g.Title {
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
