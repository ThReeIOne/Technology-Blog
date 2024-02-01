package lib

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"github.com/forgoer/openssl"
	"math/rand"
)

func MD5(text string) string {
	m := md5.New()
	m.Write([]byte(text))
	return hex.EncodeToString(m.Sum(nil))
}

func Encode(src, key string) string {
	text, err := openssl.Des3ECBEncrypt([]byte(src), []byte(key), openssl.PKCS7_PADDING)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(text)
}

func Decode(src, key string) string {
	dst, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		panic(err)
	}
	text, err := openssl.Des3ECBDecrypt(dst, []byte(key), openssl.PKCS7_PADDING)
	if err != nil {
		panic(err)
	}
	return string(text)
}

func Token(length int) string {
	source := []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	target := make([]rune, length)
	for i := range target {
		target[i] = source[rand.Intn(len(source))]
	}
	return string(target)
}

func GenerateRandomKey() string {
	key := make([]byte, 24)
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(key)[:24]
}

func GenerateSalt() (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}
