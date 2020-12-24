package ka

import (
	"net"
	"os"
)

var (
	hostname  string
	localAddr string
)

func init() {
	hostname, _ = os.Hostname()
	localAddr = getLocalIp()
}

// 获取内网ip
func getLocalIp() string {
	addrs, err := getLocalIps()
	if err == nil && len(addrs) > 0 {
		return addrs[0]
	}
	return "127.0.0.1"
}

// 获取内网列表
func getLocalIps() (ips []string, err error) {
	addrs, e := net.InterfaceAddrs()
	if e != nil {
		return ips, e
	}
	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}

		if ip == nil || ip.IsLoopback() {
			continue
		}
		ip = ip.To4()
		if ip == nil {
			// not an ipv4 address
			continue
		}
		ipStr := ip.String()
		if !isPublic(ipStr) {
			ips = append(ips, ipStr)
		}
	}
	return ips, nil
}

func isPublic(ip string) bool {
	ipAddr := net.ParseIP(ip)
	if ipAddr.IsLoopback() || ipAddr.IsLinkLocalMulticast() || ipAddr.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := ipAddr.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
}
