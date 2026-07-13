#!/bin/bash
set -euo pipefail

# Per-package READMEs are generated: godocdown renders the API body, and this
# script replaces godocdown's broken default header (`--` / `import "."`) with
# the standard expectto/be header: centered package logo, backlink to the root
# README, and the real import path.
#
# Package logos (be_ctx/logo.svg, ...) are hand-maintained SVGs, not generated
# here. Edit the READMEs via doc comments in the Go sources, then re-run this.

MAX_LOGO_WIDTH=420

gen_readme() {
    local pkg="$1" out="$2" logo="$3" back="$4" import_path="$5" alt="$6"

    local width
    width=$(sed -n 's/.*viewBox="0 0 \([0-9]*\).*/\1/p' "$(dirname "$out")/$logo")
    if [ "$width" -gt "$MAX_LOGO_WIDTH" ]; then width=$MAX_LOGO_WIDTH; fi

    {
        printf '<p align="center">\n  <img src="%s" alt="%s" width="%s">\n</p>\n\n' "$logo" "$alt" "$width"
        printf '<div align="center">\n\nPart of [`expectto/be`](%s) - composable test matchers for Go.\n\n</div>\n\n---\n\n' "$back"
        printf '```go\nimport "%s"\n```\n\n' "$import_path"
        # drop godocdown's `# pkg` / `--` / `    import "."` header + blank lines after it
        godocdown "$pkg" | tail -n +4 | awk 'NF{body=1} body'
    } > "$out"
}

for pkg in be_ctx be_http be_json be_jwt be_math be_reflected be_string be_struct be_time be_url; do
    gen_readme "$pkg" "$pkg/README.md" "logo.svg" "../README.md" "github.com/expectto/be/$pkg" "$pkg"
done

# Core matchers doc lives at the repo root next to the main logo
gen_readme "." "core-be-matchers.md" "logo.svg" "README.md" "github.com/expectto/be" "be"

# MATCHERS.md: the flat, all-packages catalog grouped by intent (the primary
# discovery surface - per-package READMEs require knowing the package exists).
go run ./internal/docgen > MATCHERS.md
