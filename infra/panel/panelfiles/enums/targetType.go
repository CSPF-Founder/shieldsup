package enums

import (
	"net"
	"net/url"

	"github.com/CSPF-Founder/shieldsup/onpremise/panel/utils/iputils"
)

type TargetType string

const (
	TargetTypeIP      TargetType = "ip"
	TargetTypeIPRange TargetType = "ip_range"
	TargetTypeURL     TargetType = "url"
	TargetTypeInvalid TargetType = ""
)

func ParseTargetType(targetAddress string) TargetType {
	if ip := net.ParseIP(targetAddress); ip != nil {
		return TargetTypeIP
	} else if _, err := url.ParseRequestURI(targetAddress); err == nil {
		return TargetTypeURL
	} else {
		ipCount, err := iputils.ConvertIPRangeToIPSize(targetAddress)
		if err != nil {
			return TargetTypeInvalid
		}
		if ipCount == nil || ipCount.Int64() > 256 {
			return TargetTypeInvalid
		} else {
			return TargetTypeIPRange
		}
	}
}

// ParseTargetType(targetAddress string) string {
// 	if ip := net.ParseIP(targetAddress); ip != nil {
// 		return "ip"
// 	} else if _, err := url.ParseRequestURI(targetAddress); err == nil {
// 		return "url"
// 	} else {
// 		ipCount, _ := ConvertIPRangeToIPSize(targetAddress)
// 		if ipCount == nil || ipCount.Int64() > 256 {
// 			return ""
// 		} else {
// 			return "ip_range"
// 		}
// 	}
// }
