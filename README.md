# Vedic Calendar PDF Generator

Converts a GCal-exported Vedic calendar text file into a printable A4 landscape PDF calendar, one page per Vedic Masa (lunar month).

## Requirements

- Go 1.21+
- LuaLaTeX (TeX Live recommended)

Check dependencies:

    make check

## Usage

    make all

This runs the full pipeline:

1. **Parse** — `Vrindavan [India].txt` → `output/calendar.json`
2. **Render** — `output/calendar.json` → `output/calendar.tex`
3. **Compile** — `output/calendar.tex` → `output/calendar.pdf`

## Custom Input

    make all INPUT="path/to/other-calendar.txt"

## Individual Stages

    make parse    # text → JSON
    make render   # JSON → LaTeX
    make pdf      # LaTeX → PDF
    make clean    # remove all generated files

## Project Structure

    cmd/parser/         Go parser CLI
    cmd/renderer/       Go renderer CLI
    internal/models/    Shared data structures
    internal/parser/    Text parsing logic
    internal/renderer/  Grid layout and template execution
    templates/          LuaLaTeX templates
    data/               Festival list and input data
    output/             Generated files
    scripts/            Helper scripts

## Calendar Features

- One page per Vedic Masa (lunar month)
- Sunday-start week grid
- Icons: ◆ Ekadasi, ★ Appearance, † Disappearance, ⊙ Festival, ▲ Sankranti
- Fasting days highlighted with grey background
- Footnotes with full event descriptions (3-column layout)
- Reference page with all abbreviations, icons, and metadata
