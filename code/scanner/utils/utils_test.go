package utils

import (
	"testing"
)

func TestGetIPCountIfRange(t *testing.T) {
	got := GetIPCountIfRange("192.168.1.1/24")
	want := 256

	if got != want {
		t.Errorf("GetIPCountIfRange() = %v, want %v", got, want)
	}
}
