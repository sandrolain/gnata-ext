// Package extcolor provides color parsing and manipulation functions for gnata.
//
// All color representations use float64 components (0–255 for RGB,
// 0–360 for hue, 0–100 for saturation/lightness/alpha expressed as 0–1).
//
// Functions:
//
//   - $colorParse(s)              – parse CSS color string → {r, g, b, a}
//   - $colorToHex(obj)            – {r,g,b} → "#rrggbb"
//   - $colorToRGB(obj)            – {r,g,b,a?} → "rgb(r,g,b)" or "rgba(r,g,b,a)"
//   - $colorToHSL(obj)            – {r,g,b} → {h, s, l}
//   - $colorMix(a, b, t)          – linear blend at t ∈ [0,1] → {r,g,b,a}
//   - $colorLighten(obj, amount)  – increase lightness by amount (0–1) → {r,g,b,a}
//   - $colorDarken(obj, amount)   – decrease lightness by amount (0–1) → {r,g,b,a}
//   - $colorContrast(fg, bg)      – WCAG contrast ratio (number)
//   - $colorLuminance(obj)        – relative luminance (number, 0–1)
package extcolor

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/recolabs/gnata"
	"github.com/sandrolain/gnata-ext/pkg/ext/extutil"
)

// All returns a map of all extcolor functions.
func All() map[string]gnata.CustomFunc {
	return map[string]gnata.CustomFunc{
		"colorParse":     ColorParse(),
		"colorToHex":     ColorToHex(),
		"colorToRGB":     ColorToRGB(),
		"colorToHSL":     ColorToHSL(),
		"colorMix":       ColorMix(),
		"colorLighten":   ColorLighten(),
		"colorDarken":    ColorDarken(),
		"colorContrast":  ColorContrast(),
		"colorLuminance": ColorLuminance(),
	}
}

// colorRGBA is an internal representation with float64 components.
type colorRGBA struct {
	R, G, B, A float64
}

// --- Public CustomFunc constructors ---

// ColorParse returns the CustomFunc for $colorParse(s).
// Supported formats: #rgb, #rrggbb, rgb(r,g,b), rgba(r,g,b,a).
func ColorParse() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		s, err := requireString1(args, "$colorParse")
		if err != nil {
			return nil, err
		}
		c, err := parseColor(s)
		if err != nil {
			return nil, fmt.Errorf("$colorParse: %w", err)
		}
		return colorToMap(c), nil
	}
}

// ColorToHex returns the CustomFunc for $colorToHex(obj).
func ColorToHex() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		c, err := requireColorMap(args, "$colorToHex")
		if err != nil {
			return nil, err
		}
		r := clampByte(c.R)
		g := clampByte(c.G)
		b := clampByte(c.B)
		return fmt.Sprintf("#%02x%02x%02x", r, g, b), nil
	}
}

// ColorToRGB returns the CustomFunc for $colorToRGB(obj).
func ColorToRGB() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		c, err := requireColorMap(args, "$colorToRGB")
		if err != nil {
			return nil, err
		}
		r := clampByte(c.R)
		g := clampByte(c.G)
		b := clampByte(c.B)
		if c.A != 1.0 {
			a := math.Round(c.A*100) / 100
			return fmt.Sprintf("rgba(%d,%d,%d,%g)", r, g, b, a), nil
		}
		return fmt.Sprintf("rgb(%d,%d,%d)", r, g, b), nil
	}
}

// ColorToHSL returns the CustomFunc for $colorToHSL(obj).
// Returns {h (0-360), s (0-100), l (0-100)}.
func ColorToHSL() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		c, err := requireColorMap(args, "$colorToHSL")
		if err != nil {
			return nil, err
		}
		h, s, l := rgbToHSL(c.R/255, c.G/255, c.B/255)
		return map[string]any{
			"h": math.Round(h*10) / 10,
			"s": math.Round(s*1000) / 10,
			"l": math.Round(l*1000) / 10,
		}, nil
	}
}

