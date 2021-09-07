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

type FontMetrics struct {
	EmSize             float64
	AscenderY          float64
	DescenderY         float64
	LineHeight         float64
	UnderlineY         float64
	UnderlineThickness float64
}

type FontGeometry struct {
	m *C.struct__fc_font_geometry_t
}

func (h *FontGeometry) free() {
	C.fc_font_geometry_free(h.m)
}

func NewFontGeometry() *FontGeometry {
	ret := &FontGeometry{m: C.fc_new_font_geometry()}
	runtime.SetFinalizer(ret, (*FontGeometry).free)
	return ret
}

func NewFontGeometryWithGlyphs(glyphs *GlyphGeometryList) *FontGeometry {
	ret := &FontGeometry{m: C.fc_new_font_geometry_with_glyphs(glyphs.m)}
	runtime.SetFinalizer(ret, (*FontGeometry).free)
	return ret
}

func (h *FontGeometry) LoadFromGlyphset(f *FontHolder, fontScale float64, charsets *Charsets) int {
	return int(C.fc_font_geometry_load_from_glyphset(h.m, f.m, C.double(fontScale), charsets.m))
}

func (h *FontGeometry) LoadFromCharset(f *FontHolder, fontScale float64, charsets *Charsets) int {
	return int(C.fc_font_geometry_load_from_charset(h.m, f.m, C.double(fontScale), charsets.m))
}

func (h *FontGeometry) LoadMetrics(f *FontHolder, fontScale float64) bool {
	return bool(C.fc_font_geometry_load_metrics(h.m, f.m, C.double(fontScale)))
}

func (h *FontGeometry) AddGlyph(glyph *GlyphGeometry) bool {
	return bool(C.fc_font_geometry_add_glyph(h.m, glyph.m))
}

func (h *FontGeometry) LoadKerning(f *FontHolder) int {
	return int(C.fc_font_geometry_load_kerning(h.m, f.m))
}

func (h *FontGeometry) SetName(name string) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	C.fc_font_geometry_set_name(h.m, cname)
}

func (h *FontGeometry) GetName() string {
	cname := C.fc_font_geometry_get_name(h.m)
	defer C.free(unsafe.Pointer(cname))
	return C.GoString(cname)
}

func (h *FontGeometry) GetGeometryScale() float64 {
	return float64(C.fc_font_geometry_geometry_scale(h.m))
}

func (h *FontGeometry) GetFontMetrics() FontMetrics {
	m := FontMetrics{}
	fm := C.fc_font_geometry_get_metrics(h.m)
	m.EmSize = float64(fm.emSize)
	m.AscenderY = float64(fm.ascenderY)
	m.DescenderY = float64(fm.descenderY)
	m.LineHeight = float64(fm.lineHeight)
	m.UnderlineY = float64(fm.underlineY)
	m.UnderlineThickness = float64(fm.underlineThickness)
	return m
}

func (h *FontGeometry) GetPreferredIdentifierType() GlyphIdentifierType {
	return GlyphIdentifierType(C.fc_font_geometry_get_preferred_identifier_type(h.m))
}

func (h *FontGeometry) GetGlyphs() *GlyphRange {
	ret := &GlyphRange{m: C.fc_font_geometry_get_glyphs(h.m)}
	runtime.SetFinalizer(ret, (*GlyphRange).free)
	return ret
}

func (h *FontGeometry) GetGlyphFromIndex(index GlyphIndex) *GlyphGeometry {
	ret := &GlyphGeometry{m: C.fc_font_geometry_get_glyph_from_index(h.m, C.fc_glyph_index_t(index))}
	runtime.SetFinalizer(ret, (*GlyphGeometry).free)
	return ret
}

func (h *FontGeometry) GetGlyphFromCodePoint(codepoint rune) *GlyphGeometry {
	ret := &GlyphGeometry{m: C.fc_font_geometry_get_glyph_from_unicode(h.m, C.fc_unicode_t(codepoint))}
	runtime.SetFinalizer(ret, (*GlyphGeometry).free)
	return ret
}

func (h *FontGeometry) GetAdvanceFromIndex(index1, index2 GlyphIndex) (bool, float64) {
	var advance C.double
	ret := bool(C.fc_font_geometry_get_advance_from_index(h.m, &advance, C.fc_glyph_index_t(index1), C.fc_glyph_index_t(index2)))
	return ret, float64(advance)
}

func (h *FontGeometry) GetAdvanceFromUnicode(codePoint1, codePoint2 rune) (bool, float64) {
	var advance C.double
	ret := bool(C.fc_font_geometry_get_advance_from_unicode(h.m, &advance, C.fc_unicode_t(codePoint1), C.fc_unicode_t(codePoint2)))
	return ret, float64(advance)
}

func (h *FontGeometry) GetKerning() *KerningMap {
	ret := &KerningMap{m: C.fc_font_geometry_get_kerning(h.m)}
	runtime.SetFinalizer(ret, (*KerningMap).free)
	return ret
}

type FontGeometryList struct {
	m *C.struct__fc_font_geometry_list_t
}

func (h *FontGeometryList) free() {
	C.fc_font_geometry_list_free(h.m)
}

func NewFontGeometryList() *FontGeometryList {
	ret := &FontGeometryList{m: C.fc_new_font_geometry_list()}
	runtime.SetFinalizer(ret, (*FontGeometryList).free)
	return ret
}

func (h *FontGeometryList) Push(g *FontGeometry) {
	C.fc_font_geometry_list_push_geometry(h.m, g.m)
}

func (h *FontGeometryList) Empty() bool {
	return bool(C.fc_font_geometry_list_empty(h.m))
}

func (h *FontGeometryList) Size() int {
	return int(C.fc_font_geometry_list_size(h.m))
}

type KerningMap struct {
	m *C.struct__fc_kerning_map_t
}

func (h *KerningMap) free() {
	C.fc_kerning_map_free(h.m)
}

func (h *KerningMap) GetKernings() []Kerning {
	var si C.size_t
	ck := C.fc_kerning_map_get_kernings(h.m, &si)

	var dSlice []C.struct__fc_kerning_t
	dHeader := (*reflect.SliceHeader)((unsafe.Pointer(&dSlice)))
	dHeader.Cap = int(si)
	dHeader.Len = int(si)
	dHeader.Data = uintptr(unsafe.Pointer(ck))
	k := make([]Kerning, int(si))
	for i := range dSlice {
		k[i] = Kerning{First: rune(dSlice[i].first), Second: rune(dSlice[i].second), Amount: float64(dSlice[i].kerning)}
	}
	return k
}

type GlyphRange struct {
	m *C.struct__fc_glyph_range_t
}

func (h *GlyphRange) free() {
	C.fc_glyph_range_free(h.m)
}

func (h *GlyphRange) Empty() bool {
	return bool(C.fc_glyph_range_empty(h.m))
}

func (h *GlyphRange) Size() int {
	return int(C.fc_glyph_range_size(h.m))
}

func (h *GlyphRange) GetGlyphs(index int) *GlyphGeometry {
	ret := &GlyphGeometry{m: C.fc_glyph_range_get(h.m, C.size_t(index))}
	runtime.SetFinalizer(ret, (*GlyphGeometry).free)
	return ret
}
