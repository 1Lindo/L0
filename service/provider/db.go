package provider

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

func ConnectDB(url string) (*sqlx.DB, error) {
	return sqlx.Open("postgres", url+"?sslmode=disable")
}

func ConnectCache(rAdrr string) (*redis.Client, error) {
	return redis.NewClient(&redis.Options{
		Addr: rAdrr,
	}), nil
}
