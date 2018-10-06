package token

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"

	"business/user/common"
	"business/user/db"
)

func init() {
	common.InitConfig(configPath)
	db.InitMySQL(common.Config.Mysql)
}

func initTestUserTokenCaseEnv(t *testing.T) {
	initUserTokenTable(t)
}

func initUserTokenTable(t *testing.T) {
	assert := assert.New(t)

	f, err := ioutil.ReadFile(TestUserTokenSql)
	assert.Nil(err)
	assert.NotNil(f)
	result, err := db.MySQL.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s;", UserTokenTableName))
	assert.Nil(err)
	assert.NotEmpty(result)
	result, err = db.MySQL.Exec(string(f))
	assert.Nil(err)
	assert.NotEmpty(result)
}

func TestSaveForeverTokenAndDeleteForeverToken(t *testing.T) {
	assert := assert.New(t)
	initTestUserTokenCaseEnv(t)

	err := saveForeverToken(context.Background(), "test_name", "test_user_id", "test_access_token")
	assert.Nil(err)

	err = deleteForeverToken(context.Background(), "test_name", "test_user_id")
	assert.Nil(err)
}