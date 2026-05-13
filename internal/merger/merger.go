package merger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/taraka/vedic-calendar/internal/models"
)

// ExtraEvent is one entry from an extra-*.json file.
type ExtraEvent struct {
	Masa        string           `json:"masa"`
	Paksa       string           `json:"paksa"`
	Tithi       string           `json:"tithi"`
	Type        models.EventType `json:"type"`
	Description string           `json:"description"`
	Highlight   bool             `json:"highlight,omitempty"`
}

// normalizePaksa maps common spellings to the single-letter codes used by the parser ("G" or "K").
// This allows extra event files to use either "G"/"K" or "shukla"/"krishna".
func normalizePaksa(p string) string {
	switch strings.ToLower(strings.TrimSpace(p)) {
	case "g", "shukla", "sukla":
		return "G"
	case "k", "krishna":
		return "K"
	default:
		return p
	}
}

// Merge injects extra events into cal in-place.
// Unmatched entries are logged to log.
// Tithi matching strips any parenthetical qualifier (e.g. "Dvitiya (suitable for fasting)" matches "Dvitiya").
func Merge(cal *models.Calendar, extras []ExtraEvent, log io.Writer) {
	// Build lookup: "masa|paksa|baseTithi" -> []*Day
	lookup := make(map[string][]*models.Day)
	for mi := range cal.Months {
		masa := cal.Months[mi].MasaName
		for di := range cal.Months[mi].Days {
			day := &cal.Months[mi].Days[di]
			base, _, _ := strings.Cut(day.Tithi, " (")
			key := masa + "|" + normalizePaksa(day.Paksa) + "|" + base
			lookup[key] = append(lookup[key], day)
		}
	}

	for _, e := range extras {
		key := e.Masa + "|" + normalizePaksa(e.Paksa) + "|" + e.Tithi
		days, ok := lookup[key]
		if !ok {
			fmt.Fprintf(log, "unmatched extra event: %s %s %s %q\n", e.Masa, e.Paksa, e.Tithi, e.Description)
			continue
		}
		evt := models.Event{
			Type:        e.Type,
			Description: e.Description,
			Highlight:   e.Highlight,
		}
		for _, day := range days {
			day.Events = append(day.Events, evt)
		}
	}
}

// LoadExtras loads and combines all extra event files.
func LoadExtras(paths []string) ([]ExtraEvent, error) {
	var all []ExtraEvent
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", p, err)
		}
		var events []ExtraEvent
		if err := json.Unmarshal(data, &events); err != nil {
			return nil, fmt.Errorf("parsing %s: %w", p, err)
		}
		all = append(all, events...)
	}
	return all, nil
}
