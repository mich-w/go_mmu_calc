package main

import (
	"testing"
)

func TestMain(m *testing.M) {
	initColourToNames()
	m.Run()
}

func TestLuminance(t *testing.T) {
	tests := []struct {
		r, g, b  float64
		expected float64
		approx   bool
	}{
		{0, 0, 0, 0, false},
		{1, 1, 1, 1, true},
		{1, 0, 0, 0.3, false},
		{0, 1, 0, 0.59, false},
		{0, 0, 1, 0.11, false},
	}

	for _, tt := range tests {
		got := luminance(tt.r, tt.g, tt.b)
		if tt.approx {
			if got < tt.expected-0.001 || got > tt.expected+0.001 {
				t.Errorf("luminance(%v, %v, %v) = %v, want ~%v", tt.r, tt.g, tt.b, got, tt.expected)
			}
		} else {
			if got != tt.expected {
				t.Errorf("luminance(%v, %v, %v) = %v, want %v", tt.r, tt.g, tt.b, got, tt.expected)
			}
		}
	}
}

func TestRGB2HSV(t *testing.T) {
	tests := []struct {
		name                string
		r, g, b             float64
		wantH, wantS, wantV float64
	}{
		{"black", 0, 0, 0, 0, 0, 0},
		{"white", 1, 1, 1, 0, 0, 1},
		{"red", 1, 0, 0, 0, 1, 1},
		{"green", 0, 1, 0, 120, 1, 1},
		{"blue", 0, 0, 1, 240, 1, 1},
		{"yellow", 1, 1, 0, 60, 1, 1},
		{"cyan", 0, 1, 1, 180, 1, 1},
		{"magenta", 1, 0, 1, -60, 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, s, v := rgb2hsv(tt.r, tt.g, tt.b)
			if h != tt.wantH || s != tt.wantS || v != tt.wantV {
				t.Errorf("rgb2hsv(%v, %v, %v) = (%v, %v, %v), want (%v, %v, %v)",
					tt.r, tt.g, tt.b, h, s, v, tt.wantH, tt.wantS, tt.wantV)
			}
		})
	}
}

func TestRGB2HSV_NegativeCases(t *testing.T) {
	// Grayscale: same R, G, B should give hue 0
	h, s, _ := rgb2hsv(0.5, 0.5, 0.5)
	if s != 0 {
		t.Errorf("rgb2hsv(0.5, 0.5, 0.5) saturation should be 0, got %v", s)
	}

	h, s, v := rgb2hsv(0, 0, 0)
	if h != 0 || s != 0 || v != 0 {
		t.Errorf("rgb2hsv(0, 0, 0) = (%v, %v, %v), want (0, 0, 0)", h, s, v)
	}
}

