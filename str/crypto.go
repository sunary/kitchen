package str

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
)

// Hash Using sha256 algorithm
func Hash(password, salt []byte) string {
	mac := hmac.New(sha256.New, salt)
	mac.Write([]byte(password))
	return hex.EncodeToString(mac.Sum(nil))
}

// Md5 ...
func Md5(input string) string {
	h := md5.New()
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}
