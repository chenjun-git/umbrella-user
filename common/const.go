package common

const (
	CurrentAccount              = "current_user"
	AuthHeader                  = "Authorization"
	AccessTokenKey              = "x-token"
	Verify2faTokenKey           = "2fa-token"
	ModuleName                  = "user"
	TwoFactorAuthVerifyTokenKey = "verify_2fa_token"
)

const (
	SignupPurpose                  = "signup"
	SigninPurpose                  = "signin"
	UpdatePhonePurpose             = "update_phone"
	UpdateEmailPurpose             = "update_email"
	ResetPasswordPurpose           = "reset_password"
	TwoFactorAuthSetupPurpose      = "setup_2fa_code"
	TwoFactorAuthVerifyCodePurpose = "verify_2fa_code"
)

const (
	MediaTypeEmail = "email"
	MediaTypePhone = "phone"
	MediaTypeTOTP  = "totp"
)

const (
	WebSource = "web"
	AppSource = "app"

	InternalSource = "internal"
)

var SourceRange []string
var PurposeRange []string
var TwoFactorAuthTypes []string

func init() {
	SourceRange = []string{WebSource, AppSource}
	PurposeRange = []string{SignupPurpose, SigninPurpose, UpdatePhonePurpose, UpdateEmailPurpose, ResetPasswordPurpose, TwoFactorAuthSetupPurpose, TwoFactorAuthVerifyCodePurpose}
	TwoFactorAuthTypes = []string{MediaTypePhone, MediaTypeEmail, MediaTypeTOTP}
}
