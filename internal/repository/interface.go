package repository

import (
	"context"
)

type User interface {
	SaveRefreshToken(ctx context.Context, token []byte, guid string) error
	GetRefreshToken(ctx context.Context, guid string) ([]byte, error)
	UpdateRefreshToken(ctx context.Context, guid string,token []byte) error
}
