package util

import (
	"encoding/hex"

	uuid2 "github.com/google/uuid"
)

func GenUUIDWithOutDash() string {
	uuid := uuid2.New()
	var buf [32]byte
	dst := buf[:]
	hex.Encode(dst, uuid[:])
	return string(buf[:])
}
