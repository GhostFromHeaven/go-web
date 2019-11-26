package models

import (
	"errors"
	"github.com/go-redis/redis"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidLogin = errors.New("invalid login")
)

func AuthenticateUser(username, password string) error {
	hash, err := GetUser(username)
	if err == redis.Nil {
		return ErrUserNotFound
	} else if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		return ErrInvalidLogin
	}
	return nil
}

func RegisterUser(username, password string) error {
	cost := bcrypt.DefaultCost
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return err
	}
	return SaveUser(username, hash)
}

func GetUser(username string) ([]byte, error) {
	return client.Get("user:" + username).Bytes()
}

func SaveUser(username string, hash []byte) error {
	return client.Set("user:"+username, hash, 0).Err()
}
