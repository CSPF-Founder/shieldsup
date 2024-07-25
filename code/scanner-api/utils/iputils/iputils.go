package iputils

import (
	"errors"
	"math"
	"math/big"
	"net"
	"regexp"
	"strings"
)

func GetIPCountIfRange(cidr string) (int, error) {

	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return 0, err
	}

	// Get the number of bits in the network mask
	ones, _ := ipnet.Mask.Size()

	// Calculate the number of host bits
	hostBits := 32 - ones

	// Calculate the number of usable IP addresses (excluding network and broadcast addresses)
	totalIPs := int(math.Pow(2, float64(hostBits)))

	return totalIPs, nil
}

func IsValidCIDRByRegex(cidrStr string) bool {
	// Define the regular expression pattern for CIDR notation
	pattern := `^(?:(?:25[0-5]|2[0-4][0-9]|[0-1]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[0-1]?[0-9][0-9]?)\/(?:[1-9]|[1-2][0-9]|3[0-2])$`

	// Use the regexp.MatchString function to check if the input matches the pattern
	match, _ := regexp.MatchString(pattern, cidrStr)
	return match
}

func ConvertIPRangeToIPSize(ipRange string) (*big.Int, error) {
	if ipRange == "" {
		return nil, nil
	}

	ipRange = strings.TrimSpace(ipRange)

	if strings.Contains(ipRange, "-") {
		return nil, errors.New("Invalid Target")
	}

	if strings.Contains(ipRange, "/") {

		totalIPs, err := countCIDR(ipRange)
		if err != nil {
			return nil, err
		}

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
