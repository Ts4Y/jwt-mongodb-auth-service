package repository

import "context"


type User interface{

	SaveRefreshToken(ctx context.Context, token []byte, guid string ) error
}


