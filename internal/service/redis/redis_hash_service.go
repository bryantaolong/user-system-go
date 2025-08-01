package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type RedisHashService struct {
	client *redis.Client
	logger *logrus.Logger
}

func NewRedisHashService(client *redis.Client, logger *logrus.Logger) *RedisHashService {
	return &RedisHashService{
		client: client,
		logger: logger,
	}
}

// Set 在 Redis 的 Hash 中存储多个键值对
func (s *RedisHashService) Set(ctx context.Context, key string, value map[string]interface{}) error {
	err := s.client.HSet(ctx, key, value).Err()
	if err != nil {
		s.logger.WithError(err).Errorf("Redis hSetAll 操作失败，key: %s, value: %v", key, value)
		return fmt.Errorf("redis hSetAll failed: %w", err)
	}
	return nil
}

// SetField 在 Redis 的 Hash 中存储一个键值对
func (s *RedisHashService) SetField(ctx context.Context, key, field string, value interface{}) error {
	err := s.client.HSet(ctx, key, field, value).Err()
	if err != nil {
		s.logger.WithError(err).Errorf("Redis hSet 操作失败，key: %s, field: %s, value: %v", key, field, value)
		return fmt.Errorf("redis hSet failed: %w", err)
	}
	return nil
}

// Get 在 Redis 的 Hash 中获取某个字段的值
func (s *RedisHashService) Get(ctx context.Context, key, field string) (string, error) {
	val, err := s.client.HGet(ctx, key, field).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	if err != nil {
		s.logger.WithError(err).Errorf("Redis hGet 操作失败，key: %s, field: %s", key, field)
		return "", fmt.Errorf("redis hGet failed: %w", err)
	}
	return val, nil
}

// Delete 从 Redis 的 Hash 中删除一个字段
func (s *RedisHashService) Delete(ctx context.Context, key, field string) (bool, error) {
	result, err := s.client.HDel(ctx, key, field).Result()
	if err != nil {
		s.logger.WithError(err).Errorf("Redis hDelete 操作失败，key: %s, field: %s", key, field)
		return false, fmt.Errorf("redis hDel failed: %w", err)
	}
	return result > 0, nil
}

// Exists 检查 Redis 的 Hash 中是否存在某个字段
func (s *RedisHashService) Exists(ctx context.Context, key, field string) (bool, error) {
	exists, err := s.client.HExists(ctx, key, field).Result()
	if err != nil {
		s.logger.WithError(err).Errorf("Redis hExists 操作失败，key: %s, field: %s", key, field)
		return false, fmt.Errorf("redis hExists failed: %w", err)
	}
	return exists, nil
}

// GetAll 获取 Redis 的 Hash 中所有字段和对应值
func (s *RedisHashService) GetAll(ctx context.Context, key string) (map[string]string, error) {
	result, err := s.client.HGetAll(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return map[string]string{}, nil
	}
	if err != nil {
		s.logger.WithError(err).Errorf("Redis hGetAll 操作失败，key: %s", key)
		return nil, fmt.Errorf("redis hGetAll failed: %w", err)
	}
	return result, nil
}

// Keys 获取 Redis 的 Hash 中所有字段
func (s *RedisHashService) Keys(ctx context.Context, key string) ([]string, error) {
	result, err := s.client.HKeys(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return []string{}, nil
	}
	if err != nil {
		s.logger.WithError(err).Errorf("Redis hKeys 操作失败，key: %s", key)
		return nil, fmt.Errorf("redis hKeys failed: %w", err)
	}
	return result, nil
}

// Values 获取 Redis 的 Hash 中所有字段的值
func (s *RedisHashService) Values(ctx context.Context, key string) ([]string, error) {
	result, err := s.client.HVals(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return []string{}, nil
	}
	if err != nil {
		s.logger.WithError(err).Errorf("Redis hVals 操作失败，key: %s", key)
		return nil, fmt.Errorf("redis hVals failed: %w", err)
	}
	return result, nil
}

// Size 获取 Redis 的 Hash 中字段数量
func (s *RedisHashService) Size(ctx context.Context, key string) (int64, error) {
	result, err := s.client.HLen(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return 0, nil
	}
	if err != nil {
		s.logger.WithError(err).Errorf("Redis hLen 操作失败，key: %s", key)
		return 0, fmt.Errorf("redis hLen failed: %w", err)
	}
	return result, nil
}
