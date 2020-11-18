package coderand

import "testing"

func TestUint64(t *testing.T) {
	for i := 0; i < 1000; i++ {
		t.Log(Uint64(7655563599020032))
	}
}

func TestAscii(t *testing.T) {
	for i := 0; i < 1000; i++ {
		t.Log(Ascii(30))
	}
}

func TestUpperNum(t *testing.T) {
	for i := 0; i < 1000; i++ {
		t.Log(UpperNum(30))
	}
}
