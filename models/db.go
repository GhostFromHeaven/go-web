package models

import "github.com/go-redis/redis"

var client *redis.Client

func init() {
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}
