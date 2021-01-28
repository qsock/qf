package qcast

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type T struct {
	A string `json:"a"`
	B string `json:"b"`
	C struct {
		D []string `json:"d"`
	} `json:"c"`
}

func Test_Decode(t *testing.T) {
	var src = map[string]interface{}{
		"a": "1",
		"b": 2,
		"c": map[string]interface{}{
			"d": []string{"1", "2", "3"},
		},
	}

	var p T

	err := Decode(src, &p)
	assert.Nil(t, err)
}
