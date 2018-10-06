package v1_0

import (
	"business/user/common"
)

func getContactByPhoneEmail(mediaType, phone, email string) string {
	if common.MediaTypeEmail == mediaType {
		return email
	} else {
		return phone
	}
}
