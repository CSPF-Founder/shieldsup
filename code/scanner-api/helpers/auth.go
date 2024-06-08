package helpers

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerValue(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no auth headers passed")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed auth header")
	}

	if vals[0] != "Bearer" {
		return "", errors.New("malformed first part of the auth header")
	}

	return vals[1], nil
}
