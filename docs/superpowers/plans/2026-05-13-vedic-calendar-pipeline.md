# Vedic Calendar PDF Pipeline — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a pipeline that converts a GCal Vedic calendar text export into a printable A4 landscape PDF with one page per Vedic Masa.

**Architecture:** Go parser reads plaintext and outputs structured JSON. Go renderer reads JSON, applies LuaLaTeX templates with TikZ grids, and outputs `.tex`. LuaLaTeX compiles the final PDF. Makefile orchestrates the chain.

**Tech Stack:** Go 1.21+, LuaLaTeX, TikZ, fontspec, Go `text/template`

**Spec:** `docs/superpowers/specs/2026-05-13-vedic-calendar-pipeline-design.md`

---

## File Structure

```
cmd/parser/main.go          — CLI entry: reads .txt path from args, writes JSON to stdout or file
cmd/renderer/main.go        — CLI entry: reads JSON path from args, writes .tex to file
internal/models/models.go   — All shared structs (Calendar, VedicMonth, Day, Event, etc.)
internal/parser/parser.go   — Line classification, day entry parsing, event attachment
internal/parser/parser_test.go — Unit tests for parser
internal/renderer/renderer.go  — Template loading, footnote numbering, grid data prep
internal/renderer/abbreviations.go — Tithi/Naksatra abbreviation maps
internal/renderer/renderer_test.go — Unit tests for renderer
data/festivals.json         — Curated festival classification list
templates/calendar.tex.tmpl — LuaLaTeX template for monthly grid pages
templates/reference.tex.tmpl — LuaLaTeX template for the reference/legend page
scripts/check-deps.sh       — Dependency checker (go, lualatex)
Makefile                    — Build automation
go.mod                      — Go module definition
```

---

### Task 1: Project scaffold and Go module

**Files:**
- Create: `go.mod`
- Create: `cmd/parser/main.go`
- Create: `cmd/renderer/main.go`
- Create: `data/festivals.json`

- [ ] **Step 1: Initialize Go module**

Run:
```bash
cd /home/taraka/Documents/thomas/calander
go mod init github.com/taraka/vedic-calendar
```
Expected: `go.mod` created

- [ ] **Step 2: Create stub parser CLI**

Create `cmd/parser/main.go`:

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: parser <input.txt> [output.json]\n")
		os.Exit(1)
	}
	fmt.Println("parser: not yet implemented")
}
```

- [ ] **Step 3: Create stub renderer CLI**

Create `cmd/renderer/main.go`:

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: renderer <input.json> [output.tex]\n")
		os.Exit(1)
	}
	fmt.Println("renderer: not yet implemented")
}
```

- [ ] **Step 4: Create festivals.json**

Create `data/festivals.json`:

```json
{
  "festivals": [
    "Sri Krsna Pusya Abhiseka",
    "Vasanta Pancami",
    "Sarasvati Puja",
    "Bhismastami",
    "Varaha Dvadasi",
    "Nityananda Trayodasi",
    "Sri Krsna Madhura Utsava",
    "Gaura Purnima",
    "Festival of Jagannatha Misra",
    "Rama Navami",
    "Aksaya Trtiya",
    "Candana Yatra",
    "Jahnu Saptami",
    "Nrsimha Caturdasi",
    "Krsna Phula Dola, Salila Vihara",
    "Rukmini Dvadasi",
    "Damanakaropana Dvadasi",
    "Sri Balarama Rasayatra",
    "Sri Krsna Vasanta Rasa",
    "Appearance of Radha Kunda",
    "Snana Yatra",
    "Panihati Cida Dahi Utsava",
    "Ganga Puja",
    "Ratha Yatra",
    "Gundica Marjana",
    "Hera Pancami",
    "Return Ratha",
    "Guru (Vyasa) Purnima",
    "Nandotsava",
    "Sri Krsna Janmastami",
    "Radhastami",
    "Sri Vamana Dvadasi",
    "Ananta Caturdasi Vrata",
    "Sri Visvarupa Mahotsava",
    "Bhadra Purnima",
    "Durga Puja",
    "Ramacandra Vijayotsava",
    "Sri Krsna Saradiya Rasayatra",
    "Laksmi Puja",
    "Dipa dana, Dipavali",
    "Go Puja, Govardhana Puja",
    "Bali Daityaraja Puja",
    "Gopastami, Gosthastami",
    "Jagaddhatri Puja",
    "Bahulastami",
    "Sri Krsna Rasayatra",
    "Tulasi-Saligrama Vivaha",
    "Ganga Sagara Mela",
    "Odana Sasthi",
    "Siva Ratri",
    "Krsna Phula Dola",
    "Salila Vihara",
    "Go Puja",
    "Govardhana Puja",
    "Gopastami",
    "Gosthastami",
    "Dipa dana",
    "Dipavali",
    "Kali Puja",
    "Bali Daityaraja Puja"
  ]
}
```

Note: some festivals appear as comma-separated pairs in the source text (e.g., "Krsna Phula Dola, Salila Vihara"). We include both the combined and individual forms so substring matching catches them.

- [ ] **Step 5: Verify the stubs build**

Run:
```bash
go build ./cmd/parser && go build ./cmd/renderer
```
Expected: both binaries compile without error

- [ ] **Step 6: Commit**

```bash
git add go.mod cmd/ data/festivals.json
git commit -m "scaffold: Go module, stub CLIs, curated festival list"
```

---

### Task 2: Data model

**Files:**
- Create: `internal/models/models.go`

- [ ] **Step 1: Write the models file**

Create `internal/models/models.go`:

```go
package models

import "time"

type Calendar struct {
	Location Location     `json:"location"`
	Months   []VedicMonth `json:"months"`
}

type Location struct {
	Name     string `json:"name"`
	Coords   string `json:"coords"`
	Timezone string `json:"timezone"`
}

type VedicMonth struct {
	MasaName     string `json:"masa_name"`
	MasaAltName  string `json:"masa_alt_name"`
	GaurabdaYear int    `json:"gaurabda_year"`
	Days         []Day  `json:"days"`
}

type Day struct {
	Date      time.Time  `json:"date"`
	Weekday   string     `json:"weekday"`
	Tithi     string     `json:"tithi"`
	Paksa     string     `json:"paksa"`
	Naksatra  string     `json:"naksatra"`
	Fasting   bool       `json:"fasting"`
	Events    []Event    `json:"events,omitempty"`
	BreakFast *BreakFast `json:"break_fast,omitempty"`
	Sankranti *Sankranti `json:"sankranti,omitempty"`
}

type Event struct {
	Type        EventType `json:"type"`
	Description string    `json:"description"`
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
	Start string `json:"start"`
	End   string `json:"end"`
	Notes string `json:"notes"`
}

type Sankranti struct {
	Name  string `json:"name"`
	Notes string `json:"notes"`
}
```

- [ ] **Step 2: Verify it compiles**

Run:
```bash
go build ./internal/models/
```
Expected: compiles without error

- [ ] **Step 3: Commit**

```bash
git add internal/models/
git commit -m "feat: add data model structs"
```

---

### Task 3: Parser — line classification

**Files:**
- Create: `internal/parser/parser.go`
- Create: `internal/parser/parser_test.go`

- [ ] **Step 1: Write failing test for line classification**

Create `internal/parser/parser_test.go`:

```go
package parser

import "testing"

func TestClassifyLine(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected LineType
	}{
		{"blank", "", LineBlank},
		{"blank spaces", "     ", LineBlank},
		{"separator", "--------------------------------------------------------------------------------", LineSeparator},
		{"masa header", "                       Pausa (Narayana) Masa, Gaurabda 539      GCal 11, Build 5", LineMasaHeader},
		{"location", "                      Vrindavan [India] (27N35 77E42, +5:30)", LineLocation},
		{"table header", " DATE            TITHI                         PAKSA NAKSATRA       FAST ", LineTableHeader},
		{"day entry", " 1 Jan 2026 Th   Trayodasi                         G Rohini", LineDayEntry},
		{"day entry 2digit", "14 Jan 2026 We   Ekadasi (suitable for fasting)    K Anuradha        *", LineDayEntry},
		{"continuation", "                 Break fast 07:09 (sunrise) - 10:38 (1/3 of daylight) LT", LineContinuation},
		{"continuation event", "                 Sri Krsna Pusya Abhiseka", LineContinuation},
		{"sankranti", "---------- Makara Sankranti (Sun enters Capricorn on 14 Jan, 14:54 LT) ----------", LineSankranti},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := classifyLine(tt.line)
			if got != tt.expected {
				t.Errorf("classifyLine(%q) = %v, want %v", tt.line, got, tt.expected)
			}
		})
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run:
```bash
go test ./internal/parser/ -v -run TestClassifyLine
```
Expected: FAIL — `classifyLine` not defined

- [ ] **Step 3: Implement line classification**

Create `internal/parser/parser.go`:

```go
package parser

