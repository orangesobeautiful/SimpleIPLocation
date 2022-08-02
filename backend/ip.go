package main

import (
	"net/http"
	"net/netip"
	"strings"
)

// ParseReqRemoteIP 從 http.Request 解析客戶端請求的 IP 地址
func ParseReqRemoteIP(r *http.Request) netip.Addr {
	var realIP netip.Addr
	var parseErr error
	for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		addresses := r.Header.Values(h)

		// march from right to left until we get a public address
		// that will be the address right before our proxy.
		for i := len(addresses) - 1; i >= 0; i-- {
			ip := strings.TrimSpace(addresses[i])
			// header can contain spaces too, strip those out.
			realIP, parseErr = netip.ParseAddr(ip)
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