// ColorMix returns the CustomFunc for $colorMix(a, b, t).
// t is the blend factor in [0,1]. Returns {r,g,b,a}.
func ColorMix() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("$colorMix: requires 3 arguments (a, b, t)")
		}
		ca, err := requireColorMap(args[:1], "$colorMix: first argument")
		if err != nil {
			return nil, err
		}
		cb, err := requireColorMap(args[1:2], "$colorMix: second argument")
		if err != nil {
			return nil, err
		}
		t, err := extutil.ToFloat(args[2])
		if err != nil {
			return nil, fmt.Errorf("$colorMix: t must be a number")
		}
		t = math.Max(0, math.Min(1, t))
		mixed := colorRGBA{
			R: lerp(ca.R, cb.R, t),
			G: lerp(ca.G, cb.G, t),
			B: lerp(ca.B, cb.B, t),
			A: lerp(ca.A, cb.A, t),
		}
		return colorToMap(mixed), nil
	}
}

// ColorLighten returns the CustomFunc for $colorLighten(obj, amount).
// Increases HSL lightness by amount (0–1).
func ColorLighten() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		return adjustLightness(args, "$colorLighten", 1)
	}
}

// ColorDarken returns the CustomFunc for $colorDarken(obj, amount).
// Decreases HSL lightness by amount (0–1).
func ColorDarken() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		return adjustLightness(args, "$colorDarken", -1)
	}
}

// ColorContrast returns the CustomFunc for $colorContrast(fg, bg).
// Returns the WCAG contrast ratio between fg and bg.
func ColorContrast() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$colorContrast: requires 2 arguments (fg, bg)")
		}
		fg, err := requireColorMap(args[:1], "$colorContrast: first argument")
		if err != nil {
			return nil, err
		}
		bg, err := requireColorMap(args[1:2], "$colorContrast: second argument")
		if err != nil {
			return nil, err
		}
		l1 := relativeLuminance(fg)
		l2 := relativeLuminance(bg)
		if l1 < l2 {
			l1, l2 = l2, l1
		}
		ratio := (l1 + 0.05) / (l2 + 0.05)
		return math.Round(ratio*100) / 100, nil
	}
}

// ColorLuminance returns the CustomFunc for $colorLuminance(obj).
// Returns the relative luminance (0–1) per WCAG 2.1.
func ColorLuminance() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		c, err := requireColorMap(args, "$colorLuminance")
		if err != nil {
			return nil, err
		}
		lum := relativeLuminance(c)
		return math.Round(lum*1e6) / 1e6, nil
	}
}

// --- helpers ---

func adjustLightness(args []any, name string, sign float64) (any, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("%s: requires 2 arguments (obj, amount)", name)
	}
	c, err := requireColorMap(args[:1], name)
	if err != nil {
		return nil, err
	}
	amount, err := extutil.ToFloat(args[1])
	if err != nil {
		return nil, fmt.Errorf("%s: amount must be a number", name)
	}
	h, s, l := rgbToHSL(c.R/255, c.G/255, c.B/255)
	l = math.Max(0, math.Min(1, l+sign*amount))
	r, g, b := hslToRGB(h, s, l)
	return colorToMap(colorRGBA{R: r * 255, G: g * 255, B: b * 255, A: c.A}), nil
}

func requireString1(args []any, name string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("%s: requires 1 argument", name)
	}
	s, ok := args[0].(string)
	if !ok {
		return "", fmt.Errorf("%s: argument must be a string", name)
	}
	return s, nil
}

func requireColorMap(args []any, name string) (colorRGBA, error) {
	if len(args) < 1 {
		return colorRGBA{}, fmt.Errorf("%s: requires 1 argument", name)
	}
	obj, ok := args[0].(map[string]any)
	if !ok {
		return colorRGBA{}, fmt.Errorf("%s: argument must be a color object {r,g,b}", name)
	}
	r, _ := extutil.ToFloat(obj["r"])
	g, _ := extutil.ToFloat(obj["g"])
	b, _ := extutil.ToFloat(obj["b"])
	a := 1.0
	if av, exists := obj["a"]; exists {
		if af, err := extutil.ToFloat(av); err == nil {
			a = af
		}
	}
	return colorRGBA{R: r, G: g, B: b, A: a}, nil
}

