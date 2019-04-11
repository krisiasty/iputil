package iputil

import (
	"net"
	"testing"
)

var convertTests = []struct {
	in  net.IP
	out uint32
}{
	{net.IPv4(0, 0, 0, 0), 0},
	{net.IPv4(0, 0, 0, 1), 1},
	{net.IPv4(0, 0, 1, 0), 256},
	{net.IPv4(0, 1, 0, 0), 65536},
	{net.IPv4(1, 0, 0, 0), 16777216},
	{net.IPv4(0, 0, 1, 1), 257},
	{net.IPv4(1, 1, 0, 0), 16842752},
	{net.IPv4(0, 1, 1, 0), 65792},
	{net.IPv4(1, 0, 0, 1), 16777217},
	{net.IPv4(1, 1, 1, 1), 16843009},
	{net.IPv4(0, 0, 0, 255), 255},
	{net.IPv4(0, 0, 255, 0), 65280},
	{net.IPv4(0, 255, 0, 0), 16711680},
	{net.IPv4(255, 0, 0, 0), 4278190080},
	{net.IPv4(0, 0, 255, 255), 65535},
	{net.IPv4(255, 255, 0, 0), 4294901760},
	{net.IPv4(0, 255, 255, 0), 16776960},
	{net.IPv4(255, 0, 0, 255), 4278190335},
	{net.IPv4(255, 255, 255, 255), 4294967295},
	{net.IPv4(127, 0, 0, 1), 2130706433},
}

var numAdrTests = []struct {
	in  string
	out uint64
}{
	{"192.168.0.0/32", 1},
	{"192.168.0.0/31", 2},
	{"192.168.0.0/30", 4},
	{"192.168.0.0/29", 8},
	{"192.168.0.0/28", 16},
	{"192.168.0.0/27", 32},
	{"192.168.0.0/26", 64},
	{"192.168.0.0/25", 128},
	{"192.168.0.0/24", 256},

	{"192.168.0.0/23", 512},
	{"192.168.0.0/22", 1024},
	{"192.168.0.0/21", 2048},
	{"192.168.0.0/20", 4096},
	{"192.168.0.0/19", 8192},
	{"192.168.0.0/18", 16384},
	{"192.168.0.0/17", 32768},
	{"192.168.0.0/16", 65536},

	{"192.168.0.0/15", 131072},
	{"192.168.0.0/14", 262144},
	{"192.168.0.0/13", 524288},
	{"192.168.0.0/12", 1048576},
	{"192.168.0.0/11", 2097152},
	{"192.168.0.0/10", 4194304},
	{"192.168.0.0/9", 8388608},
	{"192.168.0.0/8", 16777216},

	{"192.168.0.0/7", 33554432},
	{"192.168.0.0/6", 67108864},
	{"192.168.0.0/5", 134217728},
	{"192.168.0.0/4", 268435456},
	{"192.168.0.0/3", 536870912},
	{"192.168.0.0/2", 1073741824},
	{"192.168.0.0/1", 2147483648},
	{"192.168.0.0/0", 4294967296},

	{"192.168.1.1/24", 256}, // wrong network address

	{"fe80::1/24", 0}, // non-IPv4
	{"fe80::1/64", 0}, // non-IPv4
}

var firstLastAdrTests = []struct {
	in    string
	first uint32
	last  uint32
}{
	{"0.0.0.0/0", 0, 4294967295},
	{"0.0.0.1/24", 0, 255},
	{"0.0.1.0/32", 256, 256},
	{"0.1.0.0/30", 65536, 65539},
	{"1.0.0.0/28", 16777216, 16777231},
	{"0.0.1.1/26", 256, 319},
	{"1.1.0.0/16", 16842752, 16908287},
	{"0.1.1.0/17", 65536, 98303},
	{"1.0.0.1/32", 16777217, 16777217},
	{"1.1.1.1/32", 16843009, 16843009},
	{"0.0.0.255/32", 255, 255},
	{"0.0.255.0/24", 65280, 65535},
	{"0.255.0.0/16", 16711680, 16777215},
	{"255.0.0.0/8", 4278190080, 4294967295},
	{"0.0.255.255/32", 65535, 65535},
	{"255.255.0.0/16", 4294901760, 4294967295},
	{"0.255.255.0/24", 16776960, 16777215},
	{"255.0.0.255/32", 4278190335, 4278190335},
	{"255.255.255.255/32", 4294967295, 4294967295},
	{"127.0.0.1/32", 2130706433, 2130706433}}

var overlapTests = []struct {
	net1 string
	net2 string
	out  bool
}{
	{"192.168.0.0/20", "192.168.1.0/24", true},
	{"192.168.0.0/24", "192.168.1.0/20", true},
	{"192.168.0.0/16", "192.169.1.0/24", false},
	{"10.0.0.0/7", "12.1.2.3/24", false},
}

func TestIPToUint32(t *testing.T) {
	for _, tt := range convertTests {
		if out := IPToUint32(tt.in); out != tt.out {
			t.Errorf("IPToUint32(%v) = %v, want %v", tt.in, out, tt.out)
		}
	}
}

func TestUint32ToIP(t *testing.T) {
	for _, tt := range convertTests {
		if out := Uint32ToIP(tt.out); !out.Equal(tt.in) {
			t.Errorf("Uint32ToIP(%v) = %v, want %v", tt.out, out, tt.in)
		}
	}
}

func TestNumAdr(t *testing.T) {
	for _, tt := range numAdrTests {
		_, net, _ := net.ParseCIDR(tt.in)
		if out := NumAdr(net); out != tt.out {
			t.Errorf("NumAddr(%q) = %v, want %v", net, out, tt.out)
		}
	}
}

func TestFirstAddr(t *testing.T) {
	for _, tt := range firstLastAdrTests {
		_, net, _ := net.ParseCIDR(tt.in)
		if out := FirstAdr(net); out != tt.first {
			t.Errorf("FirstAdr(%v) = %v, want %v", tt.in, out, tt.first)
		}
	}
}

func TestLastAddr(t *testing.T) {
	for _, tt := range firstLastAdrTests {
		_, net, _ := net.ParseCIDR(tt.in)
		if out := LastAdr(net); out != tt.last {
			t.Errorf("FirstAdr(%v) = %v, want %v", tt.in, out, tt.last)
		}
	}
}

func TestOverlap(t *testing.T) {
	for _, tt := range overlapTests {
		_, net1, _ := net.ParseCIDR(tt.net1)
		_, net2, _ := net.ParseCIDR(tt.net2)
		if out := Overlap(net1, net2); out != tt.out {
			t.Errorf("Overlap(%v,%v) = %v, want %v", net1, net2, out, tt.out)
		}
	}
}

func TestIsIPv4(t *testing.T) {
	ipv6 := net.ParseIP("fe80::1")
	ipv4 := net.ParseIP("127.0.0.1")
	if out := IsIPv4(&ipv6); out == true {
		t.Errorf("IsIPv4(%v) = %v, want %v", ipv6, out, false)
	}
	if out := IsIPv4(&ipv4); out != true {
		t.Errorf("IsIPv4(%v) = %v, want %v", ipv4, out, true)
	}
}
