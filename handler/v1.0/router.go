package v1_0

import (
	"fmt"

	"github.com/go-chi/chi"

	"github.com/chenjun-git/umbrella-common/monitor"
)

func RegisterRouter(r chi.Router) {
	registerRouter("v1.0", r)
}

func registerRouter(version string, r chi.Router) {
	var prefix string
	// normal API
	if version == "v1" {
		prefix = fmt.Sprintf("/user/v1/")
	} else {
		prefix = fmt.Sprintf("/api/%s/user", version)
	}
	r.Route(prefix, func(r chi.Router) {
		r.Post("/sign_up", monitor.HttpHandlerWrapper("SignUpHandler", SignUpHandler))
		r.Post("/sign_in", monitor.HttpHandlerWrapper("SignInHandler", SignInHandler))
		r.Post("/sign_out", monitor.HttpHandlerWrapper("SignOutHandler", SignOutHandler))
		r.Post("/reset_password", monitor.HttpHandlerWrapper("ResetPasswordHandler", ResetPasswordHandler))

		r.Post("/send_verify_code", monitor.HttpHandlerWrapper("SendVerifyCodeHandler", SendVerifyCodeHandler))
		r.Post("/check_verify_code", monitor.HttpHandlerWrapper("CheckVerifyCodeHandler", CheckVerifyCodeHandler))

		r.Get("/get_captcha_image", monitor.HttpHandlerWrapper("GetCaptchaImageHandler", GetCaptchaImageHandler))
		r.Post("/get_captcha_token", monitor.HttpHandlerWrapper("CreateCaptchaHandler", CreateCaptchaHandler))
	})

	// internal API
	if version == "v1" {
		prefix = fmt.Sprintf("/user/v1/internal")
	} else {
		prefix = fmt.Sprintf("/internal/api/%s/user", version)
	}
	r.Route(prefix, func(r chi.Router) {
		r.Post("/sign_up", monitor.HttpHandlerWrapper("SignUpWithoutVerifyCodeHandler", SignUpWithoutVerifyCodeHandler))
		r.Post("/sign_in", monitor.HttpHandlerWrapper("SignInWithoutVerifyCodeHandler", SignInWithoutVerifyCodeHandler))
		//r.Get("/verify_token", monitor.HttpHandlerWrapper("VerifyTokenHandler", VerifyTokenHandler))
		//r.Put("/account_token", monitor.HttpHandlerWrapper("GetTokenHandler", GetTokenHandler))
		r.Get("/get_phone", monitor.HttpHandlerWrapper("GetPhoneHandler", GetPhoneHandler))
		r.Get("/get_user", monitor.HttpHandlerWrapper("GetUserHandler", GetUserHandler))
		r.Get("/get_users", monitor.HttpHandlerWrapper("GetUsersHandler", GetUsersHandler))
		r.Post("/change_password", monitor.HttpHandlerWrapper("ForceChangePasswordHandler", ForceChangePasswordHandler))
		r.Post("/update_email", monitor.HttpHandlerWrapper("ForceUpdateEmailHandler", ForceUpdateEmailHandler))
		r.Get("/check_existence", monitor.HttpHandlerWrapper("CheckExistencHandler", CheckExistencHandler))
		// r.Post("/org_setting", monitor.HttpHandlerWrapper("AddOrgSetting", AddOrgSetting))
		// r.Put("/org_setting/{setting_id}", monitor.HttpHandlerWrapper("UpdateOrgSetting", UpdateOrgSetting))
		// r.Get("/org_setting", monitor.HttpHandlerWrapper("GetOrgSeting", GetOrgSeting))
	})

	prefix = fmt.Sprintf("/3rd/api/%s/user", version)
	r.Route(prefix, func(r chi.Router) {

	})
}
