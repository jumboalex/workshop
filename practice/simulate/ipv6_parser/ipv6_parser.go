package main

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

// ParseIPv6 manually parses an IPv6 address into a 16-byte slice.
func ParseIPv6(addr string) ([]uint8, error) {
	// Check if there is an embedded IPv4 part
	ipv4Index := strings.LastIndex(addr, ".")
	var ipv4Bytes []byte
	if ipv4Index != -1 {
		// Extract IPv4 portion
		lastColon := strings.LastIndex(addr, ":")
		if lastColon == -1 {
			return nil, fmt.Errorf("invalid embedded IPv4: %s", addr)
		}
		ipv4Part := addr[lastColon+1:]
		addr = addr[:lastColon]

		// Parse IPv4 part manually
		ipv4Parts := strings.Split(ipv4Part, ".")
		if len(ipv4Parts) != 4 {
			return nil, fmt.Errorf("invalid IPv4 part: %s", ipv4Part)
		}
		for _, p := range ipv4Parts {
			v, err := strconv.Atoi(p)
			if err != nil || v < 0 || v > 255 {
				return nil, fmt.Errorf("invalid IPv4 octet: %s", p)
			}
			ipv4Bytes = append(ipv4Bytes, byte(v))
		}
	}

	// Split the IPv6 part on "::"
	var head, tail []string
	if strings.Contains(addr, "::") {
		// Check for multiple "::" (invalid)
		if strings.Count(addr, "::") > 1 {
			return nil, fmt.Errorf("multiple '::' found in address: %s", addr)
		}

		parts := strings.SplitN(addr, "::", 2)
		head = strings.Split(parts[0], ":")
		tail = strings.Split(parts[1], ":")
		if parts[0] == "" {
			head = []string{}
		}
		if parts[1] == "" {
			tail = []string{}
		}
	} else {
		head = strings.Split(addr, ":")
	}

	// Calculate how many groups of zeros to insert
	totalGroups := len(head) + len(tail)
	expectedGroups := 8
	if ipv4Bytes != nil {
		expectedGroups = 6 // IPv4 takes up 2 groups (4 bytes)
	}
	insertZeros := expectedGroups - totalGroups
	if insertZeros < 0 {
		return nil, fmt.Errorf("too many groups in address: %s", addr)
	}

	// Combine head + zero groups + tail
	var groups []string
	groups = append(groups, head...)
	for range insertZeros {
		groups = append(groups, "0")
	}
	groups = append(groups, tail...)

	// If IPv4 is present, add placeholder (counts as 2 groups)
	if ipv4Bytes != nil {
		groups = append(groups, "ipv4")
		expectedGroups = 7 // Adjust expected to account for ipv4 placeholder
	} else {
		expectedGroups = 8
	}

	if len(groups) != expectedGroups {
		return nil, fmt.Errorf("invalid IPv6 format (got %d groups, expected %d): %s", len(groups), expectedGroups, addr)
	}

	result := make([]uint8, 16)
	pos := 0

	for _, g := range groups {
		if g == "ipv4" {
			copy(result[pos:], ipv4Bytes)
			break
		}
		if g == "" {
			g = "0"
		}

		// Pad to 4 hex digits (2 bytes) for hex.DecodeString
		padded := strings.Repeat("0", 4-len(g)) + g

		bytes, err := hex.DecodeString(padded)
		if err != nil {
			return nil, fmt.Errorf("invalid group %q: %v", g, err)
		}
		copy(result[pos:], bytes)
		pos += 2
	}

	return result, nil
}

func main() {
	tests := []string{
		"1234:abcd:0000:0000:0102:0000:0000:fffe",
		"1234:abcd:0:0:102:0:0:fffe", // Same as above but with leading zeros omitted
		"::1",
		"::",
		"2001:db8::ff00:42:8329",
		"::ffff:192.168.1.1",
		"::192.0.2.1",             // IPv4-compatible
		"fe80::1",                 // Link-local
		"2001:db8::",              // Zero compression at end
		"::8a2e:370:7334",         // Zero compression at start
		"ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff", // All ones
	}

	fmt.Println("=== Valid IPv6 Addresses ===")
	for _, t := range tests {
		bytes, err := ParseIPv6(t)
		if err != nil {
			fmt.Printf("%-45s -> error: %v\n", t, err)
			continue
		}
		fmt.Printf("%-45s -> %v\n", t, bytes)
	}

	fmt.Println("\n=== Invalid IPv6 Addresses ===")
	invalidTests := []string{
		"1::2::3",                 // Multiple ::
		"gggg:1234:5678:90ab:cdef:1234:5678:90ab", // Invalid hex
		"1:2:3:4:5:6:7:8:9",       // Too many groups
		"::ffff:999.0.2.1",        // Invalid IPv4
	}

	for _, t := range invalidTests {
		bytes, err := ParseIPv6(t)
		if err != nil {
			fmt.Printf("%-45s -> error: %v\n", t, err)
		} else {
			fmt.Printf("%-45s -> %v (should have failed!)\n", t, bytes)
		}
	}
}
