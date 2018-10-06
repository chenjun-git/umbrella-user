package token

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/json-iterator/go"
	
	"business/user/common"
	"business/user/db"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	AccessTokenKeyName  = "access_token"
	RefreshTokenKeyName = "refresh_token"
)

type deviceApp struct {
	App    string `json:"app"`
	Device string `json:"device"`
}

/*
token存储在redis说明

存储：一对token会存四条数据

key: <access_token>, value: {"app":<app>,"device":<device>}
key: <refresh_token>, value: {"app":<app>,"device":<device>}
hkey: <user_id>, field: <device>-<app>-access_token, value: <access_token>
hkey: <user_id>, field: <device>-<app>-refresh_token, value: <refresh_token>

目标
1、指定用户+device+app，能拿到token
2、指定token，能拿到device+app：刷新需要拿到token中的devie+app
*/

func joinTokenRedisKey(userID, device, app string) (string, string, string, string, error) {
	userTokenKey := fmt.Sprintf("%s", userID)
	accessTokenDeviceKey := fmt.Sprintf("%s-%s-%s", device, app, AccessTokenKeyName)
	refreshTokenDeviceKey := fmt.Sprintf("%s-%s-%s", device, app, RefreshTokenKeyName)

	deviceValue, err := json.Marshal(deviceApp{App: app, Device: device})
	if err != nil {
		return "", "", "", "", common.ExtendErrorStatus(common.AccountJsonMarshalError, err)
	}

	return userTokenKey, accessTokenDeviceKey, refreshTokenDeviceKey, string(deviceValue), nil
}

func redisSaveV1Token(redisTx db.RedisExec, userID, device, app, accessToken, refreshToken string, accessTTL, refreshTTL int64) error {
	userTokenKey, accessTokenDeviceKey, refreshTokenDeviceKey, deviceValue, err := joinTokenRedisKey(userID, device, app)
	if err != nil {
		fmt.Printf("redisSaveV1Token err : %+v\n", err)
		return err
	}

	redisTx.HSet(userTokenKey, accessTokenDeviceKey, accessToken)
	redisTx.HSet(userTokenKey, refreshTokenDeviceKey, refreshToken)
	redisTx.Set(accessToken, deviceValue, time.Duration(accessTTL)*time.Minute)
	redisTx.Set(refreshToken, deviceValue, time.Duration(refreshTTL)*time.Minute)

	return nil
}

// 如果：app == "" && device == ""，那么删除这个user的所有device+app的token
// 如果：device == "" &&  <device>-<app>-access_token 以 "-<app>-<AccessTokenKeyName>" 为后缀，删除他
// 如果：<device>-<app>-access_token / <device>-<app>-refresh_token 等于 "<device>-<app>-<AccessTokenKeyName>" / "<device>-<app>-<RefreshTokenKeyName>"，删除他
func redisDeleteToken(redisTx db.RedisExec, userID, app, device string) error {
	userTokenKey, accessTokenDeviceKey, refreshTokenDeviceKey, _, _ := joinTokenRedisKey(userID, device, app)

	deviceTokenMap, err := redisTx.HGetAll(userTokenKey)
	if err != nil {
		return err
	}

	var hdels []string
	var dels []string
	deleteALL := app == "" && device == ""
	for deviceKey, token := range deviceTokenMap {
		deleteApp := device == "" && (strings.HasSuffix(deviceKey, accessTokenDeviceKey) || strings.HasSuffix(deviceKey, refreshTokenDeviceKey))
		deleteAppDevice := deviceKey == accessTokenDeviceKey || deviceKey == refreshTokenDeviceKey

		if deleteALL || deleteApp || deleteAppDevice {
			hdels = append(hdels, deviceKey)
			dels = append(dels, token)
		}
	}

	redisTx.HDel(userTokenKey, hdels...)
	redisTx.Del(dels...)

	return nil
}

func redisDeleteV1TokenByDeviceAndApp(redis db.RedisExec, userID, device, app string) error {
	err := redisDeleteToken(redis, userID, app, device)
	if err != nil {
		fmt.Printf("redisDeleteV1TokenByDeviceAndApp, err : %+v\n", err)
	}

	return err
}

func RedisDeleteALLTokenOfUser(redisTx db.RedisExec, userID string) error {
	err := redisDeleteToken(redisTx, userID, "", "")
	if err != nil {
		fmt.Printf("RedisDeleteALLTokenOfUser, err : %+v\n", err)
	}
	return err
}

func RedisDeleteV1TokenByApp(redisTx db.RedisExec, userID, app string) error {
	err := redisDeleteToken(redisTx, userID, app, "")
	if err != nil {
		fmt.Printf("RedisDeleteV1TokenByApp, err : %+v\n", err)
	}
	return err
}

// RedisCheckTokenExistByToken 检查token是否存在，包括access_token以及refresh_token
func RedisCheckTokenExistByToken(ctx context.Context, v1Token string) (bool, error) {
	redis, err := db.NewRedisClient(ctx)
	if err != nil {
		fmt.Printf("RedisCheckTokenExistByToken, err : %+v\n", err)
		return false, err
	}
	defer redis.PutPool()

	exist, err := redis.Exists(v1Token)
	if err != nil {
		fmt.Printf("RedisCheckTokenExistByToken, err : %+v\n", err)
	}

	return exist, err
}

func redisCheckAccessTokenExistAndGetDeviceApp(ctx context.Context, v1Token string) (bool, string, string, error) {
	redis, err := db.NewRedisClient(ctx)
	if err != nil {
		return false, "", "", err
	}
	defer redis.PutPool()

	deviceAndApp, err := redis.Get(v1Token)
	if err != nil {
		return false, "", "", err
	} else if deviceAndApp == "" {
		return false, "", "", nil
	}

	var deviceApp deviceApp
	if err = json.Unmarshal([]byte(deviceAndApp), &deviceApp); err != nil {
		return false, "", "", err
	}

	return true, deviceApp.Device, deviceApp.App, nil
}

func redisGetDeviceAndAppOfAccessOrRefreshToken(ctx context.Context, v1Token string) (device string, app string, err error) {
	exist, device, app, err := redisCheckAccessTokenExistAndGetDeviceApp(ctx, v1Token)
	if err != nil {
		err = common.ExtendErrorStatus(common.AccountRedisError, err)
		fmt.Printf("redisGetDeviceAndAppOfAccessOrRefreshToken, err : %+v\n", err)
		return
	} else if !exist {
		err = fmt.Errorf("token: %s not exist", v1Token)
		fmt.Printf("redisGetDeviceAndAppOfAccessOrRefreshToken, err : %+v\n", err)
		return
	} else if device == "" && app == "" {
		err = common.NormalErrorStatus(common.AccountBindNoDevice)
		fmt.Printf("redisGetDeviceAndAppOfAccessOrRefreshToken, err : %+v\n", err)
		return
	}

	return device, app, nil
}