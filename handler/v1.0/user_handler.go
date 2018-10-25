package v1_0

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	linq "github.com/ahmetb/go-linq"
	"github.com/satori/go.uuid"

	commonErrors "github.com/chenjun-git/umbrella-common/errors"
	"github.com/chenjun-git/umbrella-common/lang"
	"github.com/chenjun-git/umbrella-common/token"

	"github.com/chenjun-git/umbrella-user/common"
	"github.com/chenjun-git/umbrella-user/db"
	"github.com/chenjun-git/umbrella-user/middleware"
	"github.com/chenjun-git/umbrella-user/model"
	innerToken "github.com/chenjun-git/umbrella-user/token"
	"github.com/chenjun-git/umbrella-user/utils"
	"github.com/chenjun-git/umbrella-user/utils/cacheverify"
	"github.com/chenjun-git/umbrella-user/utils/captcha"
	"github.com/chenjun-git/umbrella-user/utils/password"
	"github.com/chenjun-git/umbrella-user/utils/render"
)

func checkVerifyCode(ctx context.Context, mediaType, id, source, purpose, code string) (bool, error) {
	// if mediaType == common.MediaTypeTOTP {
	// 	secret, err := encrypt.CBCDecrypt(id)
	// 	if err != nil {
	// 		return false, err
	// 	}
	// 	return gotp.NewDefaultTOTP(secret).Verify(code, int(time.Now().Unix())), nil
	// }

	verify, err := cacheverify.NewVerifyKey(ctx, id, purpose, source).GetCodeVerify()
	if err != nil {
		return false, err
	} else if verify == nil {
		return false, nil
	} else {
		if verify.CheckTimes >= common.Config.Verify.MaxCheckTimes {
			return false, commonErrors.NewError(common.AccountVerifyCodeNotMatch, "verify code expired")
		}
	}

	ok, err := verify.Verify(code)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest
	if err := utils.BindJSON(r, &req); err != nil {
		render.JSON(w, r, http.StatusBadRequest, render.M{
			"code":    common.AccountBindFailed,
			"message": err.Error(),
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

	strId := getContactByPhoneEmail(req.MediaType, req.Phone, req.Email)

	if common.Config.ReleaseMode {
		ok, err := checkVerifyCode(r.Context(), req.MediaType, strId, req.Source, common.SignupPurpose, req.VerifyCode)
		if err != nil {
			render.JSON(w, r, http.StatusInternalServerError, render.M{
				"code":    common.AccountVerifyCodeNotMatch,
				"message": err.Error(),
			})
			return
		} else if !ok {
			render.JSON(w, r, http.StatusBadRequest, render.M{
				"code": common.AccountVerifyCodeNotMatch,
			})
			return
		}
	}

	hashedPassword, err := password.Encrypt(req.Password)
	if err != nil {
		render.JSON(w, r, http.StatusInternalServerError, render.M{
			"code":    common.AccountInternalError,
			"message": err.Error(),
		})
		return
	}

	passwordLevel := password.PasswordStrength(req.Password)
	if passwordLevel == password.LevelIllegal {
		render.JSON(w, r, http.StatusBadRequest, render.M{
			"code": common.AccountPasswordLevelIllegal,
		})
		return
	}

	user := model.User{
		UserName:       strId,
		Tel:            utils.StringToNullString(req.Phone),
		Email:          utils.StringToNullString(req.Email),
		RegisterSource: req.Source,
		HashPasswd:     hashedPassword,
		PasswdLevel:    passwordLevel,
	}

	userID, err := utils.GenID()
	if err != nil {
		render.JSON(w, r, http.StatusInternalServerError, render.M{
			"code":    common.AccountGenerateIdFailed,
			"message": "GenID: " + err.Error(),
		})
		return
	}

	sqlExec := db.BindDBerWithContext(r.Context(), db.MySQL)
	if common.MediaTypeEmail == req.MediaType {
		_, err = model.CreateUserByEmail(sqlExec, userID, user)
	} else {
		_, err = model.CreateUserByPhone(sqlExec, userID, user)
	}

	if err != nil {
		render.JSON(w, r, http.StatusInternalServerError, render.M{
			"code":    common.AccountDBError,
			"message": "CreateAccount: " + err.Error(),
		})
		return
	}

	render.JSON(w, r, http.StatusOK, render.M{
		"code": common.OK,
	})
}

func SignUpWithoutVerifyCodeHandler(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest
	if err := utils.BindJSON(r, &req); err != nil {
		render.JSON(w, r, http.StatusBadRequest, render.M{
			"code":    common.AccountBindFailed,
			"message": err.Error(),
		})
		return
	}
	if !assertInternalSource(req.Source) {
		render.JSON(w, r, http.StatusBadRequest, render.M{
			"code": common.AccountInvalidSource,
		})
		return
	}

	if !ValidPhoneEmail(w, r, req.Phone, req.Email, req.MediaType) {
		return
	}

	hashedPassword, err := password.Encrypt(req.Password)
	if err != nil {
		render.JSON(w, r, http.StatusInternalServerError, render.M{
			"code":    common.AccountInternalError,
			"message": err.Error(),
		})
		return
	}

	passwordLevel := password.PasswordStrength(req.Password)
	if passwordLevel == password.LevelIllegal {
		render.JSON(w, r, http.StatusBadRequest, render.M{
			"code": common.AccountPasswordLevelIllegal,
		})
		return
	}

	strId := getContactByPhoneEmail(req.MediaType, req.Phone, req.Email)

	user := model.User{
		UserName:       strId,
		Tel:            utils.StringToNullString(req.Phone),
		Email:          utils.StringToNullString(req.Email),
		RegisterSource: req.Source,
		HashPasswd:     hashedPassword,
		PasswdLevel:    passwordLevel,
	}

	userID, err := utils.GenID()
	if err != nil {
		render.JSON(w, r, http.StatusInternalServerError, render.M{
			"code":    common.AccountGenerateIdFailed,
			"message": "GenID: " + err.Error(),
		})
		return
	}

	sqlExec := db.BindDBerWithContext(r.Context(), db.MySQL)
	if common.MediaTypeEmail == req.MediaType {
		_, err = model.CreateUserByEmail(sqlExec, userID, user)
	} else {
		_, err = model.CreateUserByPhone(sqlExec, userID, user)
	}

	if err != nil {
		render.JSON(w, r, http.StatusInternalServerError, render.M{
			"code":    common.AccountDBError,
			"message": "CreateAccount: " + err.Error(),
		})
		return
	}
	render.JSON(w, r, http.StatusOK, render.M{
		"code": common.OK,
		"body": map[string]interface{}{
			"id": userID,
		},
	})
	return
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	var req SignInRequest
	if err := utils.BindJSON(r, &req); err != nil {
		render.JSON(w, r, http.StatusBadRequest, render.M{
			"code":    common.AccountBindFailed,
			"message": err.Error(),
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

	sqlExec := db.BindDBerWithContext(r.Context(), db.MySQL)
	user, err := model.GetUserByPhoneEmail(sqlExec, req.MediaType, req.Phone, req.Email)
	if err != nil {
		render.JSON(w, r, http.StatusInternalServerError, render.M{
			"code":    common.AccountDBError,
			"message": fmt.Sprintf("GetAccount %s: %s", req.MediaType, err),
		})
		return
	} else if user == nil {
		render.JSON(w, r, http.StatusBadRequest, render.M{
			"code": common.AccountAccountNotExist,
		})
		return
	}
	//  else if account.EnablePasswordLifetime && int(time.Now().Sub(*account.LastUpdatePassword).Seconds()) >= account.MaxPasswordLifetime {
	// 	render.JSON(w, r, http.StatusBadRequest, render.M{
	// 		"code": common.AccountPasswordExpired,
	// 	})
	// 	return
	// }

	verifyKey := cacheverify.NewVerifyKey(r.Context(), getContactByPhoneEmail(req.MediaType, req.Phone, req.Email), common.SigninPurpose, req.Source)
	login, err := verifyKey.GetLoginVerify()
	if err != nil {
		render.JSON(w, r, http.StatusInternalServerError, render.M{
			"code":    common.AccountRedisError,
			"message": "GetLogin: " + err.Error(),
		})
		return
	} else if login == nil {
		login, err = verifyKey.CreateLoginVerify(common.Config.Login.TTL.D())
		if err != nil || login == nil {
			render.JSON(w, r, http.StatusInternalServerError, render.M{
				"code":    common.AccountRedisError,
				"message": "CreateLogin: " + err.Error(),
			})
			return
		}
	}

	languageCode := strings.ToUpper(lang.FromOutgoingContext(r.Context())[0])
	var errformat string
	// 输错密码超过指定次数之后冻结
	if login.RequestTimes >= common.Config.Login.MaxRequestTimes {
		if languageCode == "ZH-CN" {
			errformat = "账号冻结%d分钟，请稍后重试"
		} else {
			errformat = "Account suspenced for %d minutes, please try again later."
		}
		render.JSON(w, r, http.StatusBadRequest, render.M{
			"code":    common.AccountWasFreezed,
			"message": fmt.Sprintf(errformat, int(common.Config.Login.TTL.D().Minutes())),
			// TODO: 前端调整后删掉以下部分
			"body": render.M{
				"freeze_time": int(common.Config.Login.TTL.D().Minutes()),
			},
		})
		return
	}

	needCaptcha := login.RequestTimes >= common.Config.Login.MaxCaptchaTImes
	if needCaptcha {
		leftTimes := login.GetLeftTimes()
		if req.CaptchaToken == "" || req.CaptchaValue == "" {
			render.JSON(w, r, http.StatusBadRequest, render.M{
				"code": common.AccountNeedCaptcha,
				// TODO: 前端调整后删掉以下部分
				"left_times": leftTimes,
				"body": render.M{
					"left_times": leftTimes,
					"max_times":  common.Config.Login.MaxRequestTimes,
				},
			})
			return
		}

		if !captcha.VerifyCaptcha(req.CaptchaToken, req.CaptchaValue) {
			render.JSON(w, r, http.StatusBadRequest, render.M{
				"code": common.AccountCaptchaNotMatch,
				// TODO: 前端调整后删掉以下部分
				"left_times": leftTimes,
				"body": render.M{
					"left_times": leftTimes,
					"max_times":  common.Config.Login.MaxRequestTimes,
				},
			})
			return
		}
	}

	login.RequestTimes++
	if err := password.Verify(req.Password, user.HashPasswd); err != nil {
		login.Save()
		leftTimes := login.GetLeftTimes()

		// 第五次密码错误
		if leftTimes == 0 {
			if languageCode == "ZH-CN" {
				errformat = "账号冻结%d分钟，请稍后重试"
			} else {
				errformat = "Account suspenced for %d minutes, please try again later."
			}
			render.JSON(w, r, http.StatusBadRequest, render.M{
				"code":    common.AccountOrPasswordError,
				"message": fmt.Sprintf(errformat, int(common.Config.Login.TTL.D().Minutes())),
				// TODO: 前端调整后删掉以下部分
				"body": render.M{
					"left_times":  leftTimes,
					"max_times":   common.Config.Login.MaxRequestTimes,
					"freeze_time": int(common.Config.Login.TTL.D().Minutes()),
				},
			})
			return
		}
		// 第一次密码错误
		if login.RequestTimes == 1 {
			render.JSON(w, r, http.StatusBadRequest, render.M{
				"code": common.AccountOrPasswordError,
				// TODO: 前端调整后删掉以下部分
				"body": render.M{
					"left_times":  leftTimes,
					"max_times":   common.Config.Login.MaxRequestTimes,
					"freeze_time": int(common.Config.Login.TTL.D().Minutes()),
				},
			})
			return
		}

		if languageCode == "ZH-CN" {
			errformat = "账号或密码错误，你还有%d次机会"
		} else {
			errformat = "Account or password error，%d chances left."
		}
		render.JSON(w, r, http.StatusBadRequest, render.M{
			"code":    common.AccountOrPasswordError,
			"message": fmt.Sprintf(errformat, leftTimes),
			// TODO: 前端调整后删掉以下部分
			"left_times": leftTimes,
			"body": render.M{
				"left_times": leftTimes,
				"max_times":  common.Config.Login.MaxRequestTimes,
			},
		})
		return
	}

	// 二步验证
	// twofa, comerr := twoFactorAuthVerify(r.Context(), sqlExec, account.ID, req.Source)
	// if comerr != nil {
	// 	render.JSON(w, r, http.StatusOK, render.M{
	// 		"code":    comerr.GetCode(),
	// 		"message": comerr.GetDescription(),
	// 		"body":    twofa,
	// 	})
	// 	return
	// }

	login.Clean()
	accessToken, err2 := genTokenAndStore(r.Context(), user.ID, "", req.Source)
	if err2 != nil {
		render.JSON(w, r, http.StatusInternalServerError, render.M{
			"code":    err2.GetCode(),
			"message": err2.GetDescription(),
		})
		return
	}

	render.JSON(w, r, http.StatusOK, render.M{
		"code":  common.OK,
		"token": accessToken,
		"body": render.M{
			"token": accessToken,
		},
	})
}

func SignInWithoutVerifyCodeHandler(w http.ResponseWriter, r *http.Request) {
	var req SignInRequest
	if err := utils.BindJSON(r, &req); err != nil {
		render.JSON(w, r, http.StatusBadRequest, render.M{
			"code":    common.AccountBindFailed,
			"message": err.Error(),
		})
		return
	}

	if !assertInternalSource(req.Source) {
		render.JSON(w, r, http.StatusBadRequest, render.M{
			"code": common.AccountInvalidSource,
		})
		return
	}

	if !ValidPhoneEmail(w, r, req.Phone, req.Email, req.MediaType) {
		return
	}

	sqlExec := db.BindDBerWithContext(r.Context(), db.MySQL)
	user, err := model.GetUserByPhoneEmail(sqlExec, req.MediaType, req.Phone, req.Email)

	if err != nil {
		render.JSON(w, r, http.StatusInternalServerError, map[string]interface{}{
			"code":    common.AccountDBError,
			"message": "GetAccount" + req.MediaType + ":" + err.Error(),
		})
		return
	} else if user == nil {
		render.JSON(w, r, http.StatusBadRequest, render.M{
			"code": common.AccountAccountNotExist,
		})
		return
	}

	if password.Verify(req.Password, user.HashPasswd) != nil {
		render.JSON(w, r, http.StatusBadRequest, render.M{
			"code": common.AccountOrPasswordError,
		})
		return
	}

	accessToken, err2 := genTokenAndStore(r.Context(), user.ID, "", req.Source)
	if err2 != nil {
		render.JSON(w, r, http.StatusInternalServerError, render.M{
			"code":    err2.GetCode(),
			"message": err2.GetDescription(),
		})
		return
	}

	render.JSON(w, r, http.StatusOK, render.M{
		"code":  common.OK,
		"token": accessToken,
		"body": render.M{
			"token": accessToken,
		},
	})
	return
}

func SignOutHandler(w http.ResponseWriter, r *http.Request) {
}

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest

	if err := utils.BindJSON(r, &req); err != nil {
		render.JSON(w, r, http.StatusBadRequest, render.M{
			"code":    common.AccountBindFailed,
			"message": err.Error(),
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

	strId := getContactByPhoneEmail(req.MediaType, req.Phone, req.Email)

	if common.Config.ReleaseMode {
		ok, err := checkVerifyCode(r.Context(), req.MediaType, strId, req.Source, common.ResetPasswordPurpose, req.VerifyCode)
		if err != nil {
			render.JSON(w, r, http.StatusInternalServerError, render.M{
				"code":    common.AccountVerifyCodeNotMatch,
				"message": err.Error(),
			})
			return
		} else if !ok {
			render.JSON(w, r, http.StatusBadRequest, render.M{
				"code": common.AccountVerifyCodeNotMatch,
			})
			return
		}
	}

	sqlExec := db.BindDBerWithContext(r.Context(), db.MySQL)
	user, err := model.GetUserByPhoneEmail(sqlExec, req.MediaType, req.Phone, req.Email)
	if err != nil {
		render.JSON(w, r, http.StatusInternalServerError, render.M{
			"code":    common.AccountDBError,
			"message": "GetAccountByEmail: " + err.Error(),
		})
		return
	} else if user == nil {
		render.JSON(w, r, http.StatusBadRequest, render.M{
			"code": common.AccountAccountNotExist,
		})
		return
	}

	hashedPassword, err := password.Encrypt(req.NewPassword)
	if err != nil {
		render.JSON(w, r, http.StatusInternalServerError, render.M{
			"code":    common.AccountInternalError,
			"message": err.Error(),
		})
		return
	}
	if password.Verify(req.NewPassword, user.HashPasswd) == nil {
		render.JSON(w, r, http.StatusInternalServerError, render.M{
			"code": common.AccountPasswordSameWithOld,
		})
		return
	}
	passwordLevel := password.PasswordStrength(req.NewPassword)
	if passwordLevel == password.LevelIllegal {
		render.JSON(w, r, http.StatusBadRequest, render.M{
			"code": common.AccountPasswordLevelIllegal,
		})
		return
	}

	if _, err := model.UpdatePassword(sqlExec, user.ID, hashedPassword, passwordLevel); err != nil {
		render.JSON(w, r, http.StatusInternalServerError, render.M{
			"code":    common.AccountDBError,
			"message": "UpdatePassword: " + err.Error(),
		})
		return
	}

	// 清除本帐户所有的account token, 及access token/refresh token
	//token.RemoveAllTokenOfAccount(r.Context(), account.ID)
	//tenant.RemoveAllTokenOfAccount(r.Context(), account.ID)

	render.JSON(w, r, http.StatusOK, render.M{
		"code": common.OK,
	})
}

func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, http.StatusOK, render.M{
		"code": common.OK,
	})
}

func VerifyTokenHandler(w http.ResponseWriter, r *http.Request) {
	accessToken := r.URL.Query().Get("token")
	user, _, err := middleware.CheckUserToken(r.Context(), accessToken)
	if err != nil {
		render.JSON(w, r, http.StatusOK, render.M{
			"code":    common.OK,
			"message": err.Error(),
			"body": render.M{
				"result":  false,
				"user_id": "",
				//"device":  "",
			},
		})
		return
	}

	render.JSON(w, r, http.StatusOK, render.M{
		"code": common.OK,
		"body": render.M{
			"result":  true,
			"user_id": user.ID,
			//"device":     payload.LoginSource,
		},
	})
}

func GetTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID string `json:"user_id"`
		Source string `json:"source"`
	}
	if err := utils.BindJSON(r, &req); err != nil {
		render.JSON(w, r, http.StatusBadRequest, render.M{
			"code":    common.AccountBindFailed,
			"message": err.Error(),
		})
		return
	}

	if req.UserID == "" {
		render.JSON(w, r, http.StatusBadRequest, render.M{
			"code": common.AccountInvalidAccountId,
		})
		return
	}
	if req.Source != "" {
		if !linq.From(common.SourceRange).Contains(req.Source) {
			render.JSON(w, r, http.StatusBadRequest, render.M{
				"code": common.AccountInvalidSource,
			})
			return
		}
	} else {
		req.Source = common.InternalSource
	}

	sqlExec := db.BindDBerWithContext(r.Context(), db.MySQL)
	if user, err := model.GetUserById(sqlExec, req.UserID); err != nil {
		render.JSON(w, r, http.StatusInternalServerError, render.M{
			"code":    common.AccountInternalError,
			"message": err.Error(),
		})
		return
	} else if user == nil {
		render.JSON(w, r, http.StatusBadRequest, render.M{
			"code": common.AccountAccountNotExist,
		})
		return
	}

	accessToken, err2 := genTokenAndStore(r.Context(), req.UserID, "", req.Source)
	if err2 != nil {
		render.JSON(w, r, http.StatusInternalServerError, render.M{
			"code":    err2.GetCode(),
			"message": err2.GetDescription(),
		})
		return
	}

	render.JSON(w, r, http.StatusOK, render.M{
		"code":  common.OK,
		"token": accessToken,
		"body": render.M{
			"token": accessToken,
		},
	})
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		render.JSON(w, r, http.StatusBadRequest, render.M{
			"code":    common.AccountBindFailed,
			"message": "empty id",
		})
		return
	}

	sqlExec := db.BindDBerWithContext(r.Context(), db.MySQL)
	user, err := model.GetUserById(sqlExec, id)
	if err != nil {
		render.JSON(w, r, http.StatusInternalServerError, render.M{
			"code":    common.AccountInternalError,
			"message": err.Error(),
		})
		return
	} else if user == nil {
		render.JSON(w, r, http.StatusBadRequest, render.M{
			"code": common.AccountAccountNotExist,
		})
		return
	}

	info := render.M{
		"phone":                utils.NullStringToString(user.Tel),
		"email":                utils.NullStringToString(user.Email),
		"passwd_level":         user.PasswdLevel,
		"name":                 user.UserName,
		"register_source":      user.RegisterSource,
		"last_update_password": user.LastUpdatePasswd,
	}

	render.JSON(w, r, http.StatusOK, render.M{
		"code": common.OK,
		"user": info,
		"body": render.M{
			"user": info,
		},
	})
	return
}

func CheckExistencHandler(w http.ResponseWriter, r *http.Request) {
	phone := r.URL.Query().Get("phone")
	email := r.URL.Query().Get("email")

	var (
		field string
		value string
	)
	if phone != "" {
		field = "tel"
		value = phone
	} else if email != "" {
		field = "email"
		value = email
	} else {
		render.JSON(w, r, http.StatusOK, render.M{
			"code": common.AccountBindFailed,
		})
		return
	}

	sqlExec := db.BindDBerWithContext(r.Context(), db.MySQL)
	count, err := model.Count(sqlExec, map[string]interface{}{
		field: value,
	})
	if err != nil {
		render.JSON(w, r, http.StatusOK, render.M{
			"code":    common.AccountDBError,
			"message": err.Error(),
		})
		return
	}
	render.JSON(w, r, http.StatusOK, render.M{
		"code": common.OK,
		"body": render.M{
			"existence": count > 0,
		},
	})
}

func UpdateEmailHandler(w http.ResponseWriter, r *http.Request) {
}

func ForceChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req = struct {
		UserID      string `json:"user_id" binding:"required,len=32"`
		NewPassword string `json:"new_password" binding:"required,min=8,max=20"`
	}{}

	if err := utils.BindJSON(r, &req); err != nil {
		render.JSON(w, r, http.StatusOK, render.M{
			"code":    common.AccountBindFailed,
			"message": err.Error(),
		})
		return
	}
	passwordLevel := password.PasswordStrength(req.NewPassword)
	if passwordLevel == password.LevelIllegal {
		render.JSON(w, r, http.StatusOK, render.M{
			"code": common.AccountPasswordLevelIllegal,
		})
		return
	}
	hashedPassword, err := password.Encrypt(req.NewPassword)
	if err != nil {
		render.JSON(w, r, http.StatusOK, render.M{
			"code":    common.AccountInternalError,
			"message": err.Error(),
		})
		return
	}
	// 验证 account ID 有效性
	sqlExec := db.BindDBerWithContext(r.Context(), db.MySQL)
	user, err := model.GetUserById(sqlExec, req.UserID)
	if err != nil {
		render.JSON(w, r, http.StatusOK, render.M{
			"code":    common.AccountDBError,
			"message": err.Error(),
		})
		return
	}
	if user == nil {
		render.JSON(w, r, http.StatusOK, render.M{
			"code": common.AccountAccountNotExist,
		})
		return
	}

	// 修改密码
	if _, err := model.UpdatePassword(sqlExec, user.ID, hashedPassword, passwordLevel); err != nil {
		render.JSON(w, r, http.StatusOK, render.M{
			"code":    common.AccountDBError,
			"message": "UpdatePassword: " + err.Error(),
		})
		return
	}
	// 强制 account 下线
	//token.RemoveAllTokenOfAccount(r.Context(), account.ID)
	//tenant.RemoveAllTokenOfAccount(r.Context(), account.ID)

	render.JSON(w, r, http.StatusOK, render.M{
		"code": common.OK,
	})
	return
}

func ForceUpdateEmailHandler(w http.ResponseWriter, r *http.Request) {
	var req = struct {
		UserID string `json:"user_id" binding:"required,len=32"`
		Email  string `json:"email" binding:"required"`
		Source string `json:"source" binding:"required"`
	}{}
	if err := utils.BindJSON(r, &req); err != nil {
		render.JSON(w, r, http.StatusOK, render.M{
			"code":    common.AccountBindFailed,
			"message": err.Error(),
		})
		return
	}

	// 检查source
	if !assertInternalSource(req.Source) {
		render.JSON(w, r, http.StatusOK, render.M{
			"code": common.AccountInvalidSource,
		})
		return
	}
	// 检查email格式
	if !ValidPhoneEmail(w, r, "", req.Email, common.MediaTypeEmail) {
		return
	}
	// 验证 account ID 有效性
	sqlExec := db.BindDBerWithContext(r.Context(), db.MySQL)
	user, err := model.GetUserById(sqlExec, req.UserID)
	if err != nil {
		render.JSON(w, r, http.StatusOK, render.M{
			"code":    common.AccountDBError,
			"message": err.Error(),
		})
		return
	}
	if user == nil {
		render.JSON(w, r, http.StatusOK, render.M{
			"code": common.AccountAccountNotExist,
		})
		return
	}
	// 更新email
	if _, err := model.UpdateEmail(sqlExec, user.ID, req.Email); err != nil {
		render.JSON(w, r, http.StatusOK, render.M{
			"code":    common.AccountDBError,
			"message": "UpdateEmail: " + err.Error(),
		})
		return
	}

	render.JSON(w, r, http.StatusOK, render.M{
		"code": common.OK,
	})
	return
}

func GetPhoneHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		render.JSON(w, r, http.StatusBadRequest, render.M{
			"code":    common.AccountBindFailed,
			"message": "empty id",
		})
		return
	}

	sqlExec := db.BindDBerWithContext(r.Context(), db.MySQL)
	user, err := model.GetUserById(sqlExec, id)
	if err != nil {
		render.JSON(w, r, http.StatusInternalServerError, render.M{
			"code":    common.AccountInternalError,
			"message": err.Error(),
		})
		return
	} else if user == nil {
		render.JSON(w, r, http.StatusBadRequest, render.M{
			"code": common.AccountAccountNotExist,
		})
		return
	}

	render.JSON(w, r, http.StatusOK, render.M{
		"code":  common.OK,
		"phone": utils.NullStringToString(user.Tel),
		"body": render.M{
			"phone": utils.NullStringToString(user.Tel),
		},
	})
	return
}

