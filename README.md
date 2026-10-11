# Vedic Calendar

[![Go 1.25](https://img.shields.io/badge/Go-1.25-blue?logo=go)](https://go.dev)
[![license MIT](https://img.shields.io/badge/license-MIT-green)](LICENSE)
[![built with LuaLaTeX](https://img.shields.io/badge/built%20with-LuaLaTeX-orange)](https://www.luatex.org)

A Go + LuaLaTeX pipeline turning GCal-exported Vedic calendar text into a
printable A4 landscape PDF, one page per lunar month. Drop a plaintext
export in `data/`, run `make`, collect the PDF — festivals, Ekadasi fasts,
appearances/disappearances, and sankrantis on a Sunday-start grid with footnotes.

## 30-second example

Each input is a GCal-exported plaintext calendar (`data/*.txt`):

```
                       Pausa (Narayana) Masa, Gaurabda 540      GCal 11, Build 5

                      Lisbon [Portugal] (38N43 9W08, +0:00)

 DATE            TITHI                         PAKSA NAKSATRA       FAST
--------------------------------------------------------------------------------
 1 Jan 2027 Fr   Dasami                            K Citra
 2 Jan 2027 Sa   Dasami                            K Swati
 3 Jan 2027 Su   Ekadasi (suitable for fasting)    K Visakha         *
                 Fasting for Saphala Ekadasi
```

```sh
make check                                            # verify go + lualatex are installed
make FILE="data/Lisbon [Portugal]-2027.txt"           # text → JSON → LaTeX → PDF
# → output/Lisbon_Portugal_2027/calendar-Lisbon_Portugal_2027.pdf (13 pages)
```

## Sample output

| Pausa (Narayana) — Jan | Magha (Madhava) — Jan/Feb | Phalguna (Govinda) — Feb/Mar |
|---|---|---|
| ![Pausa month grid](docs/img/preview-01.png) | ![Magha month grid with footnotes](docs/img/preview-02.png) | ![Phalguna month grid with footnotes](docs/img/preview-03.png) |

## Prerequisites

- Go 1.25+ (`go`) — builds the `parser`, `renderer`, and `merger` binaries
- LuaLaTeX (`lualatex`, TeX Live) — compiles the generated `.tex` into the PDF
  (TikZ grids, footnotes, font loading)
- `Charis SIL` font (`fonts-sil-charis`) — document main font; covers the
  transliteration diacritics (ā, ī, ṭh, ṅ, ñ, etc.)
- `fontawesome5` LaTeX package (TeX Live: `texlive-fonts-extra`) — grid icons
  (`\faMoon`, `\faStar`, `\faFire`, `\faDrum`, `\faSun`, `\faCrown`)

```sh
sudo apt install texlive-luatex texlive-latex-extra texlive-fonts-extra fonts-sil-charis
```

## Commands

```sh
make                                       # build PDF from the most recently modified data/*.txt
make FILE="data/Lisbon [Portugal]-2027.txt"  # build PDF from a specific input file
make parse FILE="..."                      # text → JSON only (output/<name>/calendar.json)
make render FILE="..."                     # JSON → LaTeX only (output/<name>/calendar-<name>.tex)
make pdf FILE="..."                        # LaTeX → PDF only (output/<name>/calendar-<name>.pdf)
make clean                                 # rm -rf output/ and the parser/renderer/merger binaries
make check                                 # verify go + lualatex are installed
```

Tests: `go test ./...`

## Project layout

```
cmd/parser/main.go        # .txt + festivals.json → calendar.json
cmd/merger/main.go        # + data/extra-*.json → merged.json
cmd/renderer/main.go      # + templates/ → calendar-<name>.tex
internal/models/          # shared structs + JSON tags
internal/parser/          # line classification → model
internal/merger/          # Vedic-date lookup, injects extra events
internal/renderer/        # week grids, footnotes, template execution
templates/calendar.tex.tmpl  # TikZ grid, 3-column footnotes, legend
data/*.txt                # sample inputs (one city per file)
data/festivals.json       # Gaudiya Vaisnava festivals (major/minor)
data/extra-*.json         # optional extra events (auto-globbed)
scripts/check-deps.sh     # behind `make check`
docs/img/                 # screenshots from real output (150 dpi)
output/                   # generated JSON/LaTeX/PDF (gitignored)
```

## Extra events (`data/extra-*.json`)

Any `data/extra-*.json` is auto-merged before rendering, matched by Vedic
date (`masa` + `paksa` + `tithi`, qualifier ignored). Unmatched entries are
reported on stderr and skipped.

```json
[
  {
    "masa": "Pausa",
    "paksa": "shukla",
    "tithi": "Purnima",
    "type": "appearance",
    "description": "Example Devotee -- Appearance",
    "highlight": true
  }
]
```

| Field | Meaning |
|---|---|
| `masa` | Vedic month, e.g. `Pausa`, `Magha` (must match input headers) |
| `paksa` | `shukla`/`krishna` (or `G`/`K`) |
| `tithi` | base tithi, e.g. `Purnima`, `Ekadasi`, `Dvadasi` |
| `type` | `appearance`, `disappearance`, `festival`, `fast`, or `note` |
| `description` | footnote/cell text, rendered as-is |
| `highlight` | `true` bolds the footnote and adds a crown watermark to the cell |

## Architecture

- **Parse** — `cmd/parser`: `.txt` → `calendar.json` (line classification;
  `festivals.json` substring match sets `festival` vs `note`, majors get `Highlight`).
- **Merge** — `cmd/merger`: `+ data/extra-*.json` → `merged.json`
  (`masa|paksa|tithi` lookup).
- **Render** — `cmd/renderer`: days → Sunday-start grids, LaTeX icons,
  per-page footnotes → executes the template.
- **Compile** — `lualatex` runs in `output/<name>/`, keeping `.aux`/`.log`
  with the `.tex`.
- **Contract** — JSON between stages; each binary runs standalone
  (`make parse` / `make render` / `make pdf`).

## License

MIT — see [LICENSE](LICENSE).
