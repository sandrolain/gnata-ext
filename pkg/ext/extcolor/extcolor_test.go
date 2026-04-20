package extcolor_test

import (
	"math"
	"testing"

	"github.com/sandrolain/gnata-ext/pkg/ext/extcolor"
)

// --- helpers ---

func colorObj(r, g, b, a float64) map[string]any {
	return map[string]any{"r": r, "g": g, "b": b, "a": a}
}

func approxEq(a, b, eps float64) bool {
	return math.Abs(a-b) <= eps
}

// ---------- ColorParse ----------

func TestColorParse_Hex6(t *testing.T) {
	fn := extcolor.ColorParse()
	got, err := fn([]any{"#ff8000"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	obj := got.(map[string]any)
	if obj["r"] != float64(255) || obj["g"] != float64(128) || obj["b"] != float64(0) {
		t.Errorf("unexpected parse: %v", obj)
	}
}

func TestColorParse_Hex3(t *testing.T) {
	fn := extcolor.ColorParse()
	got, err := fn([]any{"#f80"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	obj := got.(map[string]any)
	if obj["r"] != float64(255) || obj["g"] != float64(136) || obj["b"] != float64(0) {
		t.Errorf("unexpected parse: %v", obj)
	}
}

func TestColorParse_RGB(t *testing.T) {
	fn := extcolor.ColorParse()
	got, err := fn([]any{"rgb(10, 20, 30)"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	obj := got.(map[string]any)
	if obj["r"] != float64(10) || obj["g"] != float64(20) || obj["b"] != float64(30) {
		t.Errorf("unexpected: %v", obj)
	}
}

func TestColorParse_RGBA(t *testing.T) {
	fn := extcolor.ColorParse()
	got, err := fn([]any{"rgba(10, 20, 30, 0.5)"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	obj := got.(map[string]any)
	if obj["a"] != float64(0.5) {
		t.Errorf("expected alpha=0.5, got %v", obj["a"])
	}
}

func TestColorParse_Unsupported(t *testing.T) {
	fn := extcolor.ColorParse()
	_, err := fn([]any{"hsl(200,50%,50%)"}, nil)
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestColorParse_InvalidHex(t *testing.T) {
	fn := extcolor.ColorParse()
	_, err := fn([]any{"#xyz"}, nil)
	if err == nil {
		t.Error("expected error for invalid hex")
	}
}

func TestColorParse_BadHexLength(t *testing.T) {
	fn := extcolor.ColorParse()
	_, err := fn([]any{"#12"}, nil)
	if err == nil {
		t.Error("expected error for bad hex length")
	}
}

func TestColorParse_NoArgs(t *testing.T) {
	fn := extcolor.ColorParse()
	_, err := fn([]any{}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestColorParse_WrongType(t *testing.T) {
	fn := extcolor.ColorParse()
	_, err := fn([]any{42}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

// ---------- ColorToHex ----------

func TestColorToHex_Basic(t *testing.T) {
	fn := extcolor.ColorToHex()
	got, err := fn([]any{colorObj(255, 128, 0, 1)}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != "#ff8000" {
		t.Errorf("expected '#ff8000', got %v", got)
	}
}

func TestColorToHex_Clamp(t *testing.T) {
	fn := extcolor.ColorToHex()
	got, _ := fn([]any{colorObj(300, -10, 0, 1)}, nil)
	if got != "#ff0000" {
		t.Errorf("expected '#ff0000' after clamping, got %v", got)
	}
}

func TestColorToHex_NoArgs(t *testing.T) {
	fn := extcolor.ColorToHex()
	_, err := fn([]any{}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestColorToHex_WrongType(t *testing.T) {
	fn := extcolor.ColorToHex()
	_, err := fn([]any{"not-an-object"}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

// ---------- ColorToRGB ----------

func TestColorToRGB_Opaque(t *testing.T) {
	fn := extcolor.ColorToRGB()
	got, err := fn([]any{colorObj(255, 0, 0, 1)}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != "rgb(255,0,0)" {
		t.Errorf("expected 'rgb(255,0,0)', got %v", got)
	}
}

func TestColorToRGB_WithAlpha(t *testing.T) {
	fn := extcolor.ColorToRGB()
	got, err := fn([]any{colorObj(255, 0, 0, 0.5)}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != "rgba(255,0,0,0.5)" {
		t.Errorf("expected 'rgba(255,0,0,0.5)', got %v", got)
	}
}

func TestColorToRGB_NoArgs(t *testing.T) {
	fn := extcolor.ColorToRGB()
	_, err := fn([]any{}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestColorToRGB_WrongType(t *testing.T) {
	fn := extcolor.ColorToRGB()
	_, err := fn([]any{42}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

// ---------- ColorToHSL ----------

func TestColorToHSL_Red(t *testing.T) {
	fn := extcolor.ColorToHSL()
	// pure red: h=0, s=100%, l=50%
	got, err := fn([]any{colorObj(255, 0, 0, 1)}, nil)
	if err != nil {
		t.Fatal(err)
	}
	obj := got.(map[string]any)
	h := obj["h"].(float64)
	s := obj["s"].(float64)
	l := obj["l"].(float64)
	if !approxEq(h, 0, 1) || !approxEq(s, 100, 1) || !approxEq(l, 50, 1) {
		t.Errorf("unexpected hsl for red: h=%v s=%v l=%v", h, s, l)
	}
}

func TestColorToHSL_White(t *testing.T) {
	fn := extcolor.ColorToHSL()
	got, _ := fn([]any{colorObj(255, 255, 255, 1)}, nil)
	obj := got.(map[string]any)
	if obj["l"].(float64) != 100.0 {
		t.Errorf("expected l=100 for white, got %v", obj["l"])
	}
}

func TestColorToHSL_NoArgs(t *testing.T) {
	fn := extcolor.ColorToHSL()
	_, err := fn([]any{}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestColorToHSL_WrongType(t *testing.T) {
	fn := extcolor.ColorToHSL()
	_, err := fn([]any{"bad"}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

// ---------- ColorMix ----------

func TestColorMix_Midpoint(t *testing.T) {
	fn := extcolor.ColorMix()
	black := colorObj(0, 0, 0, 1)
	white := colorObj(255, 255, 255, 1)
	got, err := fn([]any{black, white, float64(0.5)}, nil)
	if err != nil {
		t.Fatal(err)
	}
	obj := got.(map[string]any)
	r := obj["r"].(float64)
	if !approxEq(r, 127.5, 1) {
		t.Errorf("expected r≈127.5, got %v", r)
	}
}

func TestColorMix_T0(t *testing.T) {
	fn := extcolor.ColorMix()
	a := colorObj(10, 20, 30, 1)
	b := colorObj(100, 100, 100, 1)
	got, _ := fn([]any{a, b, float64(0)}, nil)
	obj := got.(map[string]any)
	if obj["r"].(float64) != float64(10) {
		t.Errorf("at t=0 expected r=10, got %v", obj["r"])
	}
}

func TestColorMix_NoArgs(t *testing.T) {
	fn := extcolor.ColorMix()
	_, err := fn([]any{colorObj(0, 0, 0, 1), colorObj(255, 255, 255, 1)}, nil)
	if err == nil {
		t.Error("expected error for missing t")
	}
}

func TestColorMix_BadT(t *testing.T) {
	fn := extcolor.ColorMix()
	_, err := fn([]any{colorObj(0, 0, 0, 1), colorObj(255, 255, 255, 1), "not-a-number"}, nil)
	if err == nil {
		t.Error("expected error for bad t")
	}
}

func TestColorMix_BadFirst(t *testing.T) {
	fn := extcolor.ColorMix()
	_, err := fn([]any{"bad", colorObj(255, 255, 255, 1), float64(0.5)}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestColorMix_BadSecond(t *testing.T) {
	fn := extcolor.ColorMix()
	_, err := fn([]any{colorObj(0, 0, 0, 1), "bad", float64(0.5)}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

// ---------- ColorLighten ----------

func TestColorLighten_Basic(t *testing.T) {
	fn := extcolor.ColorLighten()
	// #808080 = rgb(128,128,128) → l≈50%
	c := colorObj(128, 128, 128, 1)
	got, err := fn([]any{c, float64(0.1)}, nil)
	if err != nil {
		t.Fatal(err)
	}
	after := got.(map[string]any)
	// lightened: r/g/b should be higher
	if after["r"].(float64) <= 128 {
		t.Errorf("expected lighter r, got %v", after["r"])
	}
}

func TestColorLighten_NoArgs(t *testing.T) {
	fn := extcolor.ColorLighten()
	_, err := fn([]any{colorObj(128, 128, 128, 1)}, nil)
	if err == nil {
		t.Error("expected error for missing amount")
	}
}

func TestColorLighten_BadAmount(t *testing.T) {
	fn := extcolor.ColorLighten()
	_, err := fn([]any{colorObj(128, 128, 128, 1), "bad"}, nil)
	if err == nil {
		t.Error("expected error for bad amount")
	}
}

func TestColorLighten_WrongType(t *testing.T) {
	fn := extcolor.ColorLighten()
	_, err := fn([]any{"bad", float64(0.1)}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

// ---------- ColorDarken ----------

func TestColorDarken_Basic(t *testing.T) {
	fn := extcolor.ColorDarken()
	c := colorObj(200, 200, 200, 1)
	got, err := fn([]any{c, float64(0.2)}, nil)
	if err != nil {
		t.Fatal(err)
	}
	after := got.(map[string]any)
	if after["r"].(float64) >= 200 {
		t.Errorf("expected darker r, got %v", after["r"])
	}
}

func TestColorDarken_NoArgs(t *testing.T) {
	fn := extcolor.ColorDarken()
	_, err := fn([]any{colorObj(128, 128, 128, 1)}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestColorDarken_BadAmount(t *testing.T) {
	fn := extcolor.ColorDarken()
	_, err := fn([]any{colorObj(128, 128, 128, 1), "bad"}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestColorDarken_WrongType(t *testing.T) {
	fn := extcolor.ColorDarken()
	_, err := fn([]any{42, float64(0.1)}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

// ---------- ColorContrast ----------

func TestColorContrast_BlackWhite(t *testing.T) {
	fn := extcolor.ColorContrast()
	black := colorObj(0, 0, 0, 1)
	white := colorObj(255, 255, 255, 1)
	got, err := fn([]any{black, white}, nil)
	if err != nil {
		t.Fatal(err)
	}
	ratio := got.(float64)
	// WCAG: black/white contrast ≈ 21
	if !approxEq(ratio, 21.0, 0.1) {
		t.Errorf("expected ≈21, got %v", ratio)
	}
}

func TestColorContrast_Same(t *testing.T) {
	fn := extcolor.ColorContrast()
	c := colorObj(100, 100, 100, 1)
	got, _ := fn([]any{c, c}, nil)
	if got.(float64) != 1.0 {
		t.Errorf("expected 1.0 for same color, got %v", got)
	}
}

func TestColorContrast_NoArgs(t *testing.T) {
	fn := extcolor.ColorContrast()
	_, err := fn([]any{colorObj(0, 0, 0, 1)}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestColorContrast_WrongType(t *testing.T) {
	fn := extcolor.ColorContrast()
	_, err := fn([]any{"bad", colorObj(255, 255, 255, 1)}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

// ---------- ColorLuminance ----------

func TestColorLuminance_Black(t *testing.T) {
	fn := extcolor.ColorLuminance()
	got, err := fn([]any{colorObj(0, 0, 0, 1)}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got.(float64) != 0.0 {
		t.Errorf("expected 0 for black, got %v", got)
	}
}

func TestColorLuminance_White(t *testing.T) {
	fn := extcolor.ColorLuminance()
	got, _ := fn([]any{colorObj(255, 255, 255, 1)}, nil)
	if !approxEq(got.(float64), 1.0, 0.001) {
		t.Errorf("expected ≈1.0 for white, got %v", got)
	}
}

func TestColorLuminance_NoArgs(t *testing.T) {
	fn := extcolor.ColorLuminance()
	_, err := fn([]any{}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestColorLuminance_WrongType(t *testing.T) {
	fn := extcolor.ColorLuminance()
	_, err := fn([]any{"bad"}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

// ---------- All ----------

func TestAll_Keys(t *testing.T) {
	m := extcolor.All()
	expected := []string{
		"colorParse", "colorToHex", "colorToRGB", "colorToHSL",
		"colorMix", "colorLighten", "colorDarken", "colorContrast", "colorLuminance",
	}
	for _, k := range expected {
		if _, ok := m[k]; !ok {
			t.Errorf("All(): missing key %q", k)
		}
	}
	if len(m) != len(expected) {
		t.Errorf("All(): expected %d keys, got %d", len(expected), len(m))
	}
}
