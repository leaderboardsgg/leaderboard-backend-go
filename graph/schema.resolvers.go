package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"log"

	"github.com/speedrun-website/leaderboard-backend/database"
	"github.com/speedrun-website/leaderboard-backend/graph/generated"
	"github.com/speedrun-website/leaderboard-backend/graph/model"
)

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	db, err := database.GetDatabase()
	if err != nil {
		log.Println("Unable to connect to database", err)
		return nil, err
	}
	db.Find(&r.users)
	return r.users, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
