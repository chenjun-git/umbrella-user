package middleware

import (
	"context"
	"errors"
	"fmt"
	"time"
	"github.com/chenjun-git/umbrella-common/token"

	"github.com/chenjun-git/umbrella-user/common"
	"github.com/chenjun-git/umbrella-user/db"
	"github.com/chenjun-git/umbrella-user/model"
	"github.com/chenjun-git/umbrella-user/utils/cacheverify"
)

// func CheckAccessToken(ctx context.Context, accessToken string) (*model.User, error) {
// 	if accessToken == "" {
// 		return nil, errors.New("access token is empty")
// 	}
// 	// verify
// 	_, err := tenant.VerifyAccessToken(ctx, accessToken)
// 	if err != nil {
// 		return nil, err
// 	}

// 	userID, err := tenant.GetAccountIDByAccessToken(ctx, accessToken)
// 	if err != nil {
// 		return nil, err
// 	}

// 	sqlExec := db.BindDBerWithContext(ctx, db.MySQL)
// 	user, err := model.GetUserById(sqlExec, userID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if user == nil {
// 		return nil, errors.New("account by access token not found")
// 	}
// 	return user, nil
// }

func checkUserTokenFactory(ctx context.Context, authToken, purpose string) (*model.User, *token.Token, error) {
	if authToken == "" {
		return nil, nil, errors.New("account token is empty")
	}

	token, err := token.DecryptAccessToken(authToken)
	if err != nil {
		return nil, token, err
	}

	// payload.ttl 单位为分钟
	ttlDuration := time.Duration(token.TTL) * time.Minute

	// 转为秒检测超时
	if uint64(token.IssueTime)+uint64(ttlDuration.Seconds()) <= uint64(time.Now().Unix()) {
		return nil, token, errors.New("account token expired")
	}

	if purpose == common.TwoFactorAuthVerifyTokenKey {
		tokenVerify, err := cacheverify.NewVerifyKey(ctx, token.UserID, purpose, ""/*token.LoginSource*/).GetStringVerify()
		if err != nil {
			return nil, token, errors.New("redis get failed")
		}
		if token == nil || tokenVerify.Token != authToken {
			return nil, token, errors.New("2fa tmp token invalid")
		}
	} else {
		tokenVerify, err := cacheverify.NewVerifyKey(ctx, token.UserID, purpose, ""/*token.LoginSource*/).GetUserTokenVerify()
		if err != nil {
			return nil, token, fmt.Errorf("account token db get failed: %s", err)
		}
		if token == nil || tokenVerify.Token != authToken {
			return nil, token, errors.New("account token invalid")
		}
	}

	sqlExec := db.BindDBerWithContext(ctx, db.MySQL)
	user, err := model.GetUserById(sqlExec, token.UserID)
	if err != nil || user == nil {
		return nil, token, errors.New("account token payload invalid")
	}

	return user, token, nil
}

func CheckUserToken(ctx context.Context, authToken string) (*model.User, *token.Token, error) {
	return checkUserTokenFactory(ctx, authToken, "")
}