package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"text/template"
)

//Badge represents a coverage badge
type Badge struct {
	Coverage int
	Color    string
	Title    string
	TotalWidth int
	TitleWidth int
	XStartTitle int
	XStartCoverage int
}

var colors = map[string]string{
	"brightgreen": "#44cc11",
	"green":       "#97ca00",
	"yellow":      "#dfb317",
	"orange":      "#fe7d37",
	"red":         "#e05d44",
}

var _badgeTemplate string = `<svg xmlns="http://www.w3.org/2000/svg" width="{{.TotalWidth}}" height="20">
    <title>{{.Coverage}}</title>
    <desc>Generated with covbadger (https://github.com/imsky/covbadger)</desc>
    <linearGradient id="smooth" x2="0" y2="100%">
        <stop offset="0" stop-color="#bbb" stop-opacity=".1" />
        <stop offset="1" stop-opacity=".1" />
    </linearGradient>
    <rect rx="3" width="{{.TotalWidth}}" height="20" fill="#555" />
    <rect rx="3" x="{{.TitleWidth}}" width="36" height="20" fill="{{.Color}}" />
    <rect x="{{.TitleWidth}}" width="4" height="20" fill="{{.Color}}" />
    <rect rx="3" width="{{.TotalWidth}}" height="20" fill="url(#smooth)" />
    <g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,sans-serif" font-size="11">
        <text x="{{.XStartTitle}}" y="15" fill="#010101" fill-opacity=".3">{{.Title}}</text>
        <text x="{{.XStartTitle}}" y="14">{{.Title}}</text>
        <text x="{{.XStartCoverage}}" y="15" fill="#010101" fill-opacity=".3">{{.Coverage}}%</text>
        <text x="{{.XStartCoverage}}" y="14">{{.Coverage}}%</text>
    </g>
</svg>`

func RenderBadge(coverage int, title string) (string, error) {
	if coverage < 0 || coverage > 100 {
		return "", errors.New("Invalid coverage: " + strconv.Itoa(coverage))
	}

	var buffer bytes.Buffer
	badgeTemplate, _ := template.New("badge").Parse(_badgeTemplate)

	color := colors["red"]

	if coverage > 95 {
		color = colors["brightgreen"]
	} else if coverage > 80 {
		color = colors["green"]
	} else if coverage > 60 {
		color = colors["yellow"]
	} else if coverage > 40 {
		color = colors["orange"]
	}

	titleWidth := 7 * len(title)
	if titleWidth % 2 == 1 {
		titleWidth -= 1
	}
	totalWidth := titleWidth + 36
	xStartTitle := titleWidth / 2
	xStartCoverage := titleWidth + 18

	_ = badgeTemplate.Execute(&buffer, &Badge{coverage, color, title, totalWidth, titleWidth, xStartTitle, xStartCoverage})
	return buffer.String(), nil
}

func Run(args []string) {
	if len(args) != 2 {
		flag.Usage()
		return
	}

	coverage := args[0]
	title := args[1]

	if coverage == "-" {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		coverage = scanner.Text()
	}

	coverageValue, _ := strconv.ParseFloat(coverage, 64)
	badge, err := RenderBadge(int(coverageValue), string(title))

	if err != nil {
		panic(err)
	} else {
		fmt.Println(badge)
	}
}

func main() {
	flag.Usage = func() {
		fmt.Println(`Usage: covbadger [coverage] [title]`)
	}

	flag.Parse()
	Run(flag.Args())
}
