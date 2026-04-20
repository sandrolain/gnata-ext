// Package exturi provides URI parsing and construction functions for gnata.
//
// Functions:
//
//   - $uriParse(url)              – decompose URL into parts object
//   - $uriBuild(parts)            – build URL from parts object
//   - $uriJoin(base, ref)         – resolve a relative URL against a base
//   - $queryParse(qs)             – parse a query string into object
//   - $queryBuild(obj)            – serialize object to query string
//   - $uriGetPath(url)            – extract only the path component
//   - $uriGetQuery(url)           – extract and parse the query component
//   - $uriSetQuery(url, params)   – replace query string keeping the rest
package exturi

import (
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/recolabs/gnata"
)

// All returns a map of all exturi functions.
func All() map[string]gnata.CustomFunc {
	return map[string]gnata.CustomFunc{
		"uriParse":    URIParse(),
		"uriBuild":    URIBuild(),
		"uriJoin":     URIJoin(),
		"queryParse":  QueryParse(),
		"queryBuild":  QueryBuild(),
		"uriGetPath":  URIGetPath(),
		"uriGetQuery": URIGetQuery(),
		"uriSetQuery": URISetQuery(),
	}
}

// URIParse returns the CustomFunc for $uriParse(url).
// Returns an object with scheme, user, host, port, path, query (object), fragment.
func URIParse() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$uriParse: requires 1 argument")
		}
		raw, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$uriParse: argument must be a string")
		}
		u, err := url.Parse(raw)
		if err != nil {
			return nil, fmt.Errorf("$uriParse: %w", err)
		}
		queryObj := queryToObject(u.Query())
		user := ""
		password := ""
		if u.User != nil {
			user = u.User.Username()
			password, _ = u.User.Password()
		}
		return map[string]any{
			"scheme":   u.Scheme,
			"user":     user,
			"password": password,
			"host":     u.Hostname(),
			"port":     u.Port(),
			"path":     u.Path,
			"query":    queryObj,
			"fragment": u.Fragment,
		}, nil
	}
}

// URIBuild returns the CustomFunc for $uriBuild(parts).
// parts must be an object with optional keys: scheme, user, password, host, port, path, query, fragment.
func URIBuild() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$uriBuild: requires 1 argument")
		}
		parts, ok := args[0].(map[string]any)
		if !ok {
			return nil, fmt.Errorf("$uriBuild: argument must be an object")
		}
		u := &url.URL{}
		u.Scheme = stringField(parts, "scheme")
		u.Path = stringField(parts, "path")
		u.Fragment = stringField(parts, "fragment")

		host := stringField(parts, "host")
		port := stringField(parts, "port")
		if port != "" {
			u.Host = host + ":" + port
		} else {
			u.Host = host
		}

		username := stringField(parts, "user")
		if username != "" {
			password := stringField(parts, "password")
			if password != "" {
				u.User = url.UserPassword(username, password)
			} else {
				u.User = url.User(username)
			}
		}

		if q, exists := parts["query"]; exists && q != nil {
			if qMap, ok := q.(map[string]any); ok {
				vals := objectToQuery(qMap)
				u.RawQuery = vals.Encode()
			}
		}

		return u.String(), nil
	}
}

// URIJoin returns the CustomFunc for $uriJoin(base, ref).
// Resolves ref relative to base.
func URIJoin() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$uriJoin: requires 2 arguments (base, ref)")
		}
		baseStr, ok1 := args[0].(string)
		refStr, ok2 := args[1].(string)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("$uriJoin: arguments must be strings")
		}
		base, err := url.Parse(baseStr)
		if err != nil {
			return nil, fmt.Errorf("$uriJoin: invalid base: %w", err)
		}
		ref, err := url.Parse(refStr)
		if err != nil {
			return nil, fmt.Errorf("$uriJoin: invalid ref: %w", err)
		}
		return base.ResolveReference(ref).String(), nil
	}
}

// QueryParse returns the CustomFunc for $queryParse(qs).
// Parses a query string (with or without leading ?) into an object.
// Values that appear once are strings; values that appear multiple times become arrays.
func QueryParse() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$queryParse: requires 1 argument")
		}
		qs, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$queryParse: argument must be a string")
		}
		qs = strings.TrimPrefix(qs, "?")
		vals, err := url.ParseQuery(qs)
		if err != nil {
			return nil, fmt.Errorf("$queryParse: %w", err)
		}
		return queryToObject(vals), nil
	}
}

// QueryBuild returns the CustomFunc for $queryBuild(obj).
// Serializes an object into a URL-encoded query string (without leading ?).
func QueryBuild() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$queryBuild: requires 1 argument")
		}
		obj, ok := args[0].(map[string]any)
		if !ok {
			return nil, fmt.Errorf("$queryBuild: argument must be an object")
		}
		return objectToQuery(obj).Encode(), nil
	}
}

// URIGetPath returns the CustomFunc for $uriGetPath(url).
// Returns just the path component of the URL.
func URIGetPath() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$uriGetPath: requires 1 argument")
		}
		raw, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$uriGetPath: argument must be a string")
		}
		u, err := url.Parse(raw)
		if err != nil {
			return nil, fmt.Errorf("$uriGetPath: %w", err)
		}
		return u.Path, nil
	}
}

// URIGetQuery returns the CustomFunc for $uriGetQuery(url).
// Returns the parsed query string as an object.
func URIGetQuery() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$uriGetQuery: requires 1 argument")
		}
		raw, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$uriGetQuery: argument must be a string")
		}
		u, err := url.Parse(raw)
		if err != nil {
			return nil, fmt.Errorf("$uriGetQuery: %w", err)
		}
		return queryToObject(u.Query()), nil
	}
}

// URISetQuery returns the CustomFunc for $uriSetQuery(url, params).
// Replaces the query string of url with the serialized params object.
func URISetQuery() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$uriSetQuery: requires 2 arguments (url, params)")
		}
		raw, ok1 := args[0].(string)
		params, ok2 := args[1].(map[string]any)
		if !ok1 {
			return nil, fmt.Errorf("$uriSetQuery: first argument must be a string")
		}
		if !ok2 {
			return nil, fmt.Errorf("$uriSetQuery: second argument must be an object")
		}
		u, err := url.Parse(raw)
		if err != nil {
			return nil, fmt.Errorf("$uriSetQuery: %w", err)
		}
		u.RawQuery = objectToQuery(params).Encode()
		return u.String(), nil
	}
}

// --- helpers ---

func stringField(m map[string]any, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// queryToObject converts url.Values to map[string]any.
// Single-value keys become strings; multi-value keys become []any.
func queryToObject(vals url.Values) map[string]any {
	obj := make(map[string]any, len(vals))
	for k, vs := range vals {
		if len(vs) == 1 {
			obj[k] = vs[0]
		} else {
			arr := make([]any, len(vs))
			for i, v := range vs {
				arr[i] = v
			}
			obj[k] = arr
		}
	}
	return obj
}

// objectToQuery converts map[string]any to url.Values.
// Array values produce multiple entries for the same key.
func objectToQuery(obj map[string]any) url.Values {
	vals := make(url.Values)
	// Sort keys for deterministic output.
	keys := make([]string, 0, len(obj))
	for k := range obj {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := obj[k]
		switch val := v.(type) {
		case []any:
			for _, elem := range val {
				vals.Add(k, fmt.Sprintf("%v", elem))
			}
		default:
			vals.Set(k, fmt.Sprintf("%v", val))
		}
	}
	return vals
}
