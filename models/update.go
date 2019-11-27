package models

import (
	"fmt"
	"strconv"
)

type Update struct {
	Id int64
}

const (
	UpdateKeyNextId  = "update:next-id"
	UpdateKeyUpdates = "updates"
)

func updateKey(id int64) string {
	return fmt.Sprintf("update:%d", id)
}

func userUpdateKey(userId int64) string {
	return fmt.Sprintf("user:%d:update", userId)
}

func (update *Update) Key() string {
	return updateKey(update.Id)
}

func NewUpdate(userId int64, body string) (*Update, error) {
	id, err := client.Incr(UpdateKeyNextId).Result()
	if err != nil {
		return nil, err
	}

	key := updateKey(id)
	pipe := client.Pipeline()
	pipe.HSet(key, "id", id)
	pipe.HSet(key, "user_id", userId)
	pipe.HSet(key, "body", body)
	pipe.LPush(UpdateKeyUpdates, id)
	pipe.LPush(userUpdateKey(userId), id)

	_, err = pipe.Exec()
	if err != nil {
		return nil, err
	}

	return &Update{id}, nil
}

func (update *Update) GetBody() (string, error) {
	return client.HGet(update.Key(), "body").Result()
}

func (update *Update) GetUser() (*User, error) {
	userId, err := client.HGet(update.Key(), "user_id").Int64()
	if err != nil {
		return nil, err
	}

	return &User{userId}, nil
}

func queryUpdates(key string) ([]*Update, error) {
	updateIds, err := client.LRange(key, 0, 10).Result()
	if err != nil {
		return nil, err
	}

	updates := make([]*Update, len(updateIds))
	for i, strId := range updateIds {
		id, err := strconv.Atoi(strId)
		if err != nil {
			return nil, err
		}

		updates[i] = &Update{int64(id)}
	}

	return updates, nil
}

func GetAllUpdates() ([]*Update, error) {
	return queryUpdates(UpdateKeyUpdates)
}

func GetUserUpdates(userId int64) ([]*Update, error) {
	return queryUpdates(userUpdateKey(userId))
}

func SaveUpdate(userId int64, body string) (*Update, error) {
	return NewUpdate(userId, body)
}
