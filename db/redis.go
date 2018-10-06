package db

import (
	"context"
	"time"

	radix "github.com/fzzy/radix/redis"

	//"git.meiqia.com/livechat/common/log"
	"github.com/chenjun-git/umbrella-common/redis/pool"

	"business/user/common"
)


var RedisPool *pool.Pool

func InitRedis(config *common.RedisConfig) {

	df := func(network, addr string) (*radix.Client, error) {
		client, err := radix.DialTimeout(network, addr, config.Timeout.D())
		if err != nil {
			return nil, err
		}
		if config.Password != "" {
			if reply := client.Cmd("AUTH", config.Password); reply.Err != nil {
				return nil, reply.Err
			}
		}
		if reply := client.Cmd("SELECT", config.DB); reply.Err != nil {
			return nil, reply.Err
		}
		return client, nil
	}

	lazyPool := pool.NewCustomLazyPool("tcp", config.Addr, config.PoolSize, df)
	pool, err := lazyPool.GetPool()
	if err != nil {
		//log.Fatal("get redis pool failed.", zap.Error(err))
	}

	RedisPool = pool
}

func RedisDel(ctx context.Context, keys ...string) error {
	args := make([]interface{}, len(keys))
	for i, key := range keys {
		args[i] = key
	}
	_, err := RedisCmd(ctx, "del", args...)
	return err
}

func RedisSet(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if expiration < 0 {
		reply, err := RedisCmd(ctx, "set", key, value)
		return checkErr(reply, err)
	}
	if expiration < time.Second {
		reply, err := RedisCmd(ctx, "set", key, value, "px", int64(expiration/time.Millisecond))
		return checkErr(reply, err)
	}
	reply, err := RedisCmd(ctx, "set", key, value, "ex", int64(expiration/time.Second))
	return checkErr(reply, err)
}

func RedisGet(ctx context.Context, key string) (*radix.Reply, error) {
	return RedisCmd(ctx, "get", key)
}

func RedisHdel(ctx context.Context, key string, fields ...string) error {
	args := make([]interface{}, len(fields)+1)
	args[0] = key
	for i, field := range fields {
		args[i+1] = field
	}
	reply, err := RedisCmd(ctx, "hdel", args...)
	return checkErr(reply, err)
}

func RedisHset(ctx context.Context, key, field string, value interface{}) error {
	reply, err := RedisCmd(ctx, "hset", key, field, value)
	return checkErr(reply, err)
}

func RedisEX(ctx context.Context, key string, expiration time.Duration) error {
	if expiration < 0 {
		return nil
	}
	if expiration < time.Second {
		reply, err := RedisCmd(ctx, "pexpire", key, int64(expiration/time.Millisecond))
		return checkErr(reply, err)
	}
	reply, err := RedisCmd(ctx, "expire", key, int64(expiration/time.Second))
	return checkErr(reply, err)
}

func RedisHget(ctx context.Context, key, field string) (*radix.Reply, error) {
	return RedisCmd(ctx, "hget", key, field)
}

func RedisHlen(ctx context.Context, key string) (*radix.Reply, error) {
	return RedisCmd(ctx, "hlen", key)
}

func RedisHgetall(ctx context.Context, key string) (*radix.Reply, error) {
	return RedisCmd(ctx, "hgetall", key)
}

func RedisCmd(ctx context.Context, cmd string, args ...interface{}) (*radix.Reply, error) {
	client, err := RedisPool.Get()
	if err != nil {
		return nil, err
	}
	defer RedisPool.CarefullyPut(client, &err)

	reply := client.Cmd(ctx, cmd, args...)
	err = reply.Err
	return reply, err
}

func checkErr(redisReply *radix.Reply, err error) error {
	if err != nil {
		return err
	}
	if redisReply == nil {
		return nil
	}
	return redisReply.Err
}