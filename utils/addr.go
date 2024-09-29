package utils

import (
	"fmt"
	"net"
)

// ResolveAddress takes a string address and resolves it to a normalized IP address.
// It returns the resolved address as a string and any error encountered.
func ResolveAddress(address string) (string, error) {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return "", fmt.Errorf("failed to resolve TCP address: %w", err)
	}

	resolvedIP := ResolveIP(addr.IP)
	if resolvedIP == nil {
		return "", fmt.Errorf("invalid IP address: %s", addr.IP)
	}

	return fmt.Sprintf("%s:%d", resolvedIP, addr.Port), nil
}

// resolveIP takes a net.IP and returns a resolved net.IP.
// It returns nil for loopback or unspecified IP addresses.
func ResolveIP(ip net.IP) net.IP {
	if ip == nil || ip.IsLoopback() || ip.IsUnspecified() {
		return nil
	}
	return ip
}

// NormalizeIP takes a net.IP and returns a normalized string representation.
// It returns an empty string for nil, loopback, or unspecified IP addresses.
func NormalizeIP(ip net.IP) string {
	resolvedIP := ResolveIP(ip)
	if resolvedIP == nil {
		return ""
	}
	return resolvedIP.String()
}
