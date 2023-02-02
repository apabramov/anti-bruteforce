package bucket

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/apabramov/anti-bruteforce/internal/config"
	internalredis "github.com/apabramov/anti-bruteforce/internal/redis"
	"github.com/apabramov/anti-bruteforce/internal/storage"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

var client internalredis.RedisClient

func TestMain(m *testing.M) {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	c := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	client = internalredis.RedisClient{Client: c}

	code := m.Run()
	os.Exit(code)
}

func TestBucket(t *testing.T) {
	buc := New(&client, config.LimitConf{
		LimitLogin: 1,
		LimitPass:  10,
		LimitIP:    100,
	})

	authorize := storage.Authorize{
		Login:    "login",
		Password: "password",
		IP:       "127.0.0.1",
	}

	ctx := context.Background()

	t.Run("check limit", func(t *testing.T) {
		b, err := buc.CheckLimit(ctx, authorize)
		require.NoError(t, err)
		require.True(t, b)

		b, err = buc.CheckLimit(ctx, authorize)
		require.NoError(t, err)
		require.False(t, b)
	})
}
