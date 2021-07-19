package sql_driver

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/samsarahq/go/oops"
	"github.com/speedrun-website/leaderboard-backend/data"

	_ "github.com/lib/pq"
)

type SqlDriver interface {
	GetAllGames(ctx context.Context) ([]*data.Game, error)
	InsertGame(ctx context.Context, game *data.Game) error
}

type sqlDriver struct {
	db *sql.DB
}

func New() (SqlDriver, error) {
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5432, "postgres", "password", "leaderboard")
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, oops.Wrapf(err, "connecting to db")
	}

	return &sqlDriver{
		db: db,
	}, nil
}

func (s *sqlDriver) GetAllGames(ctx context.Context) ([]*data.Game, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT name FROM games")
	if err != nil {
		return nil, oops.Wrapf(err, "querying games table")
	}

	var gamesResp []*data.Game
	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			return nil, oops.Wrapf(err, "scanning row")
		}

		gamesResp = append(gamesResp, &data.Game{Title: title})
	}

	return gamesResp, nil
}

func (s *sqlDriver) InsertGame(ctx context.Context, game *data.Game) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO games(name) VALUES ($1)", game.Title)
	if err != nil {
		return oops.Wrapf(err, "inserting into games table")
	}

	return nil
}
