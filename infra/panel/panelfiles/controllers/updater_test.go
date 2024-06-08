package controllers

import "testing"

func TestParseRsyncURI(t *testing.T) {
	uri := "rsync://example.com:800/path/to/file"
	host, port, err := parseRsyncURI(uri)
	if err != nil {
		t.Fatalf("error parsing rsync URI: %v", err)
	}

	if host != "example.com" {
		t.Fatalf("invalid host. expected %s got %s", "example.com", host)
	}

	if port != "800" {
		t.Fatalf("invalid path. expected %s got %s", "/path/to/file", port)
	}

	uri2 := "rsync://example.com/path/to/file"
	host2, port2, err2 := parseRsyncURI(uri2)
	if err2 != nil {
		t.Fatalf("error parsing rsync URI: %v", err2)
	}

	if host2 != "example.com" {
		t.Fatalf("invalid host. expected %s got %s", "example.com", host2)
	}

	if port2 != "873" {
		t.Fatalf("invalid path. expected %s got %s", "/path/to/file", port2)
	}

}
