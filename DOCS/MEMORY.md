# Memory

## October 3, 2026: Shared grid, two renderers

**Decided:** `internal/city` produces a grid of typed cells. The ANSI and SVG renderers only map cells to colors and output.
**Why:** The terminal and profile versions need to look the same, and keeping layout in one place means one set of tests covers both.
**Rejected:** Separate drawing code per output. It would drift apart the first time the layout changed.

## October 3, 2026: Glyphs in the SVG, not shapes

**Decided:** The SVG draws real text characters, each placed on an explicit grid.
**Why:** The whole point is an ASCII city. Explicit x positions keep columns aligned even when the viewer's monospace font has different metrics.
**Rejected:** Drawing windows as `<rect>`s. It's crisper and doesn't depend on fonts, but it stops looking like ASCII.

## October 3, 2026: Publish to an `output` branch

**Decided:** The profile workflow force-pushes a single commit to an orphan `output` branch.
**Why:** It keeps the profile repo's main branch free of nightly image commits, and needs no third-party actions.

## October 3, 2026: Renamed to skyline, themes are built in

**Decided:** The project is called `skyline` (was `neon-skyline`), and `neon` became the default theme out of 10 built-in palettes picked with `-theme`.
**Why:** Once the city came in more than one color, "neon" described a single theme rather than the whole tool.
**Rejected:** A JSON theme file option for now. Ten presets cover plenty of variety without making anyone hand-pick a dozen hex values.
