package db

import (
	"context"
	"fmt"
	"sync"
	"time"

	radix "github.com/fzzy/radix/redis"

	//"git.meiqia.com/livechat/common/log"
	"github.com/chenjun-git/umbrella-common/redis"
	"github.com/chenjun-git/umbrella-common/redis/pool"

	"github.com/chenjun-git/umbrella-user/common"
)


var RedisPool *pool.Pool

var redisClient *pool.LazyPool

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
	redisClient = lazyPool
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

var _ RedisExec = (*RedisClient)(nil)

// RedisExec redis 接口
type RedisExec interface {
	Exists(key string) (bool, error)
	HExists(key, field string) (bool, error)
	FlushAll() error

	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	Del(keys ...string) error

	HSet(key, field string, value interface{}) error
	HGet(key, field string) (string, error)
	HDel(key string, fields ...string) error
	HGetAll(key string) (map[string]string, error)
	HVals(key string) ([]string, error)

	BeginTx() error
	Commit() error
	Discard() error

	GetContext() context.Context

	PutPool()
}

type RedisClient struct {
	pool          *pool.Pool
	ctx           context.Context
	redis         *redis.Client
	isTx          bool
	txMutex       *sync.Mutex
	redisErr      *error
	redisErrMutex *sync.Mutex
}

func NewRedisClient(ctx context.Context) (*RedisClient, error) {
	redisPool, err := redisClient.GetPool()
	if err != nil {
		return nil, common.ExtendErrorStatus(common.AccountRedisError, err)
	}

	redisClient, err := redisPool.Get()
	if err != nil {
		redisPool.CarefullyPut(redisClient, &err)
		return nil, common.ExtendErrorStatus(common.AccountRedisError, err)
	}

	return &RedisClient{
		redisPool,
		ctx,
		redisClient,
		false,
		new(sync.Mutex),
		nil,
		new(sync.Mutex),
	}, nil
}

func NewRedisClientTx(ctx context.Context) (RedisExec, error) {
	redis, err := NewRedisClient(ctx)
	if err != nil {
		return nil, err
	}

	if err := redis.BeginTx(); err != nil {
		return nil, common.ExtendErrorStatus(common.AccountRedisError, err)
	}

	return redis, nil
}

func (r *RedisClient) PutPool() {
	if r != nil && r.pool != nil && r.redis != nil {
		r.pool.CarefullyPut(r.redis, r.redisErr)
	}
}

func (r *RedisClient) checkFailed() error {
	if r.redis == nil {
		return fmt.Errorf("redis failed: maybe NewRedisClient failed")
	}
	return nil
}

func (r *RedisClient) putRedisErr(err *error) {
	if err != nil {
		r.redisErrMutex.Lock()
		r.redisErr = err
		r.redisErrMutex.Unlock()
	}
}

func (r *RedisClient) cmdWithoutTx(cmd string, args ...interface{}) (*radix.Reply, error) {
	if err := r.checkFailed(); err != nil {
		return nil, err
	}

	r.txMutex.Lock()
	defer r.txMutex.Unlock()

	if r.isTx {
		client, err := NewRedisClient(r.GetContext())
		if err != nil {
			return nil, err
		}
		return client.redis.Cmd(r.ctx, cmd, args...), nil
	}

	return r.redis.Cmd(r.ctx, cmd, args...), nil
}

func (r *RedisClient) Exists(key string) (bool, error) {
	if err := r.checkFailed(); err != nil {
		return false, err
	}

	if key == "" {
		return false, nil
	}

	reply, err := r.cmdWithoutTx("exists", key)
	if err != nil {
		return false, err
	}

	exist, err := reply.Bool()
	r.putRedisErr(&err)
	return exist, err
}

// HExists hexists
func (r *RedisClient) HExists(key, field string) (bool, error) {
	if err := r.checkFailed(); err != nil {
		return false, err
	}

	reply, err := r.cmdWithoutTx("hexists", key, field)
	if err != nil {
		return false, err
	}
	exist, err := reply.Bool()
	r.putRedisErr(&err)
	return exist, err
}

// FlushAll flushall
func (r *RedisClient) FlushAll() error {
	if err := r.checkFailed(); err != nil {
		return err
	}

	err := r.redis.Cmd(r.ctx, "flushall").Err
	r.putRedisErr(&err)
	return err
}

