package ipaddr

import (
	"os"
)

// 获取机器名
func GetHostname() string {
	h, _ := os.Hostname()
	return h
}
