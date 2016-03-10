package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var imageNamePattern *regexp.Regexp

type imageConfig struct {
	width, height int
	background    color.Color
	format        string
}

func init() {
	imageNamePattern = regexp.MustCompile(`(?i)(\d+)x(\d+)(\-[0-9a-f]{6})?.(jpg|png)`)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	host := flag.String("host", "localhost", "listen on this host")
	port := flag.Int("port", 8000, "start server on this port")
	flag.Parse()

	url := fmt.Sprintf("%s:%d", *host, *port)

	log.Printf("Starting server on %s\n", url)

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(url, nil))
}

func randomColor() color.Color {
	const upperBound = 2 << 8
	r := uint8(rand.Intn(upperBound))
	g := uint8(rand.Intn(upperBound))
	b := uint8(rand.Intn(upperBound))
	return color.RGBA{r, g, b, 0xFF}
}

func parseImageConfig(name string) (*imageConfig, error) {
	match := imageNamePattern.FindStringSubmatch(name)
	if len(match) != 5 {
		return nil, fmt.Errorf("Fail to parse name: %v", name)
	}

	var width, height int
	widthStr, heightStr, colorStr, suffix := match[1], match[2], match[3], match[4]
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
		background = randomColor()
	}

	return &imageConfig{width, height, background, strings.ToUpper(suffix)}, err
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
		return randomColor()
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
	if config.format == "JPG" {
		err = writeJPEG(w, img)
	} else {
		err = writePNG(w, img)
	}
	if err != nil {
		log.Panic(err)
	}
}

func writeJPEG(w io.Writer, img image.Image) error {
	var opt jpeg.Options
	opt.Quality = 80
	return jpeg.Encode(w, img, &opt)
}

func writePNG(w io.Writer, img image.Image) error {
	return png.Encode(w, img)
}

func createImage(width, height int, background color.Color) *image.RGBA {
	rect := image.Rect(0, 0, width, height)
	img := image.NewRGBA(rect)
	draw.Draw(img, img.Bounds(), &image.Uniform{background}, image.ZP, draw.Src)
	return img
}
