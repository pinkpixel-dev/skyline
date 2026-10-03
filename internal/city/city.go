// Package city turns a contribution calendar into a grid of cityscape cells.
// Renderers (ANSI, SVG) only decide how each cell looks, never where it goes.
package city

import (
	"math"

	"github.com/pinkpixel-dev/skyline/internal/github"
)

// Kind says what a cell is so renderers can pick colors.
type Kind int

const (
	Sky Kind = iota
	Star
	Antenna
	Roof
	Window
	Street
	Label
)

// Cell is one character of the scene.
type Cell struct {
	Rune    rune
	Kind    Kind
	Level   int  // window brightness 0-4 (0 = dark)
	Shade   int  // alternating building tone, -1 for anything outside a building
	Flicker bool // lit window that animates in the SVG
}

// Scene is the finished grid plus a little metadata for headers.
type Scene struct {
	Rows  [][]Cell
	Login string
	Total int
}

// Options control the scene size.
type Options struct {
	Height int // tallest building, in rows
	Weeks  int // how many recent weeks to draw
}

const colWidth = 2

// Build lays out one building per week. Height follows the week's total on a
// log scale, and windows cycle through that week's days, lit by GitHub level.
func Build(cal *github.Calendar, opt Options) *Scene {
	if opt.Height < 3 {
		opt.Height = 3
	}
	weeks := cal.Weeks
	if opt.Weeks > 0 && opt.Weeks < len(weeks) {
		weeks = weeks[len(weeks)-opt.Weeks:]
	}

	width := len(weeks) * colWidth
	skyRows := opt.Height + 2 // headroom for roof + antenna
	streetY, labelY := skyRows, skyRows+1
	rows := make([][]Cell, skyRows+2)
	for y := range rows {
		rows[y] = make([]Cell, width)
		for x := range rows[y] {
			rows[y][x] = skyCell(x, y, streetY)
		}
	}

	totals := make([]int, len(weeks))
	maxTotal, record := 0, -1
	for i, w := range weeks {
		for _, d := range w {
			totals[i] += d.Count
		}
		if totals[i] > maxTotal {
			maxTotal, record = totals[i], i
		}
	}

	for i, w := range weeks {
		h := heightFor(totals[i], maxTotal, opt.Height)
		x0 := i * colWidth
		shade := i % 2
		top := streetY - h

		cellIdx := 0
		for y := streetY - 1; y >= top; y-- {
			for dx := 0; dx < colWidth; dx++ {
				level := 0
				if len(w) > 0 {
					level = w[cellIdx%len(w)].Level
				}
				cellIdx++
				c := Cell{Rune: '·', Kind: Window, Shade: shade}
				if level > 0 {
					c.Rune, c.Level = '▪', level
					c.Flicker = hash(i, y, dx)%9 == 0
				}
				rows[y][x0+dx] = c
			}
		}
		for dx := 0; dx < colWidth; dx++ {
			rows[top-1][x0+dx] = Cell{Rune: '▁', Kind: Roof, Shade: -1}
		}
		if i == record {
			rows[top-2][x0] = Cell{Rune: '╻', Kind: Antenna, Shade: -1}
		}
	}

	for x := 0; x < width; x++ {
		rows[streetY][x] = Cell{Rune: '━', Kind: Street, Shade: -1}
	}
	placeLabels(rows[labelY], weeks)

	return &Scene{Rows: rows, Login: cal.Login, Total: cal.Total}
}

func heightFor(total, max, maxH int) int {
	if total == 0 || max == 0 {
		return 1
	}
	scaled := math.Log1p(float64(total)) / math.Log1p(float64(max))
	return 1 + int(math.Round(scaled*float64(maxH-1)))
}

func skyCell(x, y, streetY int) Cell {
	c := Cell{Rune: ' ', Kind: Sky, Shade: -1}
	if y < streetY {
		switch h := hash(x, y, 7); {
		case h%61 == 0:
			c.Rune, c.Kind = '✦', Star
		case h%23 == 0:
			c.Rune, c.Kind = '·', Star
		}
	}
	if y > streetY {
		c.Kind = Label
	}
	return c
}

var months = [...]string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}

// placeLabels writes a month name under the first week of each month,
// skipping any that would collide with the previous label.
func placeLabels(row []Cell, weeks [][]github.Day) {
	next := 0
	prevMonth := -1
	for i, w := range weeks {
		if len(w) == 0 {
			continue
		}
		m := int(w[len(w)-1].Date.Month()) - 1
		if m == prevMonth {
			continue
		}
		prevMonth = m
		x := i * colWidth
		name := months[m]
		if i == 0 || x < next || x+len(name) > len(row) {
			continue
		}
		for j, r := range name {
			row[x+j] = Cell{Rune: r, Kind: Label, Shade: -1}
		}
		next = x + len(name) + 1
	}
}

// hash is a tiny deterministic mixer so stars and flicker stay put between runs.
func hash(a, b, c int) uint32 {
	h := uint32(2166136261)
	for _, v := range [...]int{a, b, c} {
		h ^= uint32(v)
		h *= 16777619
		h ^= h >> 13
	}
	return h
}
