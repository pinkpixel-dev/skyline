// Package render draws a city.Scene as ANSI text or an animated SVG.
package render

import "github.com/pinkpixel-dev/neon-skyline/internal/city"

// Shared neon palette so the terminal and SVG versions match.
const (
	colSky     = "#15131c"
	colStar    = "#5e5878"
	colStarHot = "#b8b2d6"
	colRoof    = "#ff3ea5"
	colAntenna = "#2de2e6"
	colStreet  = "#ff3ea5"
	colLabel   = "#8b85a3"
	colText    = "#f1eefa"
	colDark    = "#3b3550"
)

var (
	buildingShades = [2]string{"#1f1b2b", "#26213a"}
	windowLevels   = [5]string{colDark, "#8a3b7a", "#ff3ea5", "#2de2e6", "#f4f1ff"}
)

// fg returns the foreground color for a cell.
func fg(c city.Cell) string {
	switch c.Kind {
	case city.Star:
		if c.Rune == '✦' {
			return colStarHot
		}
		return colStar
	case city.Antenna:
		return colAntenna
	case city.Roof:
		return colRoof
	case city.Window:
		return windowLevels[c.Level]
	case city.Street:
		return colStreet
	case city.Label:
		return colLabel
	}
	return colText
}

// bg returns the background color for a cell.
func bg(c city.Cell) string {
	if c.Shade >= 0 {
		return buildingShades[c.Shade]
	}
	return colSky
}
