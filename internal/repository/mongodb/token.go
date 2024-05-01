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

func (r *tokenRepo) SaveRefreshToken(ctx context.Context, token []byte, guid string) error {

	_, err := r.db.Collection("tokens").InsertOne(ctx, bson.M{
		"guid":          guid,
		"refresh_token": token,
	})

	return err
}

func (r *tokenRepo) GetRefreshToken(ctx context.Context, guid string) ([]byte, error) {
	var result struct {
		RefreshToken []byte `bson:"refresh_token"`
	}

	err := r.db.Collection("tokens").FindOne(ctx, bson.M{"guid": guid}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.RefreshToken, nil
}

func(r *tokenRepo)  UpdateRefreshToken(ctx context.Context, guid string,token []byte) error{
	_,err := r.db.Collection("tokens").UpdateOne(ctx, bson.M{"guid": guid},bson.M{"$set": bson.M{"refresh_token": token}})
	if err != nil {
		return err
	}
	return err
}