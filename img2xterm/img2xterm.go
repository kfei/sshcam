package img2xterm

import (
	"fmt"
	"math"
)

const (
	colorUndef       = iota
	colorTransparent = iota
)

var oldfg, oldbg = colorUndef, colorUndef

func floatMod(x, y float64) float64 {
	return x - y*math.Floor(x/y)
}

func floatMin(x, y float64) float64 {
	if x-y > 0 {
		return y
	}
	return x
}

func rawRGB2Pixels(raw []byte) (ret [][3]byte) {
	for cur := 0; cur < len(raw); cur += 3 {
		pixel := [3]byte{raw[cur], raw[cur+1], raw[cur+2]}
		ret = append(ret, pixel)
	}
	return
}

func rawRGB2BrightnessPixels(raw []byte) (ret []float64) {
	for cur := 0; cur < len(raw); cur += 3 {
		r, g, b := raw[cur], raw[cur+1], raw[cur+2]
		bri := (float64(r)*0.299 + float64(g)*0.587 + float64(b)*0.114) / 255.0
		ret = append(ret, bri)
	}
	return
}

func DrawRGB(raw []byte, width, height int, colorful bool) {
	var color1, color2, brightness1, brightness2 int
	if colorful {
		// Draw image with color
		pixels := rawRGB2Pixels(raw)
		for y := 0; y < height; y += 2 {
			for x := 0; x < width; x++ {
				// Compute the color of upper block
				r1 := int(pixels[y*width+x][0])
				g1 := int(pixels[y*width+x][1])
				b1 := int(pixels[y*width+x][2])
				color1 = rgb2Xterm(r1, g1, b1)

				// Compute the color of lower block
				if (y + 1) < height {
					r2 := int(pixels[(y+1)*width+x][0])
					g2 := int(pixels[(y+1)*width+x][1])
					b2 := int(pixels[(y+1)*width+x][2])
					color2 = rgb2Xterm(r2, g2, b2)
				} else {
					color2 = colorTransparent
				}

				// Draw onn pixel
				bifurcate(color1, color2)
			}
			if (y + 2) < height {
				fmt.Printf("\n")
			}
		}

	} else {
		// Draw image in grayscale
		pixels := rawRGB2BrightnessPixels(raw)
		for y := 0; y < height; y += 2 {
			for x := 0; x < width; x++ {
				brightness1 = int(pixels[y*width+x]*23) + 232
				if (y + 1) < height {
					brightness2 = int(pixels[(y+1)*width+x]*23) + 232
				} else {
					brightness2 = colorTransparent
				}
				bifurcate(brightness1, brightness2)
			}
			if (y + 2) < height {
				fmt.Printf("\n")
			}
		}
	}
}

func bifurcate(color1, color2 int) {
	fg, bg := oldfg, oldbg
	// The lower half block "▄"
	var str = "\xe2\x96\x84"

	if color1 == color2 {
		bg = color1
		str = " "
	} else if color2 == colorTransparent {
		// The upper half block "▀"
		str = "\xe2\x96\x80"
		bg, fg = color2, color1
	} else {
		bg, fg = color1, color2
	}

	if bg != oldbg {
		if bg == colorTransparent {
			fmt.Print("\033[49m")
		} else {
			fmt.Printf("\033[48;5;%dm", bg)
		}
	}

	if fg != oldfg {
		if fg == colorUndef {
			fmt.Print("\033[39m")
		} else {
			fmt.Printf("\033[38;5;%dm", fg)
		}
	}

	oldbg, oldfg = bg, fg

	fmt.Print(str)
}
