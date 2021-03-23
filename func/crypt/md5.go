package crypt

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(str string) string {
	md5_16 := md5.Sum([]byte(str))
	return hex.EncodeToString(md5_16[:])
}
