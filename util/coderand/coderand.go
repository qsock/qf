package coderand

import (
	"crypto/rand"
	"encoding/binary"
	"unsafe"
)

// copy from gf
var (
	buffers    = make(chan []byte, 8196)
	upperCase  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowerCase  = "abcdefghijklmnopqrstuvwxyz"
	numCase    = "0123456789"
	symbols    = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
	letterCase = upperCase + lowerCase
	asciiCase  = upperCase + lowerCase + numCase + symbols
)

func init() {
	go random()
}

func random() {
	for {
		buffer := make([]byte, 1024)
		if n, err := rand.Read(buffer); err != nil {
			panic(err)
		} else {
			for i := 0; i <= n-4; i += 4 {
				buffers <- buffer[i : i+4]
			}
		}

	}
}

func Uint32(n uint32) uint32 {
	if n == 0 {
		return n
	}
	buf := <-buffers
	return binary.BigEndian.Uint32(buf) % n
}

func Uint64(n uint64) uint64 {
	b := make([]byte, 8)
	copy(b[0:4], <-buffers)
	copy(b[4:], <-buffers)
	return binary.BigEndian.Uint64(b) % n
}

func b(n int) []byte {
	if n <= 0 {
		return nil
	}
	i := 0
	b := make([]byte, n)
	for {
		copy(b[i:], <-buffers)
		i += 4
		if i >= n {
			break
		}
	}
	return b
}

func str(n int, s string) string {
	if n <= 0 || len(s) == 0 {
		return ""
	}
	var (
		ret = make([]byte, n)
		bs  = b(n)
	)
	for i := range ret {
		idx := bs[i] % byte(len(s))
		ret[i] = s[idx]
	}
	return *(*string)(unsafe.Pointer(&ret))
}

func Num(n int) string {
	return str(n, numCase)
}

// 生成多少个ascii
func Ascii(n int) string {
	return str(n, asciiCase)
}

func Letter(n int) string {
	return str(n, letterCase)
}

func Upper(n int) string {
	return str(n, upperCase)
}

func Lower(n int) string {
	return str(n, lowerCase)
}

func UpperNum(n int) string {
	return str(n, upperCase+numCase)
}

func LowerNum(n int) string {
	return str(n, lowerCase+numCase)
}

func BetweenUnit32(min, max uint32) uint32 {
	if min > max {
		return min
	}
	return Uint32(max-min+1) + min
}

func BetweenUnit64(min, max uint64) uint64 {
	if min > max {
		return min
	}
	return Uint64(max-min+1) + min
}
