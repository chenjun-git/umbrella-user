package token

import (
	"fmt"
	"testing"
	"io/ioutil"

	"github.com/stretchr/testify/assert"

	"github.com/chenjun-git/umbrella-user/common"
	"github.com/chenjun-git/umbrella-user/db"
)

func init() {
	common.InitConfig(configPath)
	db.InitMySQL(common.Config.Mysql)
}

func initTestAccessTokenCaseEnv(t *testing.T) {
	initAccessTokenTable(t)
}

func initAccessTokenTable(t *testing.T) {
	assert := assert.New(t)

	f, err := ioutil.ReadFile(TestAccessTokenSql)
	assert.Nil(err)
	assert.NotNil(f)
	result, err := db.MySQL.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s;", AccessTokenTableName))
	assert.Nil(err)
	assert.NotEmpty(result)
	result, err = db.MySQL.Exec(string(f))
	assert.Nil(err)
	assert.NotEmpty(result)
}

func TestMySQLSaveV1TokenAndMySQLDeleteALLTokenOfUser(t *testing.T) {
	assert := assert.New(t)
	initTestAccessTokenCaseEnv(t)

	err := MySQLSaveV1Token("test_user_id", "huawei", "umbrella", "access_token", "refresh_token", 1, 1)
	assert.Nil(err)

	err = MySQLDeleteALLTokenOfUser(db.MySQL, "test_user_id")
	assert.Nil(err)
}

// func Test(t *testing.T) {
// 	assert := assert.New(t)
// 	initAccessTokenTable(t)
// }