func genTokenAndStore(ctx context.Context, userId, purpose, source string) (string, commonErrors.Error) {
	// TODO 不太理解修改source
	if source != common.AppSource && purpose != common.TwoFactorAuthVerifyTokenKey {
		u1, err := uuid.NewV4()
		if err != nil {
			return "", commonErrors.NewError(common.AccountInternalError, err.Error())
		}
		source += u1.String()
	}

	accessToken, err := token.EncryptAccessToken(common.Config.Token.TokenLibVersion, &token.Token{
		IssueTime: uint32(uint64(time.Now().Unix())),
		TTL:       uint16(common.Config.Token.AccessTokenTTL.Minutes()),
		UserID:    userId,
	})
	if err != nil {
		return "", commonErrors.NewError(common.AccountInternalError, err.Error())
	}

	tokenKey := cacheverify.NewVerifyKey(ctx, userId, purpose, source)
	if purpose == common.TwoFactorAuthVerifyTokenKey {
		_, err = tokenKey.CreateStringVerify(accessToken, common.Config.Token.AccessTokenTTL.D())
	} else {
		number, err := tokenKey.GetTokenNum()
		if err != nil {
			return "", commonErrors.NewError(common.AccountRedisError, err.Error())
		}

		if number > common.Config.Login.CleanInvalidTokenThreshold {
			go innerToken.RemoveInvalidToken(tokenKey.Context(), tokenKey.ID)
		}

		if number > common.Config.Login.MaxNumberLogin {
			return "", commonErrors.NewError(common.AccountTooMuchBoardingDevice, "")
		}

		_, err = tokenKey.CreateUserTokenVerify(accessToken, common.Config.Token.AccessTokenTTL.D())
	}

	if err != nil {
		return "", commonErrors.NewError(common.AccountRedisError, err.Error())
	}

	return accessToken, nil
}

func assertInternalSource(source string) bool {
	return source == common.InternalSource
}
