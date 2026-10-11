package merger

import (
	"bytes"
	"testing"
	"time"

	"github.com/bt108/vedic-calendar/internal/models"
)

func TestMergeNoExtras(t *testing.T) {
	cal := &models.Calendar{
		Months: []models.VedicMonth{
			{MasaName: "Kartika", Days: []models.Day{
				{Date: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC), Paksa: "shukla", Tithi: "Pratipat"},
			}},
		},
	}
	var log bytes.Buffer
	Merge(cal, nil, &log)
	if len(cal.Months[0].Days[0].Events) != 0 {
		t.Errorf("expected no events, got %d", len(cal.Months[0].Days[0].Events))
	}
	if log.Len() != 0 {
		t.Errorf("expected no log output, got %q", log.String())
	}
}

func TestMergeMatchingEntry(t *testing.T) {
	cal := &models.Calendar{
		Months: []models.VedicMonth{
			{MasaName: "Kartika", Days: []models.Day{
				{Date: time.Date(2026, 10, 2, 0, 0, 0, 0, time.UTC), Paksa: "shukla", Tithi: "Dvitiya"},
			}},
		},
	}
	extras := []ExtraEvent{
		{Masa: "Kartika", Paksa: "shukla", Tithi: "Dvitiya", Type: models.EventAppearance, Description: "Srila Prabhupada", Highlight: true},
	}
	var log bytes.Buffer
	Merge(cal, extras, &log)
	events := cal.Months[0].Days[0].Events
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Description != "Srila Prabhupada" {
		t.Errorf("description = %q", events[0].Description)
	}
	if events[0].Type != models.EventAppearance {
		t.Errorf("type = %q", events[0].Type)
	}
	if !events[0].Highlight {
		t.Error("expected Highlight=true")
	}
	if log.Len() != 0 {
		t.Errorf("unexpected log output: %q", log.String())
	}
}

func TestMergeNonMatchingEntry(t *testing.T) {
	cal := &models.Calendar{
		Months: []models.VedicMonth{
			{MasaName: "Kartika", Days: []models.Day{
				{Date: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC), Paksa: "shukla", Tithi: "Pratipat"},
			}},
		},
	}
	extras := []ExtraEvent{
		{Masa: "Magha", Paksa: "shukla", Tithi: "Dvitiya", Type: models.EventAppearance, Description: "Unknown Saint", Highlight: true},
	}
	var log bytes.Buffer
	Merge(cal, extras, &log)
	if len(cal.Months[0].Days[0].Events) != 0 {
		t.Errorf("expected no events, got %d", len(cal.Months[0].Days[0].Events))
	}
	if log.Len() == 0 {
		t.Error("expected log output for unmatched entry")
	}
}

func TestMergeParentheticalTithi(t *testing.T) {
	// Tithi in parsed calendar may have a parenthetical qualifier,
	// e.g. "Dvitiya (suitable for fasting)". Extra events specify just "Dvitiya".
	cal := &models.Calendar{
		Months: []models.VedicMonth{
			{MasaName: "Kartika", Days: []models.Day{
				{Date: time.Date(2026, 10, 2, 0, 0, 0, 0, time.UTC), Paksa: "shukla", Tithi: "Dvitiya (suitable for fasting)"},
			}},
		},
	}
	extras := []ExtraEvent{
		{Masa: "Kartika", Paksa: "shukla", Tithi: "Dvitiya", Type: models.EventAppearance, Description: "Saint X", Highlight: true},
	}
	var log bytes.Buffer
	Merge(cal, extras, &log)
	if len(cal.Months[0].Days[0].Events) != 1 {
		t.Fatalf("expected 1 event (parenthetical stripped), got %d", len(cal.Months[0].Days[0].Events))
	}
}

func TestMergeHighlightFalse(t *testing.T) {
	cal := &models.Calendar{
		Months: []models.VedicMonth{
			{MasaName: "Kartika", Days: []models.Day{
				{Date: time.Date(2026, 10, 2, 0, 0, 0, 0, time.UTC), Paksa: "shukla", Tithi: "Dvitiya"},
			}},
		},
	}
	extras := []ExtraEvent{
		{Masa: "Kartika", Paksa: "shukla", Tithi: "Dvitiya", Type: models.EventAppearance, Description: "Saint Y", Highlight: false},
	}
	var log bytes.Buffer
	Merge(cal, extras, &log)
	if len(cal.Months[0].Days[0].Events) == 0 {
		t.Fatal("expected 1 event, got 0")
	}
	if cal.Months[0].Days[0].Events[0].Highlight {
		t.Error("expected Highlight=false when not set in extra file")
	}
}

func TestMergeNormalizePaksa(t *testing.T) {
	// Calendar stores "G" (Gaura/Shukla); extra file uses natural-language "shukla"
	cal := &models.Calendar{
		Months: []models.VedicMonth{
			{MasaName: "Kartika", Days: []models.Day{
				{Date: time.Date(2026, 10, 2, 0, 0, 0, 0, time.UTC), Paksa: "G", Tithi: "Dvitiya"},
			}},
		},
	}
	extras := []ExtraEvent{
		{Masa: "Kartika", Paksa: "shukla", Tithi: "Dvitiya", Type: models.EventAppearance, Description: "Saint Z", Highlight: true},
	}
	var log bytes.Buffer
	Merge(cal, extras, &log)
	if len(cal.Months[0].Days[0].Events) != 1 {
		t.Fatalf("expected 1 event (shukla normalized to G), got %d", len(cal.Months[0].Days[0].Events))
	}
}
