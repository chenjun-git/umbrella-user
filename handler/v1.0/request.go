package v1_0

type SendVerifyCodeRequest struct {
	Phone        string `json:"phone" binding:"omitempty"`
	Email        string `json:"email" binding:"omitempty"`
	MediaType    string `json:"media_type" binding:"omitempty"`
	Purpose      string `json:"purpose" binding:"required"`
	Source       string `json:"source" binding:"required"`
	CaptchaToken string `json:"captcha_token" binding:"omitempty"`
	CaptchaValue string `json:"captcha_value" binding:"omitempty"`
}

type CheckVerifyCodeRequest struct {
	Phone      string `json:"phone" binding:"omitempty"`
	Email      string `json:"email" binding:"omitempty"`
	MediaType  string `json:"media_type" binding:"omitempty"`
	Purpose    string `json:"purpose" binding:"required"`
	Source     string `json:"source" binding:"required"`
	VerifyCode string `json:"verify_code" binding:"required"`
}

type SignUpRequest struct {
	Phone       string `json:"phone" binding:"omitempty"`
	Email       string `json:"email" binding:"omitempty"`
	MediaType   string `json:"media_type" binding:"omitempty"`
	Source      string `json:"source" binding:"required"`
	RegisterURL string `json:"register_url" binding:"omitempty"`
	Password    string `json:"password" binding:"required,min=8"`
	VerifyCode  string `json:"verify_code" binding:"required"`
}

type SignInRequest struct {
	Phone        string `json:"phone" binding:"omitempty"`
	Email        string `json:"email" binding:"omitempty"`
	MediaType    string `json:"media_type" binding:"omitempty"`
	Source       string `json:"source" binding:"required"`
	Password     string `json:"password" binding:"required,min=8"`
	CaptchaToken string `json:"captcha_token" binding:"omitempty"`
	CaptchaValue string `json:"captcha_value" binding:"omitempty"`
}

type ResetPasswordRequest struct {
	Phone       string `json:"phone" binding:"omitempty"`
	Email       string `json:"email" binding:"omitempty"`
	MediaType   string `json:"media_type" binding:"omitempty"`
	Source      string `json:"source" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
	VerifyCode  string `json:"verify_code" binding:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=8"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type UpdateEmailRequest struct {
	Password   string `json:"password" binding:"required,min=8"`
	Email      string `json:"email" binding:"required"`
	Source     string `json:"source" binding:"required"`
	VerifyCode string `json:"verify_code" binding:"required"`
}
