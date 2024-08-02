package database

import (
	"errors"
	"time"
)

type RefreshToken struct {
	UserID    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (db *DB) SaveRefreshToken(userID int, token string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	refreshToken := RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	}

	dbStructure.RefreshTokens[token] = refreshToken

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) UserForRefreshToken(token string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	refreshToken, ok := dbStructure.RefreshTokens[token]
	if !ok {
		return User{}, errors.New("token doesn't exit")
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return User{}, errors.New("token is expired")
	}

	user, err := db.GetUserByID(refreshToken.UserID)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) RevokeRefreshToken(token string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	delete(dbStructure.RefreshTokens, token)
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}
	return nil
}
