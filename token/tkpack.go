package token

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"io/ioutil"

	"github.com/satori/go.uuid"

	"github.com/chenjun-git/umbrella-common/token"
)

const (
	issueTimeEndIndex = 4
	ttlEndIndex       = 4 + 2
	userIDEndIndex    = 4 + 2 + 32
)

func InitKeyPem(publicKeyFile, privateKeyFile string) {
	publickKeyPem, err := ioutil.ReadFile(publicKeyFile)
	if err != nil {
		panic(err)
	}

	privateKeyPem, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		panic(err)
	}

	token.InitPublicKeys(publickKeyPem)
	token.InitPrivateKey(privateKeyPem)
}

func encodeRefreshToken(issueTime uint32, TTL uint16, userID string) (string, error) {
	issueTimeBytes := make([]byte, issueTimeEndIndex)
	ttlBytes := make([]byte, ttlEndIndex-issueTimeEndIndex)
	var datas []byte

	binary.BigEndian.PutUint32(issueTimeBytes, issueTime)
	binary.BigEndian.PutUint16(ttlBytes, TTL)

	for _, bs := range [][]byte{issueTimeBytes, ttlBytes, []byte(userID)} {
		datas = append(datas, bs...)
	}

	sign, err := uuid.NewV1()
	if err != nil {
		return "", err
	}
	datas = append(datas, sign.Bytes()...)

	return base64.URLEncoding.EncodeToString(datas), nil
}

// DecodeRefreshToken 解码refreshToken
func DecodeRefreshToken(token string) (issueTime uint32, TTL uint16, userID string, sign uuid.UUID, err error) {
	if token == "" {
		err = errors.New("DecodeRefreshToken: token is empty")
		return
	}

	data, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		err = errors.New("DecodeRefreshToken: invalid base64 token")
		return
	}

	sign, err = uuid.FromBytes(data[userIDEndIndex:])
	if err != nil {
		err = errors.New("DecodeRefreshToken: invalid uuid sign")
		return
	}

	issueTime = binary.BigEndian.Uint32(data[:issueTimeEndIndex])
	TTL = binary.BigEndian.Uint16(data[issueTimeEndIndex:ttlEndIndex])
	userID = string(data[ttlEndIndex:userIDEndIndex])
	return
}