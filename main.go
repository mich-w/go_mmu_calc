package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Constants
const (
	defaultFlushMultiplier float64 = 1.0
	maxFlushVolume         int     = 20000
	minFlushVolume         int     = 60
	maxFilaments           int     = 10
)

// Common colour names mapped to hex values.
var common = map[string]string{
	"aquamarine":           "#7fffd4",
	"black":                "#000000",
	"blue":                 "#0000ff",
	"blueviolet":           "#8a2be2",
	"brown":                "#a52a2a",
	"chartreuse":           "#7fff00",
	"coral":                "#ff7f50",
	"cornflower":           "#6495ed",
	"crimson":              "#dc143c",
	"cyan":                 "#00ffff",
	"darkblue":             "#00008b",
	"darkcyan":             "#008b8b",
	"darkgoldenrod":        "#b8860b",
	"darkgray":             "#a9a9a9",
	"darkgreen":            "#006400",
	"darkkhaki":            "#bdb76b",
	"darkmagenta":          "#8b008b",
	"darkolivegreen":       "#556b2f",
	"darkorange":           "#ff8c00",
	"darkorchid":           "#9932cc",
	"darkred":              "#8b0000",
	"darksalmon":           "#e9967a",
	"darkseagreen":         "#8fbc8f",
	"darkslateblue":        "#483d8b",
	"darkslategray":        "#2f4f4f",
	"darkturquoise":        "#00ced1",
	"darkviolet":           "#9400d3",
	"deeppink":             "#ff1493",
	"deepskyblue":          "#00bfff",
	"dimgray":              "#696969",
	"dodgerblue":           "#1e90ff",
	"firebrick":            "#b22222",
	"floralwhite":          "#fffaf0",
	"forestgreen":          "#228b22",
	"gainsboro":            "#dcdcdc",
	"ghostwhite":           "#f8f8ff",
	"gold":                 "#ffd700",
	"gray":                 "#808080",
	"green":                "#008000",
	"honeydew":             "#f0fff0",
	"hotpink":              "#ff69b4",
	"indianred":            "#cd5c5c",
	"indigo":               "#4b0082",
	"ivory":                "#fffff0",
	"khaki":                "#f0e68c",
	"lavender":             "#e6e6fa",
	"lavenderblush":        "#fff0f5",
	"lemonchiffon":         "#fffacd",
	"lightblue":            "#add8e6",
	"lightcoral":           "#f08080",
	"lightcyan":            "#e0ffff",
	"lightgoldenrodyellow": "#fafad2",
	"lightgray":            "#d3d3d3",
	"lightgreen":           "#90ee90",
	"lightpink":            "#ffb6c1",
	"lightsalmon":          "#ffa07a",
	"lightseagreen":        "#20b2aa",
	"lightskyblue":         "#87cefa",
	"lightslategray":       "#778899",
	"lightsteelblue":       "#b0c4de",
	"lightyellow":          "#ffffe0",
	"lime":                 "#00ff00",
	"limegreen":            "#32cd32",
	"linen":                "#faf0e6",
	"magenta":              "#ff00ff",
	"maroon":               "#800000",
	"mediumaquamarine":     "#66cdaa",
	"mediumblue":           "#0000cd",
	"mediumorchid":         "#ba55d3",
	"mediumpurple":         "#9370db",
	"mediumseagreen":       "#3cb371",
	"mediumslateblue":      "#7b68ee",
	"mediumspringgreen":    "#00fa9a",
	"mediumturquoise":      "#48d1cc",
	"mediumvioletred":      "#c71585",
	"midnightblue":         "#191970",
	"mintcream":            "#f5fffa",
	"mistyrose":            "#ffe4e1",
	"moccasin":             "#ffe4b5",
	"navajowhite":          "#ffdead",
	"navy":                 "#000080",
	"oldlace":              "#fdf5e6",
	"olive":                "#808000",
	"orange":               "#ffa500",
	"orangered":            "#ff4500",
	"orchid":               "#da70d6",
	"palegoldenrod":        "#eee8aa",
	"palegreen":            "#98fb98",
	"paleturquoise":        "#afeeee",
	"palevioletred":        "#db7093",
	"papayawhip":           "#ffefd5",
	"peachpuff":            "#ffdab9",
	"peru":                 "#cd853f",
	"pink":                 "#ffc0cb",
	"plum":                 "#dda0dd",
	"purple":               "#800080",
	"red":                  "#ff0000",
	"rosybrown":            "#bc8f8f",
	"royalblue":            "#4169e1",
	"saddlebrown":          "#8b4513",
	"salmon":               "#fa8072",
	"sandybrown":           "#f4a460",
	"seagreen":             "#2e8b57",
	"seashell":             "#fff5ee",
	"sienna":               "#a0522d",
	"silver":               "#c0c0c0",
	"skyblue":              "#87ceeb",
	"slateblue":            "#6a5acd",
	"slategray":            "#708090",
	"snow":                 "#fffafa",
	"springgreen":          "#00ff7f",
	"steelblue":            "#4682b4",
	"tan":                  "#d2b48c",
	"teal":                 "#008080",
	"thistle":              "#d8bfd8",
	"tomato":               "#ff6347",
	"turquoise":            "#40e0d0",
	"violet":               "#ee82ee",
	"wheat":                "#f5deb3",
	"white":                "#ffffff",
	"whitesmoke":           "#f5f5f5",
	"yellow":               "#ffff00",
	"yellowgreen":          "#9acd32",
}

