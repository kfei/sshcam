package main

import (
	"strconv"
	"sync"

	flag "github.com/docker/docker/pkg/mflag"
	"github.com/kfei/sshcam/server/sshd"
	webcam "github.com/kfei/sshcam/webcam/v4l2"
)

type Size struct {
	Width, Height int
}

var (
	h, server, color                     bool
	port                                 int
	listen, device, sizeFlag, user, pass string
	size                                 Size
)

func init() {
	flag.BoolVar(&h, []string{"h", "#help", "-help"}, false,
		"display this help message")

	flag.BoolVar(&server, []string{"s", "-server"}, false,
		"start the server")

	flag.StringVar(&listen, []string{"l", "-listen"}, "0.0.0.0",
		"start the server")

	flag.IntVar(&port, []string{"p", "-port"}, 5566,
		"port to listen")

	flag.StringVar(&device, []string{"-device"}, "/dev/video0",
		"the webcam device to open")

	flag.StringVar(&sizeFlag, []string{"-size"}, "640x480",
		"image dimension, must be supported by the device")

	flag.BoolVar(&color, []string{"c", "-color"}, false,
		"turn on color")

	flag.StringVar(&user, []string{"-user"}, "sshcam",
		"username for SSH login")

	flag.StringVar(&pass, []string{"-pass"}, "p@ssw0rd",
		"password for SSH login")

	flag.Parse()
	size = wxh2Size(sizeFlag)
}

func main() {
	switch {
	case h:
		flag.PrintDefaults()
	case server:
		sshcamArgs := []string{"--device=" + device, "--size=" + sizeFlag}
		if color {
			sshcamArgs = append(sshcamArgs, "--color")
		}
		sshd.Run(user, pass, listen, strconv.Itoa(port), sshcamArgs)
	default:
		var wg sync.WaitGroup

		// Initialize the webcam device
		webcam.OpenWebcam(device, size.Width, size.Height)
		defer webcam.CloseWebcam()

		// Start the TTY size updater goroutine
		ttyStatus := updateTTYSize()

		// Fire the drawing goroutine
		wg.Add(1)
		go draw(ttyStatus, &wg)
		wg.Wait()
	}
}
