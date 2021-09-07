package fontcatalog

// #include <stdlib.h>
// #include "fontcatalog_lib.h"
// #cgo CFLAGS: -I ./lib
// #cgo linux CXXFLAGS: -I ./lib -std=c++14
// #cgo darwin CXXFLAGS: -I ./lib  -std=gnu++14
// #cgo darwin,arm CXXFLAGS: -I ./lib  -std=gnu++14
import "C"
import (
	"reflect"
	"runtime"
	"unsafe"
)

type Charsets struct {
	m *C.struct__fc_charset_t
}

func (h *Charsets) free() {
	C.fc_charset_free(h.m)
}

func NewCharsets() *Charsets {
	ret := &Charsets{m: C.fc_new_charset()}
	runtime.SetFinalizer(ret, (*Charsets).free)
	return ret
}

func NewCharsetsASCII() *Charsets {
	ret := &Charsets{m: C.fc_new_charset_ascii()}
	runtime.SetFinalizer(ret, (*Charsets).free)
	return ret
}

func (h *Charsets) Empty() bool {
	return bool(C.fc_charset_empty(h.m))
}

func (h *Charsets) Size() int {
	return int(C.fc_charset_size(h.m))
}

func (h *Charsets) Add(code rune) {
	C.fc_charset_add(h.m, C.fc_unicode_t(code))
}

func (h *Charsets) AddRunes(codes []rune) {
	for _, c := range codes {
		C.fc_charset_add(h.m, C.fc_unicode_t(c))
	}
}

func (h *Charsets) Remove(code rune) {
	C.fc_charset_remove(h.m, C.fc_unicode_t(code))
}

func (h *Charsets) GetRunes() []rune {
	var si C.size_t
	cdata := C.fc_charset_data(h.m, &si)
	defer C.free(unsafe.Pointer(cdata))

	ret := make([]rune, int(si))

	var dSlice []rune
	dHeader := (*reflect.SliceHeader)((unsafe.Pointer(&dSlice)))
	dHeader.Cap = int(si)
	dHeader.Len = int(si)
	dHeader.Data = uintptr(unsafe.Pointer(cdata))
	copy(ret, dSlice)
	return ret
}
