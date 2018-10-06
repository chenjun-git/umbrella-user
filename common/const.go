package common

const (
	CurrentAccount              = "current_user"
	AuthHeader                  = "Authorization"
	AccessTokenKey              = "x-token"
	ModuleName                  = "user"
)

const (
	SignupPurpose                  = "signup"
	SigninPurpose                  = "signin"
	UpdatePhonePurpose             = "update_phone"
	UpdateEmailPurpose             = "update_email"
	ResetPasswordPurpose           = "reset_password"
)

const (
	MediaTypeEmail = "email"
	MediaTypePhone = "phone"
	MediaTypeTOTP  = "totp"
)

const (
	WebSource = "web"
	AppSource = "app"
)

var SourceRange []string
var PurposeRange []string
var TwoFactorAuthTypes []string

func init() {
	SourceRange = []string{WebSource, AppSource}
	PurposeRange = []string{SignupPurpose, SigninPurpose, UpdatePhonePurpose, UpdateEmailPurpose, ResetPasswordPurpose}
	TwoFactorAuthTypes = []string{MediaTypePhone, MediaTypeEmail, MediaTypeTOTP}
}
