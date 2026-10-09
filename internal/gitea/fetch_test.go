package gitea

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// Thursday, so the last week is partial.
var now = time.Date(2026, 10, 8, 15, 30, 0, 0, time.UTC)

func at(date string, hour int) int64 {
	d, _ := time.Parse("2006-01-02", date)
	return d.Add(time.Duration(hour) * time.Hour).Unix()
}

func TestBuildCalendarShape(t *testing.T) {
	cal := buildCalendar("octo", nil, now)

	first := cal.Weeks[0][0].Date
	if first.Weekday() != time.Sunday || !first.Before(now.AddDate(-1, 0, 0)) {
		t.Fatalf("first day = %s, want the Sunday on or before a year ago", first)
	}
	for i, w := range cal.Weeks[:len(cal.Weeks)-1] {
		if len(w) != 7 {
			t.Fatalf("week %d has %d days, want 7", i, len(w))
		}
	}
	last := cal.Weeks[len(cal.Weeks)-1]
	if got := last[len(last)-1].Date.Format("2006-01-02"); got != "2026-10-08" {
		t.Fatalf("last day = %s, want today", got)
	}
	if cal.Total != 0 {
		t.Fatalf("empty heatmap total = %d", cal.Total)
	}
}

func TestBuildCalendarSumsBucketsPerDay(t *testing.T) {
	entries := []entry{
		{at("2026-10-01", 9), 2},
		{at("2026-10-01", 18), 3},
		{at("2026-10-02", 12), 1},
		{at("2024-01-01", 12), 50}, // older than a year, dropped
		{at("2026-10-09", 1), 50},  // tomorrow, dropped
	}
	cal := buildCalendar("octo", entries, now)

	if cal.Total != 6 {
		t.Fatalf("total = %d, want 6", cal.Total)
	}
	byDate := map[string]int{}
	for _, w := range cal.Weeks {
		for _, d := range w {
			byDate[d.Date.Format("2006-01-02")] = d.Count
		}
	}
	if byDate["2026-10-01"] != 5 || byDate["2026-10-02"] != 1 {
		t.Fatalf("day counts = %d, %d, want 5, 1", byDate["2026-10-01"], byDate["2026-10-02"])
	}
}

func TestLevels(t *testing.T) {
	var entries []entry
	for i, n := range []int{1, 2, 3, 4, 5, 6, 7, 8} {
		entries = append(entries, entry{at("2026-09-01", 0) + int64(i)*86400, n})
	}
	cal := buildCalendar("octo", entries, now)

	got := map[int]int{}
	for _, w := range cal.Weeks {
		for _, d := range w {
			if d.Count > 0 {
				got[d.Count] = d.Level
			}
		}
	}
	want := map[int]int{1: 1, 2: 1, 3: 2, 4: 2, 5: 3, 6: 3, 7: 4, 8: 4}
	for n, l := range want {
		if got[n] != l {
			t.Errorf("count %d got level %d, want %d", n, got[n], l)
		}
	}
	if level(0, [3]int{1, 2, 3}) != 0 {
		t.Error("zero count should be level 0")
	}
}

func TestFetch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/api/v1/users/octo/heatmap":
			if got := r.Header.Get("Authorization"); got != "token secret" {
				t.Errorf("auth header = %q", got)
			}
			fmt.Fprintf(w, `[{"timestamp":%d,"contributions":4}]`, time.Now().Add(-time.Hour).Unix())
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	cal, err := Fetch(context.Background(), srv.URL+"/", "secret", "octo")
	if err != nil {
		t.Fatal(err)
	}
	if cal.Login != "octo" || cal.Total != 4 {
		t.Fatalf("login, total = %q, %d, want octo, 4", cal.Login, cal.Total)
	}

	_, err = Fetch(context.Background(), srv.URL, "", "nobody")
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Fatalf("missing user error = %v", err)
	}
}
