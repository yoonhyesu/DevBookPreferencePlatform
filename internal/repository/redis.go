package repository

import (
	"context"
	"log"
	"time"
)

// Redis 토큰 관련
func (r *CommonRepo) SaveTokenMetaData(ctx context.Context, userID, tokenUUID string, exp time.Duration) error {
	if err := r.redis.Client.Ping(ctx).Err(); err != nil {
		log.Printf("Redis 연결 확인 실패: %v", err)
		return err
	}

	err := r.redis.Client.Set(ctx, tokenUUID, userID, exp).Err()
	if err != nil {
		log.Printf("Redis 토큰 저장 실패: %v", err)
		return err
	}

	return nil
}

func (r *CommonRepo) DeleteTokenMetaData(ctx context.Context, tokenUUID string) error {
	return r.redis.Client.Del(ctx, tokenUUID).Err()
}

// Redis 채팅 관련
func (r *CommonRepo) SaveChatMessageToRedis(ctx context.Context, bookID string, message string) error {
	key := "chat:" + bookID
	err := r.redis.Client.RPush(ctx, key, message).Err()
	if err != nil {
		log.Printf("Redis 메시지 저장 실패: %v", err)
		return err
	}
	return r.redis.Client.LTrim(ctx, key, -100, -1).Err()
}

func (r *CommonRepo) GetRecentChatMessage(ctx context.Context, bookID string) ([]string, error) {
	key := "chat:" + bookID
	messages, err := r.redis.Client.LRange(ctx, key, 0, 99).Result()
	if err != nil {
		log.Printf("Redis 메시지 조회 실패: %v", err)
		return nil, err
	}
	return messages, nil
}
