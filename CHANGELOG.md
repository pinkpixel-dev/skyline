# Changelog

## 0.1.0 - October 3, 2026

### 🌃 Skyline

- First release. Draws your GitHub contribution graph as an ASCII city in the terminal
- Building height follows weekly totals on a log scale, and window colors follow GitHub's four contribution levels
- 10 color themes via `-theme`: neon, synthwave, matrix, amber, ice, sunset, toxic, vapor, crimson and mono
- `-themes` lists every theme with a color swatch
- The record week gets an antenna, and month labels run along the street
- `-height` and `-weeks` flags control the size of the city

### 🖼️ SVG

- `-svg` writes an animated SVG with flickering windows and twinkling stars, and the animation respects reduced motion
- Example GitHub Actions workflow that redraws the skyline nightly for a profile README
