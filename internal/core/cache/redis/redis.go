package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/fedotovmax/medods-test/internal/core/cache"
	"github.com/redis/go-redis/v9"
)

type Redis interface {
	Set(ctx context.Context, key string, value any, exp time.Duration) error
	SetIfNotExist(ctx context.Context, key string, value any, exp time.Duration) error
	Delete(ctx context.Context, keys ...string) error
	Get(ctx context.Context, key string) (string, error)
	GetInt64(ctx context.Context, key string) (int64, error)
	GetFloat64(ctx context.Context, key string) (float64, error)
	GetBool(ctx context.Context, key string) (bool, error)
	GetJSON(ctx context.Context, key string, dest any) error
	Stop(ctx context.Context) error
}

type Pool struct {
	*redis.Client
	log *slog.Logger
}

func New(ctx context.Context, config Config, log *slog.Logger) (*Pool, error) {

	const op = "core.cache.redis.New"

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("%s: error when validate config: %w", op, err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:            config.Addr,
		Password:        config.Password,
		DB:              config.DB,
		MaxRetries:      config.MaxRetries,
		MinRetryBackoff: config.MinRetryBackoff,
		MaxRetryBackoff: config.MaxRetryBackoff,
		PoolSize:        config.PoolSize,
		MaxIdleConns:    config.MaxIdleConns,
		ConnMaxLifetime: config.MaxConnLifetime,
		ConnMaxIdleTime: config.MaxIdleConnLifetime,
	})

	_, err := redisClient.Ping(ctx).Result()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Pool{
		Client: redisClient,
		log:    log,
	}, nil

}

func (p *Pool) Set(ctx context.Context, key string, value any, exp time.Duration) error {

	err := p.Client.Set(ctx, key, value, exp).Err()

	if err != nil {
		return err
	}

	return nil

}

func (p *Pool) SetIfNotExist(ctx context.Context, key string, value any, exp time.Duration) error {

	err := p.Client.SetArgs(ctx, key, value, redis.SetArgs{
		Mode: "NX",
		TTL:  exp,
	}).Err()

	if err != nil {
		if err == redis.Nil {
			return cache.ErrKeyExists
		}
		return err
	}

	return nil
}

func (p *Pool) Get(ctx context.Context, key string) (string, error) {

	value, err := p.Client.Get(ctx, key).Result()

	if err != nil {
		if err == redis.Nil {
			return "", cache.ErrKeyNotExists
		}
		return "", err
	}
	return value, nil
}

func (p *Pool) GetInt64(ctx context.Context, key string) (int64, error) {

	value, err := p.Client.Get(ctx, key).Int64()

	if err != nil {
		if err == redis.Nil {
			return 0, cache.ErrKeyNotExists
		}
		return 0, err
	}
	return value, nil
}

func (p *Pool) GetBool(ctx context.Context, key string) (bool, error) {

	value, err := p.Client.Get(ctx, key).Bool()

	if err != nil {
		if err == redis.Nil {
			return false, cache.ErrKeyNotExists
		}
		return false, err
	}
	return value, nil
}

func (p *Pool) GetFloat64(ctx context.Context, key string) (float64, error) {

	value, err := p.Client.Get(ctx, key).Float64()

	if err != nil {
		if err == redis.Nil {
			return 0, cache.ErrKeyNotExists
		}
		return 0, err
	}
	return value, nil
}

func (p *Pool) GetJSON(ctx context.Context, key string, dest any) error {

	value, err := p.Client.Get(ctx, key).Bytes()

	if err != nil {
		if err == redis.Nil {
			return cache.ErrKeyNotExists
		}
		return err
	}

	err = json.Unmarshal(value, dest)

	if err != nil {
		return err
	}

	return nil
}

func (p *Pool) Delete(ctx context.Context, keys ...string) error {

	if err := p.Client.Del(ctx, keys...).Err(); err != nil {
		return err
	}

	return nil

}

func (r *Pool) Stop(ctx context.Context) error {

	const op = "core.cache.redis.Pool.Stop"

	done := make(chan error, 1)

	go func() {
		err := r.Client.Close()
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("%s: %w", op, ctx.Err())
	}
}
