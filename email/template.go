package email

import (
	"strings"

	"github.com/chenjun-git/umbrella-user/common"
)

// 邮件结构内容
type EmailTplInfo struct {
	Title          string
	HtmlVerifyUrl  string
	HtmlVerifyCode string
}

var emailTemplates = make(map[string]EmailTplInfo)

func init() {
	var VerifyUrlHelper = func(do, sth string) string {
		return `您好！<BR>您正在` + do + `{{orgName}}` + sth + `，点击下面链接完成操作：<BR><A href="%s">%s</A><BR>如果您没有` + do + sth + `操作，请与管理员联系！`
	}
	var VerifyCodeHelper = func(do, sth string) string {
		return `您好！<BR>您正在` + do + `{{orgName}}` + sth + `，本次验证码为：<BR>%s<BR>如果您没有` + do + sth + `操作，请与管理员联系！`
	}

	emailTemplates[common.SignupPurpose] = EmailTplInfo{
		Title:          "注册您的{{orgName}}账户",
		HtmlVerifyUrl:  `您好！<BR>感谢选择{{orgName}}，点击下面链接输入验证码:<BR><A href="%s">%s</A><BR>如果您没有注册{{orgName}}账号操作，请与管理员联系！`,
		HtmlVerifyCode: `您好！<BR>您正在注册{{orgName}}账户，本次验证码为：<BR>%s<BR>如果您没有注册{{orgName}}账号操作，请与管理员联系！`,
	}

	emailTemplates[common.SigninPurpose] = EmailTplInfo{
		Title:          "验证您的{{orgName}}账户",
		HtmlVerifyUrl:  `您好！<BR>感谢选择{{orgName}}，点击下面链接验证邮箱：<BR><A href="%s">%s</A><BR>如果您没有注册{{orgName}}操作，请忽略此邮件。`,
		HtmlVerifyCode: `您好！<BR>您正在进行账户验证操作，本次验证码为：<BR>%s<BR>如果您没有相关操作，请忽略此邮件。`,
	}

	emailTemplates[common.UpdatePhonePurpose] = EmailTplInfo{
		Title:          "更新{{orgName}}账户手机号",
		HtmlVerifyUrl:  VerifyUrlHelper("绑定", "账号到新的手机号"),
		HtmlVerifyCode: VerifyCodeHelper("绑定", "账号到新的手机号"),
	}

	emailTemplates[common.UpdateEmailPurpose] = EmailTplInfo{
		Title:          "更新您的{{orgName}}账户邮箱",
		HtmlVerifyUrl:  VerifyUrlHelper("更新您的", "账户邮箱"),
		HtmlVerifyCode: VerifyCodeHelper("更新您的", "账户邮箱"),
	}

	emailTemplates[common.ResetPasswordPurpose] = EmailTplInfo{
		Title:          "找回您的{{orgName}}账户密码",
		HtmlVerifyUrl:  VerifyUrlHelper("找回", "密码"),
		HtmlVerifyCode: VerifyCodeHelper("找回", "密码"),
	}
	emailTemplates[common.TwoFactorAuthSetupPurpose] = EmailTplInfo{
		Title:          "设置您的{{orgName}}账户邮箱二步验证",
		HtmlVerifyUrl:  VerifyUrlHelper("设置", "邮箱二步验证方式"),
		HtmlVerifyCode: VerifyCodeHelper("设置", "邮箱二步验证方式"),
	}
	emailTemplates[common.TwoFactorAuthVerifyCodePurpose] = EmailTplInfo{
		Title:          "验证您的{{orgName}}账户邮箱二步验证",
		HtmlVerifyUrl:  VerifyUrlHelper("尝试校验", "邮箱二步验证"),
		HtmlVerifyCode: VerifyCodeHelper("尝试校验", "邮箱二步验证"),
	}
}

func GetTplInfo(purpose, orgName string) *EmailTplInfo {
	tplInfo, ok := emailTemplates[purpose]
	if !ok {
		return nil
	}
	tplInfo.Title = strings.Replace(tplInfo.Title, "{{orgName}}", orgName, -1)
	tplInfo.HtmlVerifyUrl = strings.Replace(tplInfo.HtmlVerifyUrl, "{{orgName}}", orgName, -1)
	tplInfo.HtmlVerifyCode = strings.Replace(tplInfo.HtmlVerifyCode, "{{orgName}}", orgName, -1)

	return &tplInfo
}

