package user

type User struct {
	GUID         string `json:"guid" bson:"guid"`
	AccessToken  string `json:"access_token" bson:"access_token"`
	HashedRefreshToken []byte `json:"refresh_token" bson:"refresh_token"`
}
