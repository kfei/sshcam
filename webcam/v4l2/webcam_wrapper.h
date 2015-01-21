#ifndef WEBCAM_WRAPPER_H
#define WEBCAM_WRAPPER_H

#include "webcam.h"

webcam_t* go_open_webcam(const char* dev, int width, int height) {
    webcam_t *w = webcam_open(dev);
    if (w == NULL) {
	fprintf(stderr, "[v4l2] Failed to open the webcam.\n");
	exit(EXIT_FAILURE);
    }

    if (-1 == webcam_resize(w, width, height)) {
	fprintf(stderr, "[v4l2] Failed to resize the webcam\n");
	exit(EXIT_FAILURE);
    }

    webcam_stream(w, true);
    return w;
}

buffer_t go_grab_frame(webcam_t* w) {
    buffer_t frame = { NULL, 0 };
    while(frame.length==0) {
        webcam_grab(w, &frame);
    }
    return frame;
}

int go_close_webcam(webcam_t* w) {
    webcam_stream(w, false);
    webcam_close(w);
    return 0;
}

#endif
