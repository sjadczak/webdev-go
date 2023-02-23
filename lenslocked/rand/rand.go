package rand

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func Bytes(n int) ([]byte, error) {
	bs := make([]byte, n)
	num, err := rand.Read(bs)
	if err != nil {
		return nil, fmt.Errorf("bytes: %w\n", err)
	}
	if num < n {
		return nil, fmt.Errorf("bytes: didn't read enough random bytes.")
	}

	return bs, nil
}

// n is the number of bytes being used to generate random string
func String(n int) (string, error) {
	bs, err := Bytes(n)
	if err != nil {
		return "", fmt.Errorf("string: %w\n", err)
	}

	return base64.URLEncoding.EncodeToString(bs), nil
}
