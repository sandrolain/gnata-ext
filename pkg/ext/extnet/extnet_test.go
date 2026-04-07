package extnet_test

import (
	"testing"

	"github.com/sandrolain/gnata-ext/pkg/ext/extnet"
)

func TestIPVersion(t *testing.T) {
	fn := extnet.IPVersion()
	cases := []struct {
		input any
		want  float64
	}{
		{"192.168.1.1", 4},
		{"10.0.0.1", 4},
		{"::1", 6},
		{"2001:db8::1", 6},
		{"not-an-ip", -1},
		{"", -1},
	}
	for _, c := range cases {
		got, err := fn([]any{c.input}, nil)
		if err != nil {
			t.Errorf("IPVersion(%v) unexpected error: %v", c.input, err)
			continue
		}
		if got != c.want {
			t.Errorf("IPVersion(%v): got %v, want %v", c.input, got, c.want)
		}
	}
}

func TestIsPrivateIP(t *testing.T) {
	fn := extnet.IsPrivateIP()
	cases := []struct {
		input any
		want  bool
	}{
		{"10.0.0.1", true},
		{"172.16.0.1", true},
		{"192.168.1.1", true},
		{"127.0.0.1", true},
		{"::1", true},
		{"169.254.0.1", true},
		{"8.8.8.8", false},
		{"1.1.1.1", false},
		{"not-an-ip", false},
	}
	for _, c := range cases {
		got, _ := fn([]any{c.input}, nil)
		if got != c.want {
			t.Errorf("IsPrivateIP(%v): got %v, want %v", c.input, got, c.want)
		}
	}
}

