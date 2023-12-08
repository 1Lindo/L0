package repo

import (
	"L0/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type ICacheRepo interface {
	InsertOrder(ctx context.Context, order *models.EventOrder) error
	GetOrderByID(ctx context.Context, id string) (order *models.EventOrder, err error)
	GetOrders(ctx context.Context) ([]models.EventOrder, error)
}

type cacheRepo struct {
	cache  *redis.Client
	repoDB IDBRepo
}

func (c *cacheRepo) GetOrders(ctx context.Context) ([]models.EventOrder, error) {
	orders, err := c.repoDB.GetOrders(ctx)
	if err != nil {
		return nil, err
	}
	for _, order := range orders {
		orderMf, err := json.Marshal(order)
		if err != nil {
			return nil, err
		}
		_, err = c.cache.Set(ctx, fmt.Sprintf("Order:%v", order.ID), orderMf, 0).Result()
		if err != nil {
			return nil, err
		}
	}
	return orders, nil
}
func (c *cacheRepo) InsertOrder(ctx context.Context, order *models.EventOrder) error {
	orderMf, err := json.Marshal(order)
	if err != nil {
		return c.repoDB.InsertOrder(ctx, order)
	}
	orderM, err := c.cache.Set(ctx, fmt.Sprintf("Order:%v", order.ID), orderMf, 0).Result()
	if err != nil {
		fmt.Printf("Failed to add ORDER <> %v key-value pair", orderM)
		return c.repoDB.InsertOrder(ctx, order)
	}
	return c.repoDB.InsertOrder(ctx, order)
}

func (c *cacheRepo) GetOrderByID(ctx context.Context, ID string) (order *models.EventOrder, err error) {

	value, err := c.cache.Get(ctx, fmt.Sprintf("Order:%v", ID)).Result()

	if err != nil {
		return c.repoDB.GetOrderByID(ctx, ID)
	}
	if err := json.Unmarshal([]byte(value), &order); err != nil {
		return c.repoDB.GetOrderByID(ctx, ID)
	}
	return
}
