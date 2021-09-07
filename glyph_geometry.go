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
)

type EdgeColoring uint32

const (
	EdgeColoringSimple     EdgeColoring = 0
	EdgeColoringInkTrap    EdgeColoring = 1
	EdgeColoringByDistance EdgeColoring = 2
)

type GlyphIdentifierType uint32

const (
	GLYPH_INDEX       GlyphIdentifierType = 0
	UNICODE_CODEPOINT GlyphIdentifierType = 1
)

type GlyphIndex uint32

type GlyphBox struct {
	Index   int
	Advance float64
	Bounds  [4]float64
	Rect    [4]int
}

type GlyphGeometry struct {
	m *C.struct__fc_glyph_geometry_t
}

func NewGlyphGeometryWithGlyphIndex(h *FontHolder, geometryScale float64, index GlyphIndex) *GlyphGeometry {
	handle := C.fc_new_glyph_geometry_from_glyph_index(h.m, C.double(geometryScale), C.fc_glyph_index_t(index))
	ret := &GlyphGeometry{m: handle}
	runtime.SetFinalizer(ret, (*GlyphGeometry).free)
	return ret
}

func NewGlyphGeometryWithCodePoint(h *FontHolder, geometryScale float64, codepoint rune) *GlyphGeometry {
	handle := C.fc_new_glyph_geometry_from_unicode(h.m, C.double(geometryScale), C.fc_unicode_t(codepoint))
	ret := &GlyphGeometry{m: handle}
	runtime.SetFinalizer(ret, (*GlyphGeometry).free)
	return ret
}

func (h *GlyphGeometry) free() {
	C.fc_glyph_geometry_free(h.m)
}

func (h *GlyphGeometry) EdgeColoring(ec EdgeColoring, angleThreshold float64, seed uint64) {
	C.fc_glyph_geometry_edge_coloring(h.m, C.uint(ec), C.double(angleThreshold), C.ulonglong(seed))
}

func (h *GlyphGeometry) WrapBox(scale, range_, miterLimit float64) {
	C.fc_glyph_geometry_wrap_box(h.m, C.double(scale), C.double(range_), C.double(miterLimit))
}

func (h *GlyphGeometry) PlaceBox(x, y int) {
	C.fc_glyph_geometry_place_box(h.m, C.int(x), C.int(y))
}

func (h *GlyphGeometry) GetIndex() int {
	return int(C.fc_glyph_geometry_get_index(h.m))
}

func (h *GlyphGeometry) GetGlyphIndex() GlyphIndex {
	return GlyphIndex(C.fc_glyph_geometry_get_glyph_index(h.m))
}

func (h *GlyphGeometry) GetCodePoint() rune {
	return rune(C.fc_glyph_geometry_get_codepoint(h.m))
}

func (h *GlyphGeometry) GetIdentifier(id GlyphIdentifierType) int {
	return int(C.fc_glyph_geometry_get_identifier(h.m, C.uint(id)))
}

func (h *GlyphGeometry) GetAdvance() float64 {
	return float64(C.fc_glyph_geometry_get_advance(h.m))
}

func (h *GlyphGeometry) GetBoxRect() [4]int {
	var x, y, w, he C.int
	C.fc_glyph_geometry_get_box_rect(h.m, &x, &y, &w, &he)
	return [4]int{int(x), int(y), int(w), int(he)}
}

func (h *GlyphGeometry) GetBoxSize() [2]int {
	var w, he C.int
	C.fc_glyph_geometry_get_box_size(h.m, &w, &he)
	return [2]int{int(w), int(he)}
}

func (h *GlyphGeometry) GetBoxRange() float64 {
	return float64(C.fc_glyph_geometry_get_box_range(h.m))
}

func (h *GlyphGeometry) GetBoxScale() float64 {
	return float64(C.fc_glyph_geometry_get_box_scale(h.m))
}

func (h *GlyphGeometry) GetBoxTranslate() [2]int {
	var tx, ty C.int
	C.fc_glyph_geometry_get_box_translate(h.m, &tx, &ty)
	return [2]int{int(tx), int(ty)}
}

func (h *GlyphGeometry) GetGlyphBox() GlyphBox {
	var gb GlyphBox
	cgb := C.fc_glyph_geometry_get_glyph_box(h.m)
	gb.Advance = float64(cgb.advance)
	gb.Index = int(cgb.index)
	gb.Bounds = [4]float64{
		float64(cgb.bounds.l),
		float64(cgb.bounds.b),
		float64(cgb.bounds.r),
		float64(cgb.bounds.t),
	}
	gb.Rect = [4]int{
		int(cgb.rect.x),
		int(cgb.rect.y),
		int(cgb.rect.w),
		int(cgb.rect.h),
	}
	return gb
}

func (h *GlyphGeometry) IsWhiteSpace() bool {
	return bool(C.fc_glyph_geometry_is_whitespace(h.m))
}

type GlyphGeometryList struct {
	m *C.struct__fc_glyph_geometry_list_t
}

func (h *GlyphGeometryList) free() {
	C.fc_glyph_geometry_list_free(h.m)
}

func NewGlyphGeometryList() *GlyphGeometryList {
	ret := &GlyphGeometryList{m: C.fc_new_glyph_geometry_list()}
	runtime.SetFinalizer(ret, (*GlyphGeometryList).free)
	return ret
}

func (h *GlyphGeometryList) Push(g *GlyphGeometry) {
	C.fc_glyph_geometry_list_push_geometry(h.m, g.m)
}

func (h *GlyphGeometryList) Empty() bool {
	return bool(C.fc_glyph_geometry_list_empty(h.m))
}

func (h *GlyphGeometryList) Size() int {
	return int(C.fc_glyph_geometry_list_size(h.m))
}