// colourToNames maps hex values back to their common colour names for O(1) lookup.
var colourToNames map[string]string

// FilamentColour represents a filament colour.
type FilamentColour struct {
	HexCode string
}

// FlushVolCalculator replicates the OrcaSlicer flush volume calculation.
type FlushVolCalculator struct {
	MinVol     int
	MaxVol     int
	Multiplier float64
}

// CLIArgs holds parsed command-line arguments.
type CLIArgs struct {
	ListColors    bool
	ListColorsHue bool
	Colours       []string
	Multiplier    float64
	NoColor       bool
}

// initColourToNames builds a reverse lookup map from hex values to colour names.
func initColourToNames() {
	colourToNames = make(map[string]string, len(common))
	for name, hex := range common {
		colourToNames[hex] = name
	}
}

func luminance(r, g, b float64) float64 {
	return r*0.3 + g*0.59 + b*0.11
}

func calcFlushVolInternal(srcR, srcG, srcB, srcA, dstR, dstG, dstB, dstA int, multiplier float64) int {
	if srcA == 0 {
		srcR, srcG, srcB = 255, 255, 255
	}
	if dstA == 0 {
		dstR, dstG, dstB = 255, 255, 255
	}

	srcF := [3]float64{float64(srcR) / 255.0, float64(srcG) / 255.0, float64(srcB) / 255.0}
	dstF := [3]float64{float64(dstR) / 255.0, float64(dstG) / 255.0, float64(dstB) / 255.0}

	srcH, srcS, srcV := rgb2hsv(srcF[0], srcF[1], srcF[2])
	dstH, dstS, dstV := rgb2hsv(dstF[0], dstF[1], dstF[2])

	hsDist := deltaHSBBS(srcH, srcS, srcV, dstH, dstS, dstV)

	srcLumi := luminance(srcF[0], srcF[1], srcF[2])
	dstLumi := luminance(dstF[0], dstF[1], dstF[2])

	var lumiFlush float64
	if dstLumi >= srcLumi {
		lumiFlush = math.Pow(dstLumi-srcLumi, 0.7) * 560.0
	} else {
		lumiFlush = (srcLumi - dstLumi) * 80.0
		interV := 0.67*dstV + 0.33*srcV
		if interV < hsDist {
			hsDist = interV
		}
	}

	hsFlush := 230.0 * hsDist
	a := hsFlush
	b := lumiFlush
	angle := 120.0 * math.Pi / 180.0
	flushVol := math.Sqrt(a*a + b*b - 2*a*b*math.Cos(angle))

	darkColorThres := 180.0 / 255.0
	lightColorThres := 75.0 / 255.0
	if srcLumi > darkColorThres && dstLumi < lightColorThres {
		flushVol *= 1.3
	}

	flushVol *= multiplier

	result := int(flushVol)
	result = max(result, minFlushVolume)
	result = min(result, maxFlushVolume)

	return result
}

