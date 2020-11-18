package ipaddr

import (
	"net"
	"strings"
)

// 获取内网ip
func GetLocalIp() string {
	addrs, err := GetLocalIps()
	if err == nil && len(addrs) > 0 {
		return addrs[0]
	}
	return "127.0.0.1"
}

// 获得出口ip
func GetExportIp() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()
	addr := conn.LocalAddr().String()
	idx := strings.LastIndex(addr, ":")
	return addr[:idx], nil
}

// 获取内网列表
func GetLocalIps() (ips []string, err error) {
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
		if !IsPublic(ipStr) {
			ips = append(ips, ipStr)
		}
	}
	return ips, nil
}

// 是否是公网ip
func IsPublic(ip string) bool {
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

func GetRemoteIp(host string) ([]string, error) {
	return net.LookupHost(host)
}
