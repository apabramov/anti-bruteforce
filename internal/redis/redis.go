package internalredis

import (
	"context"
	"errors"
	"net"
	"strconv"
	"time"

	"github.com/apabramov/anti-bruteforce/internal/config"
	rds "github.com/redis/go-redis/v9"
)

var ErrKeyNotFound = errors.New("key not found")

type Cache interface {
	Incr(ctx context.Context, key string, t time.Duration) error
	Get(ctx context.Context, key string) (int, error)
	Del(ctx context.Context, key string) error
	Close() error
}

type RedisClient struct {
	Client *rds.Client
}

func New(conf config.RedisConf) (*RedisClient, error) {
	rdb := rds.NewClient(&rds.Options{
		Addr:     net.JoinHostPort(conf.Host, conf.Port),
		Password: conf.Pass,
		DB:       0,
	})

	return &RedisClient{Client: rdb}, nil
}

func (c *RedisClient) Incr(ctx context.Context, key string, t time.Duration) error {
	err := c.Client.Incr(ctx, key).Err()
	if err != nil {
		return err
	}

	return c.Client.Expire(ctx, key, t).Err()
}

func (c *RedisClient) Get(ctx context.Context, key string) (int, error) {
	val, err := c.Client.Get(ctx, key).Result()
	if err != nil {
		return 0, ErrKeyNotFound
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return intVal, nil
}

func (c *RedisClient) Del(ctx context.Context, key string) error {
	return c.Client.Del(ctx, key).Err()
}

func (c *RedisClient) Close() error {
	return c.Client.Close()
}