func (f *FlushVolCalculator) CalcFlushVol(
	srcR, srcG, srcB, srcA, dstR, dstG, dstB, dstA int,
) int {
	return calcFlushVolInternal(srcR, srcG, srcB, srcA, dstR, dstG, dstB, dstA, f.Multiplier)
}

// rgbToHSV converts RGB values in [0,1] to HSV.
func rgb2hsv(r, g, b float64) (float64, float64, float64) {
	cmax := math.Max(r, math.Max(g, b))
	cmin := math.Min(r, math.Min(g, b))
	delta := cmax - cmin

	var h float64
	if delta < 1e-6 {
		h = 0.0
	} else if cmax == r {
		h = 60.0 * math.Mod((g-b)/delta, 6)
	} else if cmax == g {
		h = 60.0 * ((b-r)/delta + 2)
	} else {
		h = 60.0 * ((r-g)/delta + 4)
	}

	var s float64
	if cmax < 1e-6 {
		s = 0.0
	} else {
		s = delta / cmax
	}
	return h, s, cmax
}

func deltaHSBBS(h1, s1, v1, h2, s2, v2 float64) float64 {
	dx := math.Cos(h1*math.Pi/180.0)*s1*v1 - math.Cos(h2*math.Pi/180.0)*s2*v2
	dy := math.Sin(h1*math.Pi/180.0)*s1*v1 - math.Sin(h2*math.Pi/180.0)*s2*v2
	dxy := math.Hypot(dx, dy)
	if dxy > 1.2 {
		return 1.2
	}
	return dxy
}

// hexToRGBA parses a hex colour string and returns RGBA values in [0,255].
func hexToRGBA(hexCode string) ([4]int, error) {
	hexCode = strings.TrimPrefix(hexCode, "#")
	if len(hexCode) == 6 {
		hexCode += "ff"
	}
	if len(hexCode) != 8 {
		return [4]int{}, fmt.Errorf("invalid hex colour spec: %s", hexCode)
	}

	parts := [4]int{
		int(parseHex(hexCode[0:2])),
		int(parseHex(hexCode[2:4])),
		int(parseHex(hexCode[4:6])),
		int(parseHex(hexCode[6:8])),
	}
	return parts, nil
}

func parseHex(s string) uint64 {
	v, _ := strconv.ParseUint(s, 16, 64)
	return v
}

// parseColour resolves a colour name or hex string to a normalised hex code.
func parseColour(colour string) (string, error) {
	hexCode := nameToHex(colour)
	if hexCode == "" {
		return "", fmt.Errorf("unknown colour specifier: '%s'", colour)
	}
	if !strings.HasPrefix(hexCode, "#") {
		return "", fmt.Errorf("invalid hex colour: %s", hexCode)
	}
	if len(hexCode) == 7 {
		hexCode += "ff"
	}
	return hexCode, nil
}

func nameToHex(name string) string {
	name = strings.ToLower(name)
	if hex, ok := common[name]; ok {
		return hex
	}
	if strings.HasPrefix(name, "#") {
		return name
	}
	return ""
}

func hexToRGBTuple(hexCode string) (int, int, int) {
	rgb := strings.TrimPrefix(hexCode, "#")[:6]
	return int(parseHex(rgb[0:2])), int(parseHex(rgb[2:4])), int(parseHex(rgb[4:6]))
}

func makeSquare(hexCode string, useColor bool) string {
	if !useColor {
		return "  "
	}
	r, g, b := hexToRGBTuple(hexCode)
	return fmt.Sprintf("\033[48;2;%d;%d;%dm  \033[49m", r, g, b)
}

