package dao

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:16379",
	Password: "", // 密码
	DB:       0,  // 数据库
})

const (
	RedisErrNotFound = redis.Nil
)

func Set(ctx context.Context, keys []string, values []string) (err error) {
	pipe := rdb.Pipeline()
	for i, k := range keys {
		pipe.Set(ctx, k, values[i], 0)
	}
	_, err = pipe.Exec(ctx)
	return err
}

func Get(ctx context.Context, key string) error {
	cmd := rdb.Get(ctx, key)
	if cmd.Err() != nil {
		if cmd.Err() == redis.Nil {
			return RedisErrNotFound
		}
		return fmt.Errorf("rdb.Get failed, err: %v", cmd.Err())
	}
	return nil
}
