package encrypt

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha512"
)

func HmacSha1(key, val string) []byte {
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(val))
	return h.Sum(nil)
}

func HmacSha512(key, val string) []byte {
	h := hmac.New(sha512.New, []byte(key))
	h.Write([]byte(val))
	return h.Sum(nil)
}
