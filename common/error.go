package common

import "strings"

// 错误码
const (
	OK                              = 0
	TokenInvalid                    = 101
	TokenDeny                       = 102
	TokenExpired                    = 103
	TokenExpiredSoon                = 104 // Deprecated: use TokenExpired instead
	TokenBanned                     = 105
	SignInvalid                     = 106
	ACLConnectServerFailed          = 201
	ACLCallFailed                   = 202
	ACLObjectLevelLimit             = 203
	ACLFieldLevelLimit              = 204
	ACLRecordLevelLimit             = 205
	ACLActionTypeUnsupported        = 206
	ACLRecordAccessLevelUnsupported = 207
	ACLFunctionPermissionLimit      = 208
	AccountTokenInvalid             = 301
	AccountBindFailed               = 102100
	AccountInvalidPurpose           = 102101
	AccountInvalidSource            = 102102
	AccountAccountNotExist          = 102103
	AccountVerifyCodeNotMatch       = 102104
	AccountCaptchaNotMatch          = 102105
	AccountOrPasswordError          = 102106
	AccountNeedCaptcha              = 102107
	AccountAccountAlreadyExist      = 102108
	AccountInvalidPhone             = 102109
	AccountPasswordSameWithOld      = 102110
	AccountPasswordLevelIllegal     = 102111
	AccountInvalidEmail             = 102112
	AccountInvalidAccountId         = 102113
	AccountWasFreezed               = 102114
	AccountInvalidOAuthType         = 102115
	AccountInvalidParam             = 102116
	AccountPasswordExpired          = 102117
	AccountTwoFactorAuthMissing     = 102118
	AccountTwoFactorAuthIsEmpty     = 102119
	AccountJsonMarshalError         = 102120
	AccountRequestLimit             = 102201
	AccountInvalidToken             = 102202
	AccountTooMuchBoardingDevice    = 102203
	AccountInvalidTwoFactorAuthType = 102204
	AccountBindNoDevice             = 102205
	AccountBindNoApp                = 102206
	AccountGenerateIdFailed         = 102301
	AccountDBError                  = 102302
	AccountRedisError               = 102303
	AccountMySQLError               = 102304
	AccountInternalError            = 102401
	AccountAPIResponseIllegal       = 102402
	AccountBindNoAuthinfo           = 102403
	AccountServiceToken             = 102404
	AccountBindEmptyPersonalTokenName = 102405
	SmsServiceErr                   = 102406
)

