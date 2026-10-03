package render

import (
	"fmt"
	"html"
	"strings"

	"github.com/pinkpixel-dev/skyline/internal/city"
)

const (
	cellW    = 10
	cellH    = 18
	fontSize = 15
	padX     = 16
	padTop   = 44 // room for the header line
	padBot   = 12
)

// SVG renders the scene as a standalone animated SVG, suitable for a GitHub
// profile README. Each glyph is placed on an explicit grid so font metrics
// can't drift columns apart; some lit windows flicker and bright stars twinkle.
func SVG(s *city.Scene, t Theme) string {
	cols := 0
	if len(s.Rows) > 0 {
		cols = len(s.Rows[0])
	}
	w := cols*cellW + padX*2
	h := len(s.Rows)*cellH + padTop + padBot

	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" role="img" aria-labelledby="t">`, w, h, w, h)
	fmt.Fprintf(&b, `<title id="t">Skyline of @%s's GitHub contributions: %s in the last year</title>`,
		html.EscapeString(s.Login), commas(s.Total))
	b.WriteString(`<style>
text{font-family:ui-monospace,"SF Mono",Menlo,Consolas,"DejaVu Sans Mono",monospace;font-size:` + fmt.Sprint(fontSize) + `px;white-space:pre}
.f{animation:flick 7s infinite steps(1)}
.tw{animation:twinkle 4s infinite ease-in-out}
@keyframes flick{0%,100%{opacity:1}62%{opacity:.2}64%{opacity:1}66%{opacity:.35}70%{opacity:1}}
@keyframes twinkle{0%,100%{opacity:1}50%{opacity:.3}}
@media (prefers-reduced-motion:reduce){.f,.tw{animation:none}}
</style>`)
	fmt.Fprintf(&b, `<rect width="100%%" height="100%%" fill="%s" rx="8"/>`, t.Sky)
	fmt.Fprintf(&b, `<text x="%d" y="28" fill="%s" font-weight="700">@%s</text>`, padX, t.Roof, html.EscapeString(s.Login))
	fmt.Fprintf(&b, `<text x="%d" y="28" fill="%s">%s contributions in the last year</text>`,
		padX+(len(s.Login)+2)*cellW, t.Label, commas(s.Total))

	// Building backgrounds: one rect per horizontal run of the same shade.
	for y, row := range s.Rows {
		for x := 0; x < len(row); {
			if row[x].Shade < 0 {
				x++
				continue
			}
			end := x
			for end < len(row) && row[end].Shade == row[x].Shade {
				end++
			}
			fmt.Fprintf(&b, `<rect x="%d" y="%d" width="%d" height="%d" fill="%s"/>`,
				padX+x*cellW, padTop+y*cellH, (end-x)*cellW, cellH, t.bg(row[x]))
			x = end
		}
	}

	// Glyphs.
	for y, row := range s.Rows {
		baseline := padTop + y*cellH + cellH - 5
		for x, c := range row {
			if c.Rune == ' ' {
				continue
			}
			class := ""
			switch {
			case c.Flicker:
				class = fmt.Sprintf(` class="f" style="animation-delay:-%.1fs"`, float64((x*7+y*13)%70)/10)
			case c.Kind == city.Star && c.Rune == '✦':
				class = fmt.Sprintf(` class="tw" style="animation-delay:-%.1fs"`, float64((x*3+y*5)%40)/10)
			}
			fmt.Fprintf(&b, `<text x="%d" y="%d" fill="%s"%s>%s</text>`,
				padX+x*cellW, baseline, t.fg(c), class, html.EscapeString(string(c.Rune)))
		}
	}

	b.WriteString(`</svg>`)
	return b.String()
}
