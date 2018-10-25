package v1_0

import (
	"fmt"
	"net/http"

	linq "github.com/ahmetb/go-linq"

	"github.com/chenjun-git/umbrella-user/common"
	"github.com/chenjun-git/umbrella-user/db"
	"github.com/chenjun-git/umbrella-user/model"
	"github.com/chenjun-git/umbrella-user/utils"
	"github.com/chenjun-git/umbrella-user/utils/cacheverify"
	"github.com/chenjun-git/umbrella-user/utils/captcha"
	"github.com/chenjun-git/umbrella-user/utils/render"
	"github.com/chenjun-git/umbrella-user/sms"
)

func genCodeAndSend(w http.ResponseWriter, r *http.Request, mediaType, phone, email, purpose, source string) (haveRender bool) {
	verifyKey := cacheverify.NewVerifyKey(r.Context(), getContactByPhoneEmail(mediaType, phone, email), purpose, source)
	verify, err := verifyKey.GetCodeVerify()
	if err != nil {
		render.JSON(w, r, http.StatusInternalServerError, render.M{
			"code":    common.AccountInternalError,
			"message": err.Error(),
		})
		return true
	} else if verify == nil {
		verifyCode := utils.RandomDigits(6)
		verify, err = verifyKey.CreateCodeVerify(verifyCode, common.Config.Verify.TTL.D())
		if err != nil {
			render.JSON(w, r, http.StatusInternalServerError, render.M{
				"code":    common.AccountInternalError,
				"message": err.Error(),
			})
			return true
		}
	}

	if purpose != common.TwoFactorAuthVerifyCodePurpose {
		// 除了二步验证，其他的有请求限制
		if verify.SendTimes >= common.Config.Verify.MaxSendTimes {
			render.JSON(w, r, http.StatusBadRequest, render.M{
				"code": common.AccountRequestLimit,
			})
			return true
		} else if verify.CheckTimes >= common.Config.Verify.MaxCheckTimes {
			render.JSON(w, r, http.StatusBadRequest, render.M{
				"code": common.AccountRequestLimit,
			})
			return true
		}
	}

	if common.Config.ReleaseMode {
		if common.MediaTypeEmail == mediaType {
			// 发送邮箱校验码: 邮箱验证码在生存周期内，不管请求发送几次，都使用同一个验证码，产品需求!
			//_, err = email.SendVerifyUrlAndCode(orgSetting.Email, purpose, orgSetting.Name, userEmail, common.MeiqiaWebSiteUrl, verify.VerifyCode)
		} else {
			// 发送短信校验码: 短信验证码在生存周期内，不管请求发送几次，都使用同一个验证码，产品需求!
			err = sms.SendPhoneMsg(common.Config.Sms.Addr, verifyKey.ID, verify.VerifyCode, verifyKey.ID)
		}
		if err != nil {
			render.JSON(w, r, http.StatusInternalServerError, render.M{
				"code":    common.AccountInternalError,
				"message": err.Error(),
			})
			return true
		}
	}

	if purpose != common.TwoFactorAuthVerifyCodePurpose {
		// 除了二步验证，其他的有请求限制
		verify.SendTimes++
		verify.Save()
	}

	if !common.Config.ReleaseMode {
		render.JSON(w, r, http.StatusOK, render.M{
			"code":        common.OK,
			"verify_code": verify.VerifyCode,
			"body": render.M{
				"verify_code": verify.VerifyCode,
			},
		})
		return true
	}

	return false
}

