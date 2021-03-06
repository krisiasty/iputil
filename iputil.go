// Copyright 2019 Krzysztof Cieplucha. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package iputil contains helper functions for IPv4 address manipulations
package iputil

import (
	"encoding/binary"
	"fmt"
	"net"
)

// IPToUint32 converts IPv4 address from net.IP to uint32 format
func IPToUint32(ip *net.IP) uint32 {
	if len(*ip) == 16 {
		return binary.BigEndian.Uint32((*ip)[12:16])
	}
	return binary.BigEndian.Uint32(*ip)
}

// Uint32ToIP converts IPv4 address from uint32 to net.IP format
func Uint32ToIP(adr uint32) *net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, adr)
	return &ip
}

// NumAdr returns number of IPv4 addresses in a given subnet or zero if not IPv4
func NumAdr(subnet *net.IPNet) uint64 {
	if !IsIPv4Net(subnet) {
		return 0
	}
	size, _ := subnet.Mask.Size()
	return uint64(1 << uint64(32-size))
}

// FirstAdr returns first IP address from a given subnet in uint32 format
func FirstAdr(subnet *net.IPNet) uint32 {
	return IPToUint32(&subnet.IP)
}

// LastAdr returns first IP address from a given subnet in uint64 format
func LastAdr(subnet *net.IPNet) uint32 {
	return FirstAdr(subnet) + uint32(NumAdr(subnet)-1)
}

// Overlap checks whether two given subnets overlaps
func Overlap(subnet1, subnet2 *net.IPNet) bool {
	if LastAdr(subnet1) < FirstAdr(subnet2) || LastAdr(subnet2) < FirstAdr(subnet1) {
		return false // subnets are not overlapping
	}
	return true
}

// IsIPv4 checks whether given ip address is IPv4
func IsIPv4(ip *net.IP) bool {
	return ip.To4() != nil
}

// IsIPv4Net checks whether given ip network address is IPv4
func IsIPv4Net(ipnet *net.IPNet) bool {
	return ipnet.IP.To4() != nil
}

// IncrIP returns given IP address increased by given offset
func IncrIP(ip *net.IP, offset uint32) *net.IP {
	return Uint32ToIP(IPToUint32(ip) + offset)
}

// IPMaskToString converts giver IPMask to dotted decimal format or returns nil for non-IPv4 masks
func IPMaskToString(mask *net.IPMask) string {
	if len(*mask) != 4 {
		return "<nil>"
	}
	return fmt.Sprintf("%d.%d.%d.%d", (*mask)[0], (*mask)[1], (*mask)[2], (*mask)[3])
}
