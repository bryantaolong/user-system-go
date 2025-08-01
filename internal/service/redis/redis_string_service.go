package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type RedisStringService struct {
	client *redis.Client
	logger *logrus.Logger
}

func NewRedisStringService(client *redis.Client, logger *logrus.Logger) *RedisStringService {
	return &RedisStringService{
		client: client,
		logger: logger,
	}
}

// Set 在 Redis 中存储一个键值对（无过期时间）
func (s *RedisStringService) Set(ctx context.Context, key, value string) error {
	err := s.client.Set(ctx, key, value, 0).Err()
	if err != nil {
		s.logger.WithError(err).Errorf("Redis set 操作失败，key: %s, value: %s", key, value)
		return fmt.Errorf("redis set failed: %w", err)
	}
	return nil
}

// SetWithExpire 在 Redis 中存储一个带有过期时间的键值对
func (s *RedisStringService) SetWithExpire(ctx context.Context, key, value string, expiration time.Duration) error {
	if expiration <= 0 {
		s.logger.Warnf("Redis setWithExpire 操作：过期时间必须大于0，key: %s", key)
		return s.Set(ctx, key, value)
	}

	err := s.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		s.logger.WithError(err).Errorf("Redis setWithExpire 操作失败，key: %s, value: %s, expiration: %v", key, value, expiration)
		return fmt.Errorf("redis setWithExpire failed: %w", err)
	}
	return nil
}

// SetExpire 为 Redis 中已存在的键设置过期时间
func (s *RedisStringService) SetExpire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	if expiration <= 0 {
		s.logger.Warnf("Redis setExpire 操作：过期时间必须大于0，key: %s", key)
		return false, nil
	}

	result, err := s.client.Expire(ctx, key, expiration).Result()
	if err != nil {
		s.logger.WithError(err).Errorf("Redis setExpire 操作失败，key: %s, expiration: %v", key, expiration)
		return false, fmt.Errorf("redis expire failed: %w", err)
	}
	return result, nil
}

// Get 从 Redis 中获取与键对应的值
func (s *RedisStringService) Get(ctx context.Context, key string) (string, error) {
	val, err := s.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	if err != nil {
		s.logger.WithError(err).Errorf("Redis get 操作失败，key: %s", key)
		return "", fmt.Errorf("redis get failed: %w", err)
	}
	return val, nil
}

// Delete 从 Redis 中删除一个键
func (s *RedisStringService) Delete(ctx context.Context, key string) (bool, error) {
	result, err := s.client.Del(ctx, key).Result()
	if err != nil {
		s.logger.WithError(err).Errorf("Redis delete 操作失败，key: %s", key)
		return false, fmt.Errorf("redis del failed: %w", err)
	}
	return result > 0, nil
}

// Exists 检查 Redis 中是否存在某个键
func (s *RedisStringService) Exists(ctx context.Context, key string) (bool, error) {
	result, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		s.logger.WithError(err).Errorf("Redis exists 操作失败，key: %s", key)
		return false, fmt.Errorf("redis exists failed: %w", err)
	}
	return result > 0, nil
}

// Increment 对存储在指定键上的字符串值执行原子递增操作
func (s *RedisStringService) Increment(ctx context.Context, key string, delta int64) (int64, error) {
	result, err := s.client.IncrBy(ctx, key, delta).Result()
	if err != nil {
		s.logger.WithError(err).Errorf("Redis increment 操作失败，key: %s, delta: %d", key, delta)
		return 0, fmt.Errorf("redis incrBy failed: %w", err)
	}
	return result, nil
}

// Decrement 对存储在指定键上的字符串值执行原子递减操作
func (s *RedisStringService) Decrement(ctx context.Context, key string, delta int64) (int64, error) {
	result, err := s.client.DecrBy(ctx, key, delta).Result()
	if err != nil {
		s.logger.WithError(err).Errorf("Redis decrement 操作失败，key: %s, delta: %d", key, delta)
		return 0, fmt.Errorf("redis decrBy failed: %w", err)
	}
	return result, nil
}
