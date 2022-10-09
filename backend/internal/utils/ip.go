package utils

import (
	"net"
	"net/http"
	"net/netip"
	"strings"
	"unsafe"
)

// normalPass (only for benchmark) 初始化 net.IP 並且把 netip.Addr.As4() 的值複製到上面
func normalIPPass(p net.IP, b4Ary [4]byte) {
	// 從 go/src/net/ip.go func IPv4 複製的做法
	p[10] = 0xff
	p[11] = 0xff
	p[12] = b4Ary[0]
	p[13] = b4Ary[1]
	p[14] = b4Ary[2]
	p[15] = b4Ary[3]
}

// unsafeIPPass 初始化 net.IP 並且把 netip.Addr.As4() 的值透過 unsafe.Pointer 複製過去，
// 較清楚的邏輯可參考 func normalIPPass
func unsafeIPPass(p net.IP, b4Ary [4]byte) {
	p8Ptr := (*uint64)(unsafe.Pointer(&p[8]))
	b4AryPtr := (*uint32)(unsafe.Pointer(&b4Ary[0]))
	*p8Ptr = 0xffff<<(8*2) + uint64(*b4AryPtr)<<(8*4)
}

// NetIPAddr2netIP 將 netip.Addr 轉換為 net.IP
func NetIPAddr2netIP(addr netip.Addr) net.IP {
	if addr.Is4() {
		b4Ary := addr.As4()
		p := make(net.IP, net.IPv6len)
		unsafeIPPass(p, b4Ary)
		return p
	} else if addr.IsValid() {
		b16Ary := addr.As16()
		return net.IP(b16Ary[:])
	}
	return nil
}

// ParseReqRemoteIP 從 http.Request 解析客戶端請求的 IP
func ParseReqRemoteIP(r *http.Request) netip.Addr {
	var realIP netip.Addr
	var parseErr error

	// 嘗試解析 X-Forwarded-For 和 X-Real-Ip，從中取得客戶端 IP
	var parseHeader = []string{"X-Forwarded-For", "X-Real-Ip"}
	for hIdx := range parseHeader {
		addrStrList := r.Header.Values(parseHeader[hIdx])

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
