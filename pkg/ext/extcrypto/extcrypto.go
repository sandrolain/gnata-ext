// Package extcrypto provides cryptographic and UUID utility functions for gnata.
//
// Functions
//
//   - $uuid()                          – random UUID v4
//   - $hash(algorithm, value)          – hex-encoded hash (md5/sha1/sha256/sha384/sha512)
//   - $hmac(algorithm, key, value)     – hex-encoded HMAC
//   - $randomBytes(n)                  – n random bytes as hex string
//   - $base64url(str)                  – base64url-encode a string (no padding)
//   - $unbase64url(str)                – base64url-decode to string
package extcrypto

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"
	"strings"

	"github.com/google/uuid"
	"github.com/recolabs/gnata"
	"github.com/sandrolain/gnata-ext/pkg/ext/extutil"
)

// All returns a map of all extended crypto functions.
func All() map[string]gnata.CustomFunc {
	return map[string]gnata.CustomFunc{
		"uuid":        UUID(),
		"hash":        Hash(),
		"hmac":        HMAC(),
		"randomBytes": RandomBytes(),
		"base64url":   Base64URL(),
		"unbase64url": Unbase64URL(),
	}
}

// UUID returns the CustomFunc for $uuid().
// Generates a random UUID v4.
func UUID() gnata.CustomFunc {
	return func(_ []any, _ any) (any, error) {
		u, err := uuid.NewRandomFromReader(rand.Reader)
		if err != nil {
			return nil, fmt.Errorf("$uuid: %w", err)
		}
		return u.String(), nil
	}
}

// Hash returns the CustomFunc for $hash(algorithm, value).
func Hash() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$hash: requires 2 arguments (algorithm, value)")
		}
		algorithm, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$hash: algorithm must be a string")
		}
		value, ok := args[1].(string)
		if !ok {
			return nil, fmt.Errorf("$hash: value must be a string")
		}
		h, err := newHasher(algorithm)
		if err != nil {
			return nil, fmt.Errorf("$hash: %w", err)
		}
		h.Write([]byte(value))
		return hex.EncodeToString(h.Sum(nil)), nil
	}
}

// HMAC returns the CustomFunc for $hmac(algorithm, key, value).
func HMAC() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("$hmac: requires 3 arguments (algorithm, key, value)")
		}
		algorithm, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$hmac: algorithm must be a string")
		}
		key, ok := args[1].(string)
		if !ok {
			return nil, fmt.Errorf("$hmac: key must be a string")
		}
		value, ok := args[2].(string)
		if !ok {
			return nil, fmt.Errorf("$hmac: value must be a string")
		}
		var mac hash.Hash
		switch strings.ToLower(algorithm) {
		case "md5":
			mac = hmac.New(md5.New, []byte(key))
		case "sha1":
			mac = hmac.New(sha1.New, []byte(key))
		case "sha256":
			mac = hmac.New(sha256.New, []byte(key))
		case "sha384":
			mac = hmac.New(sha512.New384, []byte(key))
		case "sha512":
			mac = hmac.New(sha512.New, []byte(key))
		default:
			return nil, fmt.Errorf("$hmac: unsupported algorithm %q", algorithm)
		}
		mac.Write([]byte(value))
		return hex.EncodeToString(mac.Sum(nil)), nil
	}
}

// newHasher returns a hash.Hash for the given algorithm name.
func newHasher(algorithm string) (hash.Hash, error) {
	switch strings.ToLower(algorithm) {
	case "md5":
		return md5.New(), nil
	case "sha1":
		return sha1.New(), nil
	case "sha256":
		return sha256.New(), nil
	case "sha384":
		return sha512.New384(), nil
	case "sha512":
		return sha512.New(), nil
	default:
		return nil, fmt.Errorf("unsupported algorithm %q", algorithm)
	}
}

// RandomBytes returns the CustomFunc for $randomBytes(n).
// Returns n random bytes as a lowercase hex string.
func RandomBytes() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$randomBytes: requires 1 argument (n)")
		}
		n, ok := extutil.ToInt(args[0])
		if !ok || n <= 0 {
			return nil, fmt.Errorf("$randomBytes: n must be a positive integer")
		}
		buf := make([]byte, n)
		if _, err := rand.Read(buf); err != nil {
			return nil, fmt.Errorf("$randomBytes: %w", err)
		}
		return hex.EncodeToString(buf), nil
	}
}

// Base64URL returns the CustomFunc for $base64url(str).
// Base64url-encodes str without padding.
func Base64URL() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$base64url: requires 1 argument")
		}
		s, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$base64url: argument must be a string")
		}
		return base64.RawURLEncoding.EncodeToString([]byte(s)), nil
	}
}

// Unbase64URL returns the CustomFunc for $unbase64url(str).
// Decodes a base64url string (no padding) to a string.
func Unbase64URL() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$unbase64url: requires 1 argument")
		}
		s, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$unbase64url: argument must be a string")
		}
		b, err := base64.RawURLEncoding.DecodeString(s)
		if err != nil {
			return nil, fmt.Errorf("$unbase64url: %w", err)
		}
		return string(b), nil
	}
}
