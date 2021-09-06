package fontcatalog

// #include <stdlib.h>
// #include "msdfgen_lib.h"
// #cgo CFLAGS: -I ./lib
// #cgo linux CXXFLAGS: -I ./lib -std=c++14
// #cgo darwin CXXFLAGS: -I ./lib  -std=gnu++14
// #cgo darwin,arm CXXFLAGS: -I ./lib  -std=gnu++14
// #cgo darwin LDFLAGS: -L ./lib/darwin -lpng -lzlib -lharfbuzz -lfreetype -lmsdfgen -lmsdfgen_ext -llmsdf
// #cgo darwin,arm LDFLAGS: -L ./lib/darwin_arm -lpng -lzlib -lharfbuzz -lfreetype -lmsdfgen -lmsdfgen_ext -llmsdf
// #cgo linux LDFLAGS: -L ./lib/linux -Wl,--start-group -lpthread -ldl -lstdc++ -lm -lpng -lzlib -lharfbuzz -lfreetype -lmsdfgen -lmsdfgen_ext -lmsdf -Wl,--end-group
import "C"
import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"runtime"
	"strings"
	"unsafe"

	"github.com/chai2010/webp"
)

func init() {
	C.wrap_initialize_freetype()
}

const (
	MOD_HARD_MASK = "hardmask"
	MOD_SOFT_MASK = "softmask"
	MOD_SDF       = "sdf"
	MOD_PSDF      = "psdf"
	MOD_MSDF      = "msdf"
	MOD_MTSDF     = "mtsdf"
)

type PackerConfig struct {
	Size        int
	Width       int
	Height      int
	Pot         bool
	Exact       bool
	Sort        string
	Algorithm   string
	Heuristic   string
	UseWasteMap bool
}

type GenConfig struct {
	Input    string
	Inputs   []string
	Output   string
	Charset  []rune
	FontSize int
	Options  []string
	Padding  [4]int
	Spacing  [2]int
	DFSize   int
	Mode     string
	Packer   PackerConfig
}

type FontHandle struct {
	m *C.struct__font_handle_t
}

func NewFontHandle(data []byte, fontSize int) *FontHandle {
	handle := C.msdfgen_load_font_memory((*C.uchar)(unsafe.Pointer(&data[0])), C.long(len(data)), C.int(fontSize), nil)
	ret := &FontHandle{m: handle}
	runtime.SetFinalizer(ret, (*FontHandle).free)
	return ret
}

func (h *FontHandle) free() {
	C.msdfgen_free(h.m)
}

func (h *FontHandle) GetScale() float64 {
	return float64(C.msdfgen_get_scale(h.m))
}

func (h *FontHandle) GetFontName() string {
	var si C.long
	cname := C.msdfgen_get_font_name(h.m, &si)
	defer C.free(unsafe.Pointer(cname))
	return C.GoString(cname)
}

func (h *FontHandle) GenerateSDFGlyph(charcode rune, size [2]int, out []uint8, translate [2]float64, distanceRange float64, ccw bool) bool {
	ret := C.msdfgen_generate_sdf_glyph(h.m, C.int(charcode), C.int(size[0]), C.int(size[1]), (*C.uchar)(unsafe.Pointer(&out[0])), C.double(translate[0]), C.double(translate[1]), C.double(distanceRange), false, C.bool(ccw))
	return bool(ret)
}

func (h *FontHandle) GenerateMSDFGlyph(charcode rune, size [2]int, out []uint8, translate [2]float64, distanceRange float64, ccw bool) bool {
	ret := C.msdfgen_generate_msdf_glyph(h.m, C.int(charcode), C.int(size[0]), C.int(size[1]), (*C.uchar)(unsafe.Pointer(&out[0])), C.double(translate[0]), C.double(translate[1]), C.double(distanceRange), false, C.bool(ccw))
	return bool(ret)
}

func (h *FontHandle) GeneratePSDFGlyph(charcode rune, size [2]int, out []uint8, translate [2]float64, distanceRange float64, ccw bool) bool {
	ret := C.msdfgen_generate_psdf_glyph(h.m, C.int(charcode), C.int(size[0]), C.int(size[1]), (*C.uchar)(unsafe.Pointer(&out[0])), C.double(translate[0]), C.double(translate[1]), C.double(distanceRange), false, C.bool(ccw))
	return bool(ret)
}

func EncodeImage(inputName string, writer io.Writer, rgba image.Image) {
	if strings.HasSuffix(inputName, "jpg") || strings.HasSuffix(inputName, "jpeg") {
		jpeg.Encode(writer, rgba, nil)
	} else if strings.HasSuffix(inputName, "png") {
		png.Encode(writer, rgba)
	} else if strings.HasSuffix(inputName, "gif") {
		gif.Encode(writer, rgba, nil)
	} else if strings.HasSuffix(inputName, "webp") {
		webp.Encode(writer, rgba, &webp.Options{Lossless: true})
	}
}
