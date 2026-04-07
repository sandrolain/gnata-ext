package extcrypto_test

import (
	"strings"
	"testing"

	"github.com/sandrolain/gnata-ext/pkg/ext/extcrypto"
)

func invoke(f func([]any, any) (any, error), args ...any) (any, error) {
	return f(args, nil)
}

func TestUUID(t *testing.T) {
	f := extcrypto.UUID()
	got, err := invoke(f)
	if err != nil {
		t.Fatalf("uuid(): unexpected error: %v", err)
	}
	s, ok := got.(string)
	if !ok {
		t.Fatal("uuid(): expected string result")
	}
	parts := strings.Split(s, "-")
	if len(parts) != 5 {
		t.Errorf("uuid(): expected 5 parts, got %q", s)
	}
	// Version nibble should be 4
	if len(parts[2]) < 1 || parts[2][0] != '4' {
		t.Errorf("uuid(): expected version 4, got %q", s)
	}
	// Two different calls should produce different UUIDs
	got2, _ := invoke(f)
	if got.(string) == got2.(string) {
		t.Error("uuid(): two consecutive calls returned same UUID")
	}
}

func TestHash(t *testing.T) {
	f := extcrypto.Hash()
	cases := []struct {
		alg  string
		val  string
		want string
	}{
		{"md5", "", "d41d8cd98f00b204e9800998ecf8427e"},
		{"sha1", "", "da39a3ee5e6b4b0d3255bfef95601890afd80709"},
		{"sha256", "", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{"sha512", "", "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e"},
	}
	for _, c := range cases {
		got, err := invoke(f, c.alg, c.val)
		if err != nil {
			t.Errorf("hash(%s, empty): unexpected error: %v", c.alg, err)
			continue
		}
		if !strings.HasPrefix(got.(string), c.want[:8]) {
			t.Errorf("hash(%s): got %v; want prefix %v", c.alg, got, c.want[:8])
		}
	}
}

func TestHashUnsupported(t *testing.T) {
	f := extcrypto.Hash()
	_, err := invoke(f, "unknown", "data")
	if err == nil {
		t.Error("hash with unknown algorithm: expected error")
	}
}

func TestHMAC(t *testing.T) {
	f := extcrypto.HMAC()
	got, err := invoke(f, "sha256", "key", "message")
	if err != nil {
		t.Fatalf("hmac(sha256): unexpected error: %v", err)
	}
	s, ok := got.(string)
	if !ok || len(s) != 64 {
		t.Errorf("hmac(sha256): expected 64-char hex, got %q", got)
	}
}

func TestHMACUnsupported(t *testing.T) {
	f := extcrypto.HMAC()
	_, err := invoke(f, "unknown", "key", "message")
	if err == nil {
		t.Error("hmac with unknown algorithm: expected error")
	}
}

func TestAll(t *testing.T) {
	all := extcrypto.All()
	expected := []string{"uuid", "hash", "hmac"}
	for _, name := range expected {
		if _, ok := all[name]; !ok {
			t.Errorf("All(): missing function %q", name)
		}
	}
}

// --- Additional coverage tests ---

func TestUUIDErrors(t *testing.T) {
	f := extcrypto.UUID()
	// Normal call should succeed
	got, err := f([]any{}, nil)
	if err != nil {
		t.Fatalf("uuid: unexpected error: %v", err)
	}
	if _, ok := got.(string); !ok {
		t.Errorf("uuid: expected string, got %T", got)
	}
}

func TestHashErrors(t *testing.T) {
	f := extcrypto.Hash()

	// too few args
	_, err := f([]any{}, nil)
	if err == nil {
		t.Error("expected error for 0 args")
	}
	// non-string algorithm
	_, err = f([]any{42, "hello"}, nil)
	if err == nil {
		t.Error("expected error for non-string algorithm")
	}
	// non-string value
	_, err = f([]any{"sha256", 42}, nil)
	if err == nil {
		t.Error("expected error for non-string value")
	}
	// unknown algorithm
	_, err = f([]any{"xxhash", "hello"}, nil)
	if err == nil {
		t.Error("expected error for unknown algorithm")
	}
}

func TestHashAllAlgorithms(t *testing.T) {
	f := extcrypto.Hash()
	algos := []string{"md5", "sha1", "sha256", "sha384", "sha512"}
	for _, algo := range algos {
		got, err := f([]any{algo, "test"}, nil)
		if err != nil {
			t.Errorf("hash %s: unexpected error: %v", algo, err)
		}
		if got.(string) == "" {
			t.Errorf("hash %s: empty result", algo)
		}
	}
}

func TestHMACErrors(t *testing.T) {
	f := extcrypto.HMAC()

	// too few args
	_, err := f([]any{}, nil)
	if err == nil {
		t.Error("expected error for 0 args")
	}
	// non-string algorithm
	_, err = f([]any{42, "key", "val"}, nil)
	if err == nil {
		t.Error("expected error for non-string algorithm")
	}
	// non-string key
	_, err = f([]any{"sha256", 42, "val"}, nil)
	if err == nil {
		t.Error("expected error for non-string key")
	}
	// non-string value
	_, err = f([]any{"sha256", "key", 42}, nil)
	if err == nil {
		t.Error("expected error for non-string value")
	}
	// unknown algorithm
	_, err = f([]any{"xxhash", "key", "val"}, nil)
	if err == nil {
		t.Error("expected error for unknown algorithm")
	}
}

func TestHMACAllAlgorithms(t *testing.T) {
	f := extcrypto.HMAC()
	algos := []string{"md5", "sha1", "sha256", "sha384", "sha512"}
	for _, algo := range algos {
		got, err := f([]any{algo, "secret", "message"}, nil)
		if err != nil {
			t.Errorf("hmac %s: unexpected error: %v", algo, err)
		}
		if got.(string) == "" {
			t.Errorf("hmac %s: empty result", algo)
		}
	}
}
