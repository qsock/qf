package encrypt

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(b []byte) []byte {
	h := md5.New()
	h.Write(b)
	return h.Sum(nil)
}

func Md5Hex(b []byte) string {
	return hex.EncodeToString(Md5(b))
}
