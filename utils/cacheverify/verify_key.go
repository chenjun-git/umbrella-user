package cacheverify

import (
	"context"
	"fmt"
	"time"

	"github.com/fzzy/radix/redis"

	"github.com/chenjun-git/umbrella-common/json"

	"business/user/common"
	"business/user/db"
)

type VerifyKey struct {
	ID      string `json:"id"`
	Purpose string `json:"purpose"`
	Source  string `json:"source"`
	ctx     context.Context

	CreatedTimestamp int64 `json:"create_ts"`
	TTL              int64 `json:"ttl"` // 单位为秒
}

func NewVerifyKey(ctx context.Context, id, purpose, source string) *VerifyKey {
	return &VerifyKey{
		ID:               id,
		Purpose:          purpose,
		Source:           source,
		ctx:              ctx,
		CreatedTimestamp: time.Now().Unix(),
	}
}

func (r *VerifyKey) String() string {
	return fmt.Sprintf("%s:%s:%s:%s", common.ModuleName, r.ID, r.Purpose, r.Source)
}

func (r *VerifyKey) Context() context.Context {
	return r.ctx
}

func (r *VerifyKey) GetHashKey() (string, string) {
	return fmt.Sprintf("%s:%s", common.ModuleName, r.ID), r.Source
}

func (r *VerifyKey) Save(i interface{}) error {
	value, err := json.Marshal(i)
	if err != nil {
		return err
	}

	key := r.String()
	expire := r.CreatedTimestamp + r.TTL - time.Now().Unix()
	if expire <= 0 {
		return db.RedisDel(r.ctx, key)
	}

	return db.RedisSet(r.ctx, key, value, time.Duration(expire)*time.Second)
}

func (r *VerifyKey) SaveHash(i interface{}) error {
	value, err := json.Marshal(i)
	if err != nil {
		return err
	}

	key, source := r.GetHashKey()
	expire := r.CreatedTimestamp + r.TTL - time.Now().Unix()
	if expire <= 0 {
		return db.RedisDel(r.ctx, key)
	}

	if err = db.RedisHset(r.ctx, key, source, value); err != nil {
		return err
	}

	if err = db.RedisEX(r.ctx, key, time.Duration(expire)*time.Second); err != nil {
		db.RedisHdel(r.ctx, key, source)
		return err
	}

	return nil
}

func (r *VerifyKey) Clean() error {
	return db.RedisDel(r.ctx, r.String())
}

func (r *VerifyKey) CleanHash() error {
	key, source := r.GetHashKey()
	return db.RedisHdel(r.ctx, key, source)
}

func (r *VerifyKey) UnmarshalKey(data interface{}) (bool, error) {
	reply, err := db.RedisGet(r.ctx, r.String())
	if err != nil {
		return false, err
	}
	if reply.Type == redis.NilReply {
		return false, nil
	}
	bs, err := reply.Bytes()
	if err != nil {
		return false, err
	}

	if err = json.Unmarshal(bs, data); err != nil {
		return false, err
	}

	return true, nil
}

func (r *VerifyKey) UnmarshalHash(data interface{}) (bool, error) {
	key, source := r.GetHashKey()

	reply, err := db.RedisHget(r.ctx, key, source)
	if err != nil {
		return false, err
	}
	if reply.Type == redis.NilReply {
		return false, nil
	}
	bs, err := reply.Bytes()
	if err != nil {
		return false, err
	}
	if err = json.Unmarshal(bs, data); err != nil {
		return false, err
	}
	return true, nil

}

func (k *VerifyKey) GetTokenNum() (int, error) {
	key, _ := k.GetHashKey()

	reply, err := db.RedisHlen(k.ctx, key)
	if err != nil {
		return 0, err
	}
	lens, err := reply.Int()
	if err != nil {
		return 0, err
	}

	return lens, nil
}