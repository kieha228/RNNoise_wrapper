package RNNoise_wrapper

/*
#cgo CFLAGS: -I./include
#cgo LDFLAGS: -L./lib -lRNNoise
#include "rnnoise.h"
*/
import "C"

import "unsafe"

func NewDenoiser() *C.RNNoiseDenoiser {
	return C.rnnoise_create(nil)
}

func DestroyDenoiser(d *C.RNNoiseDenoiser) {
	C.rnnoise_destroy(d)
}

func ProcessFrame(d *C.RNNoiseDenoiser, in, out []float32) {
	cIn := (*C.float)(unsafe.Pointer(&in[0]))
	cOut := (*C.float)(unsafe.Pointer(&out[0]))
	C.rnnoise_process_frame(d, cOut, cIn)
}
