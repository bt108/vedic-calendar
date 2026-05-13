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
