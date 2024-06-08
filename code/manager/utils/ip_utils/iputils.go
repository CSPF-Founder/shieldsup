package iputils

import (
	"math"
	"net"
	"regexp"
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
