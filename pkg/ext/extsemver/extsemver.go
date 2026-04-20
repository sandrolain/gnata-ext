// Package extsemver provides semantic versioning functions for gnata.
//
// Functions:
//
//   - $semverParse(v)           – parse version string into object
//   - $semverCompare(a, b)      – compare two version strings (-1, 0, 1)
//   - $semverSatisfies(v, c)    – true if version satisfies constraint
//   - $semverBump(v, part)      – bump major / minor / patch
//   - $semverSort(arr)          – sort array of version strings ascending
//   - $semverMax(arr)           – return highest version string
//   - $semverMin(arr)           – return lowest version string
package extsemver

import (
	"fmt"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/recolabs/gnata"
)

// All returns a map of all extsemver functions.
func All() map[string]gnata.CustomFunc {
	return map[string]gnata.CustomFunc{
		"semverParse":     SemverParse(),
		"semverCompare":   SemverCompare(),
		"semverSatisfies": SemverSatisfies(),
		"semverBump":      SemverBump(),
		"semverSort":      SemverSort(),
		"semverMax":       SemverMax(),
		"semverMin":       SemverMin(),
	}
}

// SemverParse returns the CustomFunc for $semverParse(v).
// Returns an object {major, minor, patch, prerelease, metadata, original}.
func SemverParse() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		s, err := requireVersionString(args, "$semverParse")
		if err != nil {
			return nil, err
		}
		v, err := semver.NewVersion(s)
		if err != nil {
			return nil, fmt.Errorf("$semverParse: %w", err)
		}
		return map[string]any{
			"major":      float64(v.Major()),
			"minor":      float64(v.Minor()),
			"patch":      float64(v.Patch()),
			"prerelease": v.Prerelease(),
			"metadata":   v.Metadata(),
			"original":   v.Original(),
		}, nil
	}
}

// SemverCompare returns the CustomFunc for $semverCompare(a, b).
// Returns -1 if a < b, 0 if equal, 1 if a > b.
func SemverCompare() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$semverCompare: requires 2 arguments (a, b)")
		}
		sa, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$semverCompare: first argument must be a string")
		}
		sb, ok := args[1].(string)
		if !ok {
			return nil, fmt.Errorf("$semverCompare: second argument must be a string")
		}
		va, err := semver.NewVersion(sa)
		if err != nil {
			return nil, fmt.Errorf("$semverCompare: invalid version a: %w", err)
		}
		vb, err := semver.NewVersion(sb)
		if err != nil {
			return nil, fmt.Errorf("$semverCompare: invalid version b: %w", err)
		}
		return float64(va.Compare(vb)), nil
	}
}

// SemverSatisfies returns the CustomFunc for $semverSatisfies(v, constraint).
// Returns true if v satisfies the constraint (e.g. ">=1.0.0 <2.0.0").
func SemverSatisfies() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$semverSatisfies: requires 2 arguments (v, constraint)")
		}
		sv, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$semverSatisfies: first argument must be a string")
		}
		sc, ok := args[1].(string)
		if !ok {
			return nil, fmt.Errorf("$semverSatisfies: second argument must be a string")
		}
		v, err := semver.NewVersion(sv)
		if err != nil {
			return nil, fmt.Errorf("$semverSatisfies: invalid version: %w", err)
		}
		c, err := semver.NewConstraint(sc)
		if err != nil {
			return nil, fmt.Errorf("$semverSatisfies: invalid constraint: %w", err)
		}
		return c.Check(v), nil
	}
}

