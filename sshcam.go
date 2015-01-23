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
	h, server, colorful, asciiOnly                          bool
	port, maxFPS                                            int
	listen, device, sizeFlag, user, pass, distanceAlgorithm string
	size                                                    Size
)

func init() {
	// Arguments for ssh server
	flag.BoolVar(&h, []string{"h", "#help", "-help"}, false,
		"display this help message")
	flag.BoolVar(&server, []string{"s", "-server"}, false,
		"start the server")
	flag.StringVar(&listen, []string{"l", "-listen"}, "0.0.0.0",
		"start the server")
	flag.IntVar(&port, []string{"p", "-port"}, 5566,
		"port to listen")
	flag.StringVar(&user, []string{"-user"}, "sshcam",
		"username for SSH login")
	flag.StringVar(&pass, []string{"-pass"}, "p@ssw0rd",
		"password for SSH login")

	// Arguments for img2xterm
	flag.BoolVar(&colorful, []string{"c", "-color"}, false,
		"turn on color")
	flag.BoolVar(&asciiOnly, []string{"-ascii-only"}, false,
		"fallback to use ASCII's full block characters")
	flag.StringVar(&distanceAlgorithm, []string{"-color-algorithm"}, "yiq",
		"algorithm use to compute colors. Available options are:\n"+
			"'rgb': simple linear distance in RGB colorspace\n"+
			"'yiq': simple linear distance in YIQ colorspace (the default)\n"+
			"'cie94': use the CIE94 formula")
	flag.IntVar(&maxFPS, []string{"-max-fps"}, 4,
		"limit the maximum FPS")
	flag.StringVar(&device, []string{"-device"}, "/dev/video0",
		"the webcam device to open")
	flag.StringVar(&sizeFlag, []string{"-size"}, "640x480",
		"image dimension, must be supported by the device")

	flag.Parse()
	size = wxh2Size(sizeFlag)
}

func main() {
	switch {
	case h:
		flag.PrintDefaults()
	case server:
		// TODO: Better way to copy these arguments to sshd?
		sshcamArgs := []string{
			"--device=" + device,
			"--size=" + sizeFlag,
			"--color-algorithm=" + distanceAlgorithm,
			"--max-fps=" + strconv.Itoa(maxFPS)}
		if colorful {
			sshcamArgs = append(sshcamArgs, "--color")
		}
		if asciiOnly {
			sshcamArgs = append(sshcamArgs, "--ascii-only")
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
		go streaming(ttyStatus, &wg)
		wg.Wait()
	}
}
