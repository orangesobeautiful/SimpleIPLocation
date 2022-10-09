package utils

import (
	"net"
	"net/netip"
	"reflect"
	"testing"
)

// parseIPBench/Tests is edited from go/src/net/ip_test.go
var parseIPBench = []struct {
	in  string
	out net.IP
}{
	{"127.0.1.2", net.IPv4(127, 0, 1, 2)},
	{"127.0.0.1", net.IPv4(127, 0, 0, 1)},
	{"127.0.0.255", net.IPv4(127, 0, 0, 255)},
	{"8.8.8.8", net.IPv4(8, 8, 8, 8)},
	{"1.1.1.1", net.IPv4(1, 1, 1, 1)},
	// {"::ffff:127.1.2.3", net.IPv4(127, 1, 2, 3)},
	// {"::ffff:7f01:0203", net.IPv4(127, 1, 2, 3)},
	// {"0:0:0:0:0000:ffff:127.1.2.3", net.IPv4(127, 1, 2, 3)},
	// {"0:0:0:0:000000:ffff:127.1.2.3", net.IPv4(127, 1, 2, 3)},
	// {"0:0:0:0::ffff:127.1.2.3", net.IPv4(127, 1, 2, 3)},

	// {"2001:4860:0:2001::68", net.IP{0x20, 0x01, 0x48, 0x60, 0, 0, 0x20, 0x01, 0, 0, 0, 0, 0, 0, 0x00, 0x68}},
	// {"2001:4860:0000:2001:0000:0000:0000:0068", net.IP{0x20, 0x01, 0x48, 0x60, 0, 0, 0x20, 0x01, 0, 0, 0, 0, 0, 0, 0x00, 0x68}},
}

var parseIPTestsExt = []struct {
	in  string
	out net.IP
}{
	{"-0.0.0.0", nil},
	{"0.-1.0.0", nil},
	{"0.0.-2.0", nil},
	{"0.0.0.-3", nil},
	{"127.0.0.256", nil},
	{"abc", nil},
	{"123:", nil},
	{"a1:a2:a3:a4::b1:b2:b3:b4", nil}, // Issue 6628
	{"127.001.002.003", nil},
	{"::ffff:127.001.002.003", nil},
	{"123.000.000.000", nil},
	{"1.2..4", nil},
	{"0123.0.0.1", nil},
}
var parseIPTests = append(parseIPBench, parseIPTestsExt...)

func testFunc(t *testing.T, f func(netip.Addr) net.IP) {
	for _, tt := range parseIPTests {
		netIPAddr, parseErr := netip.ParseAddr(tt.in)
		if parseErr != nil {
			netIPAddr = netip.Addr{}
		}
		convertedIP := f(netIPAddr)
		if !reflect.DeepEqual(convertedIP, tt.out) {
			t.Errorf("IP(%q) = %v, want %v", tt.in, convertedIP, tt.out)
		}
	}
}

func TestNetIPAddr2netIP(t *testing.T) {
	testFunc(t, NetIPAddr2netIP)
}

func benchFunc(b *testing.B, f func(p net.IP, b4Ary [4]byte)) {
	var b4AryList [][4]byte
	for idx := range parseIPBench {
		netIPAddr, parseErr := netip.ParseAddr(parseIPBench[idx].in)
		if parseErr == nil {
			b4AryList = append(b4AryList, netIPAddr.As4())
		}
	}

	p := make(net.IP, net.IPv6len)

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for idx := range b4AryList {
				f(p, b4AryList[idx])
			}
		}
	})
}

func BenchmarkUnsafePtrIPPass(b *testing.B) {
	benchFunc(b, unsafeIPPass)
}
func BenchmarkNormalIPPass(b *testing.B) {
	benchFunc(b, normalIPPass)
}
