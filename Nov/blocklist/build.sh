#!/bin/bash
set -e

OUTPUT="blocklist.txt"
TEMP_DIR=$(mktemp -d)
trap 'rm -rf "$TEMP_DIR"' EXIT

# Download each source
while IFS= read -r url; do
  if [ -n "$url" ]; then
    echo "Downloading $url"
    curl -fsSL "$url" -o "$TEMP_DIR/$(basename "$url" | tr -cd 'a-zA-Z0-9.')" || echo "Failed: $url"
  fi
done < sources.txt

# Extract domains
cat "$TEMP_DIR"/* | \
  grep -v '^!' | \
  grep -v '^#' | \
  sed -E 's/^\|\|//; s/\^.*$//' | \
  sed -E 's/^0\.0\.0\.0[[:space:]]+//' | \
  grep -E '^[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$' | \
  sort -u > "$OUTPUT"

echo "Generated $OUTPUT with $(wc -l < "$OUTPUT" | awk '{print $1}') domains."
