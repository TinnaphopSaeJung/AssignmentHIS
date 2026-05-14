package service

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	maxLoginFailedAttempts = 10
	loginLockDuration      = 15 * time.Minute
	loginFailedTTL         = 15 * time.Minute
)

type LoginAttemptService struct {
	redisClient *redis.Client
}

func NewLoginAttemptService(redisClient *redis.Client) *LoginAttemptService {
	return &LoginAttemptService{
		redisClient: redisClient,
	}
}

func (s *LoginAttemptService) failedKey(username string) string {
	return fmt.Sprintf("login_failed:%s", username)
}

func (s *LoginAttemptService) lockedKey(username string) string {
	return fmt.Sprintf("login_locked:%s", username)
}

func (s *LoginAttemptService) IsLocked(ctx context.Context, username string) (bool, time.Duration, error) {
	key := s.lockedKey(username)

	ttl, err := s.redisClient.TTL(ctx, key).Result()
	if err != nil {
		return false, 0, err
	}

	if ttl > 0 {
		return true, ttl, nil
	}

	return false, 0, nil
}

func (s *LoginAttemptService) RecordFailedAttempt(ctx context.Context, username string) (int64, error) {
	key := s.failedKey(username)

	count, err := s.redisClient.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	if count == 1 {
		if err := s.redisClient.Expire(ctx, key, loginFailedTTL).Err(); err != nil {
			return 0, err
		}
	}

	if count >= maxLoginFailedAttempts {
		lockKey := s.lockedKey(username)

		if err := s.redisClient.Set(ctx, lockKey, "1", loginLockDuration).Err(); err != nil {
			return 0, err
		}
	}

	return count, nil
}

func (s *LoginAttemptService) Reset(ctx context.Context, username string) error {
	failedKey := s.failedKey(username)
	lockedKey := s.lockedKey(username)

	return s.redisClient.Del(ctx, failedKey, lockedKey).Err()
}
