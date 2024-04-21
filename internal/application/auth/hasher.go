package auth

import (
	"crypto/md5"
	"encoding/hex"
)

// Hasher generates hash for password.
type Hasher func(password string) (string, error)

// HasherMD5Hex generates a hash for a password using the MD5 algorithm and returns it as a hex-encoded string.
func HasherMD5Hex(password string) (string, error) {
	h := md5.New()
	if _, err := h.Write([]byte(password)); err != nil {
		return "", err
	}
	sum := h.Sum(nil)
	hashString := hex.EncodeToString(sum)
	return hashString, nil
}
