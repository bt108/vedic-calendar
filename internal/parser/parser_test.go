package parser

import (
	"testing"
	"time"

	"github.com/bt108/vedic-calendar/internal/models"
)

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

func TestParseMasaHeader(t *testing.T) {
	tests := []struct {
		line    string
		name    string
		altName string
		year    int
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

func TestClassifyEvent(t *testing.T) {
	major := []string{"Ratha Yatra"}
	minor := []string{"Sri Krsna Pusya Abhiseka", "Vasanta Pancami"}

	tests := []struct {
		name    string
		text    string
		evtType models.EventType
		isBF    bool
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
			evt, bf := classifyEvent(tt.text, major, minor)
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

	major := []string{}
	minor := []string{"Sri Krsna Pusya Abhiseka", "Ganga Sagara Mela"}
	cal, err := ParseCalendar(input, major, minor)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cal.Location.Name != "Vrindavan [India]" {
		t.Errorf("location name = %q", cal.Location.Name)
	}

	if len(cal.Months) != 2 {
		t.Fatalf("months = %d, want 2", len(cal.Months))
	}

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

	if m0.Days[0].BreakFast == nil {
		t.Error("day 1 should have break-fast")
	} else if m0.Days[0].BreakFast.Start != "07:09" {
		t.Errorf("break-fast start = %q", m0.Days[0].BreakFast.Start)
	}

	if len(m0.Days[2].Events) != 1 {
		t.Fatalf("day 3 events = %d, want 1", len(m0.Days[2].Events))
	}
	if m0.Days[2].Events[0].Type != models.EventFestival {
		t.Errorf("day 3 event type = %q", m0.Days[2].Events[0].Type)
	}

	m1 := cal.Months[1]
	if m1.MasaName != "Magha" {
		t.Errorf("m1 masa = %q", m1.MasaName)
	}
	if len(m1.Days) < 5 {
		t.Fatalf("m1 days = %d, want >= 5", len(m1.Days))
	}

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

	if ekDay.Sankranti == nil {
		t.Error("14 Jan should have Sankranti")
	} else if ekDay.Sankranti.Name != "Makara Sankranti" {
		t.Errorf("sankranti name = %q", ekDay.Sankranti.Name)
	}
}

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
