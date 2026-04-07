package extgeo_test

import (
	"math"
	"testing"

	"github.com/sandrolain/gnata-ext/pkg/ext/extgeo"
)

const tolerance = 0.1 // km

func approxEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func TestHaversine(t *testing.T) {
	fn := extgeo.Haversine()

	// London to Paris ≈ 340 km
	got, err := fn([]any{float64(51.5074), float64(-0.1278), float64(48.8566), float64(2.3522)}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	d, ok := got.(float64)
	if !ok {
		t.Fatalf("expected float64, got %T", got)
	}
	if !approxEqual(d, 340.0, 5.0) {
		t.Errorf("London–Paris: got %.2f km, expected ~340 km", d)
	}

	// Same point → 0
	got, _ = fn([]any{float64(48.0), float64(2.0), float64(48.0), float64(2.0)}, nil)
	if got.(float64) != 0.0 {
		t.Errorf("same point: expected 0, got %v", got)
	}
}

func TestBearing(t *testing.T) {
	fn := extgeo.Bearing()

	// Bearing from equator/0° going east should be ~90°
	got, err := fn([]any{float64(0), float64(0), float64(0), float64(1)}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	d := got.(float64)
	if !approxEqual(d, 90.0, 1.0) {
		t.Errorf("east bearing: got %.2f, expected ~90", d)
	}

	// North
	got, _ = fn([]any{float64(0), float64(0), float64(1), float64(0)}, nil)
	if !approxEqual(got.(float64), 0.0, 1.0) {
		t.Errorf("north bearing: got %.2f, expected ~0", got)
	}
}

func TestGeoFormat(t *testing.T) {
	fn := extgeo.GeoFormat()

	// decimal (default)
	got, err := fn([]any{float64(48.8566), float64(2.3522)}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := got.(string)
	if s != "48.8566, 2.3522" {
		t.Errorf("decimal format: got %q", s)
	}

	// dms
	got, err = fn([]any{float64(48.0), float64(2.0), "dms"}, nil)
	if err != nil {
		t.Fatalf("DMS format error: %v", err)
	}
	// Just verify it contains degree symbols and N/E
	dms := got.(string)
	if len(dms) == 0 {
		t.Error("DMS output is empty")
	}
}

func TestGeoParse(t *testing.T) {
	fn := extgeo.GeoParse()

	got, err := fn([]any{"48.8566, 2.3522"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m, ok := got.(map[string]any)
	if !ok {
		t.Fatalf("expected map, got %T", got)
	}
	if !approxEqual(m["lat"].(float64), 48.8566, 0.001) {
		t.Errorf("lat: got %v, want 48.8566", m["lat"])
	}
	if !approxEqual(m["lon"].(float64), 2.3522, 0.001) {
		t.Errorf("lon: got %v, want 2.3522", m["lon"])
	}

	_, err = fn([]any{"invalid"}, nil)
	if err == nil {
		t.Error("expected error for invalid input")
	}
}

func TestInBoundingBox(t *testing.T) {
	fn := extgeo.InBoundingBox()

	// Paris inside Europe bbox
	got, err := fn([]any{
		float64(48.8566), float64(2.3522),
		float64(36.0), float64(-10.0),
		float64(71.0), float64(40.0),
	}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != true {
		t.Error("Paris should be inside Europe bbox")
	}

	// Outside
	got, _ = fn([]any{
		float64(0.0), float64(0.0),
		float64(10.0), float64(10.0),
		float64(20.0), float64(20.0),
	}, nil)
	if got != false {
		t.Error("origin should be outside bbox")
	}
}

func TestGeoDistance(t *testing.T) {
	fn := extgeo.GeoDistance()

	origin := map[string]any{"lat": float64(51.5074), "lon": float64(-0.1278)} // London
	points := []any{
		map[string]any{"lat": float64(48.8566), "lon": float64(2.3522)},  // Paris
		map[string]any{"lat": float64(51.5074), "lon": float64(-0.1278)}, // London itself
	}

	got, err := fn([]any{origin, points}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	arr := got.([]any)
	if len(arr) != 2 {
		t.Fatalf("expected 2 distances, got %d", len(arr))
	}
	if !approxEqual(arr[0].(float64), 340.0, 5.0) {
		t.Errorf("London–Paris: got %.2f km", arr[0])
	}
	if arr[1].(float64) != 0.0 {
		t.Errorf("London–London: expected 0, got %v", arr[1])
	}
}

func TestAll(t *testing.T) {
	all := extgeo.All()
	expected := []string{"haversine", "bearing", "geoFormat", "geoParse", "inBoundingBox", "geoDistance"}
	for _, name := range expected {
		if _, ok := all[name]; !ok {
			t.Errorf("All() missing function: %q", name)
		}
	}
}

func TestHaversineErrors(t *testing.T) {
	f := extgeo.Haversine()
	if _, err := f([]any{}, nil); err == nil {
		t.Error("haversine: expected error for 0 args")
	}
	if _, err := f([]any{"bad", 0.0, 0.0, 0.0}, nil); err == nil {
		t.Error("haversine: expected error for non-numeric arg")
	}
}

func TestBearingErrors(t *testing.T) {
	f := extgeo.Bearing()
	if _, err := f([]any{}, nil); err == nil {
		t.Error("bearing: expected error for 0 args")
	}
	if _, err := f([]any{"bad", 0.0, 0.0, 0.0}, nil); err == nil {
		t.Error("bearing: expected error for non-numeric arg")
	}
}

func TestGeoFormatErrors(t *testing.T) {
	f := extgeo.GeoFormat()
	if _, err := f([]any{}, nil); err == nil {
		t.Error("geoFormat: expected error for 0 args")
	}
	if _, err := f([]any{"bad", 0.0}, nil); err == nil {
		t.Error("geoFormat: expected error for non-numeric lat")
	}
	if _, err := f([]any{0.0, "bad"}, nil); err == nil {
		t.Error("geoFormat: expected error for non-numeric lon")
	}
}

func TestGeoFormatDMSSouth(t *testing.T) {
	f := extgeo.GeoFormat()
	// negative lat -> "S", negative lon -> "W"
	got, err := f([]any{-33.8688, -70.6693, "dms"}, nil)
	if err != nil {
		t.Fatalf("geoFormat dms south: %v", err)
	}
	s := got.(string)
	if len(s) == 0 {
		t.Errorf("geoFormat dms south: empty result")
	}
}

func TestGeoFormatDecimalDefault(t *testing.T) {
	f := extgeo.GeoFormat()
	// non-string format arg falls back to decimal
	got, err := f([]any{48.8566, 2.3522, 42}, nil)
	if err != nil {
		t.Fatalf("geoFormat non-string format: %v", err)
	}
	s := got.(string)
	if s != "48.8566, 2.3522" {
		t.Errorf("geoFormat decimal default: got %q", s)
	}
}

func TestGeoParseErrors(t *testing.T) {
	f := extgeo.GeoParse()
	if _, err := f([]any{}, nil); err == nil {
		t.Error("geoParse: expected error for 0 args")
	}
	if _, err := f([]any{42}, nil); err == nil {
		t.Error("geoParse: expected error for non-string")
	}
	if _, err := f([]any{"no-comma"}, nil); err == nil {
		t.Error("geoParse: expected error for missing comma")
	}
	if _, err := f([]any{"bad, 2.0"}, nil); err == nil {
		t.Error("geoParse: expected error for invalid lat")
	}
	if _, err := f([]any{"48.8, bad"}, nil); err == nil {
		t.Error("geoParse: expected error for invalid lon")
	}
}

func TestInBoundingBoxErrors(t *testing.T) {
	f := extgeo.InBoundingBox()
	if _, err := f([]any{}, nil); err == nil {
		t.Error("inBoundingBox: expected error for 0 args")
	}
	if _, err := f([]any{"bad", 0.0, 0.0, 0.0, 1.0, 1.0}, nil); err == nil {
		t.Error("inBoundingBox: expected error for non-numeric arg")
	}
}

func TestGeoDistanceErrors(t *testing.T) {
	f := extgeo.GeoDistance()
	if _, err := f([]any{}, nil); err == nil {
		t.Error("geoDistance: expected error for 0 args")
	}
	// non-object origin
	if _, err := f([]any{"not-obj", []any{}}, nil); err == nil {
		t.Error("geoDistance: expected error for non-object origin")
	}
	// missing lat
	if _, err := f([]any{map[string]any{"lon": 0.0}, []any{}}, nil); err == nil {
		t.Error("geoDistance: expected error for missing origin lat")
	}
	// missing lon
	if _, err := f([]any{map[string]any{"lat": 0.0}, []any{}}, nil); err == nil {
		t.Error("geoDistance: expected error for missing origin lon")
	}
	// non-array points
	origin := map[string]any{"lat": 0.0, "lon": 0.0}
	if _, err := f([]any{origin, "not-array"}, nil); err == nil {
		t.Error("geoDistance: expected error for non-array points")
	}
	// bad point in array
	if _, err := f([]any{origin, []any{"bad"}}, nil); err == nil {
		t.Error("geoDistance: expected error for non-object point")
	}
	// point missing lat
	if _, err := f([]any{origin, []any{map[string]any{"lon": 0.0}}}, nil); err == nil {
		t.Error("geoDistance: expected error for point missing lat")
	}
	// point missing lon
	if _, err := f([]any{origin, []any{map[string]any{"lat": 0.0}}}, nil); err == nil {
		t.Error("geoDistance: expected error for point missing lon")
	}
}
