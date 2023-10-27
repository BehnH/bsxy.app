package database

import (
	"bsxy_app/logger"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	Ctx = context.TODO()
)

var RedisClient *redis.Client

func GetCachedUserProfile(handle string) (*bsky.ActorDefs_ProfileViewDetailed, error) {
	var user bsky.ActorDefs_ProfileViewDetailed
	if err := RedisClient.Ping(Ctx).Err(); err != nil {
		logger.SLog().Infof("Client dead")
		return nil, err
	}
	result, err := RedisClient.Get(Ctx, handle).Result()
	if errors.Is(err, redis.Nil) {
		logger.SLog().Infof("cache miss - key %s does not exist", handle)
		return nil, fmt.Errorf("cache miss for %s", handle)
	}
	if err != nil {
		logger.SLog().Errorf("error while looking for %s: %w", handle, err)
		return nil, err
	}

	err = json.Unmarshal([]byte(result), &user)
	if err != nil {
		logger.SLog().Errorf("failed to unwrap json to profile type: %w", err)
		return nil, err
	}

	return &user, nil
}

func GetCachedPost(post string) (*bsky.FeedGetPostThread_Output, error) {
	var user bsky.FeedGetPostThread_Output
	if err := RedisClient.Ping(Ctx).Err(); err != nil {
		logger.SLog().Infof("Client dead")
		return nil, err
	}
	result, err := RedisClient.Get(Ctx, post).Result()
	if errors.Is(err, redis.Nil) {
		logger.SLog().Infof("cache miss - key %s does not exist", post)
		return nil, fmt.Errorf("cache miss for %s", post)
	}
	if err != nil {
		logger.SLog().Errorf("error while looking for %s: %w", post, err)
		return nil, err
	}

	err = json.Unmarshal([]byte(result), &user)
	if err != nil {
		logger.SLog().Errorf("failed to unwrap json to profile type: %w", err)
		return nil, err
	}

	return &user, nil
}

func AddProfileToCache(profile bsky.ActorDefs_ProfileViewDetailed) error {
	if err := RedisClient.Ping(Ctx).Err(); err != nil {
		logger.SLog().Infof("Client ded")
		return err
	}
	out, err := json.Marshal(profile)
	if err != nil {
		return err
	}
	_, err = RedisClient.Set(Ctx, profile.Handle, out, time.Second*60*5).Result()
	if err != nil {
		return err
	}
	return nil
}

func AddPostToCache() error {
	return nil
}
