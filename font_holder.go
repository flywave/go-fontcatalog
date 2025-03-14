package fontcatalog

// #include <stdlib.h>
// #include "fontcatalog_lib.h"
// #cgo CFLAGS: -I ./lib
// #cgo linux CXXFLAGS: -I ./lib -std=c++14
// #cgo darwin CXXFLAGS: -I ./lib  -std=gnu++14
import "C"
import (
	"reflect"
	"runtime"
	"unsafe"
)

type fontInfo struct {
	UnitsPerEm   int
	Bold         bool
	Italic       bool
	LineHeight   int
	LineGap      int
	BaseLine     int
	FontHeight   int
	Ascent       int
	Descent      int
	CharacterSet []rune
}

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

func (h *FontHolder) getFontInfo() *fontInfo {
	info := &fontInfo{}

	metrics := C.fc_font_holder_get_font_info(h.m)
	defer C.free(unsafe.Pointer(metrics.characterSet))

	info.Ascent = int(metrics.ascent)
	info.Descent = int(metrics.descent)
	info.BaseLine = int(metrics.baseLine)
	info.UnitsPerEm = int(metrics.unitsPerEm)
	info.LineHeight = int(metrics.lineHeight)
	info.FontHeight = int(metrics.lineHeight)
	info.Bold = (int(metrics.flags) & 1) != 0
	info.Italic = (int(metrics.flags) & 2) != 0

	info.CharacterSet = make([]rune, int(metrics.charSize))

	var bufSlice []rune
	bufHeader := (*reflect.SliceHeader)((unsafe.Pointer(&bufSlice)))
	bufHeader.Cap = int(metrics.charSize)
	bufHeader.Len = int(metrics.charSize)
	bufHeader.Data = uintptr(unsafe.Pointer(metrics.characterSet))

	copy(info.CharacterSet, bufSlice)

	return info
}
