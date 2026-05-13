# Vedic Calendar PDF Pipeline — Design Spec

## Overview

A reusable pipeline that converts a GCal-exported Vedic calendar text file into a professionally formatted, printable PDF calendar.

**Pipeline:** `plaintext → Go parser → JSON → Go renderer → LuaLaTeX → PDF`

**Input:** Plaintext file exported from GCal (e.g., `Vrindavan [India].txt`)
**Output:** A4 landscape PDF, one page per Vedic Masa, with a reference page at the end

## Architecture & Data Flow

```
Vrindavan [India].txt
        |
        v
+-------------------+
|  cmd/parser       |  Go binary: reads .txt, classifies lines,
|                   |  builds structured data, writes JSON
+--------+----------+
         |  calendar.json
         v
+-------------------+
|  cmd/renderer     |  Go binary: reads JSON, assigns icons,
|                   |  numbers footnotes, executes .tex.tmpl
+--------+----------+
         |  calendar.tex
         v
+-------------------+
|   lualatex        |  Compiles .tex -> calendar.pdf
+--------+----------+
         |
         v
      output/calendar.pdf
```

JSON is the stable contract between parser and renderer. Each stage is independently testable and rerunnable. The Makefile tracks file modification times so it only reruns stages whose inputs have changed.

## Data Model

```go
type Calendar struct {
    Location Location
    Months   []VedicMonth
}

type Location struct {
    Name     string // "Vrindavan [India]"
    Coords   string // "27N35 77E42"
    Timezone string // "+5:30"
}

type VedicMonth struct {
    MasaName     string // "Magha"
    MasaAltName  string // "Madhava"
    GaurabdaYear int    // 539
    Days         []Day
}

type Day struct {
    Date      time.Time
    Weekday   string     // "Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"
    Tithi     string     // "Ekadasi (suitable for fasting)"
    Paksa     string     // "G" or "K"
    Naksatra  string     // "Rohini"
    Fasting   bool       // true if * marker present
    Events    []Event
    BreakFast *BreakFast
    Sankranti *Sankranti
}

type Event struct {
    Type        EventType
    Description string
}

type EventType string

const (
    EventAppearance    EventType = "appearance"
    EventDisappearance EventType = "disappearance"
    EventFestival      EventType = "festival"
    EventFast          EventType = "fast"
    EventNote          EventType = "note"
)

type BreakFast struct {
    Start string // "07:09"
    End   string // "10:38"
    Notes string // "1/3 of daylight LT"
}

type Sankranti struct {
    Name  string // "Makara Sankranti"
    Notes string // "Sun enters Capricorn on 14 Jan, 14:54 LT"
}
```

## Event Classification

Events are classified by pattern matching on continuation line text:

| Pattern | EventType |
|---|---|
| `-- Appearance` | `appearance` |
| `-- Disappearance` | `disappearance` |
| `Break fast ...` | BreakFast struct (not an Event) |
| `Fasting for ...` | `fast` |
| `(Fast till ...)` | `fast` |
| Match in `data/festivals.json` | `festival` |
| Everything else | `note` |

A curated `data/festivals.json` file ships with the project, containing all known recurring Gaudiya Vaishnava festivals. This list is used at parse time to classify events as festivals vs. notes.

### Curated Festival List

Festivals (get festival icon + footnote):

- Sri Krsna Pusya Abhiseka
- Vasanta Pancami
- Sarasvati Puja
- Bhismastami
- Varaha Dvadasi
- Nityananda Trayodasi
- Sri Krsna Madhura Utsava
- Gaura Purnima
- Festival of Jagannatha Misra
- Rama Navami
- Aksaya Trtiya
- Candana Yatra
- Jahnu Saptami
- Nrsimha Caturdasi
- Krsna Phula Dola, Salila Vihara
- Rukmini Dvadasi
- Damanakaropana Dvadasi
- Sri Balarama Rasayatra
- Sri Krsna Vasanta Rasa
- Appearance of Radha Kunda
- Snana Yatra
- Panihati Cida Dahi Utsava
- Ganga Puja
- Ratha Yatra
- Gundica Marjana
- Hera Pancami
- Return Ratha
- Guru (Vyasa) Purnima
- Nandotsava
- Sri Krsna Janmastami
- Radhastami
- Sri Vamana Dvadasi
- Ananta Caturdasi Vrata
- Sri Visvarupa Mahotsava
- Bhadra Purnima
- Durga Puja
- Ramacandra Vijayotsava
- Sri Krsna Saradiya Rasayatra
- Laksmi Puja
- Dipa dana, Dipavali
- Go Puja, Govardhana Puja
- Bali Daityaraja Puja
- Gopastami, Gosthastami
- Jagaddhatri Puja
- Bahulastami
- Sri Krsna Rasayatra
- Tulasi-Saligrama Vivaha
- Ganga Sagara Mela
- Odana Sasthi
- Siva Ratri

