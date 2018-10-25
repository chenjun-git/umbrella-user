package token

import (
	"context"
	"fmt"
	"time"

	"github.com/chenjun-git/umbrella-user/db"
)

const (
	AccessTokenTableName = "access_token"
)

func MySQLSaveV1Token(userID, device, app, accessToken, refreshToken string, accessTTL, refreshTTL int64) error {
	now := time.Now()

	SQL := `REPLACE INTO access_token (access_token, refresh_token, user_id, device, app, access_token_expired_at, refresh_token_expired_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := db.MySQL.Exec(SQL, accessToken, refreshToken, userID, device, app, now.Add(time.Duration(accessTTL)*time.Minute), now.Add(time.Duration(refreshTTL)*time.Minute))
	if err != nil {
		fmt.Printf("MySQLSaveV1Token, err : %+v\n", err)
	}
	return err
}

// MySQLDeleteALLTokenOfUser ...
func MySQLDeleteALLTokenOfUser(mysql db.MySQLExec, userID string) error {
	_, err := mysql.Exec(`DELETE FROM access_token WHERE user_id = ?`, userID)
	if err != nil {
		// log.Warn("MySQLDeleteALLTokenOfUser", err)
		fmt.Printf("MySQLDeleteALLTokenOfUser, err : %+v\n", err)
	}
	return err
}

// MySQLDeleteAllV1TokenByAccountID ...
func MySQLDeleteAllV1TokenByAccountID(mysql db.MySQLExec, accountID string) error {
	_, err := mysql.Exec(`DELETE FROM access_token WHERE account_id = ?`, accountID)
	if err != nil {
		// log.Warn("MySQLDeleteAllV1TokenByAccountID", err)
	}

	return err
}

// MySQLDeleteAllV1TokenByDeviceAndApp ...
func MySQLDeleteAllV1TokenByDeviceAndApp(mysql db.MySQLExec, userID, device, app string) error {
	_, err := mysql.Exec(`DELETE FROM access_token WHERE  user_id = ? AND device = ? AND app = ?`, userID, device, app)
	if err != nil {
		// log.Warn("MySQLDeleteAllV1TokenByDeviceAndApp", err)
	}

	return err
}

// MySQLDeleteAllV1TokenByApp ...
func MySQLDeleteAllV1TokenByApp(mysql db.MySQLExec, userID, app string) error {
	_, err := mysql.Exec(`DELETE FROM access_token WHERE user_id = ? AND app = ?`, userID, app)
	if err != nil {
		// log.Warn("MySQLDeleteAllV1TokenByApp", err)
	}
	return err
}

// MySQLDeleteV1TokenByToken ...
func MySQLDeleteV1TokenByToken(mysql db.MySQLExec, userID, accessOrRefreshToken string) error {
	_, err := mysql.Exec(`DELETE FROM access_token WHERE user_id = ? AND (access_token = ? OR refresh_token = ?)`, userID, accessOrRefreshToken, accessOrRefreshToken)
	if err != nil {
		// log.Warn("MySQLDeleteV1TokenByToken", err)
	}
	return err
}

// MySQLGetDeviceAndAppOfUser 获取用户的device和app列表
func MySQLGetDeviceAndAppOfUser(ctx context.Context, userID string) (map[string][]string, error) {
	var sources = make(map[string][]string)
	var device, app string
	rows, err := db.MySQL.QueryContext(ctx, "SELECT device, app FROM access_token WHERE user_id = ? AND refresh_token_expired_at > ?", userID, time.Now())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&device, &app); err != nil {
			return nil, err
		}
		if _, ok := sources[app]; ok {
			sources[app] = append(sources[app], device)
		} else {
			sources[app] = []string{device}
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sources, nil
}