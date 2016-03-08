package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var imageNamePattern *regexp.Regexp
var defaultColor = color.RGBA{0, 0xFF, 0, 0xCC}

type imageConfig struct {
	width, height int
	background    color.Color
}

func init() {
	imageNamePattern = regexp.MustCompile(`(?i)(\d+)x(\d+)(\-[0-9a-f]{6})?.jpg`)
}

func main() {
	port := flag.Int("port", 8000, "start server on this port")
	flag.Parse()

	url := fmt.Sprintf("localhost:%d", *port)

	log.Printf("Starting server on port %d\n", *port)

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(url, nil))
}

func parseImageConfig(name string) (*imageConfig, error) {
	match := imageNamePattern.FindStringSubmatch(name)
	if len(match) != 4 {
		msg := fmt.Sprintf("Fail to parse name: %v", name)
		return nil, errors.New(msg)
	}

	var width, height int
	widthStr, heightStr, colorStr := match[1], match[2], match[3]
	width, err := strconv.Atoi(widthStr)
	if err == nil {
		height, err = strconv.Atoi(heightStr)
	}

	if err != nil {
		return nil, err
	}

	var background color.Color
	if len(colorStr) > 0 {
		colorStr = colorStr[1:]
		background = parseColor(colorStr)
	} else {
		background = defaultColor
	}

	return &imageConfig{width, height, background}, err
}

func parseColor(s string) color.Color {
	var (
		err     error
		r, g, b uint64
	)
	r, err = strconv.ParseUint(s[0:2], 16, 8)
	if err == nil {
		g, err = strconv.ParseUint(s[2:4], 16, 8)
	}

	if err == nil {
		b, err = strconv.ParseUint(s[4:6], 16, 8)
	}

	if err != nil {
		log.Printf("Fail to parse color: %v, using default color\n%s", s, err)
		return defaultColor
	}
	return color.RGBA{uint8(r), uint8(g), uint8(b), 0xFF}
}

func handler(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimLeft(r.URL.String(), "/")

	config, err := parseImageConfig(name)
	if err != nil {
		log.Println(err)
		http.NotFound(w, r)
		return
	}

	img := createImage(config.width, config.height, config.background)
	err = writeJPEG(w, img)
	if err != nil {
		log.Panic(err)
	}
}

func writeJPEG(w io.Writer, img image.Image) error {
	var opt jpeg.Options
	opt.Quality = 80
	return jpeg.Encode(w, img, &opt)
}

func createImage(width int, height int, background color.Color) *image.RGBA {
	rect := image.Rect(0, 0, width, height)
	img := image.NewRGBA(rect)
	draw.Draw(img, img.Bounds(), &image.Uniform{background}, image.ZP, draw.Src)
	return img
}