Notes (footnote text only, no icon):

- Caturmasya month begins/ends markers
- Tulasi Jala Dan begins/ends
- Katyayani vrata begins/ends
- Jhulana Yatra begins/ends
- Radha Govinda Jhulana Yatra begins
- Fasting instruction qualifiers: "(Fasting is done yesterday...)"
- ISKCON historical notes (incorporation, Prabhupada's departure/arrival)
- Paksa vardhini / Vyanjuli Mahadvadasi labels
- Bhisma Pancaka first/last day markers

## Parser Design

### Line Classification (First Pass)

| Line Type | Detection Rule |
|---|---|
| BLANK | Empty or whitespace only |
| SEPARATOR | `^-{10,}$` (pure dashes, no content) |
| MASA_HEADER | Contains `Masa, Gaurabda` |
| LOCATION | Contains lat/long pattern like `(27N35 77E42` |
| TABLE_HEADER | Starts with `DATE` or contains `TITHI` and `PAKSA` |
| SANKRANTI | `^-+\s+\w.*Sankranti.*-+$` (dashes with content) |
| DAY_ENTRY | `^\s?\d{1,2}\s+\w{3}\s+\d{4}` |
| CONTINUATION | Starts with 17+ spaces of indentation |

### Day Entry Regex

```
^\s?(\d{1,2})\s+(\w{3})\s+(\d{4})\s+(\w{2})\s{3}(.+?)\s{2,}(G|K)\s+(\S+)\s*(\*?)\s*$
```

Fields: day, month, year, weekday, tithi, paksa, naksatra, fasting marker.

### Parsing Stages

1. Read all lines, classify each line by type
2. Walk lines top-to-bottom:
   - On MASA_HEADER: parse masa name, alt name, gaurabda year; open new VedicMonth
   - On LOCATION: parse and store location (first occurrence only)
   - On DAY_ENTRY: parse all fields; open new Day; append to current VedicMonth
   - On CONTINUATION: classify event type; attach to current Day
   - On SANKRANTI: parse date from text ("on DD Mon"); find matching Day by date; attach Sankranti struct
   - Skip BLANK, SEPARATOR, TABLE_HEADER

### Sankranti Date Extraction

Sankranti lines embed the actual event date:
```
---------- Makara Sankranti (Sun enters Capricorn on 14 Jan, 14:54 LT) ----------
```
Parser extracts "14 Jan" via regex `on\s+(\d{1,2})\s+(\w{3})` and attaches the Sankranti to the Day with that matching date.

### Output

The parser writes a single JSON file containing the full `Calendar` struct.

## Renderer & Calendar Grid Layout

### Page Layout (A4 Landscape)

```
+--------------------------------------------------------------+
|  Magha (Madhava) Masa                    Gaurabda 539        |  <- Header: 12mm
|  Vrindavan [India]                                           |
+----+----+----+----+----+----+--------------------------------+
| Sun| Mon| Tue| Wed| Thu| Fri| Sat                            |  <- Weekday bar: 6mm
+----+----+----+----+----+----+----+
| 4  | 5  | 6  | 7  | 8  | 9  |10  |
|K-Pr|K-Dv|K-Tr|K-Pa|K-Sa|K-Sp|K-Sp|  <- Row ~28mm
|Pnv |Psy |Asl |Mag |PPh |UPh |Has |
|    |    |    |**  | +  | +  |    |
|    |    |    |1,2 |3   |4   |    |
+----+----+----+----+----+----+----+
|... |    |    |    |    |    |    |
|    |    |    |    |    |    |    |  <- 5-6 rows
+----+----+----+----+----+----+----+
| 1 Sri Ramacandra Kaviraja App.   |  <- Footnotes: ~30mm
| 2 Srila Gopala Bhatta Gosvami App|
| 3 Sri Jayadeva Gosvami Dis.      |
| 4 Sri Locana Dasa Thakura Dis.   |
| Legend: * Appearance  + Dis. ... |  <- Icon legend
+----------------------------------+
```

### Cell Content (Top to Bottom)

1. Day number — large, bold, top-left
2. Paksa + Tithi abbreviation — e.g. `K-Eka` (small)
3. Naksatra abbreviation — e.g. `Anu` (smaller)
4. Icons + footnote superscripts — bottom of cell

### Icon Set (Monochrome Unicode)

| Icon | Meaning |
|------|---------|
| diamond | Ekadasi |
| star | Appearance |
| dagger | Disappearance |
| sun | Festival |
| triangle | Sankranti |

Fasting days (those with the `*` marker) are indicated by a **light grey background fill** on the entire cell rather than an icon. This provides a clear visual distinction at a glance without adding clutter to the cell content.

### Tithi Abbreviations

| Full | Abbrev |
|------|--------|
| Pratipat | Pra |
| Dvitiya | Dvi |
| Tritiya | Tri |
| Caturthi | Cat |
| Pancami | Pan |
| Sasti | Sas |
| Saptami | Sap |
| Astami | Ast |
| Navami | Nav |
| Dasami | Das |
| Ekadasi | Eka |
| Dvadasi | Dva |
| Trayodasi | Tra |
| Caturdasi | Cad |
| Purnima | Pur |
| Amavasya | Ama |

### Naksatra Abbreviations

| Full | Abbrev |
|------|--------|
| Asvini | Asv |
| Bharani | Bha |
| Krittika | Kri |
| Rohini | Roh |
| Mrigasira | Mri |
| Ardra | Ard |
| Punarvasu | Pun |
| Pusyami | Pus |
| Aslesa | Asl |
| Magha | Mag |
| Purva-phalguni | PPh |
| Uttara-phalguni | UPh |
| Hasta | Has |
| Citra | Cit |
| Swati | Swa |
| Visakha | Vis |
| Anuradha | Anu |
| Jyestha | Jye |
| Mula | Mul |
| Purva-asadha | PAs |
| Uttara-asadha | UAs |
| Sravana | Sra |
| Dhanista | Dha |
| Satabhisa | Sat |
| Purva-bhadra | PBh |
| Uttara-bhadra | UBh |
| Revati | Rev |

### Footnotes

- Numbered sequentially per page (resets each masa)
- Appearance/Disappearance shortened: "App." / "Dis."
- Break-fast info included: "Break fast 07:09-10:38 LT"
- Fasting instructions included: "Fast till noon"
- Sankranti noted: "Makara Sankranti - Sun enters Capricorn, 14:54 LT"

### Edge Cases

- Masas don't start on Sundays: first/last rows have empty cells (greyed out)
- Days with many events: all events become footnotes; only icons shown in cell
- Ekadasi marked "not suitable for fasting": no fasting icon, but Ekadasi diamond still shown

## Reference Page

The final page of the PDF contains:

- Tithi abbreviations — full table of all 16 tithis and their 3-letter codes
- Naksatra abbreviations — full table of all 27 naksatras and their short forms
- Icon legend — all 6 icons with meanings
- Paksa key — G = Gaura (bright/waxing), K = Krsna (dark/waning)
- Masa list — all Vedic months with alternate names
- Location metadata — coordinates, timezone
- Notes from source file — DST/LT explanations
- Source credit — "Generated from GCal 11, Build 5"

## Project Structure

```
calendar-generator/
├── cmd/
│   ├── parser/main.go
│   └── renderer/main.go
├── internal/
│   ├── models/models.go
│   ├── parser/parser.go
│   ├── parser/parser_test.go
│   ├── renderer/renderer.go
│   └── renderer/icons.go
├── data/
│   ├── festivals.json
│   └── Vrindavan [India].txt
├── templates/
│   ├── calendar.tex.tmpl
│   └── reference.tex.tmpl
├── output/
├── scripts/
│   └── check-deps.sh
├── Makefile
└── README.md
```

## Automation

### Makefile Targets

| Target | Action |
|--------|--------|
| `make all` | Full pipeline: parse -> render -> pdf |
| `make parse` | Run parser, produce `output/calendar.json` |
| `make render` | Run renderer, produce `output/calendar.tex` |
| `make pdf` | Compile `output/calendar.tex` -> `output/calendar.pdf` |
| `make clean` | Remove everything in `output/` |
| `make check` | Verify dependencies are installed |

### Variables

- `INPUT ?= data/Vrindavan [India].txt`
- `OUTPUT_DIR ?= output`

### Dependency Chain

```makefile
$(OUTPUT_DIR)/calendar.json: $(INPUT)
$(OUTPUT_DIR)/calendar.tex: $(OUTPUT_DIR)/calendar.json templates/*.tex.tmpl data/festivals.json
$(OUTPUT_DIR)/calendar.pdf: $(OUTPUT_DIR)/calendar.tex
```

## Technology

- Go (parser + renderer)
- LuaLaTeX (PDF compilation, Unicode-safe)
- Go `text/template` (template engine for .tex generation)
- TikZ (grid drawing in LaTeX)
- fontspec (Unicode font support in LuaLaTeX)
