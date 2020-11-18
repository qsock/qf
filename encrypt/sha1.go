package encrypt

import (
	"crypto/sha1"
	"encoding/hex"
)

func Sha1(b []byte) []byte {
	h := sha1.New()
	h.Write(b)
	return h.Sum(nil)
}

func Sha1Hex(b []byte) string {
	return hex.EncodeToString(Sha1(b))
}
