package render

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/pinkpixel-dev/skyline/internal/city"
)

// ANSI renders the scene as styled terminal text. Runs of identically styled
// cells share one style so the output stays light on escape codes.
func ANSI(s *city.Scene, t Theme) string {
	var b strings.Builder

	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(t.Roof)).Render("@" + s.Login)
	meta := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Label)).
		Render(fmt.Sprintf("  %s contributions in the last year", commas(s.Total)))
	b.WriteString(title + meta + "\n\n")

	for _, row := range s.Rows {
		start := 0
		for x := 1; x <= len(row); x++ {
			if x < len(row) && t.fg(row[x]) == t.fg(row[start]) && t.bg(row[x]) == t.bg(row[start]) {
				continue
			}
			var run strings.Builder
			for _, c := range row[start:x] {
				run.WriteRune(c.Rune)
			}
			style := lipgloss.NewStyle().
				Foreground(lipgloss.Color(t.fg(row[start]))).
				Background(lipgloss.Color(t.bg(row[start])))
			b.WriteString(style.Render(run.String()))
			start = x
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func commas(n int) string {
	s := fmt.Sprint(n)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return s
}
