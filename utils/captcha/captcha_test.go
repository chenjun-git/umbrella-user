package captcha

import (
	"github.com/stretchr/testify/assert"
)

var id string

func TestGenCaptcha(t *testing.T) {
	assert := assert.New(t)

	id = GenCaptcha(5)
	assert.Len(id, 5)
}

func TestVerifyCaptcha(t *testing.T) {
	assert := assert.New(t)

	matched, err := VerifyCaptcha(id, "12345")
	assert.NotNil(err)
	assert.False(matched)
}

func TestGetCaptchaImage(t *testing.T) {
	assert := assert.New(t)

	img, err := GetCaptchaImage(id)
	assert.NotNil(err)
	assert.Len(img, 5)
}