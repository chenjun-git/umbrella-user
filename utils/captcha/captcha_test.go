package captcha

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"business/user/common"
	"business/user/db"
)

var id string

const (
	configPath = "../../conf/config.user.toml"
)

func init() {
	common.InitConfig(configPath)
	db.InitRedis(common.Config.Redis)
	InitCaptcha(common.Config.Captcha.TTL.D())
	fmt.Printf("Captcha conf : %+v\n", common.Config.Captcha)
}

func TestSetAndGet(t *testing.T) {
	assert := assert.New(t)

	Store.Set("test", []byte{1})
	data := Store.Get("test", false)
	assert.Equal([]byte{1}, data)
}

func TestGenCaptcha(t *testing.T) {
	assert := assert.New(t)

	id = GenCaptcha(common.Config.Captcha.DefaultLength)
	assert.Len(id, 20)
}

func TestVerifyCaptcha(t *testing.T) {
	assert := assert.New(t)

	val := Store.Get(id, false)
	assert.Len(val, common.Config.Captcha.DefaultLength)

	matched := VerifyCaptcha(id, "12345")
	assert.False(matched)
}

func TestGetCaptchaImage(t *testing.T) {
	assert := assert.New(t)

	_, err := GetCaptchaImage(id, common.Config.Captcha.DefaultWidth, 
		common.Config.Captcha.DefaultHeight)
	assert.NotNil(err)
}
