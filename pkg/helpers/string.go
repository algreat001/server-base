package helpers

import (
	"crypto/rand"
	"strings"
)

func RandStr(strSize int, randType string) string {

	var dictionary string

	if randType == "alphanum" {
		dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "base32" {
		dictionary = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
	}

	if randType == "alpha" {
		dictionary = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "number" {
		dictionary = "0123456789"
	}

	var bytes = make([]byte, strSize)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}

func GetSafeString(str string) string {
	safeFilter := strings.Replace(str, "'", "", -1)
	safeFilter = strings.Replace(safeFilter, ";", "", -1)

	return safeFilter
}
