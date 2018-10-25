package db

import (
	"github.com/chenjun-git/umbrella-common/database/sql"

	"github.com/chenjun-git/umbrella-user/common"
)

// RedisCommitOrRollback ...
func RedisCommitOrRollback(e error, redis RedisExec) error {
	if e == nil {
		if err := redis.Commit(); err != nil {
			return common.ExtendErrorStatus(common.AccountRedisError, err)
		}
	} else {
		if err := redis.Discard(); err != nil {
			return common.ExtendErrorStatus(common.AccountRedisError, err)
		}
	}

	return nil
}

// MySQLCommitOrRollback ...
func MySQLCommitOrRollback(e error, mysql *sql.Tx) error {
	if e == nil {
		if err := mysql.Commit(); err != nil {
			return common.ExtendErrorStatus(common.AccountMySQLError, err)
		}
	} else {
		if err := mysql.Rollback(); err != nil {
			return common.ExtendErrorStatus(common.AccountMySQLError, err)
		}
	}

	return nil
}