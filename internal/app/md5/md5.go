package md5

import (
	"crypto/md5"
	"encoding/hex"
)

func HashString(srcString string) string{
	hasher := md5.New()
	hasher.Write([]byte(srcString))
	return hex.EncodeToString(hasher.Sum(nil))
}