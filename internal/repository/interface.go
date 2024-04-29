package repository

import "context"


type User interface{

	SaveRefreshToken(ctx context.Context, token string, guid string ) error
}


