package db

import (
	"context"
	"time"

	"github.com/stretchr/testify/assert"
	
	"business/user/common"
)

const (
	configPath = "../config/config.user.toml"
)

func init() {
	common.InitConfig(configPath)
	InitRedis(common.Config.Redis)
}

func TestRedisGetSetDel(t *testing.T) {
	assert := assert.New(t)

	err := RedisSet(context.Background(), "test", "test", 5*time.Second)
	assert.Nil(err)

	reply, err := RedisGet(context.Background(), "test")
	assert.Nil(err)
	assert.NotNil(reply)

	err = RedisDel(context.Background(), "test")
	assert.Nil(err)
}

func TestRedisHGetSetDel(t *testing.T) {
	assert := assert.New(t)

	err := RedisHset(context.Background(), "test", "test", "test", 5*time.Second)
	assert.Nil(err)

	reply, err := RedisHget(context.Background(), "test", "test")
	assert.Nil(err)
	assert.NotNil(reply)

	err = RedisHdel(context.Background(), "test", test)
	assert.Nil(err)
}