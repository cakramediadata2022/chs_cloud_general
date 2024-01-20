package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache Declaration
var DataCache CacheItf

var ctxB = context.Background()

type CacheItf interface {
	Set(ctx context.Context, companyCode string, key string, data interface{}, expiration time.Duration) error
	SetNX(ctx context.Context, companyCode string, key string, data interface{}, expiration time.Duration) error
	GetJSON(ctx context.Context, companyCode string, key string) (interface{}, error)
	Get(ctx context.Context, companyCode string, key string) ([]byte, error)
	GetString(ctx context.Context, companyCode string, key string) (string, error)
	Del(ctx context.Context, companyCode string, key string) error
}

type AppCache struct {
	client *redis.Client
}

var RedisClient *redis.Client

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:49153",
		Password: "redispw",
		DB:       0,
	})

	_, err := RedisClient.Ping(ctxB).Result()
	if err != nil {
		panic(err)
	}
}

func (r *AppCache) Set(ctx context.Context, companyCode string, key string, data interface{}, expiration time.Duration) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = r.client.Set(ctx, companyCode+"_"+key, b, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *AppCache) SetNX(ctx context.Context, companyCode string, key string, data interface{}, expiration time.Duration) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	r.client.SetNX(ctx, companyCode+"_"+key, b, expiration)
	return nil
}

func (r *AppCache) Get(ctx context.Context, companyCode string, key string) ([]byte, error) {
	res, err := r.client.Get(ctx, companyCode+"_"+key).Result()
	if err == redis.Nil {
		return nil, errors.New(companyCode + "_" + key + " does not exist")
	} else if err != nil {
		return nil, err
	}

	return []byte(res), nil
}

func (r *AppCache) GetString(ctx context.Context, companyCode string, key string) (string, error) {
	var Data string
	res, err := r.client.Get(ctx, companyCode+"_"+key).Result()
	if err == redis.Nil {
		return "", errors.New(companyCode + "_" + key + " does not exist")
	} else if err != nil {
		return "", err
	}
	if err := json.Unmarshal([]byte(res), &Data); err != nil {
		return "", err
	}

	return Data, nil
}

func (r *AppCache) GetJSON(ctx context.Context, companyCode string, key string) (DataOutput interface{}, err error) {
	res, err := r.client.Get(ctx, companyCode+"_"+key).Result()
	if err == redis.Nil {
		return nil, errors.New(companyCode + "_" + key + " does not exist")
	} else if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(res), &DataOutput); err != nil {
		return nil, err
	}
	return DataOutput, nil

}

func (r *AppCache) Del(ctx context.Context, companyCode string, key string) error {
	_, err := r.client.Del(ctx, companyCode+"_"+key).Result()
	if err == redis.Nil {
		return errors.New(companyCode + "_" + key + " does not exist")
	} else if err != nil {
		return err
	}

	return nil
}

func InitCache(client *redis.Client) {
	DataCache = &AppCache{
		client: client,
	}
}
