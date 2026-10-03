package render

import (
	"encoding/xml"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/pinkpixel-dev/skyline/internal/city"
	"github.com/pinkpixel-dev/skyline/internal/github"
)

func TestSVGIsWellFormed(t *testing.T) {
	start, _ := time.Parse("2006-01-02", "2026-01-04")
	cal := &github.Calendar{Login: `a<b&"c"`, Total: 1234}
	for w := range 6 {
		var week []github.Day
		for d := range 7 {
			n := (w * d) % 9
			week = append(week, github.Day{Date: start.AddDate(0, 0, w*7+d), Count: n, Level: min(4, n/2)})
		}
		cal.Weeks = append(cal.Weeks, week)
	}

	out := SVG(city.Build(cal, city.Options{Height: 8}), Themes[0])

	dec := xml.NewDecoder(strings.NewReader(out))
	for {
		if _, err := dec.Token(); err == io.EOF {
			break
		} else if err != nil {
			t.Fatalf("invalid XML: %v", err)
		}
	}
	for _, want := range []string{"1,234 contributions", "prefers-reduced-motion", "<title"} {
		if !strings.Contains(out, want) {
			t.Errorf("SVG missing %q", want)
		}
	}
}
