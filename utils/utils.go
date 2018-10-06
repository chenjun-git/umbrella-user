package utils

import (
	"database/sql"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"

	"github.com/json-iterator/go"
)

// NullStringToString NullString类型转换为String类型
func NullStringToString(s sql.NullString) interface{} {
	if s.Valid {
		return s.String
	} else {
		return nil
	}
}

func StringToNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

const (
	Digits   = "0123456789"
	Alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Ascii    = Alphabet + Digits + "~!@#$%^&*()-_+={}[]\\|<,>.?/\"';:`"
)

func RandomDigits(length int) string {
	return randomString(length, []byte(Digits))
}

func randomString(length int, base []byte) string {
	bytes := make([]byte, length)
	maxIndex := len(base)
	for i := 0; i < length; i++ {
		index := rand.Intn(maxIndex)
		bytes[i] = byte(base[index])
	}

	return string(bytes)
}

func BindJSON(r *http.Request, obj interface{}) error {
	defer io.Copy(ioutil.Discard, r.Body)
	if err := jsoniter.NewDecoder(r.Body).Decode(obj); err != nil {
		return err
	}

	return nil//structValidator.ValidateStruct(obj)
}