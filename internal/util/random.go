package util

import (
	crand "crypto/rand"
	"encoding/base32"
	mrand "math/rand"
	"time"
)

func init() {
	mrand.Seed(time.Now().UnixNano())
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[mrand.Intn(len(letters))]
	}
	return string(b)
}

func RandomCryptoString(n int) (string, error) {
	s := make([]byte, n)
	_, err := crand.Read(s)
	if err != nil {
		return "", err
	}

	return base32.StdEncoding.EncodeToString(s), nil
}

func RandomByteSlice(n uint32) ([]byte, error) {
	s := make([]byte, n)
	_, err := crand.Read(s)
	if err != nil {
		return nil, err
	}

	return s, nil
}