import (
	"regexp"
	"strings"
)

type LineType int

const (
	LineBlank LineType = iota
	LineSeparator
	LineMasaHeader
	LineLocation
	LineTableHeader
	LineDayEntry
	LineContinuation
	LineSankranti
)

var (
	reDayEntry  = regexp.MustCompile(`^\s?\d{1,2}\s+\w{3}\s+\d{4}`)
	reSankranti = regexp.MustCompile(`^-+\s+.*Sankranti.*-+$`)
	reSeparator = regexp.MustCompile(`^-{10,}$`)
)

func classifyLine(line string) LineType {
	trimmed := strings.TrimSpace(line)

	if trimmed == "" {
		return LineBlank
	}

	if reSankranti.MatchString(trimmed) {
		return LineSankranti
	}

	if reSeparator.MatchString(trimmed) {
		return LineSeparator
	}

	if strings.Contains(trimmed, "Masa, Gaurabda") || strings.Contains(trimmed, "Masa,  Gaurabda") {
		return LineMasaHeader
	}

	if regexp.MustCompile(`\d+N\d+\s+\d+E\d+`).MatchString(trimmed) {
		return LineLocation
	}

	if strings.HasPrefix(trimmed, "DATE") || (strings.Contains(trimmed, "TITHI") && strings.Contains(trimmed, "PAKSA")) {
		return LineTableHeader
	}

	if reDayEntry.MatchString(line) {
		return LineDayEntry
	}

	if len(line) > 16 && strings.TrimSpace(line[:16]) == "" && trimmed != "" {
		return LineContinuation
	}

	return LineBlank
}
```

- [ ] **Step 4: Run test to verify it passes**

Run:
```bash
go test ./internal/parser/ -v -run TestClassifyLine
```
Expected: PASS — all 11 cases

- [ ] **Step 5: Commit**

```bash
git add internal/parser/
git commit -m "feat: parser line classification with tests"
```

---

### Task 4: Parser — masa header and location parsing

**Files:**
- Modify: `internal/parser/parser.go`
- Modify: `internal/parser/parser_test.go`

- [ ] **Step 1: Write failing tests for masa header and location parsing**

Add to `internal/parser/parser_test.go`:

```go
func TestParseMasaHeader(t *testing.T) {
	tests := []struct {
		line     string
		name     string
		altName  string
		year     int
	}{
		{"                       Pausa (Narayana) Masa, Gaurabda 539      GCal 11, Build 5", "Pausa", "Narayana", 539},
		{"                        Magha (Madhava) Masa, Gaurabda 539      GCal 11, Build 5", "Magha", "Madhava", 539},
		{"                        Caitra (Visnu) Masa, Gaurabda 540       GCal 11, Build 5", "Caitra", "Visnu", 540},
		{"                      Purusottama-adhika Masa, Gaurabda 540     GCal 11, Build 5", "Purusottama-adhika", "", 540},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, altName, year := parseMasaHeader(tt.line)
			if name != tt.name {
				t.Errorf("name = %q, want %q", name, tt.name)
			}
			if altName != tt.altName {
				t.Errorf("altName = %q, want %q", altName, tt.altName)
			}
			if year != tt.year {
				t.Errorf("year = %d, want %d", year, tt.year)
			}
		})
	}
}

