package token

import (
	"context"
	"fmt"
	"time"

	"github.com/chenjun-git/umbrella-common/token"

	"business/user/common"
	"business/user/db"
	"business/user/utils"
	"business/user/utils/cacheverify"
)

// VerifyTokenResult 验证token的结果的类型
type VerifyTokenResult string

// AccessTokenType 验证token的token的类型
type AccessTokenType string

const (
	UserIDLength       = 32
	AccessTokenLength  = 184
	RefreshTokenLength = 116

	VerifyTokenAllow      VerifyTokenResult = "VerifyToken_Allow"
	VerifyTokenDeny       VerifyTokenResult = "VerifyToken_Deny"
	VerifyTokenExpireSoon VerifyTokenResult = "VerifyToken_ExpireSoon"
	VerifyTokenExpired    VerifyTokenResult = "VerifyToken_Expired"

	UnKnownToken AccessTokenType = "UnKnownToken"
	V1Token      AccessTokenType = "V1Token"
	ForeverToken AccessTokenType = "ForeverToken"

	TokenLibV1Version = 1
)


type VerifyResult struct {
	VerifyResult VerifyTokenResult
	Type         AccessTokenType
	UserID       string
	Device       string
	App          string
	Message      string
}

// VerifyTokenStatusToCode 将验证token返回的结果转为错误码
var VerifyTokenStatusToCode = map[VerifyTokenResult]int{
	VerifyTokenDeny:       common.TokenDeny,
	VerifyTokenExpireSoon: common.TokenExpiredSoon,
	VerifyTokenExpired:    common.TokenExpired,
}

func GenV1Token(userID string, accessTTL, refreshTTL time.Duration) (accessToken, refreshToken string, err error) {
	if len(userID) != UserIDLength {
		return "", "", fmt.Errorf("invalid length of user_id(%v)", UserIDLength)
	}

	issueTime := time.Now()
	tk := &token.Token{
		IssueTime: uint32(issueTime.Unix()),
		TTL:       uint16(accessTTL.Minutes()),
		UserID:    userID,
	}

	accessToken, err = token.EncryptAccessToken(TokenLibV1Version, tk)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = encodeRefreshToken(uint32(issueTime.Unix()), uint16(refreshTTL.Minutes()), userID)
	if err != nil {
		return "", "", err
	}

	return
}

