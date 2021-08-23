package fontcatalog

// #include <stdlib.h>
// #include "msdfgen_lib.h"
// #cgo CFLAGS: -I ./lib
// #cgo linux CXXFLAGS: -I ./lib -std=c++14
// #cgo darwin CXXFLAGS: -I ./lib  -std=gnu++14
import "C"
import (
	"io/ioutil"
	"os"
	"runtime"
	"unsafe"
)

type GlyphInfo struct {
	Char    rune
	Width   int
	Height  int
	XOffset int
	YOffset int
	Descent int
	Advance int
	IsCCW   bool
	//Rect     vec2d.Rect
	Renderer *GlyphRender
}

type GlyphRender struct {
	File         string
	Font         *FontHandle
	FontName     string
	UnitsPerEm   int
	Bold         bool
	Italic       bool
	LineHeight   int
	BaseLine     int
	FontHeight   int
	Ascent       int
	Descent      int
	GlyphMap     map[rune]*GlyphInfo
	RenderGlyphs []*GlyphInfo
}

func NewGlyphRender(path string, fontSize int) *GlyphRender {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	data, _ := ioutil.ReadAll(f)
	font := NewFontHandle(data, fontSize)
	ret := &GlyphRender{File: path, Font: font}
	var metrics C.struct__font_metrics_t
	handle := &FontHandle{m: C.msdfgen_load_font_memory((*C.uchar)(unsafe.Pointer(&data[0])), C.long(len(data)), C.int(fontSize), &metrics)}
	runtime.SetFinalizer(handle, (*FontHandle).free)

	ret.Font = handle
	ret.Ascent = int(metrics.ascent)
	ret.Descent = int(metrics.descent)
	ret.BaseLine = int(metrics.baseLine)
	ret.UnitsPerEm = int(metrics.unitsPerEm)
	ret.LineHeight = int(metrics.lineHeight)
	ret.FontHeight = int(metrics.lineHeight)
	ret.Bold = (int(metrics.flags) & 1) != 0
	ret.Italic = (int(metrics.flags) & 2) != 0

	ret.FontName = handle.GetFontName()
	ret.GlyphMap = make(map[rune]*GlyphInfo)
	ret.RenderGlyphs = []*GlyphInfo{}
	return ret
}

func (r *GlyphRender) GetChar(char rune) *GlyphInfo {
	g, ok := r.GlyphMap[char]
	if !ok {
		g = &GlyphInfo{}
		g.Renderer = r
		r.GlyphMap[char] = g
		var m C.struct__glyph_metrics_t
		if bool(C.msdfgen_get_glyph_metrics(r.Font.m, C.int(char), &m)) {
			g.Char = char
			g.Width = int(m.width)
			g.Height = int(m.height)
			g.XOffset = int(m.offsetX)
			g.YOffset = int(m.offsetY)
			g.Advance = int(m.advanceX)
			g.Descent = int(m.descent)
			g.IsCCW = bool(m.ccw)
		}
	}
	if g.Char == -1 {
		return nil
	}
	return g
}
