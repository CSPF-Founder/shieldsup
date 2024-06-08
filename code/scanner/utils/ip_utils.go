package utils

import (
	"net"
	"strings"
)

// cidrToIPList converts a CIDR notation string to a list of IP addresses.
func cidrToIPList(inputData string) []string {
	_, ipNet, err := net.ParseCIDR(inputData)
	if err != nil {
		return nil
	}

	var ipList []string
	for ip := ipNet.IP.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
		ipList = append(ipList, ip.String())
	}

	return ipList
}

// inc increments an IP address by one.
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// getIPCountIfRange returns the count of IP addresses in a given range specified in CIDR notation.
// It returns 0 if the input is empty or invalid.
func GetIPCountIfRange(inputData string) int {
	if inputData == "" || !strings.Contains(inputData, "/") {
		return 0
	}

	ipList := cidrToIPList(inputData)
	if ipList == nil {
		return 0
	}

	return len(ipList)
}
