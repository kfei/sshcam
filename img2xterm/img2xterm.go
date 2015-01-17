package img2xterm

import (
	"fmt"
	"math"
)

func floatMod(x, y float64) float64 {
	return x - y*math.Floor(x/y)
}

func floatMin(x, y float64) float64 {
	if x-y > 0 {
		return y
	} else {
		return x
	}
}

func rgb2Gray(r, g, b byte) float64 {
	return float64(r)*0.299 + float64(g)*0.587 + float64(b)*0.114
}

func rawRGB2GrayPixels(raw []byte) (ret []float64) {
	for cur := 0; cur < len(raw); cur += 3 {
		r, g, b := raw[cur], raw[cur+1], raw[cur+2]
		brightness := rgb2Gray(r, g, b) / 255.0
		ret = append(ret, brightness)
	}
	return
}

func DrawRGB(raw []byte, width, height int, color bool) {
	var chr string
	if color {
		// Draw image with color
	} else {
		// Draw image in grayscale
		pixels := rawRGB2GrayPixels(raw)
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				brightness := pixels[y*width+x]
				bg := brightness*23 + 232
				fg := floatMin(255, bg+1)
				mod := floatMod(bg, 1.0)

				switch {
				case mod < 0.2:
					chr = " "
				case mod < 0.4:
					chr = "░"
				case mod < 0.6:
					chr = "▒"
				case mod < 0.8:
					bg, fg = fg, bg
					chr = "▒"
				default:
					bg, fg = fg, bg
					chr = "░"
				}

				// TODO: Try to reduce the number of Printf call
				fmt.Printf(
					"\033[48;5;%dm\033[38;5;%dm%s", int(bg), int(fg), chr)
			}
			if y < height-1 {
				fmt.Printf("\n")
			}
		}
	}
}
