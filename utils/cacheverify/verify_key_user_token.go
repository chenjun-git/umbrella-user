package cacheverify

import (
	"time"
)
type UserTokenVerify struct {
	VerifyKey
	Token string `json:"token"`
}

func (r *VerifyKey) CreateUserTokenVerify(userToken string, ttl time.Duration) (*UserTokenVerify, error) {
	r.TTL = int64(ttl.Seconds())
	r.CreatedTimestamp = time.Now().Unix()
	t := &UserTokenVerify{
		VerifyKey: *r,
		Token:     userToken,
	}

	if err := t.VerifyKey.SaveHash(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (r *VerifyKey) GetUserTokenVerify() (*UserTokenVerify, error) {
	c := UserTokenVerify{}
	if ok, err := r.UnmarshalHash(&c); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	c.ctx = r.ctx

	return &c, nil
}