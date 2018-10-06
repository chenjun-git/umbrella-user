package v1_0

import (
	"net/http"

	linq "github.com/ahmetb/go-linq"

	"business/user/common"
	"business/user/utils"
)

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

	userID, err := utils.GenID(
	if err != nil {
		render.JSON(w, r, http.StatusInternalServerError, render.M{
			"code":    common.AccountGenerateIdFailed,
			"message": "GenID: " + err.Error(),
		})
		return
	}


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
}

func SignOutHandler(w http.ResponseWriter, r *http.Request) {
}

func GetCaptchaImageHandler(w http.ResponseWriter, r *http.Request) {
}

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {

}