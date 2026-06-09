package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/taraka/vedic-calendar/internal/models"
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

	reMasaHeader = regexp.MustCompile(`(\S+(?:-\S+)?)\s+(?:\((\w+)\)\s+)?Masa,\s+Gaurabda\s+(\d+)`)
	reLocation   = regexp.MustCompile(`(.+?)\s+\((\d+[NS]\d+\s+\d+[EW]\d+),\s*([^)]+)\)`)

	// Tithi and paksa may be separated by 2+ spaces OR zero spaces (e.g. "fasting)G Hasta")
	reDayParse = regexp.MustCompile(`^\s?(\d{1,2})\s+(\w{3})\s+(\d{4})\s+(\w{2})\s{2,3}(.+?)\s*(G|K)\s+(\S+(?:-\S+)?)\s*(\*?)\s*$`)
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

	if regexp.MustCompile(`\d+[NS]\d+\s+\d+[EW]\d+`).MatchString(trimmed) {
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

var reBreakFast = regexp.MustCompile(`^Break fast\s+(\d{2}:\d{2})\s+\([^)]+\)\s*-\s*(\d{2}:\d{2})\s+(.+)$`)

func classifyEvent(text string, major, minor []string) (models.Event, *models.BreakFast) {
	trimmed := strings.TrimSpace(text)

	// Break-fast info — structural, not an event
	if m := reBreakFast.FindStringSubmatch(trimmed); m != nil {
		return models.Event{}, &models.BreakFast{
			Start: m[1],
			End:   m[2],
			Notes: strings.TrimSpace(m[3]),
		}
	}

	// Determine TYPE by text pattern only
	var t models.EventType
	switch {
	case strings.Contains(trimmed, "-- Appearance") || strings.Contains(trimmed, "Appearance of"):
		t = models.EventAppearance
	case strings.Contains(trimmed, "-- Disappearance"):
		t = models.EventDisappearance
	case strings.HasPrefix(trimmed, "Fasting for ") || strings.HasPrefix(trimmed, "(Fast till ") || strings.HasPrefix(trimmed, "(Fasting till "):
		t = models.EventFast
	default:
		t = models.EventNote
		for _, f := range major {
			if strings.HasPrefix(trimmed, f) {
				t = models.EventFestival
				break
			}
		}
		if t == models.EventNote {
			for _, f := range minor {
				if strings.HasPrefix(trimmed, f) {
					t = models.EventFestival
					break
				}
			}
		}
	}

	// Determine HIGHLIGHT independently — any event whose description starts with a major prefix
	highlight := false
	for _, f := range major {
		if strings.HasPrefix(trimmed, f) {
			highlight = true
			break
		}
	}

	return models.Event{Type: t, Description: trimmed, Highlight: highlight}, nil
}

var reSankrantiParse = regexp.MustCompile(`^-+\s+(\w+\s+Sankranti)\s+\((.+?)\)\s+-+$`)
var reSankrantiDate = regexp.MustCompile(`on\s+(\d{1,2})\s+(\w{3})`)

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

type FestivalList struct {
	Major []string `json:"major"`
	Minor []string `json:"minor"`
}

func LoadFestivals(path string) (major, minor []string, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("read festivals: %w", err)
	}
	var fl FestivalList
	if err := json.Unmarshal(data, &fl); err != nil {
		return nil, nil, fmt.Errorf("parse festivals: %w", err)
	}
	return fl.Major, fl.Minor, nil
}

func ParseCalendarFile(inputPath, festivalsPath string) (models.Calendar, error) {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return models.Calendar{}, fmt.Errorf("read input: %w", err)
	}
	major, minor, err := LoadFestivals(festivalsPath)
	if err != nil {
		return models.Calendar{}, err
	}
	return ParseCalendar(string(data), major, minor)
}

func ParseCalendar(input string, major, minor []string) (models.Calendar, error) {
	lines := strings.Split(input, "\n")
	cal := models.Calendar{}
	var currentMonth *models.VedicMonth
	var currentDay *models.Day
	locationParsed := false

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
			evt, bf := classifyEvent(text, major, minor)
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
