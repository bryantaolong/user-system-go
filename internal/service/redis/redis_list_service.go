package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type RedisListService struct {
	client *redis.Client
	logger *logrus.Logger
}

func NewRedisListService(client *redis.Client, logger *logrus.Logger) *RedisListService {
	return &RedisListService{
		client: client,
		logger: logger,
	}
}

// LeftPush 向 Redis 列表左侧添加一个元素
func (s *RedisListService) LeftPush(ctx context.Context, key string, value interface{}) error {
	err := s.client.LPush(ctx, key, value).Err()
	if err != nil {
		s.logger.WithError(err).Errorf("Redis leftPush 操作失败，key: %s, value: %v", key, value)
		return fmt.Errorf("redis lPush failed: %w", err)
	}
	return nil
}

// LeftPushAll 向 Redis 列表左侧添加多个元素
func (s *RedisListService) LeftPushAll(ctx context.Context, key string, values ...interface{}) error {
	err := s.client.LPush(ctx, key, values...).Err()
	if err != nil {
		s.logger.WithError(err).Errorf("Redis leftPushAll 操作失败，key: %s, values: %v", key, values)
		return fmt.Errorf("redis lPushAll failed: %w", err)
	}
	return nil
}

// RightPush 向 Redis 列表右侧添加一个元素
func (s *RedisListService) RightPush(ctx context.Context, key string, value interface{}) error {
	err := s.client.RPush(ctx, key, value).Err()
	if err != nil {
		s.logger.WithError(err).Errorf("Redis rightPush 操作失败，key: %s, value: %v", key, value)
		return fmt.Errorf("redis rPush failed: %w", err)
	}
	return nil
}

// RightPushAll 向 Redis 列表右侧添加多个元素
func (s *RedisListService) RightPushAll(ctx context.Context, key string, values ...interface{}) error {
	err := s.client.RPush(ctx, key, values...).Err()
	if err != nil {
		s.logger.WithError(err).Errorf("Redis rightPushAll 操作失败，key: %s, values: %v", key, values)
		return fmt.Errorf("redis rPushAll failed: %w", err)
	}
	return nil
}

// LeftPop 从 Redis 列表左侧弹出元素
func (s *RedisListService) LeftPop(ctx context.Context, key string) (string, error) {
	val, err := s.client.LPop(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	if err != nil {
		s.logger.WithError(err).Errorf("Redis leftPop 操作失败，key: %s", key)
		return "", fmt.Errorf("redis lPop failed: %w", err)
	}
	return val, nil
}

// RightPop 从 Redis 列表右侧弹出元素
func (s *RedisListService) RightPop(ctx context.Context, key string) (string, error) {
	val, err := s.client.RPop(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	if err != nil {
		s.logger.WithError(err).Errorf("Redis rightPop 操作失败，key: %s", key)
		return "", fmt.Errorf("redis rPop failed: %w", err)
	}
	return val, nil
}

// Range 获取 Redis 列表中指定范围的元素
func (s *RedisListService) Range(ctx context.Context, key string, start, stop int64) ([]string, error) {
	val, err := s.client.LRange(ctx, key, start, stop).Result()
	if errors.Is(err, redis.Nil) {
		return []string{}, nil
	}
	if err != nil {
		s.logger.WithError(err).Errorf("Redis range 操作失败，key: %s, start: %d, stop: %d", key, start, stop)
		return nil, fmt.Errorf("redis lRange failed: %w", err)
	}
	return val, nil
}

// Length 获取 Redis 列表的长度
func (s *RedisListService) Length(ctx context.Context, key string) (int64, error) {
	val, err := s.client.LLen(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return 0, nil
	}
	if err != nil {
		s.logger.WithError(err).Errorf("Redis length 操作失败，key: %s", key)
		return 0, fmt.Errorf("redis lLen failed: %w", err)
	}
	return val, nil
}

// Index 根据索引从 Redis 列表中获取元素
func (s *RedisListService) Index(ctx context.Context, key string, index int64) (string, error) {
	val, err := s.client.LIndex(ctx, key, index).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	if err != nil {
		s.logger.WithError(err).Errorf("Redis index 操作失败，key: %s, index: %d", key, index)
		return "", fmt.Errorf("redis lIndex failed: %w", err)
	}
	return val, nil
}

// SetByIndex 根据索引更新 Redis 列表中的元素
func (s *RedisListService) SetByIndex(ctx context.Context, key string, index int64, value interface{}) error {
	err := s.client.LSet(ctx, key, index, value).Err()
	if err != nil {
		s.logger.WithError(err).Errorf("Redis setByIndex 操作失败，key: %s, index: %d, value: %v", key, index, value)
		return fmt.Errorf("redis lSet failed: %w", err)
	}
	return nil
}

// Remove 删除 Redis 列表中指定数量的元素
func (s *RedisListService) Remove(ctx context.Context, key string, count int64, value interface{}) (int64, error) {
	result, err := s.client.LRem(ctx, key, count, value).Result()
	if err != nil {
		s.logger.WithError(err).Errorf("Redis remove 操作失败，key: %s, count: %d, value: %v", key, count, value)
		return 0, fmt.Errorf("redis lRem failed: %w", err)
	}
	return result, nil
}

// Trim 对一个列表进行修剪，只保留指定区间的元素
func (s *RedisListService) Trim(ctx context.Context, key string, start, stop int64) error {
	err := s.client.LTrim(ctx, key, start, stop).Err()
	if err != nil {
		s.logger.WithError(err).Errorf("Redis trim 操作失败，key: %s, start: %d, stop: %d", key, start, stop)
		return fmt.Errorf("redis lTrim failed: %w", err)
	}
	return nil
}