func TestParseLocation(t *testing.T) {
	line := "                      Vrindavan [India] (27N35 77E42, +5:30)"
	name, coords, tz := parseLocation(line)
	if name != "Vrindavan [India]" {
		t.Errorf("name = %q, want %q", name, "Vrindavan [India]")
	}
	if coords != "27N35 77E42" {
		t.Errorf("coords = %q, want %q", coords, "27N35 77E42")
	}
	if tz != "+5:30" {
		t.Errorf("tz = %q, want %q", tz, "+5:30")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run:
```bash
go test ./internal/parser/ -v -run "TestParseMasaHeader|TestParseLocation"
```
Expected: FAIL — functions not defined

- [ ] **Step 3: Implement masa header and location parsing**

Add to `internal/parser/parser.go`:

```go
var (
	reMasaHeader = regexp.MustCompile(`(\S+(?:-\S+)?)\s+(?:\((\w+)\)\s+)?Masa,\s+Gaurabda\s+(\d+)`)
	reLocation   = regexp.MustCompile(`(.+?)\s+\((\d+N\d+\s+\d+E\d+),\s*([^)]+)\)`)
)

func parseMasaHeader(line string) (name, altName string, year int) {
	m := reMasaHeader.FindStringSubmatch(strings.TrimSpace(line))
	if m == nil {
		return "", "", 0
	}
	name = m[1]
	altName = m[2]
	y, _ := strconv.Atoi(m[3])
	return name, altName, y
}

func parseLocation(line string) (name, coords, tz string) {
	m := reLocation.FindStringSubmatch(strings.TrimSpace(line))
	if m == nil {
		return "", "", ""
	}
	return m[1], m[2], m[3]
}
```

Add `"strconv"` to the imports.

- [ ] **Step 4: Run tests to verify they pass**

Run:
```bash
go test ./internal/parser/ -v -run "TestParseMasaHeader|TestParseLocation"
```
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/parser/
git commit -m "feat: parse masa headers and location metadata"
```

---

### Task 5: Parser — day entry parsing

**Files:**
- Modify: `internal/parser/parser.go`
- Modify: `internal/parser/parser_test.go`

- [ ] **Step 1: Write failing tests for day entry parsing**

Add to `internal/parser/parser_test.go`:

```go
import (
	"testing"
	"time"
)

func TestParseDayEntry(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		day      int
		month    time.Month
		year     int
		weekday  string
		tithi    string
		paksa    string
		naksatra string
		fasting  bool
	}{
		{
			"basic day",
			" 1 Jan 2026 Th   Trayodasi                         G Rohini",
			1, time.January, 2026, "Th", "Trayodasi", "G", "Rohini", false,
		},
		{
			"fasting day",
			"14 Jan 2026 We   Ekadasi (suitable for fasting)    K Anuradha        *",
			14, time.January, 2026, "We", "Ekadasi (suitable for fasting)", "K", "Anuradha", true,
		},
		{
			"not suitable for fasting",
			"26 May 2026 Tu   Ekadasi (not suitable for fasting)G Hasta",
			26, time.May, 2026, "Tu", "Ekadasi (not suitable for fasting)", "G", "Hasta", false,
		},
		{
			"dvadasi fasting",
			"11 Jul 2026 Sa   Dvadasi (suitable for fasting)    K Krittika        *",
			11, time.July, 2026, "Sa", "Dvadasi (suitable for fasting)", "K", "Krittika", true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := parseDayEntry(tt.line)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if d.Date.Day() != tt.day {
				t.Errorf("day = %d, want %d", d.Date.Day(), tt.day)
			}
			if d.Date.Month() != tt.month {
				t.Errorf("month = %v, want %v", d.Date.Month(), tt.month)
			}
			if d.Date.Year() != tt.year {
				t.Errorf("year = %d, want %d", d.Date.Year(), tt.year)
			}
			if d.Weekday != tt.weekday {
				t.Errorf("weekday = %q, want %q", d.Weekday, tt.weekday)
			}
			if d.Tithi != tt.tithi {
				t.Errorf("tithi = %q, want %q", d.Tithi, tt.tithi)
			}
			if d.Paksa != tt.paksa {
				t.Errorf("paksa = %q, want %q", d.Paksa, tt.paksa)
			}
			if d.Naksatra != tt.naksatra {
				t.Errorf("naksatra = %q, want %q", d.Naksatra, tt.naksatra)
			}
			if d.Fasting != tt.fasting {
				t.Errorf("fasting = %v, want %v", d.Fasting, tt.fasting)
			}
		})
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run:
```bash
go test ./internal/parser/ -v -run TestParseDayEntry
```
Expected: FAIL — `parseDayEntry` not defined

- [ ] **Step 3: Implement day entry parsing**

Add to `internal/parser/parser.go`:

```go
import (
	"fmt"
	"time"

	"github.com/taraka/vedic-calendar/internal/models"
)

var reDayParse = regexp.MustCompile(`^\s?(\d{1,2})\s+(\w{3})\s+(\d{4})\s+(\w{2})\s{2,3}(.+?)\s{2,}(G|K)\s+(\S+(?:-\S+)?)\s*(\*?)\s*$`)

func parseDayEntry(line string) (models.Day, error) {
	m := reDayParse.FindStringSubmatch(line)
	if m == nil {
		return models.Day{}, fmt.Errorf("line does not match day entry pattern: %q", line)
	}

	dayNum, _ := strconv.Atoi(m[1])
	monthStr := m[2]
	yearNum, _ := strconv.Atoi(m[3])
	weekday := m[4]
	tithi := strings.TrimSpace(m[5])
	paksa := m[6]
	naksatra := m[7]
	fasting := m[8] == "*"

	monthMap := map[string]time.Month{
		"Jan": time.January, "Feb": time.February, "Mar": time.March,
		"Apr": time.April, "May": time.May, "Jun": time.June,
		"Jul": time.July, "Aug": time.August, "Sep": time.September,
		"Oct": time.October, "Nov": time.November, "Dec": time.December,
	}

	mon, ok := monthMap[monthStr]
	if !ok {
		return models.Day{}, fmt.Errorf("unknown month: %q", monthStr)
	}

	date := time.Date(yearNum, mon, dayNum, 0, 0, 0, 0, time.UTC)

	return models.Day{
		Date:     date,
		Weekday:  weekday,
		Tithi:    tithi,
		Paksa:    paksa,
		Naksatra: naksatra,
		Fasting:  fasting,
	}, nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run:
```bash
go test ./internal/parser/ -v -run TestParseDayEntry
```
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/parser/
git commit -m "feat: parse day entry lines"
```

---

### Task 6: Parser — continuation line classification (events, break-fast, fasting)

**Files:**
- Modify: `internal/parser/parser.go`
- Modify: `internal/parser/parser_test.go`

- [ ] **Step 1: Write failing tests for continuation line classification**

Add to `internal/parser/parser_test.go`:

```go
func TestClassifyEvent(t *testing.T) {
	festivals := []string{"Sri Krsna Pusya Abhiseka", "Vasanta Pancami", "Ratha Yatra"}

	tests := []struct {
		name     string
		text     string
		evtType  models.EventType
		isBF     bool
	}{
		{"appearance", "Sri Ramacandra Kaviraja -- Appearance", models.EventAppearance, false},
		{"disappearance", "Srila Gopala Bhatta Gosvami -- Disappearance", models.EventDisappearance, false},
		{"fasting for", "Fasting for Sat-tila Ekadasi", models.EventFast, false},
		{"fast till", "(Fast till noon)", models.EventFast, false},
		{"fast till midnight", "(Fast till midnight)", models.EventFast, false},
		{"fasting done", "(Fasting is done yesterday, today is feast)", models.EventNote, false},
		{"festival", "Sri Krsna Pusya Abhiseka", models.EventFestival, false},
		{"festival ratha", "Ratha Yatra", models.EventFestival, false},
		{"break fast", "Break fast 07:09 (sunrise) - 10:38 (1/3 of daylight) LT", "", true},
		{"note", "Tulasi Jala Dan begins.", models.EventNote, false},
		{"note caturmasya", "First month of Caturmasya begins [PURNIMA SYSTEM]", models.EventNote, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evt, bf := classifyEvent(tt.text, festivals)
			if tt.isBF {
				if bf == nil {
					t.Error("expected BreakFast, got nil")
				}
			} else {
				if bf != nil {
					t.Errorf("expected no BreakFast, got %+v", bf)
				}
				if evt.Type != tt.evtType {
					t.Errorf("type = %q, want %q", evt.Type, tt.evtType)
				}
			}
		})
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run:
```bash
go test ./internal/parser/ -v -run TestClassifyEvent
```
Expected: FAIL — `classifyEvent` not defined

- [ ] **Step 3: Implement classifyEvent**

Add to `internal/parser/parser.go`:

```go
var reBreakFast = regexp.MustCompile(`^Break fast\s+(\d{2}:\d{2})\s+\([^)]+\)\s*-\s*(\d{2}:\d{2})\s+(.+)$`)

func classifyEvent(text string, festivals []string) (models.Event, *models.BreakFast) {
	trimmed := strings.TrimSpace(text)

	// Break-fast info
	if m := reBreakFast.FindStringSubmatch(trimmed); m != nil {
		return models.Event{}, &models.BreakFast{
			Start: m[1],
			End:   m[2],
			Notes: strings.TrimSpace(m[3]),
		}
	}

	// Appearance
	if strings.Contains(trimmed, "-- Appearance") {
		return models.Event{Type: models.EventAppearance, Description: trimmed}, nil
	}

	// Disappearance
	if strings.Contains(trimmed, "-- Disappearance") {
		return models.Event{Type: models.EventDisappearance, Description: trimmed}, nil
	}

	// Fasting instruction
	if strings.HasPrefix(trimmed, "Fasting for ") {
		return models.Event{Type: models.EventFast, Description: trimmed}, nil
	}
	if strings.HasPrefix(trimmed, "(Fast till ") || strings.HasPrefix(trimmed, "(Fast till noon") {
		return models.Event{Type: models.EventFast, Description: trimmed}, nil
	}

	// Festival (check curated list)
	for _, f := range festivals {
		if strings.Contains(trimmed, f) {
			return models.Event{Type: models.EventFestival, Description: trimmed}, nil
		}
	}

	// Default: note
	return models.Event{Type: models.EventNote, Description: trimmed}, nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run:
```bash
go test ./internal/parser/ -v -run TestClassifyEvent
```
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/parser/
git commit -m "feat: classify continuation lines into event types"
```

---

### Task 7: Parser — Sankranti parsing

**Files:**
- Modify: `internal/parser/parser.go`
- Modify: `internal/parser/parser_test.go`

- [ ] **Step 1: Write failing test for Sankranti parsing**

Add to `internal/parser/parser_test.go`:

```go
func TestParseSankranti(t *testing.T) {
	tests := []struct {
		line  string
		name  string
		notes string
		day   int
		month time.Month
	}{
		{
			"---------- Makara Sankranti (Sun enters Capricorn on 14 Jan, 14:54 LT) ----------",
			"Makara Sankranti",
			"Sun enters Capricorn on 14 Jan, 14:54 LT",
			14, time.January,
		},
		{
			"----------- Kumbha Sankranti (Sun enters Aquarius on 13 Feb, 03:52 LT) -----------",
			"Kumbha Sankranti",
			"Sun enters Aquarius on 13 Feb, 03:52 LT",
			13, time.February,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, day, month := parseSankranti(tt.line)
			if s.Name != tt.name {
				t.Errorf("name = %q, want %q", s.Name, tt.name)
			}
			if s.Notes != tt.notes {
				t.Errorf("notes = %q, want %q", s.Notes, tt.notes)
			}
			if day != tt.day {
				t.Errorf("day = %d, want %d", day, tt.day)
			}
			if month != tt.month {
				t.Errorf("month = %v, want %v", month, tt.month)
			}
		})
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run:
```bash
go test ./internal/parser/ -v -run TestParseSankranti
```
Expected: FAIL

- [ ] **Step 3: Implement parseSankranti**

Add to `internal/parser/parser.go`:

```go
var reSankrantiParse = regexp.MustCompile(`^-+\s+(\w+\s+Sankranti)\s+\((.+?)\)\s+-+$`)
var reSankrantiDate  = regexp.MustCompile(`on\s+(\d{1,2})\s+(\w{3})`)

func parseSankranti(line string) (models.Sankranti, int, time.Month) {
	trimmed := strings.TrimSpace(line)
	m := reSankrantiParse.FindStringSubmatch(trimmed)
	if m == nil {
		return models.Sankranti{}, 0, 0
	}

	name := m[1]
	notes := m[2]

	dm := reSankrantiDate.FindStringSubmatch(notes)
	var day int
	var month time.Month
	if dm != nil {
		day, _ = strconv.Atoi(dm[1])
		monthMap := map[string]time.Month{
			"Jan": time.January, "Feb": time.February, "Mar": time.March,
			"Apr": time.April, "May": time.May, "Jun": time.June,
			"Jul": time.July, "Aug": time.August, "Sep": time.September,
			"Oct": time.October, "Nov": time.November, "Dec": time.December,
		}
		month = monthMap[dm[2]]
	}

	return models.Sankranti{Name: name, Notes: notes}, day, month
}
```

- [ ] **Step 4: Run test to verify it passes**

Run:
```bash
go test ./internal/parser/ -v -run TestParseSankranti
```
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/parser/
git commit -m "feat: parse Sankranti lines with date extraction"
```

---

### Task 8: Parser — full file orchestrator

**Files:**
- Modify: `internal/parser/parser.go`
- Modify: `internal/parser/parser_test.go`

- [ ] **Step 1: Write failing test for full parse**

Add to `internal/parser/parser_test.go`:

```go
func TestParseCalendar(t *testing.T) {
	input := `
                       Pausa (Narayana) Masa, Gaurabda 539      GCal 11, Build 5

                      Vrindavan [India] (27N35 77E42, +5:30)

 DATE            TITHI                         PAKSA NAKSATRA       FAST
--------------------------------------------------------------------------------
 1 Jan 2026 Th   Trayodasi                         G Rohini
                 Break fast 07:09 (sunrise) - 10:38 (1/3 of daylight) LT
 2 Jan 2026 Fr   Caturdasi                         G Mrigasira
 3 Jan 2026 Sa   Purnima                           G Ardra
                 Sri Krsna Pusya Abhiseka

                        Magha (Madhava) Masa, Gaurabda 539      GCal 11, Build 5

                      Vrindavan [India] (27N35 77E42, +5:30)

 DATE            TITHI                         PAKSA NAKSATRA       FAST
--------------------------------------------------------------------------------
 4 Jan 2026 Su   Pratipat                          K Punarvasu
 5 Jan 2026 Mo   Dvitiya                           K Pusyami
14 Jan 2026 We   Ekadasi (suitable for fasting)    K Anuradha        *
                 Fasting for Sat-tila Ekadasi
15 Jan 2026 Th   Dvadasi                           K Jyestha
                 Break fast 07:10 (sunrise) - 10:42 (1/3 of daylight) LT
                 Ganga Sagara Mela
---------- Makara Sankranti (Sun enters Capricorn on 14 Jan, 14:54 LT) ----------
16 Jan 2026 Fr   Trayodasi                         K Mula
`

	festivals := []string{"Sri Krsna Pusya Abhiseka", "Ganga Sagara Mela"}
	cal, err := ParseCalendar(input, festivals)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Location
	if cal.Location.Name != "Vrindavan [India]" {
		t.Errorf("location name = %q", cal.Location.Name)
	}

	// Two months
	if len(cal.Months) != 2 {
		t.Fatalf("months = %d, want 2", len(cal.Months))
	}

	// First month
	m0 := cal.Months[0]
	if m0.MasaName != "Pausa" {
		t.Errorf("m0 masa = %q", m0.MasaName)
	}
	if m0.MasaAltName != "Narayana" {
		t.Errorf("m0 alt = %q", m0.MasaAltName)
	}
	if len(m0.Days) != 3 {
		t.Fatalf("m0 days = %d, want 3", len(m0.Days))
	}

	// First day has break-fast
	if m0.Days[0].BreakFast == nil {
		t.Error("day 1 should have break-fast")
	} else if m0.Days[0].BreakFast.Start != "07:09" {
		t.Errorf("break-fast start = %q", m0.Days[0].BreakFast.Start)
	}

	// Third day has festival event
	if len(m0.Days[2].Events) != 1 {
		t.Fatalf("day 3 events = %d, want 1", len(m0.Days[2].Events))
	}
	if m0.Days[2].Events[0].Type != models.EventFestival {
		t.Errorf("day 3 event type = %q", m0.Days[2].Events[0].Type)
	}

	// Second month
	m1 := cal.Months[1]
	if m1.MasaName != "Magha" {
		t.Errorf("m1 masa = %q", m1.MasaName)
	}
	if len(m1.Days) < 5 {
		t.Fatalf("m1 days = %d, want >= 5", len(m1.Days))
	}

	// Ekadasi day has fasting
	ekDay := m1.Days[2] // 14 Jan
	if !ekDay.Fasting {
		t.Error("14 Jan should have fasting=true")
	}
	if len(ekDay.Events) != 1 {
		t.Fatalf("14 Jan events = %d", len(ekDay.Events))
	}
	if ekDay.Events[0].Type != models.EventFast {
		t.Errorf("14 Jan event type = %q", ekDay.Events[0].Type)
	}

	// 15 Jan has break-fast and festival
	d15 := m1.Days[3]
	if d15.BreakFast == nil {
		t.Error("15 Jan should have break-fast")
	}
	festFound := false
	for _, e := range d15.Events {
		if e.Type == models.EventFestival {
			festFound = true
		}
	}
	if !festFound {
		t.Error("15 Jan should have festival event (Ganga Sagara Mela)")
	}

	// Sankranti attached to 14 Jan
	if ekDay.Sankranti == nil {
		t.Error("14 Jan should have Sankranti")
	} else if ekDay.Sankranti.Name != "Makara Sankranti" {
		t.Errorf("sankranti name = %q", ekDay.Sankranti.Name)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run:
```bash
go test ./internal/parser/ -v -run TestParseCalendar
```
Expected: FAIL — `ParseCalendar` not defined

- [ ] **Step 3: Implement ParseCalendar**

Add to `internal/parser/parser.go`:

```go
func ParseCalendar(input string, festivals []string) (models.Calendar, error) {
	lines := strings.Split(input, "\n")
	cal := models.Calendar{}
	var currentMonth *models.VedicMonth
	var currentDay *models.Day
	locationParsed := false

	// Collect all days by date for Sankranti attachment
	type sankrantiPending struct {
		s     models.Sankranti
		day   int
		month time.Month
	}
	var pendingSankrantis []sankrantiPending

	for _, line := range lines {
		lt := classifyLine(line)

		switch lt {
		case LineMasaHeader:
			name, altName, year := parseMasaHeader(line)
			cal.Months = append(cal.Months, models.VedicMonth{
				MasaName:     name,
				MasaAltName:  altName,
				GaurabdaYear: year,
			})
			currentMonth = &cal.Months[len(cal.Months)-1]
			currentDay = nil

		case LineLocation:
			if !locationParsed {
				n, c, tz := parseLocation(line)
				cal.Location = models.Location{Name: n, Coords: c, Timezone: tz}
				locationParsed = true
			}

		case LineDayEntry:
			d, err := parseDayEntry(line)
			if err != nil {
				return cal, fmt.Errorf("parse day: %w", err)
			}
			if currentMonth == nil {
				return cal, fmt.Errorf("day entry before any masa header")
			}
			currentMonth.Days = append(currentMonth.Days, d)
			currentDay = &currentMonth.Days[len(currentMonth.Days)-1]

		case LineContinuation:
			if currentDay == nil {
				continue
			}
			text := strings.TrimSpace(line)
			evt, bf := classifyEvent(text, festivals)
			if bf != nil {
				currentDay.BreakFast = bf
			} else {
				currentDay.Events = append(currentDay.Events, evt)
			}

		case LineSankranti:
			s, day, month := parseSankranti(line)
			if day > 0 {
				pendingSankrantis = append(pendingSankrantis, sankrantiPending{s, day, month})
			}

		default:
			// skip blank, separator, table header
		}
	}

	// Attach Sankrantis to matching days
	for _, ps := range pendingSankrantis {
		for mi := range cal.Months {
			for di := range cal.Months[mi].Days {
				d := &cal.Months[mi].Days[di]
				if d.Date.Day() == ps.day && d.Date.Month() == ps.month {
					d.Sankranti = &ps.s
				}
			}
		}
	}

	return cal, nil
}
```

- [ ] **Step 4: Run test to verify it passes**

Run:
```bash
go test ./internal/parser/ -v -run TestParseCalendar
```
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/parser/
git commit -m "feat: full calendar parser orchestrator"
```

---

### Task 9: Parser — festival list loader and CLI wiring

**Files:**
- Modify: `internal/parser/parser.go`
- Modify: `cmd/parser/main.go`

- [ ] **Step 1: Add festival list loader**

Add to `internal/parser/parser.go`:

```go
import (
	"encoding/json"
	"os"
)

type FestivalList struct {
	Festivals []string `json:"festivals"`
}

func LoadFestivals(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read festivals: %w", err)
	}
	var fl FestivalList
	if err := json.Unmarshal(data, &fl); err != nil {
		return nil, fmt.Errorf("parse festivals: %w", err)
	}
	return fl.Festivals, nil
}

func ParseCalendarFile(inputPath, festivalsPath string) (models.Calendar, error) {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return models.Calendar{}, fmt.Errorf("read input: %w", err)
	}
	festivals, err := LoadFestivals(festivalsPath)
	if err != nil {
		return models.Calendar{}, err
	}
	return ParseCalendar(string(data), festivals)
}
```

- [ ] **Step 2: Wire up the parser CLI**

Replace `cmd/parser/main.go`:

```go
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/taraka/vedic-calendar/internal/parser"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "usage: parser <input.txt> <festivals.json> [output.json]\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	festivalsPath := os.Args[2]

	cal, err := parser.ParseCalendarFile(inputPath, festivalsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	out, err := json.MarshalIndent(cal, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	if len(os.Args) >= 4 {
		if err := os.WriteFile(os.Args[3], out, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "error writing file: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "wrote %s\n", os.Args[3])
	} else {
		fmt.Println(string(out))
	}
}
```

- [ ] **Step 3: Build and test with the actual input file**

Run:
```bash
go build -o parser ./cmd/parser && mkdir -p output && ./parser "data/Vrindavan [India].txt" data/festivals.json output/calendar.json
```
Expected: `output/calendar.json` created, no errors. Inspect with:
```bash
head -50 output/calendar.json
```
Verify JSON contains location, months, days with events.

- [ ] **Step 4: Commit**

```bash
git add cmd/parser/ internal/parser/
git commit -m "feat: parser CLI with festival loader, end-to-end parsing"
```

---

### Task 10: Renderer — abbreviation maps

**Files:**
- Create: `internal/renderer/abbreviations.go`

- [ ] **Step 1: Create abbreviation maps**

Create `internal/renderer/abbreviations.go`:

```go
package renderer

var TithiAbbrev = map[string]string{
	"Pratipat":  "Pra",
	"Dvitiya":   "Dvi",
	"Tritiya":   "Tri",
	"Caturthi":  "Cat",
	"Pancami":   "Pan",
	"Sasti":     "Sas",
	"Saptami":   "Sap",
	"Astami":    "Ast",
	"Navami":    "Nav",
	"Dasami":    "Das",
	"Ekadasi":   "Eka",
	"Dvadasi":   "Dva",
	"Trayodasi": "Tra",
	"Caturdasi": "Cad",
	"Purnima":   "Pur",
	"Amavasya":  "Ama",
}

var NaksatraAbbrev = map[string]string{
	"Asvini":          "Asv",
	"Bharani":         "Bha",
	"Krittika":        "Kri",
	"Rohini":          "Roh",
	"Mrigasira":       "Mri",
	"Ardra":           "Ard",
	"Punarvasu":       "Pun",
	"Pusyami":         "Pus",
	"Aslesa":          "Asl",
	"Magha":           "Mag",
	"Purva-phalguni":  "PPh",
	"Uttara-phalguni": "UPh",
	"Hasta":           "Has",
	"Citra":           "Cit",
	"Swati":           "Swa",
	"Visakha":         "Vis",
	"Anuradha":        "Anu",
	"Jyestha":         "Jye",
	"Mula":            "Mul",
	"Purva-asadha":    "PAs",
	"Uttara-asadha":   "UAs",
	"Sravana":         "Sra",
	"Dhanista":        "Dha",
	"Satabhisa":       "Sat",
	"Purva-bhadra":    "PBh",
	"Uttara-bhadra":   "UBh",
	"Revati":          "Rev",
}

// ShortenTithi extracts the base tithi name and abbreviates it.
// Input may contain qualifiers like "Ekadasi (suitable for fasting)".
func ShortenTithi(tithi string) string {
	for full, abbr := range TithiAbbrev {
		if len(tithi) >= len(full) && tithi[:len(full)] == full {
			return abbr
		}
	}
	if len(tithi) >= 3 {
		return tithi[:3]
	}
	return tithi
}

// ShortenNaksatra abbreviates a naksatra name.
func ShortenNaksatra(naksatra string) string {
	if abbr, ok := NaksatraAbbrev[naksatra]; ok {
		return abbr
	}
	if len(naksatra) >= 3 {
		return naksatra[:3]
	}
	return naksatra
}
```

- [ ] **Step 2: Verify it compiles**

Run:
```bash
go build ./internal/renderer/
```
Expected: compiles

- [ ] **Step 3: Commit**

```bash
git add internal/renderer/
git commit -m "feat: tithi and naksatra abbreviation maps"
```

---

### Task 11: Renderer — grid data preparation and footnote numbering

**Files:**
- Create: `internal/renderer/renderer.go`
- Create: `internal/renderer/renderer_test.go`

- [ ] **Step 1: Write failing test for grid data preparation**

Create `internal/renderer/renderer_test.go`:

```go
package renderer

import (
	"testing"
	"time"

	"github.com/taraka/vedic-calendar/internal/models"
)

func TestBuildMonthGrid(t *testing.T) {
	// A month starting on Sunday (4 Jan 2026 = Sunday)
	days := []models.Day{
		{Date: time.Date(2026, 1, 4, 0, 0, 0, 0, time.UTC), Weekday: "Su", Tithi: "Pratipat", Paksa: "K", Naksatra: "Punarvasu"},
		{Date: time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC), Weekday: "Mo", Tithi: "Dvitiya", Paksa: "K", Naksatra: "Pusyami"},
		{Date: time.Date(2026, 1, 6, 0, 0, 0, 0, time.UTC), Weekday: "Tu", Tithi: "Tritiya", Paksa: "K", Naksatra: "Aslesa"},
	}

	month := models.VedicMonth{
		MasaName:     "Magha",
		MasaAltName:  "Madhava",
		GaurabdaYear: 539,
		Days:         days,
	}

	grid := BuildMonthGrid(month)

	if grid.MasaName != "Magha" {
		t.Errorf("masa = %q", grid.MasaName)
	}

	// First day is Sunday = column 0
	if len(grid.Weeks) == 0 {
		t.Fatal("no weeks")
	}
	if grid.Weeks[0][0] == nil {
		t.Fatal("Sunday cell should not be nil")
	}
	if grid.Weeks[0][0].DayNum != 4 {
		t.Errorf("day num = %d, want 4", grid.Weeks[0][0].DayNum)
	}
}

func TestFootnoteNumbering(t *testing.T) {
	days := []models.Day{
		{
			Date: time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC), Weekday: "We",
			Tithi: "Pancami", Paksa: "K", Naksatra: "Magha",
			Events: []models.Event{
				{Type: models.EventAppearance, Description: "Srila Gopala Bhatta Gosvami -- Appearance"},
				{Type: models.EventDisappearance, Description: "Sri Ramacandra Kaviraja -- Disappearance"},
			},
		},
	}

	month := models.VedicMonth{MasaName: "Magha", Days: days}
	grid := BuildMonthGrid(month)

	if len(grid.Footnotes) != 2 {
		t.Fatalf("footnotes = %d, want 2", len(grid.Footnotes))
	}
	if grid.Footnotes[0].Number != 1 {
		t.Errorf("fn 0 number = %d", grid.Footnotes[0].Number)
	}
	if grid.Weeks[0][3].FootnoteRefs[0] != 1 {
		t.Errorf("cell fn ref = %d", grid.Weeks[0][3].FootnoteRefs[0])
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run:
```bash
go test ./internal/renderer/ -v -run "TestBuildMonthGrid|TestFootnoteNumbering"
```
Expected: FAIL

- [ ] **Step 3: Implement grid data structures and BuildMonthGrid**

Add to `internal/renderer/renderer.go`:

```go
package renderer

import (
	"fmt"
	"strings"

	"github.com/taraka/vedic-calendar/internal/models"
)

// CellData holds the rendered content for one day cell in the grid.
type CellData struct {
	DayNum       int
	Weekday      string
	TithiShort   string
	PaksaLabel   string
	NaksatraShort string
	Fasting      bool
	IsEkadasi    bool
	Icons        []string   // icon characters to show
	FootnoteRefs []int      // footnote numbers referenced by this cell
	HasSankranti bool
}

// Footnote holds one footnote entry for the page.
type Footnote struct {
	Number int
	Text   string
	Icon   string
}

// MonthGrid holds all data needed to render one month page.
type MonthGrid struct {
	MasaName     string
	MasaAltName  string
	GaurabdaYear int
	Location     string
	Weeks        [][7]*CellData // each row is Sun..Sat, nil = empty cell
	Footnotes    []Footnote
}

// Icon characters
const (
	IconEkadasi      = "◆"
	IconAppearance   = "★"
	IconDisappearance = "†"
	IconFestival     = "☀"
	IconSankranti    = "▲"
)

func iconForEventType(t models.EventType) string {
	switch t {
	case models.EventAppearance:
		return IconAppearance
	case models.EventDisappearance:
		return IconDisappearance
	case models.EventFestival:
		return IconFestival
	default:
		return ""
	}
}

func BuildMonthGrid(month models.VedicMonth) MonthGrid {
	grid := MonthGrid{
		MasaName:     month.MasaName,
		MasaAltName:  month.MasaAltName,
		GaurabdaYear: month.GaurabdaYear,
	}

	footnoteNum := 1

	// Determine which weeks we need
	type dayPos struct {
		weekIdx int
		dayCol  int
	}

	if len(month.Days) == 0 {
		return grid
	}

	// Find the weekday column (Sun=0 .. Sat=6)
	weekdayCol := func(wd string) int {
		cols := map[string]int{"Su": 0, "Mo": 1, "Tu": 2, "We": 3, "Th": 4, "Fr": 5, "Sa": 6}
		if c, ok := cols[wd]; ok {
			return c
		}
		return 0
	}

	// Build week rows
	currentWeek := 0
	firstCol := weekdayCol(month.Days[0].Weekday)
	grid.Weeks = append(grid.Weeks, [7]*CellData{})

	prevCol := firstCol - 1

	for i := range month.Days {
		day := &month.Days[i]
		col := weekdayCol(day.Weekday)

		// If we wrapped past Saturday, or col <= prevCol (new week)
		if i > 0 && col <= prevCol {
			currentWeek++
			grid.Weeks = append(grid.Weeks, [7]*CellData{})
		}
		prevCol = col

		cell := &CellData{
			DayNum:        day.Date.Day(),
			Weekday:       day.Weekday,
			TithiShort:    ShortenTithi(day.Tithi),
			PaksaLabel:    day.Paksa,
			NaksatraShort: ShortenNaksatra(day.Naksatra),
			Fasting:       day.Fasting,
			IsEkadasi:     strings.HasPrefix(day.Tithi, "Ekadasi"),
		}

		if cell.IsEkadasi {
			cell.Icons = append(cell.Icons, IconEkadasi)
		}

		if day.Sankranti != nil {
			cell.HasSankranti = true
			cell.Icons = append(cell.Icons, IconSankranti)
			grid.Footnotes = append(grid.Footnotes, Footnote{
				Number: footnoteNum,
				Text:   fmt.Sprintf("%s — %s", day.Sankranti.Name, day.Sankranti.Notes),
				Icon:   IconSankranti,
			})
			cell.FootnoteRefs = append(cell.FootnoteRefs, footnoteNum)
			footnoteNum++
		}

		for _, evt := range day.Events {
			icon := iconForEventType(evt.Type)
			if icon != "" {
				// Deduplicate icons
				found := false
				for _, ic := range cell.Icons {
					if ic == icon {
						found = true
						break
					}
				}
				if !found {
					cell.Icons = append(cell.Icons, icon)
				}
			}

			// All events with descriptions become footnotes
			grid.Footnotes = append(grid.Footnotes, Footnote{
				Number: footnoteNum,
				Text:   evt.Description,
				Icon:   icon,
			})
			cell.FootnoteRefs = append(cell.FootnoteRefs, footnoteNum)
			footnoteNum++
		}

		if day.BreakFast != nil {
			text := fmt.Sprintf("Break fast %s – %s (%s)", day.BreakFast.Start, day.BreakFast.End, day.BreakFast.Notes)
			grid.Footnotes = append(grid.Footnotes, Footnote{
				Number: footnoteNum,
				Text:   text,
			})
			cell.FootnoteRefs = append(cell.FootnoteRefs, footnoteNum)
			footnoteNum++
		}

		grid.Weeks[currentWeek][col] = cell
	}

	return grid
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run:
```bash
go test ./internal/renderer/ -v -run "TestBuildMonthGrid|TestFootnoteNumbering"
```
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/renderer/
git commit -m "feat: grid data preparation with footnote numbering"
```

---

### Task 12: LuaLaTeX calendar template

**Files:**
- Create: `templates/calendar.tex.tmpl`

- [ ] **Step 1: Create the calendar template**

Create `templates/calendar.tex.tmpl`:

```latex
\documentclass[a4paper,landscape,10pt]{article}
\usepackage[left=8mm,right=8mm,top=8mm,bottom=8mm]{geometry}
\usepackage{fontspec}
\usepackage{tikz}
\usepackage{xcolor}
\usepackage{array}
\usepackage{calc}

\setmainfont{Latin Modern Roman}

\definecolor{fastingbg}{gray}{0.90}
\definecolor{cellborder}{gray}{0.5}
\definecolor{headerbg}{gray}{0.15}
\definecolor{headertext}{gray}{1.0}
\definecolor{weekdaybg}{gray}{0.30}
\definecolor{emptycell}{gray}{0.97}

\pagestyle{empty}
\setlength{\parindent}{0pt}

\newcommand{\calheader}[3]{%
  \noindent
  \begin{tikzpicture}
    \fill[headerbg] (0,0) rectangle (\textwidth, 1.2cm);
    \node[anchor=west, text=headertext, font=\Large\bfseries] at (0.3cm, 0.6cm) {#1};
    \node[anchor=east, text=headertext, font=\large] at (\textwidth-0.3cm, 0.6cm) {Gaurabda #2};
    \node[anchor=west, text=headertext, font=\small] at (0.3cm, 0.15cm) {#3};
  \end{tikzpicture}
}

\newcommand{\weekdaybar}{%
  \noindent
  \begin{tikzpicture}
    \def\cellw{\textwidth/7}
    \fill[weekdaybg] (0,0) rectangle (\textwidth, 0.5cm);
    \foreach \d [count=\i from 0] in {Sun,Mon,Tue,Wed,Thu,Fri,Sat} {
      \node[anchor=center, text=white, font=\small\bfseries] at (\i*\cellw + \cellw/2, 0.25cm) {\d};
    }
  \end{tikzpicture}
}

\begin{document}
{{- range .Grids}}

\calheader{ {{- .MasaName}}{{if .MasaAltName}} ({{.MasaAltName}}){{end}} Masa}{ {{- .GaurabdaYear}}}{ {{- $.Location.Name}}}

\weekdaybar

\noindent
\begin{tikzpicture}[x=\textwidth/7, y=-1cm]
  \pgfmathsetmacro{\cellw}{\textwidth/7}
  \def\rowh{2.4cm}
  {{- $numWeeks := len .Weeks}}
  {{- range $wi, $week := .Weeks}}
  {{- range $ci, $cell := $week}}
  {{- if $cell}}
  % Day {{$cell.DayNum}}
  {{- if $cell.Fasting}}
  \fill[fastingbg] ({{$ci}}, {{$wi}}*\rowh) rectangle +(1, \rowh);
  {{- end}}
  \draw[cellborder, thin] ({{$ci}}, {{$wi}}*\rowh) rectangle +(1, \rowh);
  \node[anchor=north west, font=\large\bfseries, inner sep=2pt] at ({{$ci}}+0.03, {{$wi}}*\rowh+0.03) { {{- $cell.DayNum -}} };
  \node[anchor=north east, font=\tiny, inner sep=2pt] at ({{$ci}}+0.97, {{$wi}}*\rowh+0.03) { {{- range $cell.Icons}} {{- . -}} {{- end -}} };
  \node[anchor=north west, font=\scriptsize, inner sep=2pt] at ({{$ci}}+0.03, {{$wi}}*\rowh+0.55) { {{- $cell.PaksaLabel}}\cdot{{$cell.TithiShort -}} };
  \node[anchor=north west, font=\tiny, text=gray, inner sep=2pt] at ({{$ci}}+0.03, {{$wi}}*\rowh+0.85) { {{- $cell.NaksatraShort -}} };
  {{- if $cell.FootnoteRefs}}
  \node[anchor=south west, font=\tiny, inner sep=2pt] at ({{$ci}}+0.03, {{$wi}}*\rowh+\rowh-0.05) { {{- range $fi, $fn := $cell.FootnoteRefs}}{{if $fi}},{{end}}\textsuperscript{ {{- $fn -}} }{{- end -}} };
  {{- end}}
  {{- else}}
  \fill[emptycell] ({{$ci}}, {{$wi}}*\rowh) rectangle +(1, \rowh);
  \draw[cellborder, very thin, dashed] ({{$ci}}, {{$wi}}*\rowh) rectangle +(1, \rowh);
  {{- end}}
  {{- end}}
  {{- end}}
\end{tikzpicture}

\vfill

\noindent
\begin{minipage}[b]{\textwidth}
\footnotesize
{{- range .Footnotes}}
\textsuperscript{ {{- .Number}}} {{.Icon}} {{.Text}}\\
{{- end}}
\medskip
{\scriptsize ◆\,Ekadasi \quad ★\,Appearance \quad †\,Disappearance \quad ☀\,Festival \quad ▲\,Sankranti \quad \colorbox{fastingbg}{~~}\,Fasting day}
\end{minipage}

\newpage
{{- end}}

\end{document}
```

- [ ] **Step 2: Commit**

```bash
git add templates/
git commit -m "feat: LuaLaTeX calendar grid template"
```

---

### Task 13: LuaLaTeX reference page template

**Files:**
- Create: `templates/reference.tex.tmpl`

- [ ] **Step 1: Create the reference page template**

Create `templates/reference.tex.tmpl`:

```latex
% Reference page — appended after all month pages
\newpage
\begin{center}
{\Large\bfseries Reference}
\end{center}
\vspace{0.5cm}

\noindent
\begin{minipage}[t]{0.48\textwidth}
{\normalsize\bfseries Tithi Abbreviations}\\[4pt]
\begin{tabular}{ll}
Pra & Pratipat \\
Dvi & Dvitiya \\
Tri & Tritiya \\
Cat & Caturthi \\
Pan & Pancami \\
Sas & Sasti \\
Sap & Saptami \\
Ast & Astami \\
Nav & Navami \\
Das & Dasami \\
Eka & Ekadasi \\
Dva & Dvadasi \\
Tra & Trayodasi \\
Cad & Caturdasi \\
Pur & Purnima \\
Ama & Amavasya \\
\end{tabular}

\vspace{0.5cm}
{\normalsize\bfseries Paksa}\\[4pt]
\begin{tabular}{ll}
G & Gaura (bright/waxing) \\
K & Krsna (dark/waning) \\
\end{tabular}

\vspace{0.5cm}
{\normalsize\bfseries Icons}\\[4pt]
\begin{tabular}{ll}
◆ & Ekadasi \\
★ & Appearance \\
† & Disappearance \\
☀ & Festival \\
▲ & Sankranti \\
\colorbox{fastingbg}{~~} & Fasting day \\
\end{tabular}
\end{minipage}%
\hfill
\begin{minipage}[t]{0.48\textwidth}
{\normalsize\bfseries Naksatra Abbreviations}\\[4pt]
\begin{tabular}{ll}
Asv & Asvini \\
Bha & Bharani \\
Kri & Krittika \\
Roh & Rohini \\
Mri & Mrigasira \\
Ard & Ardra \\
Pun & Punarvasu \\
Pus & Pusyami \\
Asl & Aslesa \\
Mag & Magha \\
PPh & Purva-phalguni \\
UPh & Uttara-phalguni \\
Has & Hasta \\
Cit & Citra \\
Swa & Swati \\
Vis & Visakha \\
Anu & Anuradha \\
Jye & Jyestha \\
Mul & Mula \\
PAs & Purva-asadha \\
UAs & Uttara-asadha \\
Sra & Sravana \\
Dha & Dhanista \\
Sat & Satabhisa \\
PBh & Purva-bhadra \\
UBh & Uttara-bhadra \\
Rev & Revati \\
\end{tabular}

\vspace{0.5cm}
{\normalsize\bfseries Vedic Months (Masa)}\\[4pt]
\begin{tabular}{ll}
Pausa & Narayana \\
Magha & Madhava \\
Phalguna & Govinda \\
Caitra & Visnu \\
Vaisakha & Madhusudana \\
Jyestha & Trivikrama \\
Asadha & Vamana \\
Sravana & Sridhara \\
Bhadra & Hrsikesa \\
Asvina & Padmanabha \\
Kartika & Damodara \\
Margasirsa & Kesava \\
\end{tabular}

\vspace{0.5cm}
{\normalsize\bfseries Location}\\[4pt]
{{.Location.Name}}\\
Coordinates: {{.Location.Coords}}\\
Timezone: {{.Location.Timezone}}

\vspace{0.3cm}
{\normalsize\bfseries Notes}\\[4pt]
DST --- Time is in `Daylight Saving Time'\\
LT --- Time is in `Local Time'

\vspace{0.3cm}
{\scriptsize Generated from GCal 11, Build 5}
\end{minipage}
```

- [ ] **Step 2: Commit**

```bash
git add templates/reference.tex.tmpl
git commit -m "feat: reference page LuaLaTeX template"
```

---

### Task 14: Renderer — template execution and CLI

**Files:**
- Modify: `internal/renderer/renderer.go`
- Modify: `cmd/renderer/main.go`

- [ ] **Step 1: Add template execution to renderer**

Add to `internal/renderer/renderer.go`:

```go
import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"text/template"
)

// TemplateData is the top-level data passed to the LuaLaTeX template.
type TemplateData struct {
	Location models.Location
	Grids    []MonthGrid
}

func LoadCalendar(jsonPath string) (models.Calendar, error) {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return models.Calendar{}, err
	}
	var cal models.Calendar
	if err := json.Unmarshal(data, &cal); err != nil {
		return models.Calendar{}, err
	}
	return cal, nil
}

func RenderTeX(cal models.Calendar, templateDir string, w io.Writer) error {
	// Build grids for all months
	td := TemplateData{
		Location: cal.Location,
	}
	for _, m := range cal.Months {
		g := BuildMonthGrid(m)
		g.Location = cal.Location.Name
		td.Grids = append(td.Grids, g)
	}

	// Parse templates
	calTmpl, err := template.ParseFiles(
		filepath.Join(templateDir, "calendar.tex.tmpl"),
	)
	if err != nil {
		return fmt.Errorf("parse calendar template: %w", err)
	}

	refTmplPath := filepath.Join(templateDir, "reference.tex.tmpl")
	refContent, err := os.ReadFile(refTmplPath)
	if err != nil {
		return fmt.Errorf("read reference template: %w", err)
	}

	// Execute calendar template (includes \begin{document} and \end{document})
	// We need to insert the reference page before \end{document}
	// Strategy: template ends with \end{document}, we render it then inject reference before that

	// Actually, let's modify the approach: the calendar template outputs everything
	// up to the last \newpage, then we append reference, then \end{document}.

	// Simpler approach: render calendar template which does NOT include \end{document},
	// then render reference template, then write \end{document}.

	if err := calTmpl.Execute(w, td); err != nil {
		return fmt.Errorf("execute calendar template: %w", err)
	}

	// Parse and execute reference template
	refTmpl, err := template.New("reference").Parse(string(refContent))
	if err != nil {
		return fmt.Errorf("parse reference template: %w", err)
	}
	if err := refTmpl.Execute(w, td); err != nil {
		return fmt.Errorf("execute reference template: %w", err)
	}

	return nil
}
```

Note: This requires the `calendar.tex.tmpl` to NOT include `\end{document}` at the end, so the reference page can be appended. Update the calendar template: replace the last `\end{document}` with nothing — the renderer will append the reference page and then write `\end{document}`.

- [ ] **Step 2: Update calendar.tex.tmpl — remove trailing \end{document}**

Edit `templates/calendar.tex.tmpl`: replace the final `\end{document}` with a comment `% reference page and \end{document} appended by renderer`.

Then update `RenderTeX` to write `\end{document}` after the reference template:

Add at the end of `RenderTeX`:

```go
	fmt.Fprintln(w, `\end{document}`)
```

- [ ] **Step 3: Wire up the renderer CLI**

Replace `cmd/renderer/main.go`:

```go
package main

import (
	"fmt"
	"os"

	"github.com/taraka/vedic-calendar/internal/renderer"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "usage: renderer <input.json> <templates-dir> <output.tex>\n")
		os.Exit(1)
	}

	jsonPath := os.Args[1]
	templateDir := os.Args[2]
	outputPath := os.Args[3]

	cal, err := renderer.LoadCalendar(jsonPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading calendar: %v\n", err)
		os.Exit(1)
	}

	f, err := os.Create(outputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating output: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	if err := renderer.RenderTeX(cal, templateDir, f); err != nil {
		fmt.Fprintf(os.Stderr, "error rendering: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "wrote %s\n", outputPath)
}
```

- [ ] **Step 4: Build and test end-to-end (parser + renderer)**

Run:
```bash
go build -o parser ./cmd/parser && go build -o renderer ./cmd/renderer
./parser "data/Vrindavan [India].txt" data/festivals.json output/calendar.json
./renderer output/calendar.json templates output/calendar.tex
head -80 output/calendar.tex
```
Expected: valid LaTeX output with TikZ grid commands, footnotes, and the reference page.

- [ ] **Step 5: Commit**

```bash
git add cmd/renderer/ internal/renderer/ templates/
git commit -m "feat: renderer CLI with template execution"
```

---

### Task 15: Makefile and dependency checker

**Files:**
- Create: `Makefile`
- Create: `scripts/check-deps.sh`

- [ ] **Step 1: Create check-deps.sh**

Create `scripts/check-deps.sh`:

```bash
#!/usr/bin/env bash
set -euo pipefail

errors=0

check() {
  if ! command -v "$1" &>/dev/null; then
    echo "MISSING: $1 — $2"
    errors=$((errors + 1))
  else
    echo "OK: $1 ($(command -v "$1"))"
  fi
}

check go "Install Go from https://go.dev/dl/"
check lualatex "Install TeX Live: sudo apt install texlive-luatex texlive-latex-extra"

if [ "$errors" -gt 0 ]; then
  echo ""
  echo "$errors missing dependencies"
  exit 1
fi

echo ""
echo "All dependencies satisfied."
```

- [ ] **Step 2: Create Makefile**

Create `Makefile`:

```makefile
INPUT       ?= data/Vrindavan [India].txt
OUTPUT_DIR  ?= output
TEMPLATES   ?= templates
FESTIVALS   ?= data/festivals.json

JSON        := $(OUTPUT_DIR)/calendar.json
TEX         := $(OUTPUT_DIR)/calendar.tex
PDF         := $(OUTPUT_DIR)/calendar.pdf
PARSER      := parser
RENDERER    := renderer

.PHONY: all parse render pdf clean check

all: pdf

check:
	@bash scripts/check-deps.sh

$(PARSER): cmd/parser/main.go internal/parser/parser.go internal/models/models.go
	go build -o $(PARSER) ./cmd/parser

$(RENDERER): cmd/renderer/main.go internal/renderer/renderer.go internal/renderer/abbreviations.go internal/models/models.go
	go build -o $(RENDERER) ./cmd/renderer

$(OUTPUT_DIR):
	mkdir -p $(OUTPUT_DIR)

$(JSON): $(PARSER) $(FESTIVALS) $(INPUT) | $(OUTPUT_DIR)
	./$(PARSER) "$(INPUT)" $(FESTIVALS) $(JSON)

$(TEX): $(RENDERER) $(JSON) $(wildcard $(TEMPLATES)/*.tex.tmpl)
	./$(RENDERER) $(JSON) $(TEMPLATES) $(TEX)

$(PDF): $(TEX)
	cd $(OUTPUT_DIR) && lualatex -interaction=nonstopmode ../$(TEX)
	@echo "PDF written to $(PDF)"

parse: $(JSON)

render: $(TEX)

pdf: $(PDF)

clean:
	rm -rf $(OUTPUT_DIR) $(PARSER) $(RENDERER)
```

- [ ] **Step 3: Make check-deps.sh executable and test**

Run:
```bash
chmod +x scripts/check-deps.sh
make check
```
Expected: shows OK for go and lualatex (or MISSING with install instructions)

- [ ] **Step 4: Test full pipeline**

Run:
```bash
make clean && make all
```
Expected: parser runs, renderer runs, lualatex compiles, `output/calendar.pdf` created.

- [ ] **Step 5: Commit**

```bash
git add Makefile scripts/
git commit -m "feat: Makefile and dependency checker for full pipeline"
```

---

### Task 16: Template debugging and visual tuning

This task is iterative — run the pipeline, inspect the PDF, fix template issues.

**Files:**
- Modify: `templates/calendar.tex.tmpl`
- Modify: `templates/reference.tex.tmpl`
- Possibly modify: `internal/renderer/renderer.go`

- [ ] **Step 1: Run the full pipeline and open the PDF**

Run:
```bash
make clean && make all
```
Open `output/calendar.pdf` and inspect:
- Are all month pages present? (Expect ~14 pages: 13 Masas + 1 reference)
- Are day cells correctly positioned in the grid?
- Do fasting days have grey backgrounds?
- Are footnotes numbered correctly and legible?
- Does the reference page render properly?

- [ ] **Step 2: Fix any TikZ coordinate or sizing issues**

Common issues to check:
- Cell width/height ratios for A4 landscape
- TikZ coordinate units (may need `\pgfmathsetmacro` adjustments)
- Font sizes for readability
- Footnote area height vs grid area
- Empty cells at start/end of Masa properly styled

Adjust `calendar.tex.tmpl` as needed. After each change:
```bash
make render && make pdf
```

- [ ] **Step 3: Verify all 13+ Masa pages render correctly**

Scan through every page in the PDF. Each page should show:
- Correct Masa name and Gaurabda year in header
- Days in correct weekday columns
- Icons in cells where expected
- Fasting cells with grey background
- Footnotes at bottom matching cell superscripts

- [ ] **Step 4: Commit final template adjustments**

```bash
git add templates/ internal/renderer/
git commit -m "fix: template tuning for correct grid rendering"
```

---

### Task 17: README

**Files:**
- Create: `README.md`

- [ ] **Step 1: Create README.md**

Create `README.md`:

```markdown
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

1. **Parse** — `data/Vrindavan [India].txt` → `output/calendar.json`
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
    data/               Input files and festival list
    output/             Generated files
    scripts/            Helper scripts

## Calendar Features

- One page per Vedic Masa (lunar month)
- Sunday-start week grid
- Icons: ◆ Ekadasi, ★ Appearance, † Disappearance, ☀ Festival, ▲ Sankranti
- Fasting days highlighted with grey background
- Footnotes with full event descriptions
- Reference page with all abbreviations
```

- [ ] **Step 2: Commit**

```bash
git add README.md
git commit -m "docs: add README with usage instructions"
```
