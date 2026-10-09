// Command skyline draws a GitHub or Gitea contribution graph as an ASCII city.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"charm.land/lipgloss/v2"

	"github.com/pinkpixel-dev/skyline/internal/city"
	"github.com/pinkpixel-dev/skyline/internal/gitea"
	"github.com/pinkpixel-dev/skyline/internal/github"
	"github.com/pinkpixel-dev/skyline/internal/render"
)

func main() {
	height := flag.Int("height", 14, "height of the tallest building, in rows")
	weeks := flag.Int("weeks", 53, "number of recent weeks to draw")
	svgPath := flag.String("svg", "", "write an animated SVG to this path instead of printing")
	themeName := flag.String("theme", render.Themes[0].Name, "color theme (see -themes)")
	listThemes := flag.Bool("themes", false, "list the built-in themes and exit")
	giteaURL := flag.String("gitea", "", "read contributions from the Gitea or Forgejo instance at this base URL")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: skyline [flags] <username>")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *listThemes {
		lipgloss.Print(themeList())
		return
	}
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	if err := run(flag.Arg(0), *giteaURL, *themeName, *height, *weeks, *svgPath); err != nil {
		fmt.Fprintln(os.Stderr, "skyline:", err)
		os.Exit(1)
	}
}

func run(login, giteaURL, themeName string, height, weeks int, svgPath string) error {
	theme, err := render.ThemeByName(themeName)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cal, err := fetch(ctx, login, giteaURL)
	if err != nil {
		return err
	}
	scene := city.Build(cal, city.Options{Height: height, Weeks: weeks})

	if svgPath == "" {
		lipgloss.Print(render.ANSI(scene, theme))
		return nil
	}
	if err := os.WriteFile(svgPath, []byte(render.SVG(scene, theme)), 0o644); err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, "wrote", svgPath)
	return nil
}

// fetch reads the calendar from Gitea when a base URL is given, otherwise GitHub.
func fetch(ctx context.Context, login, giteaURL string) (*github.Calendar, error) {
	if giteaURL != "" {
		return gitea.Fetch(ctx, giteaURL, gitea.Token(), login)
	}
	token, err := github.Token()
	if err != nil {
		return nil, err
	}
	return github.Fetch(ctx, token, login)
}

// themeList prints each theme name next to a swatch of its roof and window colors.
func themeList() string {
	var b strings.Builder
	for _, t := range render.Themes {
		name := lipgloss.NewStyle().Width(11).Foreground(lipgloss.Color(t.Text)).Render(t.Name)
		swatch := lipgloss.NewStyle().Background(lipgloss.Color(t.Sky)).Foreground(lipgloss.Color(t.Roof)).Render(" ▁▁ ")
		for _, c := range t.Windows {
			swatch += lipgloss.NewStyle().Background(lipgloss.Color(t.Shades[0])).Foreground(lipgloss.Color(c)).Render("▪▪")
		}
		b.WriteString(name + swatch + "\n")
	}
	return b.String()
}
