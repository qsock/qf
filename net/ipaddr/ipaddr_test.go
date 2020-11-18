package ipaddr

import (
	"testing"
)

func TestGetPublicAddr(t *testing.T) {
	t.Log(GetExportIp())
	t.Log(GetLocalIp())
	t.Log(GetRemoteIp("www.qq.com"))
}
