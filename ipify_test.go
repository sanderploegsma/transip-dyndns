package main

import (
	"net"
	"testing"
)

func TestGetIPv4(t *testing.T) {
	ip, err := GetIPv4()
	if err != nil {
		t.Errorf("failed to get IPv4 address: %v", err)
	}

	parsed := net.ParseIP(ip)
	if parsed == nil {
		t.Errorf("got invalid IP address: %s", ip)
	}

	if parsed.To4() == nil {
		t.Errorf("got invalid IPv4 address: %s", ip)
	}
}

func TestGetIPv6(t *testing.T) {
	ip, err := GetIPv6()
	if err != nil {
		t.Errorf("failed to get IPv6 address: %v", err)
	}

	parsed := net.ParseIP(ip)
	if parsed == nil {
		t.Errorf("got invalid IP address: %s", ip)
	}

	if parsed.To4() != nil {
		t.Errorf("got invalid IPv6 address: %s", ip)
	}
}
