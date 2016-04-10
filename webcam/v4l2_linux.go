package webcam

import "github.com/kfei/sshcam/webcam/v4l2"

type Webcam struct{}

func (w *Webcam) OpenWebcam(path string, width, height int) error {
	v4l2.OpenWebcam(path, width, height)
	return nil
}

func (w *Webcam) GrabFrame() []byte {
	return v4l2.GrabFrame()
}

func (w *Webcam) CloseWebcam() error {
	v4l2.CloseWebcam()
	return nil
}
