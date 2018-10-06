package password

import (
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

const (
	// default cost in golang.org/x/crypto/bcrypt is 10, we use a higher cost
	cost = 12
)

// 密码级别
const (
	LevelIllegal = iota
	LevelWeak
	LevelNormal
	LevelStrong
)

func Encrypt(passwd string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(passwd), cost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func Verify(passwd string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwd))
}

func isSpecialChar(r rune) bool {
	return unicode.IsGraphic(r) && !unicode.IsNumber(r) && !unicode.IsLetter(r)
}

func PasswordStrength(password string) int {
	length := len(password)
	hasNumber := strings.IndexFunc(password, unicode.IsNumber) >= 0
	hasLetter := strings.IndexFunc(password, unicode.IsLetter) >= 0
	hasSpecialChar := strings.IndexFunc(password, isSpecialChar) >= 0

	weakLevel := length >= 8 && length <= 12 && (hasNumber && hasLetter || hasNumber && hasSpecialChar || hasLetter && hasSpecialChar)
	normalLevel := length > 12 && (hasNumber && hasLetter || hasNumber && hasSpecialChar || hasLetter && hasSpecialChar) ||
		length >= 8 && length <= 12 && hasNumber && hasLetter && hasSpecialChar
	strongLevel := length > 12 && hasNumber && hasLetter && hasSpecialChar

	if strongLevel {
		return LevelStrong
	} else if normalLevel {
		return LevelNormal
	} else if weakLevel {
		return LevelWeak
	}

	return LevelIllegal
}

