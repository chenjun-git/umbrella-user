package cacheverify

import (
	"time"
	
	"github.com/chenjun-git/umbrella-user/common"
)

type LoginVerify struct {
	VerifyKey
	RequestTimes int `json:"request_times"` // 生存周期内请求次数记录
}

func (r *LoginVerify) Save() error {
	return r.VerifyKey.Save(r)
}

func (r *LoginVerify) GetLeftTimes() int {
	if r.RequestTimes < common.Config.Login.MaxRequestTimes {
		return common.Config.Login.MaxRequestTimes - r.RequestTimes
	}
	return 0
}

// 从redis中恢复 LoginVerify
func (r *VerifyKey) GetLoginVerify() (*LoginVerify, error) {
	c := LoginVerify{}
	if ok, err := r.UnmarshalKey(&c); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	c.ctx = r.ctx

	return &c, nil
}

// 创建 LoginVerify 并保存到redis
func (r *VerifyKey) CreateLoginVerify(ttl time.Duration) (*LoginVerify, error) {
	r.TTL = int64(ttl.Seconds())
	t := &LoginVerify{
		VerifyKey:    *r,
		RequestTimes: 0,
	}

	if err := t.Save(); err != nil {
		return nil, err
	}
	return t, nil
}