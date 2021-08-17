package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

var expected string = `<svg xmlns="http://www.w3.org/2000/svg" width="126" height="20">
    <title>90</title>
    <desc>Generated with covbadger (https://github.com/imsky/covbadger)</desc>
    <linearGradient id="smooth" x2="0" y2="100%">
        <stop offset="0" stop-color="#bbb" stop-opacity=".1" />
        <stop offset="1" stop-opacity=".1" />
    </linearGradient>
    <rect rx="3" width="126" height="20" fill="#555" />
    <rect rx="3" x="90" width="36" height="20" fill="#97ca00" />
    <rect x="90" width="4" height="20" fill="#97ca00" />
    <rect rx="3" width="126" height="20" fill="url(#smooth)" />
    <g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,sans-serif" font-size="11">
        <text x="45" y="15" fill="#010101" fill-opacity=".3">coverage:test</text>
        <text x="45" y="14">coverage:test</text>
        <text x="108" y="15" fill="#010101" fill-opacity=".3">90%</text>
        <text x="108" y="14">90%</text>
    </g>
</svg>`

func TestRenderBadge(t *testing.T) {
	var err error
	badge, _ := RenderBadge(90, "coverage:test")

	if badge != expected {
		t.Errorf("RenderBadge output is incorrect")
	}

	badge, _ = RenderBadge(100, "coverage:test")

	if strings.Contains(badge, colors["brightgreen"]) != true {
		t.Errorf("Incorrect color for coverage badge, expected brightgreen")
	}

	badge, _ = RenderBadge(70, "coverage:test")

	if strings.Contains(badge, colors["yellow"]) != true {
		t.Errorf("Incorrect color for coverage badge, expected yellow")
	}

	badge, _ = RenderBadge(60, "coverage:test")

	if strings.Contains(badge, colors["orange"]) != true {
		t.Errorf("Incorrect color for coverage badge, expected orange")
	}

	badge, _ = RenderBadge(40, "coverage:test")

	if strings.Contains(badge, colors["red"]) != true {
		t.Errorf("Incorrect color for coverage badge, expected red")
		t.Errorf(badge)
	}

	badge, err = RenderBadge(101, "coverage:test")

	if err == nil {
		t.Errorf("Invalid coverage: greater than 100%%")
	}

	badge, err = RenderBadge(-1, "coverage:test")

	if err == nil {
		t.Errorf("Invalid coverage: less than 0%%")
	}
}

func TestCovbadger(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Invalid coverage value did not throw an error")
		}
	}()

	main()
	Run([]string{"99", "coverage:test"})

	userInput := []byte("90.09")
	tmpfile, _ := ioutil.TempFile("", "covbadger-stdin")
	defer os.Remove(tmpfile.Name())
	tmpfile.Write(userInput)
	tmpfile.Seek(0, 0)

	_stdin := os.Stdin
	defer func() { os.Stdin = _stdin }()
	os.Stdin = tmpfile
	Run([]string{"-", "coverage:test"})
	tmpfile.Close()

	Run([]string{"-99", "coverage:test"})
}
