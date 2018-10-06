package utils

import (
	"context"
	"database/sql"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"


	"github.com/json-iterator/go"
)

// NullStringToString NullString类型转换为String类型
func NullStringToString(s sql.NullString) interface{} {
	if s.Valid {
		return s.String
	} else {
		return nil
	}
}

func StringToNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

const (
	Digits   = "0123456789"
	Alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Ascii    = Alphabet + Digits + "~!@#$%^&*()-_+={}[]\\|<,>.?/\"';:`"
)

func RandomDigits(length int) string {
	return randomString(length, []byte(Digits))
}

func randomString(length int, base []byte) string {
	bytes := make([]byte, length)
	maxIndex := len(base)
	for i := 0; i < length; i++ {
		index := rand.Intn(maxIndex)
		bytes[i] = byte(base[index])
	}

	return string(bytes)
}

func BindJSON(r *http.Request, obj interface{}) error {
	defer io.Copy(ioutil.Discard, r.Body)
	if err := jsoniter.NewDecoder(r.Body).Decode(obj); err != nil {
		return err
	}

	return nil//structValidator.ValidateStruct(obj)
}

////////////////////////////////////////////////////////////////////////
type contextKeyType string

const (
	currentAccessToken contextKeyType = "current_access_token"
	currentUserID      contextKeyType = "current_user_id"
)

// SetCurrentUserAuth 设置context中的tenant鉴权信息
func SetCurrentUserAuth(ctx context.Context, accessToken, userID string) context.Context {
	ctx = context.WithValue(ctx, currentAccessToken, accessToken)
	ctx = context.WithValue(ctx, currentUserID, userID)
	return ctx
}

// GetCurrentUserAuth 获取context中的tenant鉴权信息
func GetCurrentUserAuth(ctx context.Context) (userID, accessToken string) {
	userID = ""
	accessToken = ""

	userID, ok := ctx.Value(currentUserID).(string)
	if !ok || userID == "" {
		return
	}

	accessToken, ok = ctx.Value(currentAccessToken).(string)
	if !ok || accessToken == "" {
		return
	}

	return
}