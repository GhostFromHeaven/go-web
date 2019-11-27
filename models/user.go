package models

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrInvalidLogin  = errors.New("invalid login")
	ErrUsernameTaken = errors.New("username taken")
)

type User struct {
	Id int64
}

const (
	UserKeyByUsername = "user:by-username"
	UserKeyNextId     = "user:next-id"
)

func userKey(id int64) string {
	return fmt.Sprintf("user:%d", id)
}

func (user *User) Key() string {
	return userKey(user.Id)
}

func NewUser(username string, hash []byte) (*User, error) {
	exists, err := client.HExists(UserKeyByUsername, username).Result()
	if exists {
		return nil, ErrUsernameTaken
	}
	//if err != nil {
	//	return nil, err
	//}

	id, err := client.Incr(UserKeyNextId).Result()
	if err != nil {
		return nil, err
	}

	key := userKey(id)
	pipe := client.Pipeline()
	pipe.HSet(key, "id", id)
	pipe.HSet(key, "username", username)
	pipe.HSet(key, "hash", hash)
	pipe.HSet(UserKeyByUsername, username, id)
	_, err = pipe.Exec()
	if err != nil {
		return nil, err
	}

	return &User{id}, nil
}

func (user *User) GetId() (int64, error) {
	return user.Id, nil
}

func (user *User) GetUsername() (string, error) {
	return client.HGet(user.Key(), "username").Result()
}

func (user *User) GetHash() ([]byte, error) {
	return client.HGet(user.Key(), "hash").Bytes()
}

func (user *User) Authenticate(password string) error {
	hash, err := user.GetHash()
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return ErrInvalidLogin
	}

	return err
}

func GetUserByUsername(username string) (*User, error) {
	id, err := client.HGet(UserKeyByUsername, username).Int64()
	if err == redis.Nil {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	return &User{id}, nil
}

func AuthenticateUser(username, password string) (*User, error) {
	user, err := GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	return user, user.Authenticate(password)
}

func RegisterUser(username, password string) (*User, error) {
	cost := bcrypt.DefaultCost
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return nil, err
	}

	return NewUser(username, hash)
}
