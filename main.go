package main

import (
	"errors"
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

func init() {
	imageNamePattern = regexp.MustCompile(`(\d+)x(\d+).jpg`)
}

func parseDimension(name string) (int, int, error) {
	match := imageNamePattern.FindStringSubmatch(name)
	if len(match) == 0 {
		msg := fmt.Sprintf("Fail to parse name: %v", name)
		return 0, 0, errors.New(msg)
	}

	var width, height int
	widthStr, heightStr := match[1], match[2]
	width, err := strconv.Atoi(widthStr)
	if err != nil {
		return width, height, err
	}
	height, err = strconv.Atoi(heightStr)
	return width, height, err
}

func handler(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimLeft(r.URL.String(), "/")

	width, height, err := parseDimension(name)
	if err != nil {
		log.Println(err)
		http.NotFound(w, r)
		return
	}

	background := color.RGBA{0, 0xFF, 0, 0xCC}

	img := createImage(width, height, background)
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

func main() {
	port := 8000
	url := fmt.Sprintf("localhost:%d", port)

	log.Printf("Starting server on port %d\n", port)

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(url, nil))
}

func createImage(width int, height int, background color.RGBA) *image.RGBA {
	rect := image.Rect(0, 0, width, height)
	img := image.NewRGBA(rect)
	draw.Draw(img, img.Bounds(), &image.Uniform{background}, image.ZP, draw.Src)
	return img
}
