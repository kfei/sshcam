package v4l2

import (
	"log"
)

// #include "webcam_wrapper.h"
import "C"
import "unsafe"

var w *C.webcam_t

func OpenWebcam(path string, width, height int) {
	dev := C.CString(path)
	defer C.free(unsafe.Pointer(dev))
	w = C.go_open_webcam(dev, C.int(width), C.int(height))
	// The following defer statement introduces a `double free or corruption`
	// error since it's already freezed in C:
	//
	// defer C.free(unsafe.Pointer(w))

	// Now open the device
	log.Println("Webcam opened")
}

func GrabFrame() []byte {
	buf := C.go_grab_frame(w)
	result := C.GoBytes(unsafe.Pointer(buf.start), C.int(buf.length))
	// Free the buffer (better way for this?)
	if unsafe.Pointer(buf.start) != unsafe.Pointer(uintptr(0)) {
		C.free(unsafe.Pointer(buf.start))
	}
	return result
}

func CloseWebcam() {
	if C.go_close_webcam(w) == 0 {
		log.Println("Webcam closed")
	}
}
