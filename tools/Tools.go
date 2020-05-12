package tools

import (
	"crypto/md5"
	"encoding/hex"
)

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func ShortString(s string, i int) string {
	runes := []rune( s )
	if len(runes) > i {
		return string(runes[:i])
	}
	return s
}
