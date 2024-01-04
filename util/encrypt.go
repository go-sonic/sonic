package util

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(s string) string {
	d := []byte(s)
	m := md5.New()
	m.Write(d)
	return hex.EncodeToString(m.Sum(nil))
}
