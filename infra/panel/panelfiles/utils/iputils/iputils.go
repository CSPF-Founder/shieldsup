package iputils

import (
	"errors"
	"math"
	"math/big"
	"net"
	"strconv"
	"strings"
)

func ConvertIPRangeToIPSize(ipRange string) (*big.Int, error) {
	if ipRange == "" {
		return nil, nil
	}

	ipRange = strings.TrimSpace(ipRange)

	if strings.Contains(ipRange, "-") {
		return nil, errors.New("Invalid Target")
	}

	if strings.Contains(ipRange, "/") {
		// ipRangeSplitted := strings.Split(ipRange, "/")
		// cidrPrefix := ipRangeSplitted[0]
		// cidrSize, err := strconv.Atoi(ipRangeSplitted[1])
		// if err != nil {
		// 	return nil, nil
		// }

		// if cidrSize < 24 {
		// 	return nil, errors.New("Please provide single ip range")
		// }

		// if ip := net.ParseIP(cidrPrefix); ip == nil {
		// 	return nil, errors.New("Invalid IP address")
		// }

		totalIPs, err := countCIDR(ipRange)
		if err != nil {
			return nil, err
		}

		// if totalIPs.Int64() > 256 {
		// 	return nil, errors.New("Please provide single ip range")
		// }
		return totalIPs, nil
	}

	return nil, nil
}

// CIDRRangeSize calculates the total number of IP addresses in a CIDR range.
func countCIDR(cidr string) (*big.Int, error) {
	var totalIPs *big.Int
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return totalIPs, err
	}

	// Get the number of bits in the mask
	ones, bits := ipNet.Mask.Size()

	// Calculate the number of possible IPs
	// The formula is 2^(bits - ones)
	totalIPs = big.NewInt(0).Exp(big.NewInt(2), big.NewInt(int64(bits-ones)), nil)

	return totalIPs, nil
}

// CIDRToIPList converts CIDR notation to a list of IP addresses
func CIDRToIPList(inputData string) ([]string, error) {
	ipList := []string{}
	ipRange, err := CIDRToRange(inputData)
	if err != nil {
		return ipList, err
	}

	if ipRange != nil {
		ipList = ipRangeToList(ipRange["start"], ipRange["end"])
	}

	return ipList, nil
}

// CIDRToRange converts CIDR notation to IP range
func CIDRToRange(cidr string) (map[string]string, error) {
	rangeMap := make(map[string]string)
	parts := strings.Split(cidr, "/")
	cidrSize, err := strconv.Atoi(parts[1])
	if err != nil || cidrSize < 1 || cidrSize > 32 {
		return nil, errors.New("Incorrect CIDR")
	}

	ip := net.ParseIP(parts[0])
	if ip == nil {
		return nil, errors.New("Invalid IP address")
	}

	mask := net.CIDRMask(cidrSize, 32)
	network := ip.Mask(mask)
	startIP := network.String()
	endIP := incrementIP(network, uint32(math.Pow(2, float64(32-cidrSize)))-1)

	rangeMap["start"] = startIP
	rangeMap["end"] = endIP.String()

	return rangeMap, nil
}

// ipRangeToList converts IP range to a list of IP addresses
func ipRangeToList(startingPart, endingPart string) []string {
	var ipList []string

	startIP := net.ParseIP(startingPart)
	endIP := net.ParseIP(endingPart)
	if startIP == nil || endIP == nil {
		return nil
	}

	for ip := startIP; ip.String() <= endIP.String(); ip = incrementIP(ip, 1) {
		ipList = append(ipList, ip.String())
	}

	return ipList
}

// incrementIP increments IP address by a given value
func incrementIP(ip net.IP, val uint32) net.IP {
	ip = ip.To4()
	if ip == nil {
		return nil
	}

	num := uint32(ip[0])<<24 + uint32(ip[1])<<16 + uint32(ip[2])<<8 + uint32(ip[3])
	num += val

	return net.IPv4(byte(num>>24), byte(num>>16), byte(num>>8), byte(num))
}