func TestIPToInt(t *testing.T) {
	fn := extnet.IPToInt()

	got, err := fn([]any{"0.0.0.0"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != float64(0) {
		t.Errorf("0.0.0.0: got %v, want 0", got)
	}

	got, err = fn([]any{"255.255.255.255"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != float64(4294967295) {
		t.Errorf("255.255.255.255: got %v, want 4294967295", got)
	}

	got, err = fn([]any{"192.168.1.1"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 192.168.1.1 == 0xC0A80101 == 3232235777
	if got != float64(3232235777) {
		t.Errorf("192.168.1.1: got %v, want 3232235777", got)
	}

	_, err = fn([]any{"::1"}, nil)
	if err == nil {
		t.Error("expected error for IPv6 address")
	}
}

func TestIntToIP(t *testing.T) {
	fn := extnet.IntToIP()

	got, err := fn([]any{float64(0)}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "0.0.0.0" {
		t.Errorf("got %v, want 0.0.0.0", got)
	}

	got, err = fn([]any{float64(3232235777)}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "192.168.1.1" {
		t.Errorf("got %v, want 192.168.1.1", got)
	}

	_, err = fn([]any{float64(-1)}, nil)
	if err == nil {
		t.Error("expected error for negative value")
	}
}

func TestIPInCIDR(t *testing.T) {
	fn := extnet.IPInCIDR()

	got, err := fn([]any{"192.168.1.100", "192.168.1.0/24"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != true {
		t.Error("expected 192.168.1.100 to be in 192.168.1.0/24")
	}

	got, err = fn([]any{"10.0.0.1", "192.168.1.0/24"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != false {
		t.Error("expected 10.0.0.1 to NOT be in 192.168.1.0/24")
	}
}

func TestExpandCIDR(t *testing.T) {
	fn := extnet.ExpandCIDR()

	got, err := fn([]any{"192.168.1.0/24"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := got.(map[string]any)
	if m["network"] != "192.168.1.0" {
		t.Errorf("network: got %v, want 192.168.1.0", m["network"])
	}
	if m["broadcast"] != "192.168.1.255" {
		t.Errorf("broadcast: got %v, want 192.168.1.255", m["broadcast"])
	}
	if m["count"] != float64(256) {
		t.Errorf("count: got %v, want 256", m["count"])
	}

	// /32 single host
	got, err = fn([]any{"10.0.0.1/32"}, nil)
	if err != nil {
		t.Fatalf("/32 unexpected error: %v", err)
	}
	m = got.(map[string]any)
	if m["count"] != float64(1) {
		t.Errorf("/32 count: got %v, want 1", m["count"])
	}
}

func TestExpandCIDRIPv6(t *testing.T) {
	fn := extnet.ExpandCIDR()

	got, err := fn([]any{"2001:db8::/32"}, nil)
	if err != nil {
		t.Fatalf("IPv6 CIDR unexpected error: %v", err)
	}
	m := got.(map[string]any)
	if _, ok := m["network"]; !ok {
		t.Error("IPv6 result missing 'network'")
	}
	if _, ok := m["first"]; !ok {
		t.Error("IPv6 result missing 'first'")
	}
	if _, ok := m["last"]; !ok {
		t.Error("IPv6 result missing 'last'")
	}
	// IPv6 result should not have broadcast or count keys
	if _, ok := m["count"]; ok {
		t.Error("IPv6 result should not have 'count'")
	}
}

func TestAll(t *testing.T) {
	all := extnet.All()
	expected := []string{"ipVersion", "isPrivateIP", "ipToInt", "intToIP", "ipInCIDR", "expandCIDR"}
	for _, name := range expected {
		if _, ok := all[name]; !ok {
			t.Errorf("All() missing function: %q", name)
		}
	}
}

// --- additional coverage tests ---

func TestIPVersionNoArgs(t *testing.T) {
	fn := extnet.IPVersion()
	got, err := fn([]any{}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != float64(-1) {
		t.Errorf("got %v, want -1", got)
	}
}

func TestIPVersionNonString(t *testing.T) {
	fn := extnet.IPVersion()
	got, _ := fn([]any{42}, nil)
	if got != float64(-1) {
		t.Errorf("got %v, want -1", got)
	}
}

func TestIsPrivateIPNoArgs(t *testing.T) {
	fn := extnet.IsPrivateIP()
	got, err := fn([]any{}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != false {
		t.Errorf("got %v, want false", got)
	}
}

func TestIsPrivateIPNonString(t *testing.T) {
	fn := extnet.IsPrivateIP()
	got, _ := fn([]any{42}, nil)
	if got != false {
		t.Errorf("got %v, want false", got)
	}
}

func TestIPToIntNoArgs(t *testing.T) {
	fn := extnet.IPToInt()
	_, err := fn([]any{}, nil)
	if err == nil {
		t.Error("expected error for 0 args")
	}
}

func TestIPToIntNonString(t *testing.T) {
	fn := extnet.IPToInt()
	_, err := fn([]any{42}, nil)
	if err == nil {
		t.Error("expected error for non-string arg")
	}
}

func TestIPToIntInvalidIP(t *testing.T) {
	fn := extnet.IPToInt()
	_, err := fn([]any{"not-an-ip"}, nil)
	if err == nil {
		t.Error("expected error for invalid IP")
	}
}

func TestIntToIPNoArgs(t *testing.T) {
	fn := extnet.IntToIP()
	_, err := fn([]any{}, nil)
	if err == nil {
		t.Error("expected error for 0 args")
	}
}

func TestIntToIPNonNumeric(t *testing.T) {
	fn := extnet.IntToIP()
	_, err := fn([]any{"not-a-number"}, nil)
	if err == nil {
		t.Error("expected error for non-numeric arg")
	}
}

func TestIntToIPOutOfRange(t *testing.T) {
	fn := extnet.IntToIP()
	_, err := fn([]any{float64(5000000000)}, nil)
	if err == nil {
		t.Error("expected error for value exceeding uint32 max")
	}
}

func TestIPInCIDRNoArgs(t *testing.T) {
	fn := extnet.IPInCIDR()
	_, err := fn([]any{}, nil)
	if err == nil {
		t.Error("expected error for 0 args")
	}
}

func TestIPInCIDRNonStringIP(t *testing.T) {
	fn := extnet.IPInCIDR()
	_, err := fn([]any{42, "192.168.1.0/24"}, nil)
	if err == nil {
		t.Error("expected error for non-string ip")
	}
}

func TestIPInCIDRNonStringCIDR(t *testing.T) {
	fn := extnet.IPInCIDR()
	_, err := fn([]any{"192.168.1.1", 42}, nil)
	if err == nil {
		t.Error("expected error for non-string cidr")
	}
}

func TestIPInCIDRInvalidIP(t *testing.T) {
	fn := extnet.IPInCIDR()
	_, err := fn([]any{"not-an-ip", "192.168.1.0/24"}, nil)
	if err == nil {
		t.Error("expected error for invalid IP")
	}
}

func TestIPInCIDRInvalidCIDR(t *testing.T) {
	fn := extnet.IPInCIDR()
	_, err := fn([]any{"192.168.1.1", "not-a-cidr"}, nil)
	if err == nil {
		t.Error("expected error for invalid CIDR")
	}
}

func TestExpandCIDRNoArgs(t *testing.T) {
	fn := extnet.ExpandCIDR()
	_, err := fn([]any{}, nil)
	if err == nil {
		t.Error("expected error for 0 args")
	}
}

func TestExpandCIDRNonString(t *testing.T) {
	fn := extnet.ExpandCIDR()
	_, err := fn([]any{42}, nil)
	if err == nil {
		t.Error("expected error for non-string arg")
	}
}

func TestExpandCIDRInvalid(t *testing.T) {
	fn := extnet.ExpandCIDR()
	_, err := fn([]any{"not-a-cidr"}, nil)
	if err == nil {
		t.Error("expected error for invalid CIDR")
	}
}
