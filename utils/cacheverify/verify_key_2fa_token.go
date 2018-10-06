package cacheverify

import (
	"time"
)

type StringVerify struct {
	VerifyKey
	Token string `json:"token"`
}

func (r *VerifyKey) CreateStringVerify(s string, ttl time.Duration) (*StringVerify, error) {
	r.TTL = int64(ttl.Seconds())
	t := &StringVerify{
		VerifyKey: *r,
		Token:     s,
	}

	if err := t.VerifyKey.Save(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (r *VerifyKey) GetStringVerify() (*StringVerify, error) {
	c := StringVerify{}
	if ok, err := r.UnmarshalKey(&c); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	c.ctx = r.ctx

	return &c, nil
}
