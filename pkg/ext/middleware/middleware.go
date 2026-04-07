// Package middleware provides wrappers that add cross-cutting behaviour
// (logging, memoization, rate limiting) to maps of gnata custom functions.
//
// Each wrapper returns a new map; the original map is left unchanged.
//
//	funcs := ext.AllFuncs()
//	funcs = middleware.WithLogging(funcs, slog.Default())
//	funcs = middleware.WithMemoize(funcs, "uuid")
//	env := gnata.NewCustomEnv(funcs)
package middleware

import (
	"fmt"
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