func SendVerifyCodeHandler(w http.ResponseWriter, r *http.Request) {
	var req SendVerifyCodeRequest
	fmt.Printf("send")
	if err := utils.BindJSON(r, &req); err != nil {
		render.JSON(w, r, http.StatusBadRequest, render.M{
			"code":    common.AccountBindFailed,
			"message": err.Error,
		})
		return
	}

	if !linq.From(common.PurposeRange).Contains(req.Purpose) {
		render.JSON(w, r, http.StatusBadRequest, render.M{
			"code": common.AccountInvalidPurpose,
		})
		return
	}

	if !linq.From(common.SourceRange).Contains(req.Source) {
		render.JSON(w, r, http.StatusBadRequest, render.M{
			"code": common.AccountInvalidSource,
		})
		return
	}

	if !ValidPhoneEmail(w, r, req.Phone, req.Email, req.MediaType) {
		return
	}

	// todo:user, logined := checkLogin(r)
	logined := false // todo
	var user = &model.User{}
	user = nil
	var err error

	// 验证登录状态
	if !logined {
		switch req.Purpose {
		case common.TwoFactorAuthSetupPurpose, common.UpdateEmailPurpose, common.UpdatePhonePurpose:
			// 开启二次验证的手机号不需要绑定，但是需要登录
			render.JSON(w, r, http.StatusBadRequest, render.M{
				"code": common.AccountInvalidToken,
			})
			return
		default:
			sqlExec := db.BindDBerWithContext(r.Context(), db.MySQL)
			user, err = model.GetUserByPhoneEmail(sqlExec, req.MediaType, req.Phone, req.Email)
			if err != nil {
				render.JSON(w, r, http.StatusInternalServerError, render.M{
					"code":    common.AccountInternalError,
					"message": err.Error(),
				})
				return
			}
		}

		// 图形验证码校验
		// 登录状态不需要图形验证码
		if common.Config.ReleaseMode {
			if !captcha.VerifyCaptcha(req.CaptchaToken, req.CaptchaValue) {
				render.JSON(w, r, http.StatusBadRequest, render.M{
					"code": common.AccountCaptchaNotMatch,
				})
				return
			}
		}
	}

	// 二次验证：账号存在与否无所谓
	// 注册账号: 账号不能存在
	// 其他：账号不能不存在
	switch req.Purpose {
	case common.TwoFactorAuthSetupPurpose:
		// pass
	case common.SignupPurpose:
		if user != nil {
			render.JSON(w, r, http.StatusBadRequest, render.M{
				"code": common.AccountAccountAlreadyExist,
			})
			return
		}
	default:
		if user == nil {
			render.JSON(w, r, http.StatusBadRequest, render.M{
				"code": common.AccountAccountNotExist,
			})
			return
		}
	}

	haveRender := genCodeAndSend(w, r, req.MediaType, req.Phone, req.Email, req.Purpose, req.Source)
	if haveRender {
		return
	}

	render.JSON(w, r, http.StatusOK, render.M{
		"code": common.OK,
	})
	return
}

func CheckVerifyCodeHandler(w http.ResponseWriter, r *http.Request) {
	var req CheckVerifyCodeRequest
	if err := utils.BindJSON(r, &req); err != nil {
		render.JSON(w, r, http.StatusBadRequest, render.M{
			"code":    common.AccountInternalError,
			"message": err.Error(),
		})
		return
	}
	// purpose校验
	if !linq.From(common.PurposeRange).Contains(req.Purpose) {
		render.JSON(w, r, http.StatusBadRequest, render.M{
			"code": common.AccountInvalidPurpose,
		})
		return
	}
	// source校验
	if !linq.From(common.SourceRange).Contains(req.Source) {
		render.JSON(w, r, http.StatusBadRequest, render.M{
			"code": common.AccountInvalidSource,
		})
		return
	}

	if !ValidPhoneEmail(w, r, req.Phone, req.Email, req.MediaType) {
		return
	}

	verifyKey := cacheverify.NewVerifyKey(r.Context(), getContactByPhoneEmail(req.MediaType, req.Phone, req.Email), req.Purpose, req.Source)

	verify, err := verifyKey.GetCodeVerify()
	if err != nil {
		render.JSON(w, r, http.StatusInternalServerError, render.M{
			"code":    common.AccountInternalError,
			"message": err.Error(),
		})
		return
	} else if verify == nil {
		render.JSON(w, r, http.StatusBadRequest, render.M{
			"code": common.AccountVerifyCodeNotMatch,
		})
		return
	}

	if verify.CheckTimes >= common.Config.Verify.MaxCheckTimes-1 { // -1 是为其他需要验证码的业务接口(如 signup/reset_password)预留一次
		render.JSON(w, r, http.StatusBadRequest, render.M{
			"code": common.AccountRequestLimit,
		})
		return
	}

	ok, err := verify.Verify(req.VerifyCode)
	if err != nil {
		render.JSON(w, r, http.StatusInternalServerError, render.M{
			"code":    common.AccountInternalError,
			"message": err.Error(),
		})
		return
	} else if !ok {
		render.JSON(w, r, http.StatusBadRequest, render.M{
			"code": common.AccountVerifyCodeNotMatch,
		})
		return
	}

	render.JSON(w, r, http.StatusOK, render.M{
		"code": common.OK,
	})
	return
}
