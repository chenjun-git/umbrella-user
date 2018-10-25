package token

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/chenjun-git/umbrella-user/common"
	"github.com/chenjun-git/umbrella-user/db"
)

var redisExec db.RedisExec

func init() {
	common.InitConfig(configPath)
	db.InitRedis(common.Config.Redis)
	redis, _ := db.NewRedisClient(context.Background())
	redisExec = redis
}

func TestRedisSaveV1Token(t *testing.T) {
	assert := assert.New(t)

	err := redisSaveV1Token(redisExec, "test_id", "huawei", "test_app", "access_token", "refresh_token", 1, 1)
	assert.Nil(err)
}


func TestRedisDeleteToken(t *testing.T) {
	assert := assert.New(t)

	err := redisDeleteToken(redisExec, "test_id", "huawei", "test_app")
	assert.Nil(err)
}