// Set set
func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	if err := r.checkFailed(); err != nil {
		return err
	}
	if expiration < time.Second || expiration%time.Second != 0 {
		err := r.redis.Cmd(r.ctx, "set", key, value, "px", int64(expiration/time.Millisecond)).Err
		r.putRedisErr(&err)
		return err
	}
	err := r.redis.Cmd(r.ctx, "set", key, value, "ex", int64(expiration/time.Second)).Err
	r.putRedisErr(&err)
	return err
}

// Get get
func (r *RedisClient) Get(key string) (string, error) {
	if err := r.checkFailed(); err != nil {
		return "", err
	}

	if key == "" {
		return "", nil
	}

	reply, err := r.cmdWithoutTx("get", key)
	if err != nil {
		return "", err
	}

	if reply.Type == radix.NilReply {
		return "", nil
	}
	str, err := reply.Str()
	r.putRedisErr(&err)
	return str, err
}

// Del del
func (r *RedisClient) Del(keys ...string) error {
	if err := r.checkFailed(); err != nil {
		return err
	}

	var args []interface{}
	for _, key := range keys {
		if key != "" {
			args = append(args, key)
		}
	}

	if len(args) == 0 {
		return nil
	}

	err := r.redis.Cmd(r.ctx, "del", args...).Err
	r.putRedisErr(&err)
	return err
}

// HSet hset
func (r *RedisClient) HSet(key, field string, value interface{}) error {
	if err := r.checkFailed(); err != nil {
		return err
	}

	err := r.redis.Cmd(r.ctx, "hset", key, field, value).Err
	r.putRedisErr(&err)
	return err
}

// HGet hget
func (r *RedisClient) HGet(key, field string) (string, error) {
	if err := r.checkFailed(); err != nil {
		return "", err
	}

	if key == "" || field == "" {
		return "", nil
	}

	reply, err := r.cmdWithoutTx("hget", key, field)
	if err != nil {
		return "", err
	}

	if reply.Type == radix.NilReply {
		return "", nil
	}
	str, err := reply.Str()
	r.putRedisErr(&err)
	return str, err
}

// HDel hdel
func (r *RedisClient) HDel(key string, fields ...string) error {
	if err := r.checkFailed(); err != nil {
		return err
	}

	var args []interface{}
	for _, field := range fields {
		if field != "" {
			args = append(args, field)
		}
	}

	if len(args) == 0 {
		return nil
	}

	a := append([]interface{}{key}, args...)

	err := r.redis.Cmd(r.ctx, "hdel", a...).Err
	r.putRedisErr(&err)
	return err
}

// HGetAll hgetall
func (r *RedisClient) HGetAll(key string) (map[string]string, error) {
	if err := r.checkFailed(); err != nil {
		return nil, err
	}

	if key == "" {
		return nil, nil
	}

	reply, err := r.cmdWithoutTx("hgetall", key)
	if err != nil {
		return nil, err
	}
	hash, err := reply.Hash()
	r.putRedisErr(&err)
	return hash, err
}

// HVals hvals
func (r *RedisClient) HVals(key string) ([]string, error) {
	if err := r.checkFailed(); err != nil {
		return nil, err
	}

	if key == "" {
		return nil, nil
	}

	reply, err := r.cmdWithoutTx("hvals", key)
	if err != nil {
		return nil, err
	}

	ls, err := reply.List()
	r.putRedisErr(&err)
	return ls, err
}

func (r *RedisClient) BeginTx() error {
	if err := r.checkFailed(); err != nil {
		return err
	}

	r.txMutex.Lock()
	defer r.txMutex.Unlock()

	if r.isTx {
		return fmt.Errorf("ERR MULTI calls can not be nested")
	}

	if err := r.redis.Cmd(r.ctx, "multi").Err; err != nil {
		r.putRedisErr(&err)
		return err
	}

	r.isTx = true
	return nil
}

func (r *RedisClient) Commit() error {
	if err := r.checkFailed(); err != nil {
		return err
	}

	r.txMutex.Lock()
	defer r.txMutex.Unlock()

	if !r.isTx {
		return fmt.Errorf("not redis transaction client")
	}

	r.isTx = false
	err := r.redis.Cmd(r.ctx, "exec").Err
	r.putRedisErr(&err)
	return err
}

// Discard 放弃事务
func (r *RedisClient) Discard() error {
	if err := r.checkFailed(); err != nil {
		return err
	}

	r.txMutex.Lock()
	defer r.txMutex.Unlock()

	if !r.isTx {
		return fmt.Errorf("not redis transaction client")
	}

	r.isTx = false
	err := r.redis.Cmd(r.ctx, "discard").Err
	r.putRedisErr(&err)
	return err
}

// GetContext 获取context
func (r *RedisClient) GetContext() context.Context {
	return r.ctx
}