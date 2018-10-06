package db

import (
	"fmt"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	
	"business/user/common"
)

const (
	configPath = "../conf/config.user.toml"
)

func init() {
	common.InitConfig(configPath)
	fmt.Printf("redis config : %+v\n", common.Config.Redis)
	InitRedis(common.Config.Redis)
}

func TestRedisGetSetDel(t *testing.T) {
	assert := assert.New(t)

	err := RedisSet(context.Background(), "test", "test_str", 5*time.Second)
	assert.Nil(err)

	reply, err := RedisGet(context.Background(), "test")
	assert.Nil(err)
	assert.NotNil(reply)
	data, err := reply.Str()
	assert.Nil(err)
	assert.Equal("test_str", data)

	err = RedisDel(context.Background(), "test")
	assert.Nil(err)
}

func TestRedisHGetSetDel(t *testing.T) {
	assert := assert.New(t)

	err := RedisHset(context.Background(), "test", "test", "test_hash")
	assert.Nil(err)

	reply, err := RedisHget(context.Background(), "test", "test")
	assert.Nil(err)
	assert.NotNil(reply)
	data, err := reply.Str()
	assert.Nil(err)
	assert.Equal("test_hash", data)

	err = RedisHdel(context.Background(), "test", "test")
	assert.Nil(err)
}