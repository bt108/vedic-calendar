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
