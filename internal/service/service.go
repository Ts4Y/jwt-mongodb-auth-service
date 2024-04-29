package service

import "golang.org/x/crypto/bcrypt"

func hashRefreshToken(token string) ([]byte, error) {
	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hashedRefreshToken, nil
}
