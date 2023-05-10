package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/redis/go-redis/v9"
)

// RedisClient redis client type
type RedisClient struct{ *redis.Client }

var once sync.Once
var redisClient *RedisClient


type Record struct {
	ID int64
	Data Data
}
type Data struct {
	Login   string `json:"login"`
	Password string `json:"password"`
}

func init() {
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))

	if err != nil {
		log.Fatal(err)
	}

	once.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_DB_URL"),
			Password: os.Getenv("REDIS_DB_PASSWORD"),
			DB:       db,
		})

		redisClient = &RedisClient{client}
	})

	ctx := context.Background()
	pong, err := redisClient.Ping(ctx).Result()

	fmt.Println(pong)

	if err != nil {
		log.Fatalf("Could not connect to redis %v", err)
	}
}

func SetKey(ctx context.Context, r Record) error {
	err := redisClient.HSet(ctx, strconv.FormatInt(r.ID, 10), r.Data.Login, r.Data.Password).Err()

	if err != nil {
		return err
	}

	return nil
}

func GetKey(ctx context.Context, r Record) (string, error) {
	val, err := redisClient.HGet(ctx, strconv.FormatInt(r.ID, 10), r.Data.Login).Result()

	if err != nil {
		return "", err
	}
	return val, nil
}

func DeleteKey(ctx context.Context, r Record) error {
	_, err := redisClient.HDel(ctx, strconv.FormatInt(r.ID, 10), r.Data.Login).Result()

	if err != nil {
		return err
	}
	return nil
}