// SemverBump returns the CustomFunc for $semverBump(v, part).
// part must be "major", "minor", or "patch".
func SemverBump() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$semverBump: requires 2 arguments (v, part)")
		}
		sv, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$semverBump: first argument must be a string")
		}
		part, ok := args[1].(string)
		if !ok {
			return nil, fmt.Errorf("$semverBump: second argument must be a string")
		}
		v, err := semver.NewVersion(sv)
		if err != nil {
			return nil, fmt.Errorf("$semverBump: invalid version: %w", err)
		}
		var bumped semver.Version
		switch part {
		case "major":
			bumped = v.IncMajor()
		case "minor":
			bumped = v.IncMinor()
		case "patch":
			bumped = v.IncPatch()
		default:
			return nil, fmt.Errorf("$semverBump: part must be 'major', 'minor', or 'patch'")
		}
		return bumped.Original(), nil
	}
}

// SemverSort returns the CustomFunc for $semverSort(arr).
// Sorts an array of version strings ascending.
func SemverSort() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		arr, err := requireVersionArray(args, "$semverSort")
		if err != nil {
			return nil, err
		}
		versions := make([]*semver.Version, 0, len(arr))
		originals := make([]string, 0, len(arr))
		for _, raw := range arr {
			s, ok := raw.(string)
			if !ok {
				return nil, fmt.Errorf("$semverSort: all elements must be version strings")
			}
			v, err := semver.NewVersion(s)
			if err != nil {
				return nil, fmt.Errorf("$semverSort: invalid version %q: %w", s, err)
			}
			versions = append(versions, v)
			originals = append(originals, s)
		}
		// Sort by version; keep original string.
		type pair struct {
			v   *semver.Version
			raw string
		}
		pairs := make([]pair, len(versions))
		for i := range versions {
			pairs[i] = pair{versions[i], originals[i]}
		}
		sort.Slice(pairs, func(i, j int) bool {
			return pairs[i].v.LessThan(pairs[j].v)
		})
		out := make([]any, len(pairs))
		for i, p := range pairs {
			out[i] = p.raw
		}
		return out, nil
	}
}

// SemverMax returns the CustomFunc for $semverMax(arr).
// Returns the highest version string in the array.
func SemverMax() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		arr, err := requireVersionArray(args, "$semverMax")
		if err != nil {
			return nil, err
		}
		if len(arr) == 0 {
			return nil, fmt.Errorf("$semverMax: array must not be empty")
		}
		var best *semver.Version
		var bestRaw string
		for _, raw := range arr {
			s, ok := raw.(string)
			if !ok {
				return nil, fmt.Errorf("$semverMax: all elements must be version strings")
			}
			v, err := semver.NewVersion(s)
			if err != nil {
				return nil, fmt.Errorf("$semverMax: invalid version %q: %w", s, err)
			}
			if best == nil || v.GreaterThan(best) {
				best = v
				bestRaw = s
			}
		}
		return bestRaw, nil
	}
}

// SemverMin returns the CustomFunc for $semverMin(arr).
// Returns the lowest version string in the array.
func SemverMin() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		arr, err := requireVersionArray(args, "$semverMin")
		if err != nil {
			return nil, err
		}
		if len(arr) == 0 {
			return nil, fmt.Errorf("$semverMin: array must not be empty")
		}
		var best *semver.Version
		var bestRaw string
		for _, raw := range arr {
			s, ok := raw.(string)
			if !ok {
				return nil, fmt.Errorf("$semverMin: all elements must be version strings")
			}
			v, err := semver.NewVersion(s)
			if err != nil {
				return nil, fmt.Errorf("$semverMin: invalid version %q: %w", s, err)
			}
			if best == nil || v.LessThan(best) {
				best = v
				bestRaw = s
			}
		}
		return bestRaw, nil
	}
}

// --- helpers ---

func requireVersionString(args []any, name string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("%s: requires 1 argument", name)
	}
	s, ok := args[0].(string)
	if !ok {
		return "", fmt.Errorf("%s: argument must be a string", name)
	}
	return s, nil
}

func requireVersionArray(args []any, name string) ([]any, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("%s: requires 1 argument", name)
	}
	arr, ok := args[0].([]any)
	if !ok {
		return nil, fmt.Errorf("%s: argument must be an array", name)
	}
	return arr, nil
}
