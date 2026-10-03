// Command neon-skyline draws a GitHub contribution graph as a neon city.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"charm.land/lipgloss/v2"

	"github.com/pinkpixel-dev/neon-skyline/internal/city"
	"github.com/pinkpixel-dev/neon-skyline/internal/github"
	"github.com/pinkpixel-dev/neon-skyline/internal/render"
)

func main() {
	height := flag.Int("height", 14, "height of the tallest building, in rows")
	weeks := flag.Int("weeks", 53, "number of recent weeks to draw")
	svgPath := flag.String("svg", "", "write an animated SVG to this path instead of printing")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: neon-skyline [flags] <github-username>")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	if err := run(flag.Arg(0), *height, *weeks, *svgPath); err != nil {
		fmt.Fprintln(os.Stderr, "neon-skyline:", err)
		os.Exit(1)
	}
}

func run(login string, height, weeks int, svgPath string) error {
	token, err := github.Token()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cal, err := github.Fetch(ctx, token, login)
	if err != nil {
		return err
	}
	scene := city.Build(cal, city.Options{Height: height, Weeks: weeks})

	if svgPath == "" {
		lipgloss.Print(render.ANSI(scene))
		return nil
	}
	if err := os.WriteFile(svgPath, []byte(render.SVG(scene)), 0o644); err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, "wrote", svgPath)
	return nil
}
