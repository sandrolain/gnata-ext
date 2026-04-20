// Package middleware provides wrappers that add cross-cutting behaviour
// (logging, memoization, rate limiting, timeout, retry, circuit breaking, audit
// logging) to maps of gnata custom functions.
//
// Each wrapper returns a new map; the original map is left unchanged.
//
//	funcs := ext.AllFuncs()
//	funcs = middleware.WithLogging(funcs, slog.Default())
//	funcs = middleware.WithTimeout(funcs, 500*time.Millisecond)
//	funcs = middleware.WithRetry(funcs, middleware.MaxAttempts(3))
//	env := gnata.NewCustomEnv(funcs)
package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"time"

	"github.com/recolabs/gnata"
)

// WithLogging wraps each function in funcs to log the function name, arguments,
// execution duration, and any error using the provided slog.Logger.
// The result is a new map; the original is not modified.
func WithLogging(funcs map[string]gnata.CustomFunc, logger *slog.Logger) map[string]gnata.CustomFunc {
	out := make(map[string]gnata.CustomFunc, len(funcs))
	for name, fn := range funcs {
		name, fn := name, fn
		out[name] = func(args []any, focus any) (any, error) {
			start := time.Now()
			result, err := fn(args, focus)
			dur := time.Since(start)
			if err != nil {
				logger.Error("gnata-ext function error",
					"func", name,
					"args", fmt.Sprintf("%v", args),
					"duration", dur,
					"error", err,
				)
			} else {
				logger.Debug("gnata-ext function called",
					"func", name,
					"args", fmt.Sprintf("%v", args),
					"duration", dur,
				)
			}
			return result, err
		}
	}
	return out
}

// WithMemoize wraps each function in funcs with a simple in-memory cache.
// Results are keyed by the function name and a string representation of the
// arguments. Functions listed in exclude are not memoized (e.g. "uuid").
// The cache is per-wrapped-map and is not shared between calls to WithMemoize.
func WithMemoize(funcs map[string]gnata.CustomFunc, exclude ...string) map[string]gnata.CustomFunc {
	skip := make(map[string]bool, len(exclude))
	for _, n := range exclude {
		skip[n] = true
	}

	type cacheEntry struct {
		result any
		err    error
	}
	var mu sync.Mutex
	cache := make(map[string]cacheEntry)

	out := make(map[string]gnata.CustomFunc, len(funcs))
	for name, fn := range funcs {
		name, fn := name, fn
		if skip[name] {
			out[name] = fn
			continue
		}
		out[name] = func(args []any, focus any) (any, error) {
			key := name + ":" + fmt.Sprintf("%v", args)
			mu.Lock()
			if entry, ok := cache[key]; ok {
				mu.Unlock()
				return entry.result, entry.err
			}
			mu.Unlock()

			result, err := fn(args, focus)

			mu.Lock()
			cache[key] = cacheEntry{result, err}
			mu.Unlock()

			return result, err
		}
	}
	return out
}

// WithRateLimit wraps each function in funcs with a token-bucket rate limiter
// that allows at most rps calls per second across all wrapped functions sharing
// the same bucket. Calls that exceed the limit block until a token is available.
func WithRateLimit(funcs map[string]gnata.CustomFunc, rps float64) map[string]gnata.CustomFunc {
	if rps <= 0 {
		// No limiting — return a plain copy.
		out := make(map[string]gnata.CustomFunc, len(funcs))
		for k, v := range funcs {
			out[k] = v
		}
		return out
	}

	interval := time.Duration(float64(time.Second) / rps)
	ticker := time.NewTicker(interval)

	out := make(map[string]gnata.CustomFunc, len(funcs))
	for name, fn := range funcs {
		name, fn := name, fn
		out[name] = func(args []any, focus any) (any, error) {
			<-ticker.C
			return fn(args, focus)
		}
	}
	return out
}

// WithTimeout wraps each function so that calls exceeding timeout return an
// error. The underlying function continues running in a goroutine after the
// timeout but its result is discarded.
func WithTimeout(funcs map[string]gnata.CustomFunc, timeout time.Duration) map[string]gnata.CustomFunc {
	out := make(map[string]gnata.CustomFunc, len(funcs))
	for name, fn := range funcs {
		name, fn := name, fn
		out[name] = func(args []any, focus any) (any, error) {
			type result struct {
				val any
				err error
			}
			ch := make(chan result, 1)
			go func() {
				v, e := fn(args, focus)
				ch <- result{v, e}
			}()
			select {
			case r := <-ch:
				return r.val, r.err
			case <-time.After(timeout):
				return nil, fmt.Errorf("gnata-ext: %q exceeded timeout %s", name, timeout)
			}
		}
	}
	return out
}

// RetryOption configures WithRetry behaviour.
type RetryOption func(*retryConfig)

type retryConfig struct {
	maxAttempts int
	backoff     time.Duration
	shouldRetry func(error) bool
}

// MaxAttempts sets the maximum number of call attempts (first call + retries).
// Defaults to 3.
func MaxAttempts(n int) RetryOption {
	return func(c *retryConfig) { c.maxAttempts = n }
}

// Backoff sets the initial wait duration between retries. Each subsequent
// retry doubles the wait (exponential backoff). Defaults to 100ms.
func Backoff(d time.Duration) RetryOption {
	return func(c *retryConfig) { c.backoff = d }
}