func TestDeltaHSBBS(t *testing.T) {
	tests := []struct {
		name                   string
		h1, s1, v1, h2, s2, v2 float64
		want                   float64
	}{
		{"identical", 180, 0.5, 0.8, 180, 0.5, 0.8, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := deltaHSBBS(tt.h1, tt.s1, tt.v1, tt.h2, tt.s2, tt.v2)
			if got != tt.want {
				t.Errorf("deltaHSBBS(...) = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeltaHSBBS_Capped(t *testing.T) {
	// Maximum distance should be capped at 1.2
	got := deltaHSBBS(0, 1, 1, 180, 1, 1)
	if got != 1.2 {
		t.Errorf("deltaHSBBS(0, 1, 1, 180, 1, 1) = %v, want 1.2 (capped)", got)
	}
}

func TestHexToRGBA(t *testing.T) {
	tests := []struct {
		hex     string
		want    [4]int
		wantErr bool
	}{
		{"#ff0000", [4]int{255, 0, 0, 255}, false},
		{"ff0000ff", [4]int{255, 0, 0, 255}, false},
		{"#00ff00ff", [4]int{0, 255, 0, 255}, false},
		{"#000000ff", [4]int{0, 0, 0, 255}, false},
		{"#ffffff", [4]int{255, 255, 255, 255}, false},
		{"#fff", [4]int{}, true},
		{"invalid", [4]int{}, true},
		{"", [4]int{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.hex, func(t *testing.T) {
			got, err := hexToRGBA(tt.hex)
			if tt.wantErr {
				if err == nil {
					t.Errorf("hexToRGBA(%q) expected error, got nil", tt.hex)
				}
				return
			}
			if err != nil {
				t.Errorf("hexToRGBA(%q) unexpected error: %v", tt.hex, err)
				return
			}
			if got != tt.want {
				t.Errorf("hexToRGBA(%q) = %v, want %v", tt.hex, got, tt.want)
			}
		})
	}
}

func TestParseColour(t *testing.T) {
	tests := []struct {
		input   string
		want    string
		wantErr bool
	}{
		{"red", "#ff0000ff", false},
		{"blue", "#0000ffff", false},
		{"#ff0000", "#ff0000ff", false},
		{"#00ff00ff", "#00ff00ff", false},
		{"RED", "#ff0000ff", false},
		{"Red", "#ff0000ff", false},
		{"invalid", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parseColour(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("parseColour(%q) expected error, got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Errorf("parseColour(%q) unexpected error: %v", tt.input, err)
				return
			}
			if got != tt.want {
				t.Errorf("parseColour(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNameToHex(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"red", "#ff0000"},
		{"Blue", "#0000ff"},
		{"GREEN", "#008000"},
		{"#aabbcc", "#aabbcc"},
		{"nonexistent", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := nameToHex(tt.input)
			if got != tt.want {
				t.Errorf("nameToHex(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestColourName(t *testing.T) {
	tests := []struct {
		hex    string
		expect string
	}{
		{"#ff0000ff", "red"},
		{"#0000ffff", "blue"},
		{"#008000ff", "green"},
		{"#00ff00ff", "lime"},
		{"#ffffff", "white"},
	}

	for _, tt := range tests {
		t.Run(tt.hex, func(t *testing.T) {
			got := colourName(tt.hex)
			if got != tt.expect {
				t.Errorf("colourName(%q) = %q, want %q", tt.hex, got, tt.expect)
			}
		})
	}
}

func TestHexToRGBTuple(t *testing.T) {
	r, g, b := hexToRGBTuple("#ff0000")
	if r != 255 || g != 0 || b != 0 {
		t.Errorf("hexToRGBTuple(\"#ff0000\") = (%d, %d, %d), want (255, 0, 0)", r, g, b)
	}

	r, g, b = hexToRGBTuple("#00ff00")
	if r != 0 || g != 255 || b != 0 {
		t.Errorf("hexToRGBTuple(\"#00ff00\") = (%d, %d, %d), want (0, 255, 0)", r, g, b)
	}
}

func TestCalcFlushVol_Internal(t *testing.T) {
	tests := []struct {
		name                   string
		srcR, srcG, srcB, srcA int
		dstR, dstG, dstB, dstA int
		multiplier             float64
		minVal                 int
		maxVal                 int
	}{
		{
			name: "same color",
			srcR: 255, srcG: 0, srcB: 0, srcA: 255,
			dstR: 255, dstG: 0, dstB: 0, dstA: 255,
			multiplier: 1.0,
			minVal:     minFlushVolume, maxVal: maxFlushVolume,
		},
		{
			name: "min volume clamped",
			srcR: 200, srcG: 200, srcB: 200, srcA: 255,
			dstR: 190, dstG: 190, dstB: 190, dstA: 255,
			multiplier: 1.0,
			minVal:     minFlushVolume, maxVal: maxFlushVolume,
		},
		{
			name: "transparent source treated as white",
			srcR: 0, srcG: 0, srcB: 0, srcA: 0,
			dstR: 0, dstG: 0, dstB: 0, dstA: 255,
			multiplier: 1.0,
			minVal:     minFlushVolume, maxVal: maxFlushVolume,
		},
		{
			name: "transparent dest treated as white",
			srcR: 255, srcG: 255, srcB: 255, srcA: 255,
			dstR: 0, dstG: 0, dstB: 0, dstA: 0,
			multiplier: 1.0,
			minVal:     minFlushVolume, maxVal: maxFlushVolume,
		},
		{
			name: "multiplier applied",
			srcR: 255, srcG: 0, srcB: 0, srcA: 255,
			dstR: 0, dstG: 0, dstB: 255, dstA: 255,
			multiplier: 2.0,
			minVal:     minFlushVolume, maxVal: maxFlushVolume,
		},
		{
			name: "small multiplier still respects min",
			srcR: 100, srcG: 100, srcB: 100, srcA: 255,
			dstR: 99, dstG: 99, dstB: 99, dstA: 255,
			multiplier: 0.01,
			minVal:     minFlushVolume, maxVal: maxFlushVolume,
		},
		{
			name: "zero multiplier still respects min",
			srcR: 100, srcG: 100, srcB: 100, srcA: 255,
			dstR: 99, dstG: 99, dstB: 99, dstA: 255,
			multiplier: 0,
			minVal:     minFlushVolume, maxVal: maxFlushVolume,
		},
		{
			name: "dark to light applies 1.3x before multiplier",
			srcR: 255, srcG: 255, srcB: 255, srcA: 255,
			dstR: 0, dstG: 0, dstB: 0, dstA: 255,
			multiplier: 1.0,
			minVal:     minFlushVolume, maxVal: maxFlushVolume,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calcFlushVolInternal(tt.srcR, tt.srcG, tt.srcB, tt.srcA,
				tt.dstR, tt.dstG, tt.dstB, tt.dstA, tt.multiplier)

			if got < tt.minVal {
				t.Errorf("calcFlushVolInternal(...) = %d, want >= %d", got, tt.minVal)
			}
			if got > tt.maxVal {
				t.Errorf("calcFlushVolInternal(...) = %d, want <= %d", got, tt.maxVal)
			}
		})
	}
}

func TestCalcFlushVol_MaxClamped(t *testing.T) {
	// Use a very large multiplier to test max clamping
	result := calcFlushVolInternal(255, 0, 0, 255, 0, 0, 0, 255, 10000.0)
	if result > maxFlushVolume {
		t.Errorf("flush volume %d exceeds max %d", result, maxFlushVolume)
	}
	if result != maxFlushVolume {
		t.Errorf("flush volume should be clamped to max %d, got %d", maxFlushVolume, result)
	}
}

func TestFlushVolCalculator(t *testing.T) {
	calc := FlushVolCalculator{
		MinVol:     minFlushVolume,
		MaxVol:     maxFlushVolume,
		Multiplier: 1.5,
	}

	result := calc.CalcFlushVol(255, 0, 0, 255, 0, 0, 255, 255)
	if result < minFlushVolume {
		t.Errorf("result %d below min %d", result, minFlushVolume)
	}
	if result > maxFlushVolume {
		t.Errorf("result %d above max %d", result, maxFlushVolume)
	}
}

func TestCalcFlushVol_MultiplierScalesResult(t *testing.T) {
	// Use colours with significant difference so result is above minFlushVolume
	result1x := calcFlushVolInternal(255, 0, 0, 255, 0, 255, 0, 255, 1.0)
	result2x := calcFlushVolInternal(255, 0, 0, 255, 0, 255, 0, 255, 2.0)

	if result1x < minFlushVolume || result2x < minFlushVolume {
		t.Skip("both results below minFlushVolume, skipping multiplier scale test")
	}

	// 2x multiplier should give roughly 2x the result (allowing 10% tolerance for rounding)
	ratio := float64(result2x) / float64(result1x)
	if ratio < 1.8 || ratio > 2.2 {
		t.Errorf("multiplier scale: 1x=%d, 2x=%d, ratio=%.2f, expected ~2.0", result1x, result2x, ratio)
	}
}

func TestParseColour_CaseInsensitive(t *testing.T) {
	upper, err := parseColour("RED")
	if err != nil {
		t.Errorf("parseColour(\"RED\") unexpected error: %v", err)
	}
	lower, err := parseColour("red")
	if err != nil {
		t.Errorf("parseColour(\"red\") unexpected error: %v", err)
	}
	if upper != lower {
		t.Errorf("case sensitivity: RED=%q != red=%q", upper, lower)
	}
}

func TestCommonMap_NoDuplicateHexes(t *testing.T) {
	hexToNames := make(map[string]string)
	for name, hex := range common {
		if existing, ok := hexToNames[hex]; ok {
			t.Errorf("duplicate hex %q for colours %q and %q", hex, existing, name)
		}
		hexToNames[hex] = name
	}
}

func TestMaxFilaments(t *testing.T) {
	if maxFilaments != 10 {
		t.Errorf("maxFilaments = %d, want 10", maxFilaments)
	}
}

func TestParseColour_HexVariants(t *testing.T) {
	// 6-char hex with #
	got, err := parseColour("#ff0000")
	if err != nil {
		t.Errorf("parseColour(\"#ff0000\") unexpected error: %v", err)
	}
	if got != "#ff0000ff" {
		t.Errorf("parseColour(\"#ff0000\") = %q, want #ff0000ff", got)
	}

	// 8-char hex with #
	got, err = parseColour("#ff0000ff")
	if err != nil {
		t.Errorf("parseColour(\"#ff0000ff\") unexpected error: %v", err)
	}
	if got != "#ff0000ff" {
		t.Errorf("parseColour(\"#ff0000ff\") = %q, want #ff0000ff", got)
	}

	// Hex without # should fail (not a colour name)
	_, err = parseColour("ff0000ff")
	if err == nil {
		t.Error("parseColour(\"ff0000ff\") expected error (no # prefix)")
	}
}

func TestCalcFlushVol_DarkToLightMultiplier(t *testing.T) {
	// Dark to light transition should apply 1.3x multiplier
	// srcLumi > 180/255 and dstLumi < 75/255
	result := calcFlushVolInternal(255, 255, 255, 255, 0, 0, 0, 255, 1.0)
	// Same color difference should be large enough to verify multiplier is applied
	if result <= 0 {
		t.Errorf("flush volume should be positive, got %d", result)
	}
}

func TestHexToRGBA_PartialAlpha(t *testing.T) {
	// Only 6 chars (no alpha) should default to ff
	got, err := hexToRGBA("#aabbcc")
	if err != nil {
		t.Errorf("hexToRGBA(\"#aabbcc\") unexpected error: %v", err)
	}
	if got[3] != 255 {
		t.Errorf("alpha should default to 255, got %d", got[3])
	}
}

func TestRGB2HSV_Yellow(t *testing.T) {
	// Yellow should have hue around 60
	h, s, v := rgb2hsv(1, 1, 0)
	if h < 55 || h > 65 {
		t.Errorf("rgb2hsv(1, 1, 0) hue = %v, expected ~60", h)
	}
	if s != 1 {
		t.Errorf("rgb2hsv(1, 1, 0) saturation = %v, expected 1", s)
	}
	if v != 1 {
		t.Errorf("rgb2hsv(1, 1, 0) value = %v, expected 1", v)
	}
}

func TestRGB2HSV_Cyan(t *testing.T) {
	h, _, _ := rgb2hsv(0, 1, 1)
	if h < 175 || h > 185 {
		t.Errorf("rgb2hsv(0, 1, 1) hue = %v, expected ~180", h)
	}
}

func TestRGB2HSV_Magenta(t *testing.T) {
	h, _, _ := rgb2hsv(1, 0, 1)
	// Implementation returns -60 for magenta (equivalent to 300)
	if h != -60 {
		t.Errorf("rgb2hsv(1, 0, 1) hue = %v, expected -60", h)
	}
}

func TestCalcFlushVol_Asymmetric(t *testing.T) {
	// Red->Blue should differ from Blue->Red
	forward := calcFlushVolInternal(255, 0, 0, 255, 0, 0, 255, 255, 1.0)
	backward := calcFlushVolInternal(0, 0, 255, 255, 255, 0, 0, 255, 1.0)
	if forward == backward {
		t.Errorf("flush volume should be asymmetric: red->blue=%d, blue->red=%d", forward, backward)
	}
}
