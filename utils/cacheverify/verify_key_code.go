package cacheverify

import (
	"time"
)

type CodeVerify struct {
	VerifyKey
	SendTimes  int    `json:"send_times"`
	CheckTimes int    `json:"check_times"`
	VerifyCode string `json:"verify_code"`
}

func (r *CodeVerify) Verify(verifyCode string) (bool, error) {
	if r.VerifyCode == verifyCode {
		return true, nil
	}

	r.CheckTimes++
	return false, r.Save()
}

func (r *CodeVerify) Save() error {
	return r.VerifyKey.Save(r)
}

func (r *VerifyKey) GetCodeVerify() (*CodeVerify, error) {
	c := CodeVerify{}
	if ok, err := r.UnmarshalKey(&c); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	c.ctx = r.ctx

	return &c, nil
}

func (r *VerifyKey) CreateCodeVerify(verifyCode string, ttl time.Duration) (*CodeVerify, error) {
	r.TTL = int64(ttl.Seconds())
	t := &CodeVerify{
		VerifyKey:  *r,
		SendTimes:  0,
		CheckTimes: 0,
		VerifyCode: verifyCode,
	}
	if err := t.Save(); err != nil {
		return nil, err
	}
	return t, nil
}