// RetryOn sets the predicate that decides whether to retry after an error.
// Defaults to IsTransient (retry all errors).
func RetryOn(f func(error) bool) RetryOption {
	return func(c *retryConfig) { c.shouldRetry = f }
}

// IsTransient is a RetryOn predicate that retries all errors.
func IsTransient(_ error) bool { return true }

// WithRetry wraps each function in funcs so that errors are retried with
// exponential backoff. Use RetryOption helpers to configure behaviour.
func WithRetry(funcs map[string]gnata.CustomFunc, opts ...RetryOption) map[string]gnata.CustomFunc {
	cfg := &retryConfig{
		maxAttempts: 3,
		backoff:     100 * time.Millisecond,
		shouldRetry: IsTransient,
	}
	for _, o := range opts {
		o(cfg)
	}

	out := make(map[string]gnata.CustomFunc, len(funcs))
	for name, fn := range funcs {
		_, fn := name, fn
		out[name] = func(args []any, focus any) (any, error) {
			var lastErr error
			wait := cfg.backoff
			for attempt := 0; attempt < cfg.maxAttempts; attempt++ {
				v, err := fn(args, focus)
				if err == nil {
					return v, nil
				}
				lastErr = err
				if !cfg.shouldRetry(err) || attempt == cfg.maxAttempts-1 {
					break
				}
				time.Sleep(wait)
				wait *= 2
			}
			return nil, lastErr
		}
	}
	return out
}

// CircuitBreakerOption configures WithCircuitBreaker behaviour.
type CircuitBreakerOption func(*cbConfig)

type cbConfig struct {
	threshold  int
	resetAfter time.Duration
}

// Threshold sets the number of consecutive failures that open the circuit.
// Defaults to 5.
func Threshold(n int) CircuitBreakerOption {
	return func(c *cbConfig) { c.threshold = n }
}

// ResetAfter sets the duration before an open circuit is probed again.
// Defaults to 30 seconds.
func ResetAfter(d time.Duration) CircuitBreakerOption {
	return func(c *cbConfig) { c.resetAfter = d }
}

// ErrCircuitOpen is returned when a circuit breaker is in the open state.
var ErrCircuitOpen = errors.New("gnata-ext: circuit breaker open")

type circuitState int

const (
	stateClosed circuitState = iota
	stateOpen
)

type breaker struct {
	mu       sync.Mutex
	state    circuitState
	failures int
	openedAt time.Time
}

// WithCircuitBreaker wraps each function with an independent per-function
// circuit breaker. After Threshold consecutive failures the circuit opens and
// subsequent calls return ErrCircuitOpen without executing the function. After
// ResetAfter has elapsed the circuit closes and normal execution resumes.
func WithCircuitBreaker(funcs map[string]gnata.CustomFunc, opts ...CircuitBreakerOption) map[string]gnata.CustomFunc {
	cfg := &cbConfig{
		threshold:  5,
		resetAfter: 30 * time.Second,
	}
	for _, o := range opts {
		o(cfg)
	}

	out := make(map[string]gnata.CustomFunc, len(funcs))
	for name, fn := range funcs {
		name, fn := name, fn
		b := &breaker{}
		out[name] = func(args []any, focus any) (any, error) {
			b.mu.Lock()
			if b.state == stateOpen {
				if time.Since(b.openedAt) < cfg.resetAfter {
					b.mu.Unlock()
					return nil, fmt.Errorf("%w: %s", ErrCircuitOpen, name)
				}
				// half-open: allow one probe attempt
				b.state = stateClosed
				b.failures = 0
			}
			b.mu.Unlock()

			v, err := fn(args, focus)

			b.mu.Lock()
			if err != nil {
				b.failures++
				if b.failures >= cfg.threshold {
					b.state = stateOpen
					b.openedAt = time.Now()
				}
			} else {
				b.failures = 0
			}
			b.mu.Unlock()

			return v, err
		}
	}
	return out
}

// AuditRecord is the structure written to the audit log for each function call.
type AuditRecord struct {
	Time     string `json:"time"`
	Func     string `json:"func"`
	Args     any    `json:"args,omitempty"`
	Duration string `json:"duration"`
	Error    string `json:"error,omitempty"`
}

// WithAuditLog wraps each function in funcs so that every call appends an
// AuditRecord as a JSON line (NDJSON) to w. Writes are serialised with a mutex
// so w does not need to be goroutine-safe.
func WithAuditLog(funcs map[string]gnata.CustomFunc, w io.Writer) map[string]gnata.CustomFunc {
	var mu sync.Mutex
	enc := json.NewEncoder(w)

	out := make(map[string]gnata.CustomFunc, len(funcs))
	for name, fn := range funcs {
		name, fn := name, fn
		out[name] = func(args []any, focus any) (any, error) {
			start := time.Now()
			v, err := fn(args, focus)
			dur := time.Since(start)

			rec := AuditRecord{
				Time:     start.UTC().Format(time.RFC3339Nano),
				Func:     name,
				Args:     args,
				Duration: dur.String(),
			}
			if err != nil {
				rec.Error = err.Error()
			}

			mu.Lock()
			_ = enc.Encode(rec)
			mu.Unlock()

			return v, err
		}
	}
	return out
}
