INPUT       ?= Vrindavan [India].txt
OUTPUT_DIR  ?= output
TEMPLATES   ?= templates
FESTIVALS   ?= data/festivals.json

JSON        := $(OUTPUT_DIR)/calendar.json
TEX         := $(OUTPUT_DIR)/calendar.tex
PDF         := $(OUTPUT_DIR)/calendar.pdf
PARSER      := parser
RENDERER    := renderer

.PHONY: all parse render pdf clean check

all: pdf

check:
	@bash scripts/check-deps.sh

$(PARSER): $(wildcard cmd/parser/*.go) $(wildcard internal/parser/*.go) $(wildcard internal/models/*.go)
	go build -o $(PARSER) ./cmd/parser

$(RENDERER): $(wildcard cmd/renderer/*.go) $(wildcard internal/renderer/*.go) $(wildcard internal/models/*.go)
	go build -o $(RENDERER) ./cmd/renderer

$(OUTPUT_DIR):
	mkdir -p $(OUTPUT_DIR)

$(JSON): $(PARSER) $(FESTIVALS) | $(OUTPUT_DIR)
	./$(PARSER) "$(INPUT)" $(FESTIVALS) $(JSON)

$(TEX): $(RENDERER) $(JSON) $(wildcard $(TEMPLATES)/*.tex.tmpl)
	./$(RENDERER) $(JSON) $(TEMPLATES) $(TEX)

$(PDF): $(TEX)
	cd $(OUTPUT_DIR) && lualatex -interaction=nonstopmode calendar.tex
	@echo "PDF written to $(PDF)"

parse: $(JSON)

render: $(TEX)

pdf: $(PDF)

clean:
	rm -rf $(OUTPUT_DIR) $(PARSER) $(RENDERER)
