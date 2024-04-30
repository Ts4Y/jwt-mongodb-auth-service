package token

import "github.com/dgrijalva/jwt-go"


type Token struct {
	AccessToken  string `json:"access_token" bson:"access_token"`
	RefreshToken string `json:"refresh_token" bson:"refresh_token"`
}

type Claims struct {
	GUID string `json:"guid"`
	jwt.StandardClaims
}


