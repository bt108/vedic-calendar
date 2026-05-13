package renderer

import (
	"testing"
	"time"

	"github.com/taraka/vedic-calendar/internal/models"
)

func TestBuildMonthGrid(t *testing.T) {
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
	if len(grid.Weeks) == 0 {
		t.Fatal("no weeks")
	}
	if grid.Weeks[0][0] == nil {
		t.Fatal("Sunday cell should not be nil")
	}
	if grid.Weeks[0][0].DayNum != 4 {
		t.Errorf("day num = %d, want 4", grid.Weeks[0][0].DayNum)
	}
	if grid.Weeks[0][0].Tithi != "Pratipat" {
		t.Errorf("tithi = %q, want Pratipat", grid.Weeks[0][0].Tithi)
	}
	if grid.Weeks[0][0].Naksatra != "Punarvasu" {
		t.Errorf("naksatra = %q, want Punarvasu", grid.Weeks[0][0].Naksatra)
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
	if grid.Footnotes[0].DayNum != 7 {
		t.Errorf("fn 0 DayNum = %d, want 7", grid.Footnotes[0].DayNum)
	}
	if !grid.Footnotes[0].ShowDayNum {
		t.Error("fn 0 ShowDayNum should be true")
	}
	if grid.Footnotes[1].DayNum != 7 {
		t.Errorf("fn 1 DayNum = %d, want 7", grid.Footnotes[1].DayNum)
	}
	if grid.Footnotes[1].ShowDayNum {
		t.Error("fn 1 ShowDayNum should be false")
	}
	if grid.Weeks[0][3] == nil {
		t.Fatal("Wednesday cell should not be nil")
	}
	if !grid.Weeks[0][3].HasFootnotes {
		t.Error("Wednesday cell should have HasFootnotes=true")
	}
}

func TestSplitTithi(t *testing.T) {
	tests := []struct{ input, wantMain, wantSub string }{
		{"Ekadasi", "Ekadasi", ""},
		{"Ekadasi (suitable for fasting)", "Ekadasi", "(suitable for fasting)"},
		{"Dvadasi (suitable for fasting)", "Dvadasi", "(suitable for fasting)"},
		{"Purnima", "Purnima", ""},
	}
	for _, tt := range tests {
		gotMain, gotSub := splitTithi(tt.input)
		if gotMain != tt.wantMain || gotSub != tt.wantSub {
			t.Errorf("splitTithi(%q) = %q, %q; want %q, %q", tt.input, gotMain, gotSub, tt.wantMain, tt.wantSub)
		}
	}
}
