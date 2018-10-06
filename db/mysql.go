package db

import (
	"context"
	rawsql "database/sql"

	_ "github.com/go-sql-driver/mysql" // 注册 mysql driver

	"github.com/chenjun-git/umbrella-common/database/sql"

	"business/user/common"
)

// MySQL 全局的MySQL连接池
var MySQL *sql.DB

// InitMySQL 根据配置信息连接MySQL并设置参数
func InitMySQL(config *common.MysqlConfig) {
	var middlewares []sql.DBMiddleware

	pool, err := sql.Open("mysql", config.Dsn, middlewares...)
	if err != nil {
		panic(err)
	}

	if config.MaxIdle > 0 {
		pool.SetMaxIdleConns(config.MaxIdle)
	}
	if config.MaxOpen > 0 {
		pool.SetMaxOpenConns(config.MaxOpen)
	}

	MySQL = pool
}

// MySQLExec mysql接口
type MySQLExec interface {
	Exec(query string, args ...interface{}) (rawsql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)

	ExecContext(ctx context.Context, query string, args ...interface{}) (rawsql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type dberWithContext struct {
	MySQLExec
	ctx context.Context
}

// Exec exec
func (db *dberWithContext) Exec(query string, args ...interface{}) (rawsql.Result, error) {
	return db.MySQLExec.ExecContext(db.ctx, query, args...)
}

// Query query
func (db *dberWithContext) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.MySQLExec.QueryContext(db.ctx, query, args...)
}

// QueryRow query one row
func (db *dberWithContext) QueryRow(query string, args ...interface{}) *sql.Row {
	return db.MySQLExec.QueryRowContext(db.ctx, query, args...)
}

// BindDBerWithContext 生成带context的client
func BindDBerWithContext(ctx context.Context, db MySQLExec) MySQLExec {
	return &dberWithContext{db, ctx}
}
