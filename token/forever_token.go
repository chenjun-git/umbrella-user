package token

import (
	"context"
	//"time"

	"business/user/common"
	"business/user/db"
	"business/user/utils"
)

const (
	UserTokenTableName = "user_token"
)

func saveForeverToken(ctx context.Context, name, userID, accessToken string) error {
	_, err := db.MySQL.ExecContext(ctx, `INSERT INTO user_token (name, token, user_id) VALUES (?, ?, ?)`, name, accessToken, userID)
	return err
}

func deleteForeverToken(ctx context.Context, name, userID string) error {
	_, err := db.MySQL.ExecContext(ctx, "DELETE FROM user_token WHERE name = ? AND user_id = ?", name, userID)
	return err
}

// DeleteForeverTokenOfUser 删除 foreverToken
func DeleteForeverTokenOfUser(mysql db.MySQLExec, userID string) error {
	_, err := mysql.Exec("DELETE FROM user_token WHERE user_id = ?", userID)
	return err
}

func isExistForeverToken(ctx context.Context, foreverToken string) (bool, error) {
	var c int
	if err := db.MySQL.QueryRowContext(ctx, `SELECT COUNT(1) FROM personal_token WHERE token IN (?, ?)`, foreverToken, utils.Hash(foreverToken)).Scan(&c); err != nil {
		return false, common.ExtendErrorStatus(common.AccountMySQLError, err)
	}

	return c > 0, nil
}