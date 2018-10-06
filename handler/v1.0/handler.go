package v1_0

import (
	"net/http"

	"business/user/common"
	"business/user/utils"
	"business/user/utils/render"
)

func ValidPhoneEmail(w http.ResponseWriter, r *http.Request, phone string, email string, mediaType string) bool {
	if common.MediaTypeEmail == mediaType {
		if err := utils.ValidEmail(email); err != nil {
			render.JSON(w, r, http.StatusBadRequest, render.M{
				"code":    common.AccountInvalidEmail,
				"message": err.Error(),
			})
			return false
		}
	} else {
		if err := utils.ValidPhone(phone); err != nil {
			render.JSON(w, r, http.StatusBadRequest, render.M{
				"code":    common.AccountInvalidPhone,
				"message": err.Error(),
			})
			return false
		}
	}
	return true
}
