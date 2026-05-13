package renderer

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/taraka/vedic-calendar/internal/models"
)

type CellData struct {
	DayNum       int
	Weekday      string
	Tithi        string
	TithiSub     string // parenthetical qualifier, e.g. "(suitable for fasting)"
	PaksaLabel   string
	Naksatra     string
	Fasting      bool
	IsEkadasi    bool
	Icons        []string
	FootnoteRefs []int
	HasSankranti bool
}

type Footnote struct {
	Number int
	Text   string
	Icon   string
}

type MonthGrid struct {
	MasaName        string
	MasaAltName     string
	GaurabdaYear    int
	Location        string
	GregorianMonths string
	Weeks           [][7]*CellData
	Footnotes       []Footnote
}

const (
	IconEkadasi       = `\faMoon`
	IconAppearance    = `\faStar`
	IconDisappearance = `\faFire`
	IconFestival      = `\faDrum`
	IconSankranti     = `\faSun`
)

func splitTithi(tithi string) (main, sub string) {
	if idx := strings.Index(tithi, "("); idx >= 0 {
		return strings.TrimSpace(tithi[:idx]), strings.TrimSpace(tithi[idx:])
	}
	return tithi, ""
}

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

func gregorianMonthRange(days []models.Day) string {
	seen := make(map[time.Month]bool)
	var ordered []time.Month
	for _, d := range days {
		m := d.Date.Month()
		if !seen[m] {
			seen[m] = true
			ordered = append(ordered, m)
		}
	}
	parts := make([]string, len(ordered))
	for i, m := range ordered {
		parts[i] = m.String()[:3]
	}
	return strings.Join(parts, "/")
}

func BuildMonthGrid(month models.VedicMonth) MonthGrid {
	grid := MonthGrid{
		MasaName:        month.MasaName,
		MasaAltName:     month.MasaAltName,
		GaurabdaYear:    month.GaurabdaYear,
		GregorianMonths: gregorianMonthRange(month.Days),
	}

	footnoteNum := 1

	if len(month.Days) == 0 {
		return grid
	}

	weekdayCol := func(wd string) int {
		cols := map[string]int{"Su": 0, "Mo": 1, "Tu": 2, "We": 3, "Th": 4, "Fr": 5, "Sa": 6}
		if c, ok := cols[wd]; ok {
			return c
		}
		return 0
	}

	currentWeek := 0
	grid.Weeks = append(grid.Weeks, [7]*CellData{})
	prevCol := weekdayCol(month.Days[0].Weekday) - 1

	for i := range month.Days {
		day := &month.Days[i]
		col := weekdayCol(day.Weekday)

		if i > 0 && col <= prevCol {
			currentWeek++
			grid.Weeks = append(grid.Weeks, [7]*CellData{})
		}
		prevCol = col

		tithiMain, tithiSub := splitTithi(day.Tithi)
		cell := &CellData{
			DayNum:     day.Date.Day(),
			Weekday:    day.Weekday,
			Tithi:      tithiMain,
			TithiSub:   tithiSub,
			PaksaLabel: day.Paksa,
			Naksatra:   day.Naksatra,
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
				Text:   fmt.Sprintf("%s --- %s", day.Sankranti.Name, day.Sankranti.Notes),
				Icon:   IconSankranti,
			})
			cell.FootnoteRefs = append(cell.FootnoteRefs, footnoteNum)
			footnoteNum++
		}

		for _, evt := range day.Events {
			icon := iconForEventType(evt.Type)
			if icon != "" {
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

			grid.Footnotes = append(grid.Footnotes, Footnote{
				Number: footnoteNum,
				Text:   evt.Description,
				Icon:   icon,
			})
			cell.FootnoteRefs = append(cell.FootnoteRefs, footnoteNum)
			footnoteNum++
		}

		if day.BreakFast != nil {
			text := fmt.Sprintf("Break fast %s -- %s (%s)", day.BreakFast.Start, day.BreakFast.End, day.BreakFast.Notes)
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
	td := TemplateData{
		Location: cal.Location,
	}
	for _, m := range cal.Months {
		g := BuildMonthGrid(m)
		g.Location = cal.Location.Name
		td.Grids = append(td.Grids, g)
	}

	calTmpl, err := template.ParseFiles(filepath.Join(templateDir, "calendar.tex.tmpl"))
	if err != nil {
		return fmt.Errorf("parse calendar template: %w", err)
	}

	if err := calTmpl.Execute(w, td); err != nil {
		return fmt.Errorf("execute calendar template: %w", err)
	}

	refContent, err := os.ReadFile(filepath.Join(templateDir, "reference.tex.tmpl"))
	if err != nil {
		return fmt.Errorf("read reference template: %w", err)
	}

	refTmpl, err := template.New("reference").Parse(string(refContent))
	if err != nil {
		return fmt.Errorf("parse reference template: %w", err)
	}
	if err := refTmpl.Execute(w, td); err != nil {
		return fmt.Errorf("execute reference template: %w", err)
	}

	fmt.Fprintln(w, `\end{document}`)

	return nil
}