func hexDisplay(hexCode string) string {
	clean := hexCode
	if len(clean) > 7 {
		clean = clean[:7]
	}
	return fmt.Sprintf("(%s)", clean)
}

func colourHue(hexCode string) float64 {
	rgb := strings.TrimPrefix(hexCode, "#")[:6]
	rN := float64(parseHex(rgb[0:2])) / 255.0
	gN := float64(parseHex(rgb[2:4])) / 255.0
	bN := float64(parseHex(rgb[4:6])) / 255.0
	h, _, _ := rgb2hsv(rN, gN, bN)
	return h
}

func colourName(hexCode string) string {
	rgb := strings.TrimPrefix(hexCode, "#")[:6]
	if name, ok := colourToNames["#"+rgb]; ok {
		return name
	}
	if strings.HasSuffix(hexCode, "ff") {
		return hexCode[:7]
	}
	return hexCode
}

// buildFlushMatrix computes the pairwise flush volumes for all colour pairs.
func buildFlushMatrix(colours []FilamentColour, multiplier float64) []int {
	n := len(colours)
	matrix := make([]int, n*n)
	calc := FlushVolCalculator{
		MinVol:     minFlushVolume,
		MaxVol:     maxFlushVolume,
		Multiplier: multiplier,
	}

	for fromIdx := 0; fromIdx < n; fromIdx++ {
		srcRGBA, err := hexToRGBA(colours[fromIdx].HexCode)
		if err != nil {
			continue
		}

		for toIdx := 0; toIdx < n; toIdx++ {
			if fromIdx == toIdx {
				matrix[fromIdx*n+toIdx] = 0
				continue
			}

			dstRGBA, err := hexToRGBA(colours[toIdx].HexCode)
			if err != nil {
				continue
			}
			matrix[fromIdx*n+toIdx] = calc.CalcFlushVol(
				srcRGBA[0], srcRGBA[1], srcRGBA[2], srcRGBA[3],
				dstRGBA[0], dstRGBA[1], dstRGBA[2], dstRGBA[3],
			)
		}
	}
	return matrix
}

func printMatrix(matrix []int, colours []FilamentColour, useColor bool) {
	n := len(colours)
	maxVal := 0
	for _, v := range matrix {
		if v > maxVal {
			maxVal = v
		}
	}
	colWidth := 4
	if w := len(strconv.Itoa(maxVal)); w > colWidth {
		colWidth = w
	}
	idxWidth := 1
	if n-1 > 0 {
		idxWidth = len(strconv.Itoa(n - 1))
	}
	leftColWidth := idxWidth + 5
	cellWidth := colWidth + 2

	sep := func() string {
		parts := []string{"+" + strings.Repeat("-", leftColWidth)}
		for i := 0; i < n; i++ {
			parts = append(parts, "+"+strings.Repeat("-", cellWidth))
		}
		parts = append(parts, "+")
		return strings.Join(parts, "")
	}

	fmt.Println()
	fmt.Printf("Flush-volume matrix (%dx%d):\n", n, n)
	fmt.Println(sep())

	// Header row
	headerLeft := fmt.Sprintf("%"+strconv.Itoa(leftColWidth)+"s|", "")
	headerCells := make([]string, n)
	for c := 0; c < n; c++ {
		sq := makeSquare(colours[c].HexCode, useColor)
		headerCells[c] = fmt.Sprintf(" %d %s ", c, sq)
	}
	fmt.Println("|" + headerLeft + strings.Join(headerCells, "|") + "|")
	fmt.Println(sep())

	// Body rows
	for r := 0; r < n; r++ {
		sq := makeSquare(colours[r].HexCode, useColor)
		left := fmt.Sprintf(" %d %s ", r, sq)
		left = fmt.Sprintf("%"+strconv.Itoa(leftColWidth)+"s", left)

		cells := make([]string, n)
		for i, v := range matrix[r*n : r*n+n] {
			cells[i] = fmt.Sprintf(" %*d ", colWidth, v)
		}
		fmt.Printf("|%s|%s|\n", left, strings.Join(cells, "|"))
	}
	fmt.Println(sep())
	fmt.Println()
}

