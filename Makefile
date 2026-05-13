# Vedic Calendar PDF Generator
#
# USAGE
#   make                        build PDF from the most recently modified .txt in data/
#   make FILE=data/Tokyo\ [Japan].txt   build PDF from a specific input file
#   make parse                  text → JSON only
#   make render                 JSON → .tex only
#   make pdf                    .tex → PDF only (same as default)
#   make clean                  remove all output/ dirs and compiled binaries
#   make check                  verify go + lualatex are installed
#
# INPUT FILES
#   Place .txt calendar files (exported from GCal) in data/.
#   The default target picks the most recently modified one automatically.
#   Switch inputs by either touching/replacing the file, or passing FILE=.
#
# OUTPUT
#   Each input produces its own subdirectory under output/, e.g.:
#     data/Tokyo [Japan].txt  →  output/Tokyo_Japan/calendar-Tokyo_Japan.pdf
#   Non-alphanumeric characters (spaces, brackets, hyphens) become underscores.

TEMPLATES ?= templates
FESTIVALS ?= data/festivals.json

# Auto-detect most recently modified .txt in data/, or override with FILE=
FILE ?= $(shell ls -t data/*.txt 2>/dev/null | head -1)

ifeq ($(FILE),)
$(error No .txt files found in data/ — place a calendar file there or use FILE=path/to/file.txt)
endif

# Derive a sanitized name from the input filename (no dir, no extension).
# Runs of non-alphanumeric characters are collapsed to a single underscore,
# and any leading/trailing underscores are trimmed.
# e.g. "Vrindavan [India]-2027" → "Vrindavan_India_2027"
_BASENAME := $(notdir $(basename $(FILE)))
NAME      := $(shell printf '%s' '$(_BASENAME)' | tr -cs 'A-Za-z0-9' '_' | sed 's/^_*//;s/_*$$//')

OUTPUT_DIR := output/$(NAME)
JSON       := $(OUTPUT_DIR)/calendar.json
TEX        := $(OUTPUT_DIR)/calendar-$(NAME).tex
PDF        := $(OUTPUT_DIR)/calendar-$(NAME).pdf
PARSER     := parser
RENDERER   := renderer

.PHONY: all parse render pdf clean check

all: pdf

check:
	@bash scripts/check-deps.sh

# Build the parser binary from source (rebuilt only when Go source changes)
$(PARSER): $(wildcard cmd/parser/*.go) $(wildcard internal/parser/*.go) $(wildcard internal/models/*.go)
	go build -o $(PARSER) ./cmd/parser

# Build the renderer binary from source (rebuilt only when Go source changes)
$(RENDERER): $(wildcard cmd/renderer/*.go) $(wildcard internal/renderer/*.go) $(wildcard internal/models/*.go)
	go build -o $(RENDERER) ./cmd/renderer

$(OUTPUT_DIR):
	mkdir -p $(OUTPUT_DIR)

# Stage 1: parse .txt → calendar.json
$(JSON): $(PARSER) $(FESTIVALS) | $(OUTPUT_DIR)
	@echo "Input:  $(FILE)"
	@echo "Output: $(PDF)"
	./$(PARSER) "$(FILE)" $(FESTIVALS) $(JSON)

# Stage 2: render calendar.json → .tex (using templates/)
$(TEX): $(RENDERER) $(JSON) $(wildcard $(TEMPLATES)/*.tex.tmpl)
	./$(RENDERER) $(JSON) $(TEMPLATES) $(TEX)

# Stage 3: compile .tex → .pdf via lualatex
# lualatex is run from within the output dir so auxiliary files (.aux, .log)
# stay co-located with the .tex rather than polluting the project root.
$(PDF): $(TEX)
	cd $(OUTPUT_DIR) && lualatex -interaction=nonstopmode calendar-$(NAME).tex
	@echo "PDF written to $(PDF)"

parse: $(JSON)

render: $(TEX)

pdf: $(PDF)

# Removes all generated output and compiled binaries.
# Each input file gets its own output/ subdirectory, so this wipes them all.
clean:
	rm -rf output $(PARSER) $(RENDERER)
