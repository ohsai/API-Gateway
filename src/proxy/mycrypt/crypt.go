package mycrypt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func CreateMAC(message, key string) []byte {
	key_encoded := base64.StdEncoding.EncodeToString([]byte(key))
	mac := hmac.New(sha256.New, []byte(key_encoded))
	mac.Write([]byte(message))
	outMAC := mac.Sum(nil)
	return outMAC
}
func CheckMAC(message string, messageMAC []byte, key string) bool {
	expectedMAC := CreateMAC(message, key)
	return hmac.Equal([]byte(messageMAC), expectedMAC)
}
