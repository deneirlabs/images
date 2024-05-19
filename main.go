package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"

	"github.com/nfnt/resize"
)

func main() {
	originalFile, err := os.Open("npc/noke.png")
	if err != nil {
		panic(err)
	}
	defer originalFile.Close()

	originalImage, _, err := image.Decode(originalFile)
	if err != nil {
		panic(err)
	}

	// Open the ring mask image file
	maskFile, err := os.Open("mask.png")
	if err != nil {
		panic(err)
	}
	defer maskFile.Close()

	// Decode the ring mask image
	maskImage, _, err := image.Decode(maskFile)
	if err != nil {
		panic(err)
	}

	bounds := originalImage.Bounds()
	destImage := image.NewRGBA(bounds)

	draw.Draw(destImage, bounds, originalImage, bounds.Min, draw.Src)

	centerX := float64(bounds.Dx()) / 2.0
	centerY := float64(bounds.Dy()) / 2.0

	radius := math.Min(centerX, centerY) - (math.Min(centerX, centerY))*0.1

	// Iterate over all pixels and make pixels outside of the circle transparent
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dx := float64(x) - centerX
			dy := float64(y) - centerY
			distance := math.Sqrt(dx*dx + dy*dy)
			if distance > radius {
				destImage.Set(x, y, color.Transparent)
			}
		}
	}

	maskImage = resize.Resize(uint(originalImage.Bounds().Dx()), uint(originalImage.Bounds().Dy()), maskImage, resize.NearestNeighbor)

	// Overlay the mask image onto the original image
	draw.DrawMask(destImage, bounds, maskImage, image.Point{}, maskImage, image.Point{}, draw.Over)

	outfile, err := os.Create("output.png")
	if err != nil {
		panic(err)
	}

	defer outfile.Close()

	err = png.Encode(outfile, destImage)
	if err != nil {
		panic(err)
	}

}

// Function to resize an image
func resizeImage(img image.Image, width, height int) image.Image {
	newImg := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(newImg, newImg.Bounds(), img, image.Point{}, draw.Src)

	return newImg
}
