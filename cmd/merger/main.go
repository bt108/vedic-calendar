package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/taraka/vedic-calendar/internal/merger"
	"github.com/taraka/vedic-calendar/internal/models"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "usage: merger <calendar.json> <output.json> [extra1.json ...]\n")
		os.Exit(1)
	}

	calPath := os.Args[1]
	outPath := os.Args[2]
	extraPaths := os.Args[3:]

	data, err := os.ReadFile(calPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading calendar: %v\n", err)
		os.Exit(1)
	}
	var cal models.Calendar
	if err := json.Unmarshal(data, &cal); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing calendar: %v\n", err)
		os.Exit(1)
	}

	extras, err := merger.LoadExtras(extraPaths)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading extras: %v\n", err)
		os.Exit(1)
	}

	merger.Merge(&cal, extras, os.Stdout)

	out, err := json.MarshalIndent(cal, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshaling: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(outPath, out, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error writing output: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "wrote %s\n", outPath)
}
