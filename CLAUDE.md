# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Vedic calendar PDF generator. Converts GCal-exported plaintext Vedic calendar files into A4 landscape PDFs with one page per Vedic Masa (lunar month).

## Pipeline

```
*.txt → cmd/parser → JSON → cmd/renderer → .tex → lualatex → .pdf
```

The JSON is the stable contract between parser and renderer. Each stage is independent.

## Build & Run

```bash
make all                          # full pipeline (parse → render → compile PDF)
make all INPUT="Tokyo [Japan].txt"  # different input file
make parse                        # text → JSON only
make render                       # JSON → .tex only
make pdf                          # .tex → PDF only
make clean                        # remove output/ and binaries
make check                        # verify go + lualatex are installed
```

## Tests

```bash
go test ./...                              # all tests
go test ./internal/parser/ -v              # parser tests only
go test ./internal/renderer/ -v            # renderer tests only
go test ./internal/parser/ -run TestParseDayEntry -v  # single test
```

No external dependencies — stdlib only.

## Architecture

**`internal/models/`** — Shared structs (`Calendar`, `VedicMonth`, `Day`, `Event`, `BreakFast`, `Sankranti`). Used by both parser and renderer. JSON serialization tags define the intermediate format.

**`internal/parser/`** — Reads plaintext, classifies lines by type (day entry, continuation, masa header, sankranti, etc.), then walks them to build the model. Key functions: `classifyLine`, `parseDayEntry`, `classifyEvent`, `parseSankranti`, `ParseCalendar`. Loads `data/festivals.json` to distinguish festivals from notes via substring matching.

**`internal/renderer/`** — Converts the model into grid data (`MonthGrid` with `[7]*CellData` week rows), assigns icons (LaTeX amssymb commands, not Unicode), numbers footnotes sequentially per page, then executes Go `text/template` against the `.tex.tmpl` files.

**`templates/`** — LuaLaTeX templates using TikZ for the calendar grid. `calendar.tex.tmpl` produces all month pages (does NOT include `\end{document}`). `reference.tex.tmpl` produces the abbreviation/legend page. The renderer appends `\end{document}` after both.

**`data/festivals.json`** — Curated list of recurring Gaudiya Vaishnava festival names. Parser uses substring matching against this list to classify events as `festival` vs `note`.

## Key Design Decisions

- Icons are LaTeX commands (`$\blacklozenge$`, `$\bigstar$`, `\dag`, `$\odot$`, `$\blacktriangle$`), not Unicode — Latin Modern Roman lacks the Unicode glyphs.
- Fasting days use grey cell background, not an icon.
- Sankranti lines in the source appear between days but contain an embedded date (`on DD Mon`). The parser extracts this date and attaches the Sankranti to the matching Day struct.
- The day entry regex handles edge cases like zero space before paksa character (`fasting)G`) by using `\s*` instead of `\s{2,}` for the paksa separator.
- Footnotes are 3-column `\scriptsize` to fit pages with many events.
- Row height is 1.8cm (reduced from 2.4cm to leave room for footnotes).
