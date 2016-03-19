package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"strings"
	"testing"
	"time"
)

type imgConfig struct {
	name   string
	w, h   int
	bg     color.Color
	format string
}

var r = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

func randImgConfig(r *rand.Rand, withColor bool) *imgConfig {
	var c imgConfig
	c.w = r.Intn(1000) + 1
	c.h = r.Intn(1000) + 1
	c.bg = randomColor()
	if n := r.Float32(); n < 0.5 {
		c.format = "PNG"
	} else {
		if n > 0.25 {
			c.format = "JPEG"
		} else {
			c.format = "JPG"
		}
	}
	if withColor {
		r, b, g, _ := c.bg.RGBA()
		c.name = fmt.Sprintf(
			"%dx%d-%0x%0x%0x.%s",
			c.w, c.h, uint8(r), uint8(g), uint8(b), strings.ToLower(c.format),
		)
	} else {
		c.name = fmt.Sprintf(
			"%dx%d.%s",
			c.w, c.h, strings.ToLower(c.format),
		)
	}
	return &c
}

func TestParseImageSizeOnly(t *testing.T) {
	for i := 0; i < 10; i++ {
		s := randImgConfig(r, false)
		width, height, _, _, _ := parseImageConfig(s.name)
		if width != s.w || height != s.h {
			t.Errorf("parseImageConfig(%q): (%q, %q)", s.name, width, height)
		}
	}
}

func TestParseImageFormat(t *testing.T) {
	for i := 0; i < 10; i++ {
		s := randImgConfig(r, false)
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
