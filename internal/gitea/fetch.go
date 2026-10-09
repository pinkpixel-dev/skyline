// Package gitea fetches a user's contribution heatmap from a Gitea (or
// Forgejo) instance and shapes it like a GitHub calendar.
package gitea

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/pinkpixel-dev/skyline/internal/github"
)

// entry is one heatmap bucket. Gitea groups activity into short time slots,
// so a single day can show up several times.
type entry struct {
	Timestamp     int64 `json:"timestamp"`
	Contributions int   `json:"contributions"`
}

// Token returns GITEA_TOKEN, or "" since public heatmaps need no auth.
func Token() string {
	return strings.TrimSpace(os.Getenv("GITEA_TOKEN"))
}

// Fetch downloads the last year of contributions for login from baseURL.
func Fetch(ctx context.Context, baseURL, token, login string) (*github.Calendar, error) {
	endpoint := strings.TrimRight(baseURL, "/") + "/api/v1/users/" + url.PathEscape(login) + "/heatmap"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	switch {
	case res.StatusCode == http.StatusNotFound:
		return nil, fmt.Errorf("user %q not found, or the heatmap is disabled on this instance", login)
	case res.StatusCode != http.StatusOK:
		return nil, fmt.Errorf("gitea api: %s", res.Status)
	}

	var entries []entry
	if err := json.NewDecoder(res.Body).Decode(&entries); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return buildCalendar(login, entries, time.Now()), nil
}

// buildCalendar sums entries into UTC days and lays them out like GitHub's
// graph: Sunday-first weeks running from a year back up to today.
func buildCalendar(login string, entries []entry, now time.Time) *github.Calendar {
	today := now.UTC().Truncate(24 * time.Hour)
	start := today.AddDate(-1, 0, 0)
	start = start.AddDate(0, 0, -int(start.Weekday()))

	counts := map[time.Time]int{}
	for _, e := range entries {
		day := time.Unix(e.Timestamp, 0).UTC().Truncate(24 * time.Hour)
		if !day.Before(start) && !day.After(today) {
			counts[day] += e.Contributions
		}
	}
	q := quartiles(counts)

	cal := &github.Calendar{Login: login}
	var week []github.Day
	for d := start; !d.After(today); d = d.AddDate(0, 0, 1) {
		n := counts[d]
		cal.Total += n
		week = append(week, github.Day{Date: d, Count: n, Level: level(n, q)})
		if d.Weekday() == time.Saturday {
			cal.Weeks = append(cal.Weeks, week)
			week = nil
		}
	}
	if len(week) > 0 {
		cal.Weeks = append(cal.Weeks, week)
	}
	return cal
}

// quartiles returns the 25th, 50th and 75th percentile of the active days,
// standing in for the levels GitHub computes server-side.
func quartiles(counts map[time.Time]int) [3]int {
	var active []int
	for _, n := range counts {
		if n > 0 {
			active = append(active, n)
		}
	}
	if len(active) == 0 {
		return [3]int{}
	}
	slices.Sort(active)
	at := func(p int) int { return active[(len(active)-1)*p/100] }
	return [3]int{at(25), at(50), at(75)}
}

// level maps a day's count to 0-4 against the quartile cutoffs.
func level(n int, q [3]int) int {
	if n <= 0 {
		return 0
	}
	l := 1
	for _, cut := range q {
		if n > cut {
			l++
		}
	}
	return l
}