// ErrorMap 错误码与错误信息map
var ErrorMap = map[int]map[string]string{
	0: map[string]string{
		"EN-US": "OK",
		"ZH-CN": "操作成功",
	},
	101: map[string]string{
		"EN-US": "Token Invalid",
		"ZH-CN": "",
	},
	102: map[string]string{
		"EN-US": "Token Deny",
		"ZH-CN": "",
	},
	103: map[string]string{
		"EN-US": "Token Expired",
		"ZH-CN": "",
	},
	104: map[string]string{
		"EN-US": "Token Expired Soon",
		"ZH-CN": "",
	},
	105: map[string]string{
		"EN-US": "Token Banned",
		"ZH-CN": "",
	},
	106: map[string]string{
		"EN-US": "Sign Invalid",
		"ZH-CN": "",
	},
	201: map[string]string{
		"EN-US": "ACL Connect Server Failed",
		"ZH-CN": "",
	},
	202: map[string]string{
		"EN-US": "ACL Call Failed",
		"ZH-CN": "",
	},
	203: map[string]string{
		"EN-US": "ACL Object Level Limit",
		"ZH-CN": "",
	},
	204: map[string]string{
		"EN-US": "ACL Field Level Limit",
		"ZH-CN": "",
	},
	205: map[string]string{
		"EN-US": "ACL Record Level Limit",
		"ZH-CN": "",
	},
	206: map[string]string{
		"EN-US": "ACL Action Type Unsupported",
		"ZH-CN": "",
	},
	207: map[string]string{
		"EN-US": "ACL Record Access Level Unsupported",
		"ZH-CN": "",
	},
	208: map[string]string{
		"EN-US": "ACL Function Permission Limit",
		"ZH-CN": "",
	},
	301: map[string]string{
		"EN-US": "Account Token Invalid",
		"ZH-CN": "",
	},
	102100: map[string]string{
		"EN-US": "Account Bind Failed",
		"ZH-CN": "参数解析错误",
	},
	102101: map[string]string{
		"EN-US": "Account Invalid Purpose",
		"ZH-CN": "无效目的操作",
	},
	102102: map[string]string{
		"EN-US": "Account Invalid Source",
		"ZH-CN": "无效来源",
	},
	102103: map[string]string{
		"EN-US": "Account Account Not Exist",
		"ZH-CN": "账户不存在",
	},
	102104: map[string]string{
		"EN-US": "Account VerifyCode Not Match",
		"ZH-CN": "校验码不匹配",
	},
	102105: map[string]string{
		"EN-US": "Account Captcha Not Match",
		"ZH-CN": "图形验证码不匹配",
	},
	102106: map[string]string{
		"EN-US": "Account or password error.",
		"ZH-CN": "账号或密码错误",
	},
	102107: map[string]string{
		"EN-US": "Account Need Captcha",
		"ZH-CN": "需要图形验证码",
	},
	102108: map[string]string{
		"EN-US": "Account Account Already Exist",
		"ZH-CN": "账户已存在",
	},
	102109: map[string]string{
		"EN-US": "Account Invalid Phone",
		"ZH-CN": "电话号码无效",
	},
	102110: map[string]string{
		"EN-US": "Account Password Same With Old",
		"ZH-CN": "与旧密码相同",
	},
	102111: map[string]string{
		"EN-US": "Account Password Level Illegal",
		"ZH-CN": "密码等级非法",
	},
	102112: map[string]string{
		"EN-US": "Account Invalid Email",
		"ZH-CN": "无效邮箱",
	},
	102113: map[string]string{
		"EN-US": "Account Invalid Account Id",
		"ZH-CN": "无效accountId",
	},
	102114: map[string]string{
		"EN-US": "Account Was Freezed",
		"ZH-CN": "账号被冻结",
	},
	102115: map[string]string{
		"EN-US": "Account Invalid OAuth Type",
		"ZH-CN": "无效的OAuth类型",
	},
	102116: map[string]string{
		"EN-US": "Account Invalid Param",
		"ZH-CN": "无效的参数",
	},
	102117: map[string]string{
		"EN-US": "Password Expired",
		"ZH-CN": "密码已过期",
	},
	102118: map[string]string{
		"EN-US": "Account Two Factor Auth Missing",
		"ZH-CN": "账号缺少二步验证",
	},
	102119: map[string]string{
		"EN-US": "Account Two Factor Auth IsEmpty",
		"ZH-CN": "账号没有设置二步验证方式",
	},
	102120: map[string]string{
		"EN-US": "Account Json Marshal",
		"ZH-CN": "账号Json序列化",
	},
	102201: map[string]string{
		"EN-US": "Account Request Limit",
		"ZH-CN": "达到最大发送限制，请稍后重试",
	},
	102202: map[string]string{
		"EN-US": "Account Invalid Token",
		"ZH-CN": "无效token",
	},
	102203: map[string]string{
		"EN-US": "Account Too Much Boarding Device",
		"ZH-CN": "该账号已达到同时登录数量上限",
	},
	102204: map[string]string{
		"EN-US": "Account Invalid Two Factor Auth Type",
		"ZH-CN": "无效二次验证类型",
	},
	102205: map[string]string{
		"EN-US": "Account Bind No Device",
		"ZH-CN": "无效终端设备代码",
	},
	102206: map[string]string{
		"EN-US": "Account Bind No App",
		"ZH-CN": "无效应用代码",
	},
	102301: map[string]string{
		"EN-US": "Account Generate Id Failed",
		"ZH-CN": "生成ID失败",
	},
	102302: map[string]string{
		"EN-US": "Account DB Error",
		"ZH-CN": "数据库错误",
	},
	102303: map[string]string{
		"EN-US": "Account Redis Error",
		"ZH-CN": "Redis错误",
	},
	102304: map[string]string{
		"EN-US": "Account MySQL Error",
		"ZH-CN": "MySQL错误",
	},
	102401: map[string]string{
		"EN-US": "Account Internal Error",
		"ZH-CN": "内部错误",
	},
	102402: map[string]string{
		"EN-US": "Account API Response Illegal",
		"ZH-CN": "接口返回值不合法",
	},
	102403: map[string]string{
		"EN-US": "Account No Bind Auth Info",
		"ZH-CN": "未包涵AUTHINFO",
	},
	102404: map[string]string{
		"EN-US": "Account Service Token",
		"ZH-CN": "调用token服务失败",
	},
	102405: map[string]string{
		"EN-US": "Account Personal Token Name Empty",
		"ZH-CN": "personal token 名字为空",
	},
	102406: map[string]string{
		"EN-US": "Sms Service Err",
		"ZH-CN": "短信服务 错误",
	},
}

// GetMsg 错误码转各个语言的错误信息
func GetMsg(code int, languages []string) string {
	msgMap, ok := ErrorMap[code]
	if !ok {
		return "Unknown error"
	}
	for _, lang := range languages {
		if msg, ok := msgMap[strings.ToUpper(lang)]; ok {
			if msg != "" {
				return msg
			}
		}
	}
	return "Unknown error"
}
