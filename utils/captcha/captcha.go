package captcha

import (
	"context"
	"bytes"
	"fmt"
	"time"
	
	cap "github.com/dchest/captcha"
	"github.com/fzzy/radix/redis"

	"github.com/chenjun-git/umbrella-user/db"
)

type CaptchaStore struct {
	ttl time.Duration
}

var Store CaptchaStore

const (
	CaptcharKeyPrefix = "captcha"
)

func (c CaptchaStore) Set(id string, digits []byte) {
	key := fmt.Sprintf("%s:%s", CaptcharKeyPrefix, id)
	if err := db.RedisSet(context.Background(), key, digits, c.ttl); err != nil {
		fmt.Printf("set captcha err, key : %s, digits : %s, err : %+v\n", key, digits, err)
		//log.Error("CaptchaStore.Set set id failed", zap.Error(err))
		return
	}
}

func (c CaptchaStore) Get(id string, clear bool) []byte {
	key := fmt.Sprintf("%s:%s", CaptcharKeyPrefix, id)
	reply, err := db.RedisGet(context.Background(), key)
	if err != nil {
		fmt.Printf("get captcha err, err : %+v\n", err)
		//log.Error("CaptchaStore.Get get key from redis failed", zap.Error(err))
		return nil
	}
	if reply.Type == redis.NilReply {
		return nil
	}

	digits, err := reply.Bytes()
	if err != nil {
		fmt.Printf("get captcha err, err : %+v\n", err)
		//log.Error("CaptchaStore.Get parse result failed", zap.Error(err))
		return nil
	}

	if clear {
		if err = db.RedisDel(context.Background(), key); err != nil {
			fmt.Printf("get captcha err, err : %+v\n", err)
			//log.Error("CaptchaStore.Get delete key failed", zap.Error(err))
		}
	}
	return digits
}

func InitCaptcha(ttl time.Duration) {
	if ttl > 0 {
		Store.ttl = ttl
	} else {
		Store.ttl = 0
	}
	cap.SetCustomStore(cap.Store(Store))
}

func GenCaptcha(length int) (id string) {
	return cap.NewLen(length)
}

func VerifyCaptcha(id, value string) bool {
	defer func() {
		cap.Reload(id)
	}()

	return cap.VerifyString(id, value)
}

func GetCaptchaImage(id string, width, height int) ([]byte, error) {
	var buf bytes.Buffer
	if err := cap.WriteImage(&buf, id, width, height); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// func GetCaptchaAudio(id string) ([]byte, error) {
// 	var buf bytes.Buffer
// 	if err := cap.WriteAudioWithoutNoise(&buf, id, "zh"); err != nil {
// 		return nil, err
// 	}
// 	return buf.Bytes(), nil
// }

func ReloadCaptcha(id string) {
	cap.Reload(id)
}
