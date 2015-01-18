package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kfei/sshcam/img2xterm"
	webcam "github.com/kfei/sshcam/webcam/v4l2"
)

func wxh2Size(s string) Size {
	// TODO: Also reads in "w*h" and "w h" format
	splits := strings.Split(s, "x")
	w, _ := strconv.Atoi(splits[0])
	h, _ := strconv.Atoi(splits[1])
	return Size{w, h}
}

func clearScreen() {
	// TODO: Use terminfo
	fmt.Printf("\033[2J")
}

func resetCursor() {
	// TODO: Use terminfo
	fmt.Printf("\033[00H")
}

func updateTTYSize() <-chan string {
	ttyStatus := make(chan string)
	go func() {
		for {
			// TODO: Use syscall.Syscall?
			cmd := exec.Command("stty", "size")
			cmd.Stdin = os.Stdin
			out, err := cmd.Output()
			if err != nil {
				log.Fatal(err)
			}
			ttyStatus <- strings.TrimSuffix(string(out), "\n")
			time.Sleep(300 * time.Millisecond)
		}
	}()
	return ttyStatus
}

func grabRGBPixels(ttySize Size) (ret []byte) {
	rgbArray := webcam.GrabFrame()
	// Check the image size actually captured by webcam
	if size.Width*size.Height*3 > len(rgbArray) {
		log.Fatal("Pixels conversion failed. Did you specified a size " +
			"which is not supported by the webcam?")
	}

	// Assuming the captured image is larger than terminal size
	// TODO: Scale up the image when termial size is bigger
	// TODO: Improve this inefficient and loosy algorithm
	skipX, skipY := size.Width/ttySize.Width, size.Height/(ttySize.Height*2)
	for y := 0; y < ttySize.Height*2; y++ {
		for x := 0; x < ttySize.Width; x++ {
			cur := size.Width*3*y*skipY + 3*x*skipX
			ret = append(ret, rgbArray[cur], rgbArray[cur+1], rgbArray[cur+2])
		}
	}
	return
}

func draw(ttyStatus <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Println("Start streaming, press Ctrl-c to exit...")
	time.Sleep(1500 * time.Millisecond)
	clearScreen()

	ttySize := Size{0, 0}

	for {
		// Update TTY size before every draw (synchronous)
		curSize := strings.Split(<-ttyStatus, " ")
		ttySize.Height, _ = strconv.Atoi(curSize[0])
		ttySize.Width, _ = strconv.Atoi(curSize[1])

		// Fetch image from webcam and call img2xterm to draw
		rgbRaw := grabRGBPixels(ttySize)
		resetCursor()
		img2xterm.DrawRGB(rgbRaw, ttySize.Width, ttySize.Height*2, false)
	}
}
