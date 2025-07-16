package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"orders-api/model"
)

type RedisRepository struct {
	Redis *redis.Client
}

func orderIDKey(id uint64) string { // function 2 find the order id key
	return fmt.Sprintf("orders:id:%d", id)
}

func (r *RedisRepository) Insert(ctx context.Context, order model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	key := orderIDKey(order.OrderId)

	txn := r.Redis.TxPipeline()

	res := txn.SetNX(ctx, key, string(data), 0)
	if res.Err() != nil {
		txn.Discard()
		return res.Err()
	}

	if err := txn.SAdd(ctx, "orders", key).Err(); err != nil { // store all the order keys
		txn.Discard()
		return err
	}

	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("txn exec: %w", err)
	}
	return nil
}

func (r *RedisRepository) Find(ctx context.Context, id uint64) (model.Order, error) {
	key := orderIDKey(id)
	res := r.Redis.Get(ctx, key)
	if res.Err() != nil {
		return model.Order{}, res.Err()
	}

	var order model.Order
	if err := json.Unmarshal([]byte(res.Val()), &order); err != nil {
		return model.Order{}, fmt.Errorf("json unmarshal: %w", err)
	}
	return order, nil
}

func (r *RedisRepository) Delete(ctx context.Context, id uint64) error {
	key := orderIDKey(id)

	txn := r.Redis.TxPipeline()

	res := txn.Del(ctx, key)
	if errors.Is(res.Err(), redis.Nil) {
		txn.Discard()
		return fmt.Errorf("Order Not Exist")
	} else if errors.Is(res.Err(), redis.Nil) {
		txn.Discard()
		return res.Err()
	}

	if err := txn.SRem(ctx, "orders", key).Err(); err != nil {
		txn.Discard()
		return err
	}

	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("txn exec: %w", err)
	}
	return nil
}

func (r *RedisRepository) Update(ctx context.Context, order model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	key := orderIDKey(order.OrderId)
	res := r.Redis.SetXX(ctx, key, string(data), 0)
	if errors.Is(res.Err(), redis.Nil) {
		return fmt.Errorf("order not found")
	} else if err != nil {
		return res.Err()
	}
	return nil
}

type FindAllPage struct {
	Size   uint
	Offset uint64
}

type FindResult struct {
	Orders []model.Order
	Cursor uint64
}

func (r *RedisRepository) FindAll(ctx context.Context, page FindAllPage) (FindResult, error) {
	res := r.Redis.SScan(ctx, "orders", page.Offset, "*", int64(page.Size))
	keys, cursor, err := res.Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("Failed to get order ids from redis: %w", err)
	}

	if len(keys) == 0 {
		return FindResult{}, nil
	}

	xs, err := r.Redis.MGet(ctx, keys...).Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("Failed to get Orders: %w", err)
	}

	orders := make([]model.Order, len(xs))
	for i, x := range xs {
		x := x.(string)
		var order model.Order
		if err := json.Unmarshal([]byte(x), &order); err != nil {
			return FindResult{}, fmt.Errorf("Failed to unmarshal order: %w", err)
		}

		orders[i] = order
	}

	return FindResult{orders, cursor}, nil
}
