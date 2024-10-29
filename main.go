package rnnoisew

/*
#cgo LDFLAGS: -lrnnoise
#include <stdlib.h>
#include <rnnoise.h>

// Объявление функции для создания RNNoise
void *create_denoise_state() {
    return rnnoise_create(NULL);
}

// Функция для удаления RNNoise
void destroy_denoise_state(void *st) {
    rnnoise_destroy(st);
}

// Функция для обработки фрейма
float process_frame(void *st, short *inout) {
    return rnnoise_process_frame(st, inout, inout);
}
*/
import "C"
import (
	"errors"
	"unsafe"
)

type Denoise struct {
	state *C.struct_DenoiseState
}

// NewDenoise инициализирует новый объект для подавления шума
func NewDenoise() (*Denoise, error) {
	state := C.create_denoise_state()
	if state == nil {
		return nil, errors.New("failed to create RNNoise state")
	}
	return &Denoise{state: (*C.struct_DenoiseState)(state)}, nil
}

// Process обрабатывает фрейм данных (480 сэмплов для 10 мс при 48 кГц)
func (d *Denoise) Process(frame []int16) (float32, error) {
    if len(frame) != 480 {
        return 0, errors.New("frame must contain exactly 480 samples")
    }

    // Конвертация int16 в float32
    floatFrame := make([]float32, len(frame))
    for i, sample := range frame {
        floatFrame[i] = float32(sample) / 32768.0
    }

    inout := (*C.float)(unsafe.Pointer(&floatFrame[0]))
    vadProb := C.process_frame(unsafe.Pointer(d.state), inout)

    // Конвертация float32 обратно в int16
    for i := range frame {
        frame[i] = int16(floatFrame[i] * 32768.0)
    }
    return float32(vadProb), nil
}

// Close освобождает память, связанную с объектом Denoise
func (d *Denoise) Close() {
	if d.state != nil {
		C.destroy_denoise_state(unsafe.Pointer(d.state))
		d.state = nil
	}
}