func printHelp() {
	fmt.Fprintf(os.Stderr, `Usage: go_mmu_calc [OPTIONS] [COLOURS...]

Calculate flush volumes for MMU filament changes.

Options:
  -l, --list-colors          List all available colours (alphabetical)
  -lh, --list-colors-hue     List all available colours (by hue)
  -m, --multiplier FLOAT     Flush volume multiplier (default: %.1f)
      --no-color             Disable coloured output
  -h, --help                 Show this help message

If no colours are provided, colours are read interactively.

Examples:
  go_mmu_calc red green blue
  go_mmu_calc -l
  go_mmu_calc "#ff0000" "#0000ff" -m 1.5
`, defaultFlushMultiplier)
}

func printColorList(useColor bool, sortByHue bool, compact bool) {
	sortedColors := make([]string, 0, len(common))
	for name := range common {
		sortedColors = append(sortedColors, name)
	}
	if sortByHue {
		sort.Slice(sortedColors, func(i, j int) bool {
			return colourHue(common[sortedColors[i]]) < colourHue(common[sortedColors[j]])
		})
	} else {
		sort.Strings(sortedColors)
	}

	if compact {
		parts := make([]string, len(sortedColors))
		for i, name := range sortedColors {
			hexCode := common[name]
			square := makeSquare(hexCode, useColor)
			parts[i] = fmt.Sprintf("%s %s (%s)", square, name, hexCode)
		}
		fmt.Println(strings.Join(parts, ", "))
		return
	}

	maxNameLen := 0
	for _, name := range sortedColors {
		if l := len(name); l > maxNameLen {
			maxNameLen = l
		}
	}
	fmt.Printf("Available colours (%d total):\n", len(sortedColors))
	for _, name := range sortedColors {
		hexCode := common[name]
		square := makeSquare(hexCode, useColor)
		fmt.Printf("  %s%-*s (%s)\n", square, maxNameLen, name, hexCode)
	}
}

func resolveColours(reader *bufio.Reader, names []string) ([]FilamentColour, error) {
	fmt.Printf("Enter the number of filaments (1-%d):\n", maxFilaments)
	n, err := readPositiveInt(reader, maxFilaments)
	if err != nil {
		return nil, err
	}

	resolved := make([]FilamentColour, n)
	for i := 0; i < n; i++ {
		resolved[i], err = resolveSingleColour(reader, i+1)
		if err != nil {
			return nil, err
		}
	}
	return resolved, nil
}

func readPositiveInt(reader *bufio.Reader, max int) (int, error) {
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			return 0, err
		}
		input = strings.TrimSpace(input)
		n, err := strconv.Atoi(input)
		if err != nil || n <= 0 {
			fmt.Println("  Please enter a positive integer.")
			continue
		}
		if n > max {
			fmt.Printf("  Maximum %d filaments allowed.\n", max)
			continue
		}
		return n, nil
	}
}

func resolveSingleColour(reader *bufio.Reader, index int) (FilamentColour, error) {
	for {
		fmt.Printf("  Colour #%d (name or hex, e.g. 'red' or '#ff0000'): ", index)
		input, err := reader.ReadString('\n')
		if err != nil {
			return FilamentColour{}, err
		}
		input = strings.TrimSpace(input)
		if input == "" {
			fmt.Println("  Colour cannot be empty.")
			continue
		}
		hexCode, err := parseColour(input)
		if err != nil {
			fmt.Printf("  %v\n", err)
			fmt.Println("  Available colours: ")
			printColorList(true, true, true)
			continue
		}
		return FilamentColour{HexCode: hexCode}, nil
	}
}

func resolveMultiplier(reader *bufio.Reader, defaultVal float64) float64 {
	fmt.Printf("Enter flush multiplier or press Enter to skip (default %v):\n", defaultVal)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			return defaultVal
		}
		input = strings.TrimSpace(input)
		if input == "" {
			return defaultVal
		}
		val, err := strconv.ParseFloat(input, 64)
		if err != nil || val <= 0 {
			fmt.Println("  Please enter a positive number.")
			continue
		}
		return val
	}
}

