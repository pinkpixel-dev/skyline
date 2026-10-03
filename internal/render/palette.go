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
		Name: "matrix", Sky: "#0f1410", Star: "#3a5a44", StarHot: "#9fd8ad",
		Roof: "#3ddc84", Antenna: "#b6ff9e", Street: "#2fbf6a", Label: "#5f8a6b", Text: "#e4f5e8",
		Shades:  [2]string{"#142017", "#182619"},
		Windows: [5]string{"#22382a", "#1f6b3a", "#2fbf6a", "#6dff9a", "#d6ffe0"},
	},
	{
		Name: "amber", Sky: "#15120d", Star: "#5a4c33", StarHot: "#d8c08a",
		Roof: "#ffb000", Antenna: "#ffd27a", Street: "#ff9a00", Label: "#8f7a55", Text: "#f7ecd6",
		Shades:  [2]string{"#1f1a12", "#251f15"},
		Windows: [5]string{"#3a2f1c", "#7a5212", "#c98410", "#ffb000", "#ffe2a8"},
	},
	{
		Name: "ice", Sky: "#10161c", Star: "#4a5d70", StarHot: "#c5d8e8",
		Roof: "#7fd4ff", Antenna: "#e6f6ff", Street: "#4fb8ff", Label: "#7a8fa3", Text: "#e8f2fa",
		Shades:  [2]string{"#17212b", "#1c2833"},
		Windows: [5]string{"#2a3a48", "#2f5d80", "#4fb8ff", "#a6e3ff", "#f0fbff"},
	},
	{
		Name: "sunset", Sky: "#1c1418", Star: "#5e464c", StarHot: "#e3c4b8",
		Roof: "#ff6b5b", Antenna: "#ffd166", Street: "#ff8c42", Label: "#9a7f85", Text: "#f8ebe6",
		Shades:  [2]string{"#261a1f", "#2d1f25"},
		Windows: [5]string{"#42292f", "#8c3b3b", "#ff6b5b", "#ff9f5a", "#ffd9a0"},
	},
	{
		Name: "toxic", Sky: "#13150f", Star: "#4c5236", StarHot: "#d3dbad",
		Roof: "#c6ff3d", Antenna: "#f5ff7a", Street: "#a4e600", Label: "#8a9370", Text: "#f0f5e0",
		Shades:  [2]string{"#1b1e14", "#20241a"},
		Windows: [5]string{"#33381f", "#5a7a1a", "#9ad11f", "#c6ff3d", "#f2ffc7"},
	},
	{
		Name: "vapor", Sky: "#17151f", Star: "#5a5578", StarHot: "#d5d0f0",
		Roof: "#ffb3d9", Antenna: "#a0f0ed", Street: "#b8a9ff", Label: "#9690b3", Text: "#f3f0fc",
		Shades:  [2]string{"#221f2e", "#282436"},
		Windows: [5]string{"#3a3550", "#6e6299", "#b8a9ff", "#ffb3d9", "#e8fffe"},
	},
	{
		Name: "crimson", Sky: "#170f11", Star: "#523a3f", StarHot: "#d6b8bd",
		Roof: "#ff2e4d", Antenna: "#ff9aa8", Street: "#d61f3c", Label: "#8f6f75", Text: "#f7e9eb",
		Shades:  [2]string{"#211417", "#28181c"},
		Windows: [5]string{"#3d2328", "#6b1d29", "#b3203a", "#ff2e4d", "#ffd1d8"},
	},
	{
		Name: "mono", Sky: "#141518", Star: "#4a4e57", StarHot: "#c4c8d0",
		Roof: "#e8eaf0", Antenna: "#f4f5f8", Street: "#b8bcc6", Label: "#7d828c", Text: "#eef0f4",
		Shades:  [2]string{"#1c1e22", "#222429"},
		Windows: [5]string{"#33363d", "#555a66", "#8a909c", "#c4c8d0", "#f4f5f8"},
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
