package main

import (
	"math"

	"github.com/fogleman/gg"
)

// Greyscale turns the image grey using luminosity method
// 0.21 R + 0.72 G + 0.07 B
func Greyscale(image *gg.Context, percent float64) {

	percentGrey := math.Max(math.Min(percent, 1), 0)
	percentColor := 1.0 - percentGrey

	for x := 0; x < image.Width(); x++ {
		for y := 0; y < image.Height(); y++ {
			r, g, b, _ := image.Image().At(x, y).RGBA()
			pix := (((float64(r) / 65536.0) * 0.21) + ((float64(g) / 65536.0) * 0.72) + ((float64(b) / 65536.0) * 0.07)) * percentGrey
			image.SetRGB(((float64(r)/65536.0)*percentColor)+pix, ((float64(g)/65536.0)*percentColor)+pix, ((float64(b)/65536.0)*percentColor)+pix)
			image.SetPixel(x, y)
		}
	}
}

// Shift moves an images pixels from current color to another
func Shift(image *gg.Context, percent, rTowards, gTowards, bTowards float64) {

	percentGrey := math.Max(math.Min(percent, 1), 0)
	percentColor := 1.0 - percentGrey

	for x := 0; x < image.Width(); x++ {
		for y := 0; y < image.Height(); y++ {
			r, g, b, _ := image.Image().At(x, y).RGBA()
			image.SetRGB(
				((float64(r)/65536.0)*percentColor)+(rTowards*percentGrey),
				((float64(g)/65536.0)*percentColor)+(gTowards*percentGrey),
				((float64(b)/65536.0)*percentColor)+(bTowards*percentGrey))
			image.SetPixel(x, y)
		}
	}
}
