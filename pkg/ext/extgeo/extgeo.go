// Package extgeo provides geospatial utility functions for gnata.
// All calculations use the WGS-84 mean Earth radius (6371 km) and operate
// on decimal degrees. No external dependencies are required.
//
// Functions
//
//   - $haversine(lat1, lon1, lat2, lon2)               – great-circle distance in km
//   - $bearing(lat1, lon1, lat2, lon2)                 – initial bearing in degrees (0–360)
//   - $geoFormat(lat, lon [, format])                  – format as "decimal" or "dms" string
//   - $geoParse(str)                                   – parse "lat, lon" string → {lat, lon}
//   - $inBoundingBox(lat, lon, minLat, minLon, maxLat, maxLon) – point-in-bbox test
//   - $geoDistance(point, points)                      – distances from point to each in points
package extgeo

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/recolabs/gnata"
	"github.com/sandrolain/gnata-ext/pkg/ext/extutil"
)

const earthRadiusKm = 6371.0

// All returns a map of all geospatial functions.
func All() map[string]gnata.CustomFunc {
	return map[string]gnata.CustomFunc{
		"haversine":    Haversine(),
		"bearing":      Bearing(),
		"geoFormat":    GeoFormat(),
		"geoParse":     GeoParse(),
		"inBoundingBox": InBoundingBox(),
		"geoDistance":  GeoDistance(),
	}
}

func degToRad(d float64) float64 { return d * math.Pi / 180.0 }
func radToDeg(r float64) float64 { return r * 180.0 / math.Pi }

// haversineCalc computes the great-circle distance between two points in km.
func haversineCalc(lat1, lon1, lat2, lon2 float64) float64 {
	φ1 := degToRad(lat1)
	φ2 := degToRad(lat2)
	Δφ := degToRad(lat2 - lat1)
	Δλ := degToRad(lon2 - lon1)
	a := math.Sin(Δφ/2)*math.Sin(Δφ/2) +
		math.Cos(φ1)*math.Cos(φ2)*math.Sin(Δλ/2)*math.Sin(Δλ/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusKm * c
}

// Haversine returns the CustomFunc for $haversine(lat1, lon1, lat2, lon2).
func Haversine() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("$haversine: requires 4 arguments (lat1, lon1, lat2, lon2)")
		}
		vals, err := extutil.ToFloatSlice(args[:4])
		if err != nil {
			return nil, fmt.Errorf("$haversine: %w", err)
		}
		return haversineCalc(vals[0], vals[1], vals[2], vals[3]), nil
	}
}

// Bearing returns the CustomFunc for $bearing(lat1, lon1, lat2, lon2).
// Returns the initial bearing in degrees (0–360, clockwise from north).
func Bearing() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("$bearing: requires 4 arguments (lat1, lon1, lat2, lon2)")
		}
		vals, err := extutil.ToFloatSlice(args[:4])
		if err != nil {
			return nil, fmt.Errorf("$bearing: %w", err)
		}
		φ1 := degToRad(vals[0])
		φ2 := degToRad(vals[2])
		Δλ := degToRad(vals[3] - vals[1])
		θ := math.Atan2(
			math.Sin(Δλ)*math.Cos(φ2),
			math.Cos(φ1)*math.Sin(φ2)-math.Sin(φ1)*math.Cos(φ2)*math.Cos(Δλ),
		)
		return math.Mod(radToDeg(θ)+360, 360), nil
	}
}

// decimalToDMS converts a decimal degree value to a DMS component tuple.
func decimalToDMS(deg float64) (int, int, float64) {
	d := math.Abs(deg)
	degrees := int(d)
	minutes := int((d - float64(degrees)) * 60)
	seconds := (d - float64(degrees) - float64(minutes)/60) * 3600
	return degrees, minutes, seconds
}

