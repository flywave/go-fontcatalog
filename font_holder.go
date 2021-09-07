package fontcatalog

// #include <stdlib.h>
// #include "fontcatalog_lib.h"
// #cgo CFLAGS: -I ./lib
// #cgo linux CXXFLAGS: -I ./lib -std=c++14
// #cgo darwin CXXFLAGS: -I ./lib  -std=gnu++14
// #cgo darwin,arm CXXFLAGS: -I ./lib  -std=gnu++14
import "C"
import (
	"runtime"
	"unsafe"
)

type FontHolder struct {
	m *C.struct__fc_font_holder_t
}

func NewFontHolder(data []byte) *FontHolder {
	handle := C.fc_font_holder_load_font_memory((*C.uchar)(unsafe.Pointer(&data[0])), C.long(len(data)))
	ret := &FontHolder{m: handle}
	runtime.SetFinalizer(ret, (*FontHolder).free)
	return ret
}

func (h *FontHolder) free() {
	C.fc_font_holder_free(h.m)
}
