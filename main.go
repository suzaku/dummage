package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Please specify the path to save the generated image")
	}
	imgPath := os.Args[1]

	out, err := os.Create(imgPath)
	if err != nil {
		log.Fatal(err)
	}

	width, height := 100, 100
	background := color.RGBA{0, 0xFF, 0, 0xCC}

	img := createImage(width, height, background)

	if strings.HasSuffix(strings.ToLower(imgPath), ".jpg") {
		var opt jpeg.Options
		opt.Quality = 80
		err = jpeg.Encode(out, img, &opt)
	} else {
		err = png.Encode(out, img)
	}

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Image saved to %s\n", imgPath)
}

func createImage(width int, height int, background color.RGBA) *image.RGBA {
	rect := image.Rect(0, 0, width, height)
	img := image.NewRGBA(rect)
	draw.Draw(img, img.Bounds(), &image.Uniform{background}, image.ZP, draw.Src)
	return img
}
