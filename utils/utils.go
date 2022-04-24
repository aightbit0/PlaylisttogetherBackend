package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func GetMD5Hash(text string) string {

	if text == "" {
		return text
	}

	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}
