package fontcatalog

// #include <stdlib.h>
// #include "fontcatalog_lib.h"
// #cgo CFLAGS: -I ./lib
// #cgo linux CXXFLAGS: -I ./lib -std=c++14
// #cgo darwin CXXFLAGS: -I ./lib  -std=gnu++14
// #cgo darwin,arm CXXFLAGS: -I ./lib  -std=gnu++14
import "C"
import (
	"errors"
	"runtime"
	"unsafe"
)

const (
	MOD_HARD_MASK = "hardmask"
	MOD_SOFT_MASK = "softmask"
	MOD_SDF       = "sdf"
	MOD_PSDF      = "psdf"
	MOD_MSDF      = "msdf"
	MOD_MTSDF     = "mtsdf"
)

type ErrorCorrection uint32

const (
	EC_DISABLED       ErrorCorrection = 0
	EC_INDISCRIMINATE ErrorCorrection = 1
	EC_EDGE_PRIORITY  ErrorCorrection = 2
	EC_EDGE_ONLY      ErrorCorrection = 3
)

type DistanceCheckMode uint32

const (
	DO_NOT_CHECK_DISTANCE  DistanceCheckMode = 0
	CHECK_DISTANCE_AT_EDGE DistanceCheckMode = 1
	ALWAYS_CHECK_DISTANCE  DistanceCheckMode = 2
)

type GeneratorAttributes struct {
	m *C.struct__fc_generator_attributes_t
}

func (h *GeneratorAttributes) free() {
	C.fc_generator_attributes_free(h.m)
}

func NewGeneratorAttributes() *GeneratorAttributes {
	ret := &GeneratorAttributes{m: C.fc_new_generator_attributes()}
	runtime.SetFinalizer(ret, (*GeneratorAttributes).free)
	return ret
}

func (h *GeneratorAttributes) SetMinDeviationRatio(ratio float64) {
	C.fc_generator_attributes_set_min_deviation_ratio(h.m, C.double(ratio))
}

func (h *GeneratorAttributes) SetMinImproveRatio(ratio float64) {
	C.fc_generator_attributes_set_min_improve_ratio(h.m, C.double(ratio))
}

func (h *GeneratorAttributes) SetMode(mode ErrorCorrection) {
	C.fc_generator_attributes_set_mode(h.m, C.uint(mode))
}

func (h *GeneratorAttributes) SetDistanceCheckMode(mode DistanceCheckMode) {
	C.fc_generator_attributes_set_distance_check_mode(h.m, C.uint(mode))
}

func (h *GeneratorAttributes) SetBuffer(buffer []byte) {
	C.fc_generator_attributes_set_buffer(h.m, (*C.uchar)(unsafe.Pointer(&buffer[0])))
}

func (h *GeneratorAttributes) SetOverlapSupport(overlapSupport bool) {
	C.fc_generator_attributes_set_overlap_support(h.m, C.bool(overlapSupport))
}

func (h *GeneratorAttributes) SetScanlinePass(scanlinePass bool) {
	C.fc_generator_attributes_set_scanline_pass(h.m, C.bool(scanlinePass))
}

func glyphGenerator(mode string, output *Bitmap, glyph *GlyphGeometry, attr *GeneratorAttributes) error {
	switch mode {
	case MOD_HARD_MASK:
		if output.GetChannels() == 1 {
			C.fc_scanline_generator(output.m, glyph.m, attr.m)
			return nil
		} else {
			return errors.New("bitmap hardmask channels must 1")
		}
	case MOD_SOFT_MASK:
		if output.GetChannels() == 1 {
			C.fc_sdf_generator(output.m, glyph.m, attr.m)
			return nil
		} else {
			return errors.New("bitmap softmask channels must 1")
		}
	case MOD_SDF:
		if output.GetChannels() == 1 {
			C.fc_sdf_generator(output.m, glyph.m, attr.m)
			return nil
		} else {
			return errors.New("bitmap sdf channels must 1")
		}
	case MOD_PSDF:
		if output.GetChannels() == 1 {
			C.fc_psdf_generator(output.m, glyph.m, attr.m)
			return nil
		} else {
			return errors.New("bitmap psdf channels must 1")
		}
	case MOD_MSDF:
		if output.GetChannels() == 3 {
			C.fc_msdf_generator(output.m, glyph.m, attr.m)
			return nil
		} else {
			return errors.New("bitmap psdf channels must 3")
		}
	case MOD_MTSDF:
		if output.GetChannels() == 4 {
			C.fc_mtsdf_generator(output.m, glyph.m, attr.m)
			return nil
		} else {
			return errors.New("bitmap psdf channels must 4")
		}
	}
	return errors.New("bitmap MOD error")
}
