package qhttp

import (
	"testing"
)

func TestGet(t *testing.T) {
	cfg := new(Config)
	cfg.Host = "www.baidu.com"
	cfg.Name = "baidu"
	cfg.Schema = "https"
	err := Add(cfg)
	if err != nil {
		t.Error(err)
		return
	}
	c, err := GetClient("baidu")
	if err != nil {
		t.Error(err)
		return
	}
	code, resp, err := c.Get("/", nil)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%d||%s", code, string(resp))
}
