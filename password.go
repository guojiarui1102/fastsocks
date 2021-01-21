package fastsocks

import (
	"encoding/base64"
	"errors"
	"math/rand"
	"strings"
	"time"
)

const passwordLength = 256

type password [passwordLength]byte

func init() {
	rand.Seed(time.Now().Unix())
}

func (pw *password) toString() string {
	return base64.StdEncoding.EncodeToString(pw[:])
}

// ParsePassword 解析获取密码
func ParsePassword(pwString string) ([]byte, error) {
	bs, err := base64.StdEncoding.DecodeString(strings.TrimSpace(pwString))
	if err != nil || len(bs) != passwordLength {
		return nil, errors.New("不符合规定的密码")
	}
	var pw []byte
	copy(pw, bs)
	return pw, nil
}

// RandPassword 产生一个随机密码
func RandPassword() string {
	intArr := rand.Perm(passwordLength)
	pw := &password{}
	for i, v := range intArr {
		pw[i] = byte(v)
		if i == v {
			return RandPassword()
		}
	}
	return pw.toString()
}