func printSummary(colours []FilamentColour, multiplier float64, useColor bool) {
	fmt.Println()
	fmt.Printf("Colours (%d):\n", len(colours))
	maxNameLen := 0
	for _, c := range colours {
		if l := len(colourName(c.HexCode)); l > maxNameLen {
			maxNameLen = l
		}
	}
	for i, c := range colours {
		square := ""
		if useColor {
			square = makeSquare(c.HexCode, true)
		}
		fmt.Printf(" %d %s %-*s %s\n", i, square, maxNameLen, colourName(c.HexCode), hexDisplay(c.HexCode))
	}
	fmt.Printf("Multiplier: %.1f\n", multiplier)
}

func parseArgs() (*CLIArgs, error) {
	args := &CLIArgs{
		Multiplier: defaultFlushMultiplier,
	}

	for _, arg := range os.Args[1:] {
		switch arg {
		case "--help", "-h":
			printHelp()
			os.Exit(0)
		case "--list-colors", "-l":
			args.ListColors = true
			return args, nil
		case "--list-colors-hue", "-lh":
			args.ListColorsHue = true
			return args, nil
		}
	}

	i := 0
	rawArgs := os.Args[1:]
	for i < len(rawArgs) {
		switch rawArgs[i] {
		case "-m", "--multiplier":
			if i+1 < len(rawArgs) {
				val, err := strconv.ParseFloat(rawArgs[i+1], 64)
				if err != nil {
					return nil, fmt.Errorf("invalid multiplier: %s", rawArgs[i+1])
				}
				args.Multiplier = val
				i += 2
			} else {
				return nil, fmt.Errorf("missing value for %s", rawArgs[i])
			}
		case "--no-color":
			args.NoColor = true
			i++
		default:
			if !strings.HasPrefix(rawArgs[i], "-") {
				args.Colours = append(args.Colours, rawArgs[i])
				i++
			} else {
				return nil, fmt.Errorf("unknown option: %s", rawArgs[i])
			}
		}
	}

	return args, nil
}

func main() {
	initColourToNames()
	args, err := parseArgs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	useColor := !args.NoColor

	if args.ListColors {
		printColorList(useColor, false, false)
		os.Exit(0)
	}

	if args.ListColorsHue {
		printColorList(useColor, true, false)
		os.Exit(0)
	}

	var resolved []FilamentColour
	interactive := false

	if len(args.Colours) > 0 {
		if len(args.Colours) > maxFilaments {
			fmt.Fprintf(os.Stderr, "Error: maximum %d filaments allowed\n", maxFilaments)
			os.Exit(1)
		}
		resolved, err = resolveColoursFromArgs(args.Colours)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			fmt.Println("Available colours: ")
			printColorList(useColor, true, true)
			os.Exit(1)
		}
	} else {
		interactive = true
		reader := bufio.NewReader(os.Stdin)
		resolved, err = resolveColours(reader, nil)
		if err != nil {
			fmt.Println("\nAborted.")
			os.Exit(1)
		}
	}

	multiplier := args.Multiplier
	if interactive {
		reader := bufio.NewReader(os.Stdin)
		multiplier = resolveMultiplier(reader, defaultFlushMultiplier)
	}

	printSummary(resolved, multiplier, useColor)

	matrix := buildFlushMatrix(resolved, multiplier)
	printMatrix(matrix, resolved, useColor)
}

func resolveColoursFromArgs(colours []string) ([]FilamentColour, error) {
	resolved := make([]FilamentColour, len(colours))
	for i, col := range colours {
		hexCode, err := parseColour(col)
		if err != nil {
			return nil, fmt.Errorf("invalid colour '%s': %w", col, err)
		}
		resolved[i] = FilamentColour{HexCode: hexCode}
	}
	return resolved, nil
}
