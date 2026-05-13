package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/taraka/vedic-calendar/internal/parser"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "usage: parser <input.txt> <festivals.json> [output.json]\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	festivalsPath := os.Args[2]

	cal, err := parser.ParseCalendarFile(inputPath, festivalsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	out, err := json.MarshalIndent(cal, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	if len(os.Args) >= 4 {
		if err := os.WriteFile(os.Args[3], out, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "error writing file: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "wrote %s\n", os.Args[3])
	} else {
		fmt.Println(string(out))
	}
}
