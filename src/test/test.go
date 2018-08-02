package test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func test() {
	data := "Hello World!"
	sEnc := base64.StdEncoding.EncodeToString([]byte(data))
	fmt.Println(sEnc)
	sDec, _ := base64.StdEncoding.DecodeString(sEnc)
	fmt.Println(string(sDec))
	key := "private key"
	key_encoded := base64.StdEncoding.EncodeToString([]byte(key))
	mac := hmac.New(sha256.New, []byte(key_encoded))
	mac.Write([]byte(data))
	fmt.Println("hmac Encrypted", mac.Sum(nil))
	data2 := "Hello World283284542!"
	mac2 := hmac.New(sha256.New, []byte(key_encoded))
	mac2.Write([]byte(data2))
	fmt.Println("hmac2 Encrypted", mac2.Sum(nil))
	fmt.Println("Equal? : ", hmac.Equal(mac.Sum(nil), mac2.Sum(nil)))
	fmt.Println(string([]byte(data)))

}
