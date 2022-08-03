package utils

import (
	"net"
	"net/http"
	"net/netip"
	"strings"
)

// NetIPAddr2netIP 將 netip.Addr 轉換為 net.IP
func NetIPAddr2netIP(addr netip.Addr) net.IP {
	if addr.Is4() {
		p := make(net.IP, net.IPv6len)
		p[10] = 0xff
		p[11] = 0xff
		b4Ary := addr.As4()
		p[12] = b4Ary[0]
		p[13] = b4Ary[1]
		p[14] = b4Ary[2]
		p[15] = b4Ary[3]
		return p
	}
	return net.IP(addr.AsSlice())
}

// ParseReqRemoteIP 從 http.Request 解析客戶端請求的 IP
func ParseReqRemoteIP(r *http.Request) netip.Addr {
	var realIP netip.Addr
	var parseErr error

	// 嘗試解析 X-Forwarded-For 和 X-Real-Ip，從中取得客戶端 IP
	for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		addrStrList := r.Header.Values(h)

		for i := len(addrStrList) - 1; i >= 0; i-- {
			ipStr := strings.TrimSpace(addrStrList[i])
			realIP, parseErr = netip.ParseAddr(ipStr)
			if parseErr == nil {
				if !realIP.IsLoopback() && !realIP.IsLinkLocalUnicast() &&
					!realIP.IsLinkLocalMulticast() && !realIP.IsPrivate() {
					return realIP
				}
			}
		}
	}

	if r.RemoteAddr != "" {
		var addrPort netip.AddrPort
		addrPort, _ = netip.ParseAddrPort(r.RemoteAddr)
		return addrPort.Addr()
	}

	return realIP
}
