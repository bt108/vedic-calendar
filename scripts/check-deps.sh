#!/usr/bin/env bash
set -euo pipefail

errors=0

check() {
  if ! command -v "$1" &>/dev/null; then
    echo "MISSING: $1 — $2"
    errors=$((errors + 1))
  else
    echo "OK: $1 ($(command -v "$1"))"
  fi
}

check go "Install Go from https://go.dev/dl/"
check lualatex "Install TeX Live: sudo apt install texlive-luatex texlive-latex-extra"

if [ "$errors" -gt 0 ]; then
  echo ""
  echo "$errors missing dependencies"
  exit 1
fi

echo ""
echo "All dependencies satisfied."
