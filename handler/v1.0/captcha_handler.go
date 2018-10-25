package v1_0

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/chenjun-git/umbrella-user/common"
	"github.com/chenjun-git/umbrella-user/utils/captcha"
	"github.com/chenjun-git/umbrella-user/utils/render"
)

func GetCaptchaValue(token string) string {
	bs := captcha.Store.Get(token, false)
	for i, _ := range bs {
		bs[i] += 48
	}
	fmt.Printf("bs:%s, len:%d\n", bs, len(bs))
	return string(bs)
}

func CreateCaptchaHandler(w http.ResponseWriter, r *http.Request) {
	token := captcha.GenCaptcha(common.Config.Captcha.DefaultLength)
	render.JSON(w, r, http.StatusOK, render.M{
		"code":          common.OK,
		"captcha_token": token,
		"body": render.M{
			"captcha_token": token,
		},
	})
	return
}

func GetCaptchaImageHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("captcha_token")

	width, err := strconv.ParseInt(r.URL.Query().Get("width"), 10, 64)
	if err != nil {
		width = int64(common.Config.Captcha.DefaultWidth)
	}

	height, err := strconv.ParseInt(r.URL.Query().Get("height"), 10, 64)
	if err != nil {
		height = int64(common.Config.Captcha.DefaultHeight)
	}

	data, err := captcha.GetCaptchaImage(token, int(width), int(height))
	if err != nil {
		render.JSON(w, r, http.StatusInternalServerError, render.M{
			"code":    common.AccountInternalError,
			"message": "GetCaptchaImage " + err.Error(),
		})
		return
	}

	if !common.Config.ReleaseMode {
		render.JSON(w, r, http.StatusOK, render.M{
			"code":          common.OK,
			"captcha_value": GetCaptchaValue(token),
			"body": render.M{
				"captcha_value": GetCaptchaValue(token),
			},
		})
		return
	}

	render.PNG(w, r, http.StatusOK, data)
	return
}
