# Overview

neon-skyline is a small Go CLI that draws a GitHub contribution graph as a neon city, either in the terminal or as an animated SVG.

## Flow

1. `main.go` parses flags and grabs a token (`GITHUB_TOKEN`, `GH_TOKEN`, then `gh auth token`).
2. `internal/github` sends one GraphQL query for `contributionCalendar` and returns weeks of days, each with a count and a 0 to 4 level (mapped from GitHub's `contributionLevel` enum).
3. `internal/city` builds a `Scene`: a 2D grid of `Cell`s. This is the only place that decides layout.
4. `internal/render` turns the scene into ANSI (`ANSI`) or SVG (`SVG`). Both read colors from `palette.go`, so the two outputs match.

## Layout rules (`internal/city`)

- Each week is a building 2 columns wide. Weeks alternate between two background shades so neighbors stay readable.
- Height is `1 + round(log1p(total) / log1p(max) * (height - 1))`. A week with zero contributions is a 1-row stub.
- Windows fill bottom to top, left to right, cycling through the week's days. Lit windows carry the day's level.
- About 1 in 9 lit windows gets `Flicker`, chosen by a deterministic hash. Stars use the same hash.
- The week with the highest total gets an antenna.
- Month labels go under the first week that ends in a new month, skipping the partial first week and anything that would overlap.

Grid rows from top to bottom: `height + 2` sky rows (headroom for roof and antenna), one street row, one label row.

## SVG output

Every glyph is its own `<text>` element at an explicit x position, so font width differences can't knock columns out of line. Building backgrounds are merged into one `<rect>` per horizontal run. The flicker and twinkle animations are plain CSS keyframes inside the SVG, turned off under `prefers-reduced-motion`.

## Profile automation

`examples/profile-workflow.yml` runs nightly in a user's profile repo. It runs `go run github.com/pinkpixel-dev/neon-skyline@latest` and force-pushes `skyline.svg` to an orphan `output` branch. It uses `SKYLINE_TOKEN` when that secret exists and falls back to the workflow's `GITHUB_TOKEN`.

## Tests

- `internal/city/city_test.go` covers heights, the antenna, window levels, the `-weeks` trim, month labels, and determinism.
- `internal/render/svg_test.go` checks that the SVG parses as XML, including a login with characters that need escaping.
