package mongodb

import (
	"context"
	"jwt-mongo-auth/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type tokenRepo struct {
	db *mongo.Database
}

func NewTokenRepo(db *mongo.Database) repository.User {
	return &tokenRepo{
		db: db,
	}
}

func (r *tokenRepo) SaveRefreshToken(ctx context.Context, token string, guid string) error {

	_, err := r.db.Collection("tokens").UpdateOne(ctx, bson.M{
		"guid": guid,
	}, token)

	return err

}
