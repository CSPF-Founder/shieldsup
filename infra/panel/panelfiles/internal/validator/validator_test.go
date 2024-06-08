package validator

import "testing"

func TestIsValidUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		want     bool
	}{
		{
			name:     "valid username",
			username: "john_doe",
			want:     true,
		},
		{
			name:     "invalid username",
			username: "john.doe",
			want:     false,
		},
		{
			name:     "empty username",
			username: "",
			want:     false,
		},
		{
			name:     "username too short",
			username: "j",
			want:     false,
		},
		{
			name:     "username too long",
			username: "john_doe_john_doe_john_doe_john_doe",
			want:     false,
		},
		{
			name:     "username starts with a number",
			username: "1john_doe",
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidUsername(tt.username); got != tt.want {
				t.Errorf("IsValidUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}