func GenV1TokenAndSaveRedis(ctx context.Context, userID, device, app string, accessTTL, refreshTTL int64) (accessToken, refreshToken string, err error) {
	if device == "" {
		return "", "", fmt.Errorf(common.GetMsg(common.AccountBindNoDevice, []string{}))
	} else if app == "" {
		return "", "", fmt.Errorf(common.GetMsg(common.AccountBindNoApp, []string{}))
	}

	accessToken, refreshToken, err = GenV1Token(userID, time.Duration(accessTTL)*time.Minute, time.Duration(refreshTTL)*time.Minute)
	if err != nil {
		return "", "", err
	}

	redis, err := db.NewRedisClient(ctx)
	if err != nil {
		return "", "", err
	}
	defer redis.PutPool()

	if err := redisSaveV1Token(redis, userID, device, app, accessToken, refreshToken, accessTTL, refreshTTL); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func RedisDeleteV1TokenByToken(ctx context.Context, userID, accessOrRefreshToken string) error {
	redis, err := db.NewRedisClient(ctx)
	if err != nil {
		return err
	}
	defer redis.PutPool()

	device, app, err := redisGetDeviceAndAppOfAccessOrRefreshToken(redis.GetContext(), accessOrRefreshToken)
	if err != nil {
		return err
	}

	return redisDeleteV1TokenByDeviceAndApp(redis, userID, device, app)
}

// RefreshToken 刷新token
func RefreshToken(ctx context.Context, oldRefreshToken string) (accessToken, refreshToken string, err error) {
	issueTime, TTL, userID, _, err := DecodeRefreshToken(oldRefreshToken)
	if err != nil {
		return
	}

	now := time.Now().Unix()
	if expired := issueTime + uint32(TTL)*60; int64(expired) < now {
		err = fmt.Errorf("refresh token expired: %d, now=%d", expired, now)
		return
	}

	accessTTL := int64(common.Config.Token.AccessTokenTTL.Minutes())
	refreshTTL := int64(common.Config.Token.RefreshTokenTTL.Minutes())

	device, app, err := redisGetDeviceAndAppOfAccessOrRefreshToken(ctx, oldRefreshToken)
	if err != nil {
		return "", "", err
	}

	redis, err := db.NewRedisClient(ctx)
	if err != nil {
		return "", "", err
	}
	defer redis.PutPool()

	// 刷新会replace into token，不需要删除
	if err := redisDeleteV1TokenByDeviceAndApp(redis, userID, device, app); err != nil {
		return "", "", err
	}

	accessToken, refreshToken, err = GenV1Token(userID, time.Duration(accessTTL)*time.Minute, time.Duration(refreshTTL)*time.Minute)
	if err != nil {
		return "", "", err
	}

	if err := redisSaveV1Token(redis, userID, device, app, accessToken, refreshToken, accessTTL, refreshTTL); err != nil {
		return "", "", err
	}
	MySQLSaveV1Token(userID, device, app, accessToken, refreshToken, accessTTL, refreshTTL)

	return accessToken, refreshToken, nil
}

func VerifyToken(ctx context.Context, accessToken string) (*VerifyResult, error) {
	result := &VerifyResult{
		VerifyResult: VerifyTokenDeny,
		Type:         UnKnownToken,
	}

	version, err := token.GetTokenVersion(accessToken)
	if err != nil {
		result.Message = fmt.Sprintf("get token version failed: %v %v", accessToken, err)
		return result, nil
	}

	tk, err := token.DecryptAccessToken(accessToken)
	if err != nil {
		result.Message = fmt.Sprintf("decrypt token failed: %v %v", accessToken, err)
		return result, nil
	}
	result.UserID = tk.UserID

	switch version {
	case TokenLibV1Version:
		if tk.TTL == 0 {
			result.Type = ForeverToken
			// TTL 有效时间设置为0即为无过期时间，只有GenV1ForeverToken接口可生成，其余情况限制TTL至少为1
			if exist, err := isExistForeverToken(ctx, accessToken); err != nil {
				return nil, err
			} else if !exist {
				return result, nil
			}
		} else {
			result.Type = V1Token

			expireTime := int64(tk.IssueTime + uint32(tk.TTL)*60)

			// 减5的原因是：请求的时间比较巧，token里面编码的时间还没有过期，所以不会返回过期，但是去redis/mysql查不到（这个时候过期了），所以认为是deny
			expireLeftTime := expireTime - time.Now().Unix() - 5 // 单位秒

			if expireLeftTime < 0 {
				result.VerifyResult = VerifyTokenExpired
				return result, nil
			}

			isExist := false
			isExist, result.Device, result.App, err = redisCheckAccessTokenExistAndGetDeviceApp(ctx, accessToken)
			if err != nil {
				err = fmt.Errorf("check access token: %s if in redis failed error: %v", accessToken, err)
				return nil, err
			} else if !isExist {
				// redis不存在说明是退出了，因为前面已经做超时检测了
				return result, nil
			}

			if expireLeftTime <= int64(common.Config.Token.AccessTokenExpireSoon.D().Seconds()) {
				result.VerifyResult = VerifyTokenExpireSoon
				return result, nil
			}
		}
	default:
		result.Message = fmt.Sprintf("token: %v have no valid token version: %v", accessToken, version)
		return result, nil
	}

	result.VerifyResult = VerifyTokenAllow
	return result, nil
}

func genV1ForeverToken(userID string) (accessToken string, err error) {
	if len(userID) != UserIDLength {
		return "", fmt.Errorf("invalid length of user_id(%v)", UserIDLength)
	}
	var mask1, mask2 int64

	issueTime := time.Now()
	tk := &token.Token{
		Mask1:     mask1,
		Mask2:     mask2,
		IssueTime: uint32(issueTime.Unix()),
		TTL:       uint16(0),
		UserID:    userID,
	}

	accessToken, err = token.EncryptAccessToken(TokenLibV1Version, tk)
	if err != nil {
		return "", err
	}

	return
}

// CreateV1ForeverToken 创建foreverToken
func CreateV1ForeverToken(ctx context.Context, name, userID string) (accessToken string, err error) {
	accessToken, err = genV1ForeverToken(userID)
	if err != nil {
		return "", err
	}

	err = saveForeverToken(ctx, name, userID, utils.Hash(accessToken))
	if err != nil {
		return
	}

	return accessToken, nil
}

// RemoveForeverToken 删除foreverToken
func RemoveForeverToken(ctx context.Context, name, userID string) error {
	return deleteForeverToken(ctx, name, userID)
}

func RemoveInvalidToken(ctx context.Context, userID string) {
	// 结合genTokenAndStore函数的修改source理解
	key, _ := cacheverify.NewVerifyKey(ctx, userID, "", "").GetHashKey()
	replyList, err := db.RedisHgetall(ctx, key)
	if err != nil {
		return
	}
	for _, s := range common.SourceRange {
		tokenKey := cacheverify.NewVerifyKey(ctx, userID, "", s)
		tokenKey.CleanHash()
	}
	var sourceList []string
	var source string
	for index, value := range replyList.Elems {
		if index%2 == 0 {
			source = value.String()
		} else {
			if value.String() == "" {
				sourceList = append(sourceList, source)
				continue
			}
			// 解析 token
			bs, err := value.Bytes()
			if err != nil {
				sourceList = append(sourceList, source)
				continue
			}
			t := &cacheverify.UserTokenVerify{}
			err = json.Unmarshal(bs, t)
			if err != nil {
				sourceList = append(sourceList, source)
				continue
			}
			// 检测超时
			if int64(t.CreatedTimestamp)+int64(t.TTL) <= time.Now().Unix() {
				sourceList = append(sourceList, source)
				continue
			}
		}
	}

	if len(sourceList) > 0 {
		db.RedisHdel(ctx, key, sourceList...)
	}
}

func RemoveAllTokenOfUser(ctx context.Context, userId string) {
	tokenKey := cacheverify.NewVerifyKey(ctx, userId, "", "")
	key, _ := tokenKey.GetHashKey()
	db.RedisDel(tokenKey.Context(), key)
}