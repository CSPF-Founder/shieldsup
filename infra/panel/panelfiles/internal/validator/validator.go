package validator

import "regexp"

func IsValidUsername(username string) bool {
	if username != "" {
		match, _ := regexp.MatchString("^[A-Za-z][a-zA-Z0-9_]{1,30}$", username)
		return match
	}
	return false
}
