package render

import (
	"regexp"
	"testing"
)

var hexColor = regexp.MustCompile(`^#[0-9a-f]{6}$`)

func TestThemesAreComplete(t *testing.T) {
	if len(Themes) != 12 {
		t.Errorf("got %d themes, want 12", len(Themes))
	}
	seen := map[string]bool{}
	for _, th := range Themes {
		if seen[th.Name] {
			t.Errorf("duplicate theme name %q", th.Name)
		}
		seen[th.Name] = true

		colors := []string{th.Sky, th.Star, th.StarHot, th.Roof, th.Antenna, th.Street, th.Label, th.Text}
		colors = append(colors, th.Shades[:]...)
		colors = append(colors, th.Windows[:]...)
		for _, c := range colors {
			if !hexColor.MatchString(c) {
				t.Errorf("theme %q has bad color %q", th.Name, c)
			}
		}
	}
}

func TestThemeByName(t *testing.T) {
	if th, err := ThemeByName("Matrix"); err != nil || th.Name != "matrix" {
		t.Errorf("ThemeByName(Matrix) = %q, %v", th.Name, err)
	}
	if _, err := ThemeByName("nope"); err == nil {
		t.Error("expected an error for an unknown theme")
	}
}
