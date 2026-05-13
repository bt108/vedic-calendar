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
	if grid.Weeks[0][0].TithiShort != "Pra" {
		t.Errorf("tithi short = %q, want Pra", grid.Weeks[0][0].TithiShort)
	}
	if grid.Weeks[0][0].NaksatraShort != "Pun" {
		t.Errorf("naksatra short = %q, want Pun", grid.Weeks[0][0].NaksatraShort)
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
	if grid.Weeks[0][3] == nil {
		t.Fatal("Wednesday cell should not be nil")
	}
	if len(grid.Weeks[0][3].FootnoteRefs) != 2 {
		t.Fatalf("cell fn refs = %d, want 2", len(grid.Weeks[0][3].FootnoteRefs))
	}
	if grid.Weeks[0][3].FootnoteRefs[0] != 1 {
		t.Errorf("cell fn ref 0 = %d", grid.Weeks[0][3].FootnoteRefs[0])
	}
}

func TestShortenTithi(t *testing.T) {
	tests := []struct{ input, want string }{
		{"Ekadasi", "Eka"},
		{"Ekadasi (suitable for fasting)", "Eka"},
		{"Purnima", "Pur"},
		{"Amavasya", "Ama"},
	}
	for _, tt := range tests {
		got := ShortenTithi(tt.input)
		if got != tt.want {
			t.Errorf("ShortenTithi(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestShortenNaksatra(t *testing.T) {
	tests := []struct{ input, want string }{
		{"Rohini", "Roh"},
		{"Purva-phalguni", "PPh"},
		{"Uttara-bhadra", "UBh"},
	}
	for _, tt := range tests {
		got := ShortenNaksatra(tt.input)
		if got != tt.want {
			t.Errorf("ShortenNaksatra(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
