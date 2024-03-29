package models

import (
	"fmt"
)

type Update struct {
	key string
}

func NewUpsate(userId int64, body string) (*Update, error) {
	id, err := client.Incr("update:next-id").Result()
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("update:%d", id)
	pipe := client.Pipeline()
	pipe.HSet(key, "id", id)
	pipe.HSet(key, "user_id", userId)
	pipe.HSet(key, "body", body)
	pipe.LPush("updates", id)
	pipe.LPush(fmt.Sprintf("user:%d:updates", userId), id)
	_, err = pipe.Exec()
	if err != nil {
		return nil, err
	}
	return &Update{key}, nil
}

func (update *Update) GetBody() (string, error) {
	return client.HGet(update.key, "body").Result()
}

func (update *Update) GetUser() (*User, error) {
	userId, err := client.HGet(update.key, "user_id").Int64()
	if err != nil {
		return nil, err
	}
	return GetUserById(userId)
}

func GetAllUpdates() ([]*Update, error) {
	updateIds, err := client.LRange("updates", 0, 10).Result()
	if err != nil {
		return nil, err
	}
	updates := make([]*Update, len(updateIds))
	for i, id := range updateIds {
		key := "update:" + id
		updates[i] = &Update{key}
	}
	return updates, nil
}

func GetUpdates(userId int64) ([]*Update, error) {
	key := fmt.Sprintf("user:%d:updates", userId)
	updateIds, err := client.LRange(key, 0, 10).Result()
	if err != nil {
		return nil, err
	}
	updates := make([]*Update, len(updateIds))
	for i, id := range updateIds {
		key := "update:" + id
		updates[i] = &Update{key}
	}
	return updates, nil
}

func PostUpdate(userId int64, body string) error {
	_, err := NewUpsate(userId, body)
	return err
}