// GeoFormat returns the CustomFunc for $geoFormat(lat, lon [, format]).
// format: "decimal" (default) → "48.8566, 2.3522"; "dms" → degrees/minutes/seconds.
func GeoFormat() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$geoFormat: requires at least 2 arguments (lat, lon)")
		}
		lat, err := extutil.ToFloat(args[0])
		if err != nil {
			return nil, fmt.Errorf("$geoFormat: lat: %w", err)
		}
		lon, err := extutil.ToFloat(args[1])
		if err != nil {
			return nil, fmt.Errorf("$geoFormat: lon: %w", err)
		}
		format := "decimal"
		if len(args) >= 3 {
			if s, ok := args[2].(string); ok {
				format = s
			}
		}
		switch format {
		case "dms":
			latD, latM, latS := decimalToDMS(lat)
			lonD, lonM, lonS := decimalToDMS(lon)
			latDir := "N"
			if lat < 0 {
				latDir = "S"
			}
			lonDir := "E"
			if lon < 0 {
				lonDir = "W"
			}
			return fmt.Sprintf("%d°%d'%.2f\"%s %d°%d'%.2f\"%s",
				latD, latM, latS, latDir, lonD, lonM, lonS, lonDir), nil
		default: // "decimal"
			return fmt.Sprintf("%.4f, %.4f", lat, lon), nil
		}
	}
}

// GeoParse returns the CustomFunc for $geoParse(str).
// Parses a "lat, lon" decimal string and returns {lat: float, lon: float}.
func GeoParse() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$geoParse: requires 1 argument (str)")
		}
		s, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$geoParse: argument must be a string")
		}
		parts := strings.SplitN(s, ",", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("$geoParse: expected 'lat, lon' format, got %q", s)
		}
		lat, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		if err != nil {
			return nil, fmt.Errorf("$geoParse: invalid latitude: %w", err)
		}
		lon, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		if err != nil {
			return nil, fmt.Errorf("$geoParse: invalid longitude: %w", err)
		}
		return map[string]any{"lat": lat, "lon": lon}, nil
	}
}

// InBoundingBox returns the CustomFunc for $inBoundingBox(lat, lon, minLat, minLon, maxLat, maxLon).
func InBoundingBox() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 6 {
			return nil, fmt.Errorf("$inBoundingBox: requires 6 arguments (lat, lon, minLat, minLon, maxLat, maxLon)")
		}
		vals, err := extutil.ToFloatSlice(args[:6])
		if err != nil {
			return nil, fmt.Errorf("$inBoundingBox: %w", err)
		}
		lat, lon := vals[0], vals[1]
		minLat, minLon, maxLat, maxLon := vals[2], vals[3], vals[4], vals[5]
		return lat >= minLat && lat <= maxLat && lon >= minLon && lon <= maxLon, nil
	}
}

// GeoDistance returns the CustomFunc for $geoDistance(point, points).
// point: {lat, lon}; points: array of {lat, lon}.
// Returns an array of distances in km.
func GeoDistance() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$geoDistance: requires 2 arguments (point, points)")
		}
		origin, err := extutil.ToObject(args[0])
		if err != nil {
			return nil, fmt.Errorf("$geoDistance: point: %w", err)
		}
		lat1, err := extutil.ToFloat(origin["lat"])
		if err != nil {
			return nil, fmt.Errorf("$geoDistance: point.lat: %w", err)
		}
		lon1, err := extutil.ToFloat(origin["lon"])
		if err != nil {
			return nil, fmt.Errorf("$geoDistance: point.lon: %w", err)
		}
		points, err := extutil.ToArray(args[1])
		if err != nil {
			return nil, fmt.Errorf("$geoDistance: points: %w", err)
		}
		result := make([]any, len(points))
		for i, p := range points {
			pm, err := extutil.ToObject(p)
			if err != nil {
				return nil, fmt.Errorf("$geoDistance: points[%d]: %w", i, err)
			}
			lat2, err := extutil.ToFloat(pm["lat"])
			if err != nil {
				return nil, fmt.Errorf("$geoDistance: points[%d].lat: %w", i, err)
			}
			lon2, err := extutil.ToFloat(pm["lon"])
			if err != nil {
				return nil, fmt.Errorf("$geoDistance: points[%d].lon: %w", i, err)
			}
			result[i] = haversineCalc(lat1, lon1, lat2, lon2)
		}
		return result, nil
	}
}
