package main

import (
	"image/color"
	"testing"
)

func TestParseImageSizeOnly(t *testing.T) {
	tests := []struct {
		name          string
		width, height int
	}{
		{"200x300.jpg", 200, 300},
		{"100x100.png", 100, 100},
		{"500x450.jpg", 500, 450},
	}
	for _, s := range tests {
		width, height, _, _, _ := parseImageConfig(s.name)
		if width != s.width || height != s.height {
			t.Errorf("parseImageConfig(%q): (%q, %q)", s.name, width, height)
		}
	}
}

func TestParseImageFormat(t *testing.T) {
	tests := []struct {
		name   string
		format string
	}{
		{"200x300.jpg", "JPG"},
		{"100x100.png", "PNG"},
		{"500x450.jpeg", "JPEG"},
	}
	for _, s := range tests {
		_, _, _, format, _ := parseImageConfig(s.name)
		if format != s.format {
			t.Errorf("parseImageConfig(%q): %q", s.name, format)
		}
	}
}

func TestParseInvalidImageName(t *testing.T) {
	tests := []string{
		"randomname",
		"/",
		"-769x100.png",
		"abcx200.gif",
	}
	for _, s := range tests {
		w, h, bg, fm, err := parseImageConfig(s)
		if err == nil {
			t.Errorf("parseImageConfig(%q): (%d, %d, %v, %s)", s, w, h, bg, fm)
		}
	}
}

func TestParseColor(t *testing.T) {
	tests := []struct {
		name    string
		r, g, b uint8
	}{
		{"010101", 1, 1, 1},
		{"FF0000", 255, 0, 0},
		{"0f1020", 15, 16, 32},
	}
	for _, s := range tests {
		bg := parseColor(s.name).(color.RGBA)
		if bg.R != s.r || bg.G != s.g || bg.B != s.b {
			t.Errorf("parseColor(%q): (%d, %d, %d)", s.name, bg.R, bg.G, bg.B)
		}
	}
}
