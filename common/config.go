package common

import (
	"time"
	
	"github.com/BurntSushi/toml"
)

// Configs 全局配置信息
type Configs struct {
	Listen      string
	ReleaseMode bool
	Sms         *smsConfig
	Monitor     *monitorConfig
	Verify      *VerifyConfig
	Captcha     *CaptchaConfig
	Login       *LoginConfig
	Mysql       *MysqlConfig
	Redis       *RedisConfig
	HTTP        *httpConfig
	Token       *tokenConfig
}

type smsConfig struct {
	Addr string
}

type monitorConfig struct {
	Namespace string
	Subsystem string
}

type VerifyConfig struct {
	DefaultLength int      // 短信验证码长度，默认为6位
	MaxSendTimes  int      // 周期内最大发送次数，默认为5次
	MaxCheckTimes int      // 周期内最大检验次数，默认为5次
	TTL           Duration // 检测周期，默认为10分钟
}

type CaptchaConfig struct {
	DefaultWidth  int // image width
	DefaultHeight int // image heigth
	DefaultLength int // value length
	TTL           Duration
}

type LoginConfig struct {
	TTL                        Duration
	MaxRequestTimes            int // 周期内最大错误登录次数
	MaxCaptchaTImes            int // 周期内N次后需要验证码
	CleanInvalidTokenThreshold int // 清理无效token的阈值
	MaxNumberLogin             int // 允许最多登入的web端
}

// MysqlConfig MySQL相关配置信息
type MysqlConfig struct {
	Dsn     string
	MaxIdle int
	MaxOpen int
}

// RedisConfig Redis相关配置信息
type RedisConfig struct {
	Addr     string
	Password string
	PoolSize int
	DB       int
	Timeout  Duration
}

type httpConfig struct {
	Listen string
}

type tokenConfig struct {
	TokenLibVersion 	  int
	AccessTokenTTL        Duration
	RefreshTokenTTL       Duration
	AccessTokenExpireSoon Duration // 过期前多长时间算expire soon

	PrivateKeyPath string
	PublicKeyPath  string
}

// Config 全局配置信息
var Config *Configs

// InitConfig 加载配置
func InitConfig(path string) {
	config, err := loadConfig(path)
	if err != nil {
		panic(err)
	}
	Config = config
}

func loadConfig(path string) (*Configs, error) {
	config := new(Configs)
	if _, err := toml.DecodeFile(path, config); err != nil {
		return nil, err
	}
	return config, nil
}

// Duration 配置中使用的时长
type Duration struct {
	time.Duration
}

// UnmarshalText 将字符串形式的时长信息转换为Duration类型
func (d *Duration) UnmarshalText(text []byte) (err error) {
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

// D 从Duration struct中取出time.Duration类型的值
func (d *Duration) D() time.Duration {
	return d.Duration
}
