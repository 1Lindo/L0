package repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

func NewDBRepo(db *sqlx.DB) IDBRepo {
	return &dbRepo{db: db}
}

func NewCacheRepo(cache *redis.Client, repoDB IDBRepo) ICacheRepo {
	return &cacheRepo{cache: cache, repoDB: repoDB}
}

func NewCache(cache *redis.Client, repoDB IDBRepo) ICacheRepo {
	return &cacheRepo{repoDB: repoDB, cache: cache}
}
