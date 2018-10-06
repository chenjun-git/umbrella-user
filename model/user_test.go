package model

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"business/user/common"
	"business/user/db"
)
const (
	TestConfigFilePath = "../conf/config.dev.toml"

	TestUserSQLPath = "../sql/user.sql"
)

func init() {
	common.InitConfig(TestConfigFilePath)
	db.InitMySQL(common.Config.Mysql)
}

func TestInitUserTable(t *testing.T) {
	assert := assert.New(t)

	f, err := ioutil.ReadFile(TestUserSQLPath)
	assert.Nil(err)
	assert.NotNil(f)
	result, err := db.MySQL.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s;", UserTableName))
	assert.Nil(err)
	assert.NotEmpty(result)
	result, err = db.MySQL.Exec(string(f))
	assert.Nil(err)
	assert.NotEmpty(result)
}

func TestSqlInsert(t *testing.T) {
	assert := assert.New(t)

	rowsAffected, err := sqlInsert(db.MySQL, UserTableName)
	assert.Nil(err)
}

func TestSqlQueryRow(t *testing.T) {
	assert := assert.New(t)

	rows, err := sqlQueryRow(db.MySQL, UserTableName)
	assert.Nil(err)
}

func TestSqlUpdate(t *testing.T) {
	assert := assert.New(t)

	rowsAffected, err := sqlUpdate(db.MySQL, UserTableName)
	assert.Nil(err)
}