# go_mmu_calc

A small command-line tool that calculates **flush (purge) volumes** for multi-material
3D printing (MMU / AMS-style setups). Give it a list of filament colours and it prints
the full matrix of flush volumes needed when switching from one colour to another.

The calculation replicates the colour-based flush volume algorithm used by
[OrcaSlicer](https://github.com/SoftFever/OrcaSlicer), so the results can be pasted
straight into your slicer's flushing volume table.

<p align="center">
  <img src="docs/screenshot.png" alt="go_mmu_calc output: colour list with swatches and a 5x5 flush-volume matrix" width="423">
</p>

## Features

- Pairwise flush-volume matrix for up to 10 filaments
- Colours by name (`red`, `navy`, `darkorange`, …) or hex (`#ff8800`, `#ff8800ff`)
- Adjustable flush multiplier
- Coloured swatches in the terminal (24-bit colour), with a `--no-color` fallback
- Interactive mode when no colours are given
- No dependencies outside the Go standard library

## Installation

Requires Go 1.27 or newer.

```sh
go install github.com/mich-w/go_mmu_calc@latest
```

Or build from source:

```sh
git clone https://github.com/mich-w/go_mmu_calc.git
cd go_mmu_calc
go build
```

## Usage

```
go_mmu_calc [OPTIONS] [COLOURS...]
```

| Option                     | Description                                    |
| -------------------------- | ---------------------------------------------- |
| `-l`, `--list-colors`      | List all available colour names (alphabetical) |
| `-lh`, `--list-colors-hue` | List all available colour names (by hue)       |
| `-m`, `--multiplier FLOAT` | Flush volume multiplier (default: `1.0`)       |
| `--no-color`               | Disable coloured terminal output               |
| `-h`, `--help`             | Show help                                      |

If no colours are passed, the tool asks for the number of filaments, each colour,
and the multiplier interactively.

### Examples

```sh
go_mmu_calc white black red "#00ff00"
go_mmu_calc "#ff0000" "#0000ff" -m 1.5
go_mmu_calc -lh
```

Example output:

```
Colours (4):
 0  white (#ffffff)
 1  black (#000000)
 2  red   (#ff0000)
 3  lime  (#00ff00)
Multiplier: 1.0

Flush-volume matrix (4x4):
+------+------+------+------+------+
|      | 0    | 1    | 2    | 3    |
+------+------+------+------+------+
| 0    |    0 |  103 |  262 |  248 |
| 1    |  559 |    0 |  408 |  540 |
| 2    |  586 |   90 |    0 |  443 |
| 3    |  460 |  107 |  242 |    0 |
+------+------+------+------+------+
```

Rows are the colour being flushed **from**, columns the colour being switched **to**,
and values are in mm³. Results are clamped to the range 60–20000.

## How it works

For each pair of colours the tool:

1. Converts both colours to HSV and measures their hue/saturation distance.
2. Compares their luminance — going from dark to light needs much more flushing
   than going from light to dark.
3. Combines both factors, adds an extra 30% when switching from a light colour to
   a very dark one, and applies the multiplier.

## Running tests

```sh
go test ./...
```

## License

[MIT](LICENSE)
