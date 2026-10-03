package city

import (
	"testing"
	"time"

	"github.com/pinkpixel-dev/skyline/internal/github"
)

// calendar builds weeks starting on a Sunday, each day with the given counts.
func calendar(start string, weeks [][]int) *github.Calendar {
	d, _ := time.Parse("2006-01-02", start)
	cal := &github.Calendar{Login: "octo"}
	for _, counts := range weeks {
		var week []github.Day
		for _, n := range counts {
			level := 0
			if n > 0 {
				level = min(4, 1+n/5)
			}
			week = append(week, github.Day{Date: d, Count: n, Level: level})
			cal.Total += n
			d = d.AddDate(0, 0, 1)
		}
		cal.Weeks = append(cal.Weeks, week)
	}
	return cal
}

func zeros() []int { return []int{0, 0, 0, 0, 0, 0, 0} }

func buildingHeight(s *Scene, week int) int {
	street := len(s.Rows) - 2
	h := 0
	for y := street - 1; y >= 0 && s.Rows[y][week*colWidth].Kind == Window; y-- {
		h++
	}
	return h
}

func TestHeightsFollowTotals(t *testing.T) {
	cal := calendar("2026-01-04", [][]int{zeros(), {1, 0, 0, 0, 0, 0, 0}, {20, 20, 20, 20, 20, 20, 20}})
	s := Build(cal, Options{Height: 10})

	if got := buildingHeight(s, 0); got != 1 {
		t.Errorf("empty week height = %d, want 1", got)
	}
	if got := buildingHeight(s, 2); got != 10 {
		t.Errorf("busiest week height = %d, want 10", got)
	}
	if a, b := buildingHeight(s, 1), buildingHeight(s, 2); a >= b {
		t.Errorf("quiet week (%d) should be shorter than busy week (%d)", a, b)
	}
}

func TestAntennaMarksRecordWeek(t *testing.T) {
	cal := calendar("2026-01-04", [][]int{{3, 0, 0, 0, 0, 0, 0}, {9, 9, 0, 0, 0, 0, 0}, {1, 0, 0, 0, 0, 0, 0}})
	s := Build(cal, Options{Height: 8})

	antennas := 0
	for _, row := range s.Rows {
		for x, c := range row {
			if c.Kind == Antenna {
				antennas++
				if x/colWidth != 1 {
					t.Errorf("antenna on week %d, want week 1", x/colWidth)
				}
			}
		}
	}
	if antennas != 1 {
		t.Errorf("found %d antennas, want 1", antennas)
	}
}

func TestWindowsUseDayLevels(t *testing.T) {
	// Only Sunday has commits, so lit windows must carry Sunday's level.
	cal := calendar("2026-01-04", [][]int{{12, 0, 0, 0, 0, 0, 0}})
	s := Build(cal, Options{Height: 6})

	street := len(s.Rows) - 2
	first := s.Rows[street-1][0] // first window = first day of the week
	if first.Level != 3 || first.Rune != '▪' {
		t.Errorf("first window = %+v, want lit level 3", first)
	}
	if second := s.Rows[street-1][1]; second.Level != 0 || second.Rune != '·' {
		t.Errorf("second window = %+v, want dark", second)
	}
}

func TestWeeksOptionKeepsMostRecent(t *testing.T) {
	cal := calendar("2026-01-04", [][]int{zeros(), zeros(), zeros(), {5, 0, 0, 0, 0, 0, 0}})
	s := Build(cal, Options{Height: 5, Weeks: 2})

	if got, want := len(s.Rows[0]), 2*colWidth; got != want {
		t.Fatalf("width = %d, want %d", got, want)
	}
	if buildingHeight(s, 1) <= 1 {
		t.Error("last kept week should be the busy one")
	}
}

func TestMonthLabels(t *testing.T) {
	var weeks [][]int
	for range 10 {
		weeks = append(weeks, zeros())
	}
	// Dec 28 2025 is a Sunday; Jan shows up in week 0, Feb in week 5 (ends Feb 7).
	s := Build(calendar("2025-12-28", weeks), Options{Height: 4})
	labels := s.Rows[len(s.Rows)-1]

	text := func(x int) string { return string([]rune{labels[x].Rune, labels[x+1].Rune, labels[x+2].Rune}) }
	if got := text(5 * colWidth); got != "Feb" {
		t.Errorf("label at week 5 = %q, want Feb", got)
	}
	if labels[0].Rune != ' ' {
		t.Error("first partial week should not get a label")
	}
}

func TestDeterministic(t *testing.T) {
	cal := calendar("2026-01-04", [][]int{{1, 2, 3, 4, 5, 6, 7}, {7, 6, 5, 4, 3, 2, 1}})
	a, b := Build(cal, Options{Height: 8}), Build(cal, Options{Height: 8})
	for y := range a.Rows {
		for x := range a.Rows[y] {
			if a.Rows[y][x] != b.Rows[y][x] {
				t.Fatalf("scene differs at %d,%d", x, y)
			}
		}
	}
}
