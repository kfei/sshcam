package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kfei/sshcam/img2xterm"
	webcam "github.com/kfei/sshcam/webcam/v4l2"
)

func wxh2Size(s string) Size {
	splits := [...]string{"x", "*", " "}
	for i := range splits {
		splitted := strings.Split(s, splits[i])
		if len(splitted) != 2 {
			continue
		}
		w, err1 := strconv.Atoi(splitted[0])
		h, err2 := strconv.Atoi(splitted[1])
		if err1 == nil && err2 == nil {
			return Size{w, h}
		}
	}
	log.Println("Invalid argument: --size, fallback to default...")
	return wxh2Size("640x480")
}

func clearScreen() {
	// TODO: Use terminfo
	fmt.Printf("\033[2J")
}

func restoreScreen() {
	fmt.Print("\033[39m\033[49m")
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

func grabRGBPixels(ttySize Size, wInc, hInc int) (ret []byte) {
	rgbArray := webcam.GrabFrame()
	// Check the image size actually captured by webcam
	if size.Width*size.Height*3 > len(rgbArray) {
		log.Fatal("Pixels conversion failed. Did you specified a size " +
			"which is not supported by the webcam?")
	}

	// Assuming the captured image is larger than terminal size
	if ttySize.Width*ttySize.Height*hInc > len(rgbArray)/3 {
		log.Fatal("Capture size too small.")
	}

	// TODO: Improve this inefficient and loosy algorithm
	skipX, skipY := size.Width/ttySize.Width, size.Height/(ttySize.Height*hInc)
	for y := 0; y < ttySize.Height*hInc; y++ {
		for x := 0; x < ttySize.Width; x++ {
			cur := size.Width*3*y*skipY + 3*x*skipX
			ret = append(ret, rgbArray[cur], rgbArray[cur+1], rgbArray[cur+2])
		}
	}
	return
}

func draw(ttyStatus <-chan string, wg *sync.WaitGroup) {
	var interrupt bool = false
	defer wg.Done()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	go func() {
		<-sigchan
		interrupt = true
	}()

	log.Println("Start streaming, press Ctrl-c to exit...")
	time.Sleep(1500 * time.Millisecond)
	clearScreen()

	ttySize := Size{0, 0}

	for !interrupt {
		// Update TTY size before every draw (synchronous)
		curSize := strings.Split(<-ttyStatus, " ")
		ttySize.Height, _ = strconv.Atoi(curSize[0])
		ttySize.Width, _ = strconv.Atoi(curSize[1])

		resetCursor()

		// Fetch image from webcam and call img2xterm to draw
		if asciiOnly {
			rgbRaw := grabRGBPixels(ttySize, 1, 1)
			img2xterm.AsciiDrawRGB(rgbRaw, ttySize.Width, ttySize.Height)
		} else {
			rgbRaw := grabRGBPixels(ttySize, 1, 2)
			img2xterm.DrawRGB(rgbRaw, ttySize.Width, ttySize.Height*2, color)
		}
	}

	restoreScreen()
	log.Println("Exiting...")
}
