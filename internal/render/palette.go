// Package render draws a city.Scene as ANSI text or an animated SVG.
package render

import (
	"fmt"
	"strings"

	"github.com/pinkpixel-dev/skyline/internal/city"
)

// Theme is a full city palette. Both renderers read from it, so the terminal
// and SVG versions of a theme always match.
type Theme struct {
	Name    string
	Sky     string
	Star    string
	StarHot string // bright ✦ stars
	Roof    string
	Antenna string
	Street  string
	Label   string
	Text    string
	Shades  [2]string // alternating building backgrounds
	Windows [5]string // dark window, then GitHub levels 1-4
}

// Themes lists every built-in palette in display order. The first is the default.
var Themes = []Theme{
	{
		Name: "neon", Sky: "#15131c", Star: "#5e5878", StarHot: "#b8b2d6",
		Roof: "#ff3ea5", Antenna: "#2de2e6", Street: "#ff3ea5", Label: "#8b85a3", Text: "#f1eefa",
		Shades:  [2]string{"#1f1b2b", "#26213a"},
		Windows: [5]string{"#3b3550", "#8a3b7a", "#ff3ea5", "#2de2e6", "#f4f1ff"},
	},
	{
		Name: "synthwave", Sky: "#1a1330", Star: "#6f5f9a", StarHot: "#d9ccf5",
		Roof: "#ff7edb", Antenna: "#fede5d", Street: "#f97e72", Label: "#8d7fb3", Text: "#f4eefe",
		Shades:  [2]string{"#241a3f", "#2b1f4a"},
		Windows: [5]string{"#3d2d5c", "#6b4a9e", "#ff7edb", "#fede5d", "#fff4c9"},
	},
	{
		Name: "matrix", Sky: "#050805", Star: "#1f3a26", StarHot: "#8dffb0",
		Roof: "#00ff66", Antenna: "#d4ff3a", Street: "#00c853", Label: "#4c8058", Text: "#dcffe4",
		Shades:  [2]string{"#0b120c", "#0f1810"},
		Windows: [5]string{"#163020", "#11803d", "#00e676", "#c6ff00", "#eaffef"},
	},
	{
		Name: "amber", Sky: "#0b0e14", Star: "#363b46", StarHot: "#c9c6bd",
		Roof: "#ff8f40", Antenna: "#ffd580", Street: "#e6b450", Label: "#6c7380", Text: "#e6e1cf",
		Shades:  [2]string{"#12161e", "#171c25"},
		Windows: [5]string{"#252b36", "#a3474f", "#f07178", "#ff8f40", "#ffd580"},
	},
	{
		Name: "ice", Sky: "#0c0f14", Star: "#34404f", StarHot: "#cfe6ff",
		Roof: "#9ccfff", Antenna: "#f2fbff", Street: "#3d8bff", Label: "#6b7d93", Text: "#e9f3fc",
		Shades:  [2]string{"#131922", "#18202b"},
		Windows: [5]string{"#212d3c", "#1f4e99", "#3d8bff", "#7fe3ff", "#f2fbff"},
	},
	{
		Name: "sunset", Sky: "#0f0a0d", Star: "#4a333d", StarHot: "#f2d0c4",
		Roof: "#ff5e62", Antenna: "#ffe066", Street: "#ff9e44", Label: "#8a6c78", Text: "#fbeee8",
		Shades:  [2]string{"#1a1116", "#21141b"},
		Windows: [5]string{"#311e28", "#8c2f66", "#ff4f6d", "#ff9e44", "#ffe066"},
	},
	{
		Name: "toxic", Sky: "#070906", Star: "#2f3a22", StarHot: "#e2f5b0",
		Roof: "#b8ff00", Antenna: "#c77dff", Street: "#8fd400", Label: "#748060", Text: "#f1f8dc",
		Shades:  [2]string{"#0f130b", "#14190f"},
		Windows: [5]string{"#212a17", "#5c8a00", "#b8ff00", "#b56cff", "#f2ffc7"},
	},
	{
		Name: "vapor", Sky: "#0e0c14", Star: "#3e3858", StarHot: "#e2dcff",
		Roof: "#ff9de2", Antenna: "#7df9ff", Street: "#b29dff", Label: "#8580a6", Text: "#f5f2ff",
		Shades:  [2]string{"#17141f", "#1d1927"},
		Windows: [5]string{"#2b263b", "#6d5cb0", "#ff9de2", "#7df9ff", "#fff3b8"},
	},
	{
		Name: "crimson", Sky: "#0a0607", Star: "#3d2328", StarHot: "#e8b8bf",
		Roof: "#ff2e4d", Antenna: "#ffb199", Street: "#c4142f", Label: "#7f5a60", Text: "#f9e8ea",
		Shades:  [2]string{"#140b0d", "#1a0e11"},
		Windows: [5]string{"#2a1519", "#7a1426", "#e0193a", "#ff6b3d", "#ffd9c7"},
	},
	{
		Name: "mono", Sky: "#0c0d0f", Star: "#363940", StarHot: "#d0d4db",
		Roof: "#f0f2f5", Antenna: "#fafbfc", Street: "#9aa0aa", Label: "#6b707a", Text: "#eef0f4",
		Shades:  [2]string{"#15171a", "#1b1d21"},
		Windows: [5]string{"#272a30", "#4d525b", "#8b919c", "#c9ced6", "#fafbfc"},
	},
	{
		Name: "prism", Sky: "#121118", Star: "#4d4866", StarHot: "#d6d0f0",
		Roof: "#b14dff", Antenna: "#ff7a00", Street: "#00e5ff", Label: "#8a84a3", Text: "#f3f0fb",
		Shades:  [2]string{"#1b1924", "#221f2d"},
		Windows: [5]string{"#2f2b3e", "#ff2e88", "#ffe600", "#39ff14", "#00e5ff"},
	},
	{
		Name: "rainbow", Sky: "#050505", Star: "#38383e", StarHot: "#d4d4dc",
		Roof: "#9b5de5", Antenna: "#ffd23f", Street: "#ff7f2a", Label: "#76767e", Text: "#f2f2f2",
		Shades:  [2]string{"#0e0e10", "#141416"},
		Windows: [5]string{"#222228", "#e63946", "#ffb627", "#3fbf5f", "#3d9bff"},
	},
}

// ThemeByName looks up a built-in theme, case-insensitively.
func ThemeByName(name string) (Theme, error) {
	names := make([]string, len(Themes))
	for i, t := range Themes {
		if strings.EqualFold(t.Name, name) {
			return t, nil
		}
		names[i] = t.Name
	}
	return Theme{}, fmt.Errorf("unknown theme %q (try: %s)", name, strings.Join(names, ", "))
}

// fg returns the foreground color for a cell.
func (t Theme) fg(c city.Cell) string {
	switch c.Kind {
	case city.Star:
		if c.Rune == '✦' {
			return t.StarHot
		}
		return t.Star
	case city.Antenna:
		return t.Antenna
	case city.Roof:
		return t.Roof
	case city.Window:
		return t.Windows[c.Level]
	case city.Street:
		return t.Street
	case city.Label:
		return t.Label
	}
	return t.Text
}

// bg returns the background color for a cell.
func (t Theme) bg(c city.Cell) string {
	if c.Shade >= 0 {
		return t.Shades[c.Shade]
	}
	return t.Sky
}