func colorToMap(c colorRGBA) map[string]any {
	return map[string]any{
		"r": math.Round(c.R*100) / 100,
		"g": math.Round(c.G*100) / 100,
		"b": math.Round(c.B*100) / 100,
		"a": c.A,
	}
}

func clampByte(v float64) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(math.Round(v))
}

func lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

// parseColor parses CSS color strings: #rgb, #rrggbb, rgb(...), rgba(...).
func parseColor(s string) (colorRGBA, error) {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "#") {
		return parseHex(s)
	}
	if strings.HasPrefix(s, "rgba(") || strings.HasPrefix(s, "rgb(") {
		return parseRGBFunc(s)
	}
	return colorRGBA{}, fmt.Errorf("unsupported color format %q", s)
}

func parseHex(s string) (colorRGBA, error) {
	s = strings.TrimPrefix(s, "#")
	switch len(s) {
	case 3:
		s = string([]byte{s[0], s[0], s[1], s[1], s[2], s[2]})
	case 6:
		// ok
	default:
		return colorRGBA{}, fmt.Errorf("invalid hex color length")
	}
	r, err1 := strconv.ParseUint(s[0:2], 16, 8)
	g, err2 := strconv.ParseUint(s[2:4], 16, 8)
	b, err3 := strconv.ParseUint(s[4:6], 16, 8)
	if err1 != nil || err2 != nil || err3 != nil {
		return colorRGBA{}, fmt.Errorf("invalid hex color %q", "#"+s)
	}
	return colorRGBA{R: float64(r), G: float64(g), B: float64(b), A: 1.0}, nil
}

func parseRGBFunc(s string) (colorRGBA, error) {
	s = strings.TrimSuffix(strings.TrimSpace(s), ")")
	s = strings.TrimPrefix(s, "rgba(")
	s = strings.TrimPrefix(s, "rgb(")
	parts := strings.Split(s, ",")
	if len(parts) < 3 {
		return colorRGBA{}, fmt.Errorf("invalid rgb() color")
	}
	r, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	g, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	b, err3 := strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
	if err1 != nil || err2 != nil || err3 != nil {
		return colorRGBA{}, fmt.Errorf("invalid rgb() values")
	}
	a := 1.0
	if len(parts) >= 4 {
		av, err := strconv.ParseFloat(strings.TrimSpace(parts[3]), 64)
		if err == nil {
			a = av
		}
	}
	return colorRGBA{R: r, G: g, B: b, A: a}, nil
}

// rgbToHSL converts r,g,b ∈ [0,1] to h ∈ [0,1], s ∈ [0,1], l ∈ [0,1].
func rgbToHSL(r, g, b float64) (h, s, l float64) {
	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))
	l = (max + min) / 2
	if max == min {
		return 0, 0, l
	}
	d := max - min
	if l > 0.5 {
		s = d / (2 - max - min)
	} else {
		s = d / (max + min)
	}
	switch max {
	case r:
		h = (g - b) / d
		if g < b {
			h += 6
		}
	case g:
		h = (b-r)/d + 2
	case b:
		h = (r-g)/d + 4
	}
	h /= 6
	return h, s, l
}

// hslToRGB converts h,s,l ∈ [0,1] to r,g,b ∈ [0,1].
func hslToRGB(h, s, l float64) (r, g, b float64) {
	if s == 0 {
		return l, l, l
	}
	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q
	r = hueToRGB(p, q, h+1.0/3)
	g = hueToRGB(p, q, h)
	b = hueToRGB(p, q, h-1.0/3)
	return
}

func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t++
	}
	if t > 1 {
		t--
	}
	switch {
	case t < 1.0/6:
		return p + (q-p)*6*t
	case t < 0.5:
		return q
	case t < 2.0/3:
		return p + (q-p)*(2.0/3-t)*6
	default:
		return p
	}
}

// relativeLuminance computes WCAG 2.1 relative luminance.
func relativeLuminance(c colorRGBA) float64 {
	linearize := func(ch float64) float64 {
		v := ch / 255
		if v <= 0.04045 {
			return v / 12.92
		}
		return math.Pow((v+0.055)/1.055, 2.4)
	}
	return 0.2126*linearize(c.R) + 0.7152*linearize(c.G) + 0.0722*linearize(c.B)
}
