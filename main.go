package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"log"
	"os"
)

func main() {
	exampleImageFile := "./mouse.jpg"
	file, err := os.Open(exampleImageFile)
	if err != nil {
		panic("error opening file")
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("Error decoding image %v", err)
	}
	fmt.Printf("Resizing image of size: %d %d\n", img.Bounds().Max.X, img.Bounds().Max.Y)
	newImg := resizeAreaAverage(img, 100)
	chars := "@@@##***++===---::..  "

	// Loop rows
	bounds := newImg.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		line := ""
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := newImg.At(x, y)
			b := brightness(c)
			idx := int(b) * (len(chars) - 1) / 255
			line += string(chars[idx])
		}
		fmt.Println(line)
	}

}

func brightness(c color.Color) uint8 {
	r, g, b, _ := c.RGBA()
	y := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
	return uint8(y / 257.0) // 65535/255 ≈ 257
}

func resizeAreaAverage(img image.Image, newWidth int) *image.RGBA {
	originalWidth := img.Bounds().Dx()
	originalHeight := img.Bounds().Dy()

	aspectCorrection := 0.5
	newHeight := int(float64(originalHeight) / float64(originalWidth) * float64(newWidth) * aspectCorrection)
	scaleX := originalWidth / newWidth
	scaleY := originalHeight / newHeight

	newImage := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			startX := x * scaleX
			endX := (x + 1) * scaleX
			startY := y * scaleY
			endY := (y + 1) * scaleY

			var rSum, gSum, bSum, aSum uint32
			var count uint32

			for sy := startY; sy < endY; sy++ {
				for sx := startX; sx < endX; sx++ {
					r, g, b, a := img.At(sx, sy).RGBA()
					rSum += r
					gSum += g
					bSum += b
					aSum += a
					count++
				}
			}

			newColor := color.RGBA{
				R: uint8(rSum / count >> 8),
				G: uint8(gSum / count >> 8),
				B: uint8(bSum / count >> 8),
				A: uint8(aSum / count >> 8),
			}
			newImage.Set(x, y, newColor)
		}
	}

	return newImage
}
