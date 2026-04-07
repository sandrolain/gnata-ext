// Package extnet provides IP-address and network utility functions for gnata.
// All IP operations use the Go standard library net package.
//
// Functions
//
//   - $ipVersion(str)       – returns 4.0, 6.0, or -1.0
//   - $isPrivateIP(str)     – true if RFC1918 / loopback / link-local
//   - $ipToInt(str)         – IPv4 address string → uint32 as float64
//   - $intToIP(n)           – uint32 as float64 → IPv4 address string
//   - $ipInCIDR(ip, cidr)   – true if ip is contained in CIDR block
//   - $expandCIDR(cidr)     – returns network info object for the CIDR block
package extnet

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"

	"github.com/recolabs/gnata"
	"github.com/sandrolain/gnata-ext/pkg/ext/extutil"
)

// privateBlocks lists RFC1918, loopback, and link-local CIDR ranges.
var privateBlocks []*net.IPNet

func init() {
	cidrs := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"127.0.0.0/8",
		"::1/128",
		"169.254.0.0/16",
		"fe80::/10",
	}
	for _, cidr := range cidrs {
		_, block, err := net.ParseCIDR(cidr)
		if err == nil {
			privateBlocks = append(privateBlocks, block)
		}
	}
}

// All returns a map of all network functions.
func All() map[string]gnata.CustomFunc {
	return map[string]gnata.CustomFunc{
		"ipVersion":   IPVersion(),
		"isPrivateIP": IsPrivateIP(),
		"ipToInt":     IPToInt(),
		"intToIP":     IntToIP(),
		"ipInCIDR":    IPInCIDR(),
		"expandCIDR":  ExpandCIDR(),
	}
}

// IPVersion returns the CustomFunc for $ipVersion(str).
// Returns 4.0 for IPv4, 6.0 for IPv6, -1.0 if the string is not a valid IP.
func IPVersion() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return float64(-1), nil
		}
		s, ok := args[0].(string)
		if !ok {
			return float64(-1), nil
		}
		ip := net.ParseIP(s)
		if ip == nil {
			return float64(-1), nil
		}
		if ip.To4() != nil {
			return float64(4), nil
		}
		return float64(6), nil
	}
}

// IsPrivateIP returns the CustomFunc for $isPrivateIP(str).
// Returns true for RFC1918, loopback, and link-local addresses.
func IsPrivateIP() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return false, nil
		}
		s, ok := args[0].(string)
		if !ok {
			return false, nil
		}
		ip := net.ParseIP(s)
		if ip == nil {
			return false, nil
		}
		for _, block := range privateBlocks {
			if block.Contains(ip) {
				return true, nil
			}
		}
		return false, nil
	}
}

// IPToInt returns the CustomFunc for $ipToInt(str).
// Converts an IPv4 address string to a 32-bit unsigned integer as float64.
func IPToInt() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$ipToInt: requires 1 argument (ip)")
		}
		s, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$ipToInt: argument must be a string")
		}
		ip := net.ParseIP(s)
		if ip == nil {
			return nil, fmt.Errorf("$ipToInt: invalid IP address %q", s)
		}
		ip4 := ip.To4()
		if ip4 == nil {
			return nil, fmt.Errorf("$ipToInt: only IPv4 addresses are supported")
		}
		n := binary.BigEndian.Uint32(ip4)
		return float64(n), nil
	}
}

// IntToIP returns the CustomFunc for $intToIP(n).
// Converts a 32-bit unsigned integer (as float64) to an IPv4 address string.
func IntToIP() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$intToIP: requires 1 argument (n)")
		}
		f, err := extutil.ToFloat(args[0])
		if err != nil {
			return nil, fmt.Errorf("$intToIP: %w", err)
		}
		if f < 0 || f > math.MaxUint32 {
			return nil, fmt.Errorf("$intToIP: value %v out of IPv4 range", f)
		}
		n := uint32(f)
		b := make([]byte, 4)
		binary.BigEndian.PutUint32(b, n)
		return net.IP(b).String(), nil
	}
}

// IPInCIDR returns the CustomFunc for $ipInCIDR(ip, cidr).
// Returns true if the IP address is within the CIDR block.
func IPInCIDR() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$ipInCIDR: requires 2 arguments (ip, cidr)")
		}
		ipStr, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$ipInCIDR: ip must be a string")
		}
		cidrStr, ok := args[1].(string)
		if !ok {
			return nil, fmt.Errorf("$ipInCIDR: cidr must be a string")
		}
		ip := net.ParseIP(ipStr)
		if ip == nil {
			return nil, fmt.Errorf("$ipInCIDR: invalid IP address %q", ipStr)
		}
		_, network, err := net.ParseCIDR(cidrStr)
		if err != nil {
			return nil, fmt.Errorf("$ipInCIDR: invalid CIDR %q: %w", cidrStr, err)
		}
		return network.Contains(ip), nil
	}
}

// ExpandCIDR returns the CustomFunc for $expandCIDR(cidr).
// For IPv4: returns {network, broadcast, first, last, count}.
// For IPv6: returns {network, first, last} (no broadcast concept; count omitted due to scale).
func ExpandCIDR() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$expandCIDR: requires 1 argument (cidr)")
		}
		cidrStr, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$expandCIDR: argument must be a string")
		}
		ip, network, err := net.ParseCIDR(cidrStr)
		if err != nil {
			return nil, fmt.Errorf("$expandCIDR: invalid CIDR %q: %w", cidrStr, err)
		}

		// Determine if IPv4
		if ip.To4() != nil || network.IP.To4() != nil {
			return expandIPv4CIDR(network)
		}
		return expandIPv6CIDR(network)
	}
}

func expandIPv4CIDR(network *net.IPNet) (map[string]any, error) {
	ip4 := network.IP.To4()
	if ip4 == nil {
		return nil, fmt.Errorf("not an IPv4 network")
	}
	mask := network.Mask
	ones, bits := mask.Size()
	// first usable (network + 1) and last usable (broadcast - 1)
	// network address
	netAddr := make(net.IP, 4)
	copy(netAddr, ip4)

	// broadcast: network | ^mask
	broadcast := make(net.IP, 4)
	for i := range ip4 {
		broadcast[i] = ip4[i] | ^mask[i]
	}

	// first host
	first := make(net.IP, 4)
	copy(first, netAddr)
	first[3]++

	// last host
	last := make(net.IP, 4)
	copy(last, broadcast)
	last[3]--

	count := math.Pow(2, float64(bits-ones))

	return map[string]any{
		"network":   netAddr.String(),
		"broadcast": broadcast.String(),
		"first":     first.String(),
		"last":      last.String(),
		"count":     count,
	}, nil
}

func expandIPv6CIDR(network *net.IPNet) (map[string]any, error) {
	ip6 := network.IP

	// first address is the network address itself
	first := make(net.IP, len(ip6))
	copy(first, ip6)

	// last address: network | ^mask
	mask := network.Mask
	last := make(net.IP, len(ip6))
	for i := range ip6 {
		last[i] = ip6[i] | ^mask[i]
	}

	return map[string]any{
		"network": network.IP.String(),
		"first":   first.String(),
		"last":    last.String(),
	}, nil
}
