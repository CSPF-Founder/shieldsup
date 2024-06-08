// apierror.go

package utils

import (
	"fmt"
)

// ApiError represents an API error
type ApiErrorModule struct {
	Messages   []string
	StatusCode int
}

// NewApiError creates a new ApiError instance
func ApiError(message string, statusCode int) *ApiErrorModule {
	return &ApiErrorModule{
		Messages:   []string{message},
		StatusCode: statusCode,
	}
}

// NewApiErrorWithList creates a new ApiError instance with a list of messages
func NewApiErrorWithList(messages []string, statusCode int) *ApiErrorModule {
	return &ApiErrorModule{
		Messages:   messages,
		StatusCode: statusCode,
	}
}

// Error returns the error message
func (e *ApiErrorModule) Error() string {
	return fmt.Sprintf("API error: %v", e.Messages)
}
