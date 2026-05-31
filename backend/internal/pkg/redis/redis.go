package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/bryan/user-system/internal/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisClient struct {
	client *redis.Client
	logger *zap.Logger
}

func NewRedisClient(cfg *config.RedisConfig, logger *zap.Logger) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return &RedisClient{
		client: client,
		logger: logger,
	}
}

func (s *RedisClient) GetClient() *redis.Client {
	return s.client
}

func (s *RedisClient) Set(ctx context.Context, key, value string) bool {
	err := s.client.Set(ctx, key, value, 0).Err()
	if err != nil {
		s.logger.Error("Redis set 操作失败", zap.String("key", key), zap.Error(err))
		return false
	}
	return true
}

func (s *RedisClient) SetWithExpire(ctx context.Context, key, value string, seconds int64) bool {
	if seconds <= 0 {
		return s.Set(ctx, key, value)
	}
	err := s.client.Set(ctx, key, value, time.Duration(seconds)*time.Second).Err()
	if err != nil {
		s.logger.Error("Redis setWithExpire 操作失败", zap.String("key", key), zap.Error(err))
		return false
	}
	return true
}

func (s *RedisClient) SetExpire(ctx context.Context, key string, seconds int64) bool {
	if seconds <= 0 {
		s.logger.Warn("Redis setExpire: 过期时间必须大于0", zap.String("key", key))
		return false
	}
	ok, err := s.client.Expire(ctx, key, time.Duration(seconds)*time.Second).Result()
	if err != nil {
		s.logger.Error("Redis setExpire 操作失败", zap.String("key", key), zap.Error(err))
		return false
	}
	return ok
}

func (s *RedisClient) Get(ctx context.Context, key string) string {
	val, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return ""
	}
	if err != nil {
		s.logger.Error("Redis get 操作失败", zap.String("key", key), zap.Error(err))
		return ""
	}
	return val
}

func (s *RedisClient) Delete(ctx context.Context, key string) bool {
	err := s.client.Del(ctx, key).Err()
	if err != nil {
		s.logger.Error("Redis delete 操作失败", zap.String("key", key), zap.Error(err))
		return false
	}
	return true
}

func (s *RedisClient) HasKey(ctx context.Context, key string) bool {
	ok, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		s.logger.Error("Redis hasKey 操作失败", zap.String("key", key), zap.Error(err))
		return false
	}
	return ok > 0
}

func (s *RedisClient) Increment(ctx context.Context, key string, delta int64) (int64, error) {
	val, err := s.client.IncrBy(ctx, key, delta).Result()
	if err != nil {
		s.logger.Error("Redis increment 操作失败", zap.String("key", key), zap.Error(err))
		return 0, fmt.Errorf("Redis increment 失败: %w", err)
	}
	return val, nil
}

func (s *RedisClient) Decrement(ctx context.Context, key string, delta int64) (int64, error) {
	val, err := s.client.DecrBy(ctx, key, delta).Result()
	if err != nil {
		s.logger.Error("Redis decrement 操作失败", zap.String("key", key), zap.Error(err))
		return 0, fmt.Errorf("Redis decrement 失败: %w", err)
	}
	return val, nil
}
