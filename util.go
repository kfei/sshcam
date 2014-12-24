package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	webcam "github.com/kfei/sshcam/webcam/v4l2"
)

func wxh2Size(s string) Size {
	// TODO: Also reads in "w*h" and "w h" format
	splits := strings.Split(s, "x")
	w, _ := strconv.Atoi(splits[0])
	h, _ := strconv.Atoi(splits[1])
	return Size{w, h}
}

func updateTTYSize() <-chan string {
	ttyStatus := make(chan string)
	go func() {
		for {
			cmd := exec.Command("stty", "size")
			cmd.Stdin = os.Stdin
			out, err := cmd.Output()
			if err != nil {
				log.Fatal(err)
			}
			ttyStatus <- strings.TrimSuffix(string(out), "\n")
			time.Sleep(100)
		}
	}()
	return ttyStatus
}

func rgb2Gray(r, g, b byte) float64 {
	gray := (float64(r)*0.299 + float64(g)*0.587 + float64(b)*0.114)
	return gray
}

func scaleRGBArrayToGrayPixels(from []byte, tSize Size) []float64 {
	// Check the image size actually captured by webcam
	if size.Width*size.Height*3 > len(from) {
		log.Fatal("Pixels conversion failed. Did you specified a size " +
			"which is not supported by the webcam?")
	}
	// TODO: Improve this inefficient and loosy algorithm
	var to []float64
	skipX := size.Width / tSize.Width
	skipY := size.Height / tSize.Height
	for y := 0; y < tSize.Height; y++ {
		for x := 0; x < tSize.Width; x++ {
			cur := size.Width*3*y*skipY + 3*x*skipX
			r, g, b := from[cur], from[cur+1], from[cur+2]
			brightness := rgb2Gray(r, g, b) / 255.0
			to = append(to, brightness)
		}
	}
	return to
}

func fetchGrayPixels(s Size) []float64 {
	rgbArray := webcam.GrabFrame()
	return scaleRGBArrayToGrayPixels(rgbArray, s)
}

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

func draw(ttyStatus <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Println("Start streaming, press Ctrl-c to exit...")
	time.Sleep(3 * time.Second)

	var chr string
	tty := Size{0, 0}

	for {
		// Update TTY size before every draw (synchronous)
		curSize := strings.Split(<-ttyStatus, " ")
		h, _ := strconv.Atoi(curSize[0])
		w, _ := strconv.Atoi(curSize[1])
		tty.Width, tty.Height = w, h

		// Fetch image from webcam
		pixels := fetchGrayPixels(tty)

		// Move cursor to top right and start to draw image
		fmt.Printf("\033[00H")
		for y := 0; y < tty.Height; y++ {
			for x := 0; x < tty.Width; x++ {
				brightness := pixels[y*tty.Width+x]
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
			if y < tty.Height-1 {
				fmt.Printf("\n")
			}
		}
	}
}
