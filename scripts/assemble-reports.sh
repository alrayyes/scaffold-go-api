#!/usr/bin/env bash
# Lays the test and coverage reports out the way they're published, under
# site/reports/ (rules/published-reports.md has the layout):
#
#   reports/index.html
#   reports/tests/unit.xml
#   reports/coverage/{index.html,coverage.xml,coverage.out}
#
# Needs junit.xml and coverage.out in the working directory, which the test
# job writes. Runs on pull requests too, so a broken conversion fails before
# the merge; only the upload and the deploy are default-branch-only.
set -euo pipefail

cd "$(dirname "$0")/.."

out="${1:-site/reports}"
rm -rf "$out"
mkdir -p "$out/tests" "$out/coverage"

cp junit.xml "$out/tests/unit.xml"
cp coverage.out "$out/coverage/coverage.out"
go tool cover -html=coverage.out -o "$out/coverage/index.html"

# Cobertura is the one coverage format every viewer reads; Go can't write it
# natively. Pinned exactly. Invisible to Dependabot here, since it's
# fetched by `go run` rather than listed in go.mod: bump it by hand.
go run github.com/boumenot/gocover-cobertura@v1.5.0 \
  <coverage.out >"$out/coverage/coverage.xml"
grep -q '<coverage ' "$out/coverage/coverage.xml"

sha="${GITHUB_SHA:-$(git rev-parse HEAD)}"
date="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
cat >"$out/index.html" <<HTML
<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Reports</title>
  </head>
  <body>
    <h1>Reports</h1>
    <p>Commit <code>${sha:0:7}</code>, ${date}.</p>
    <ul>
      <li><a href="tests/unit.xml">Test results</a> (JUnit XML)</li>
      <li><a href="coverage/">Coverage</a> (HTML)</li>
      <li><a href="coverage/coverage.xml">Coverage</a> (Cobertura XML)</li>
      <li><a href="coverage/coverage.out">Coverage</a> (Go profile)</li>
    </ul>
  </body>
</html>
HTML
