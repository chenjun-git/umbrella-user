package v1_0

import (
	"fmt"

	"github.com/go-chi/chi"
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
	router.Route(prefix, func(r chi.Router) {
		r.Post("/sign_up", monitor.HttpHandlerWrapper("SignUpHandler", SignUpHandler))
		r.Post("/sign_in", monitor.HttpHandlerWrapper("SignInHandler", SignInHandler))
		r.Post("/sign_out", monitor.HttpHandlerWrapper("SignOutHandler", SignOutHandler))
		r.Post("/reset_password", monitor.HttpHandlerWrapper("ResetPasswordHandler", ResetPasswordHandler))
		
		r.Get("/get_captcha_image", monitor.HttpHandlerWrapper("GetCaptchaImageHandler", GetCaptchaImageHandler))
		r.Post("/get_captcha_token", monitor.HttpHandlerWrapper("CreateCaptchaHandler", CreateCaptchaHandler))
	})

	prefix = fmt.Sprintf("/3rd/api/%s/user", version)
	router.Route(prefix, func(r chi.Router) {
		
	})
}