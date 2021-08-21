package fontcatalog

// #define STB_TRUETYPE_IMPLEMENTATION
// #include <stb_truetype.h>
import "C"
import (
	"fmt"
	"image"
	"unsafe"
)

type (
	PackContext C.stbtt_pack_context
	PackRange   C.stbtt_pack_range
	PackedChar  C.stbtt_packedchar
)

type FontInfo struct {
	info C.stbtt_fontinfo
	data unsafe.Pointer
}

type AlignedQuad struct {
	X0, Y0, S0, T0 float64
	X1, Y1, S1, T1 float64
}

func makeAlignedQuad(q C.stbtt_aligned_quad) AlignedQuad {
	return AlignedQuad{
		X0: float64(q.x0),
		Y0: float64(q.y0),
		S0: float64(q.s0),
		T0: float64(q.t0),
		X1: float64(q.x1),
		Y1: float64(q.y1),
		S1: float64(q.s1),
		T1: float64(q.t1),
	}
}

func NewPackContext() *PackContext {
	return (*PackContext)(C.calloc(1, C.sizeof_stbtt_pack_context))
}

func NewFontInfo() *FontInfo {
	return (*FontInfo)(C.calloc(1, C.sizeof_stbtt_fontinfo))
}

func NewFontInfoFromData(data []byte, index int) *FontInfo {
	ret := (*FontInfo)(C.calloc(1, C.sizeof_stbtt_fontinfo))
	offset := GetFontOffsetForIndex(data, index)
	ret.Init(data, offset)
	return ret
}

func MakePackedChars(n int) []PackedChar {
	p := C.calloc(C.size_t(n), C.sizeof_stbtt_packedchar)
	s := ((*[1 << 26]PackedChar)(unsafe.Pointer(p)))[:n:n]
	return s
}

func MakeRects(n int) []Rect {
	p := C.calloc(C.size_t(n), C.sizeof_stbrp_rect)
	s := ((*[1 << 26]Rect)(unsafe.Pointer(p)))[:n:n]
	return s
}

func MakePackRanges(n int) []PackRange {
	p := C.calloc(C.size_t(n), C.sizeof_stbtt_pack_range)
	s := ((*[1 << 26]PackRange)(unsafe.Pointer(p)))[:n:n]
	return s
}

func FreePackContext(p *PackContext) {
	C.free(unsafe.Pointer(p))
}

func FreeFontInfo(p *FontInfo) {
	C.free(unsafe.Pointer(p))
}

func FreePackedChars(p []PackedChar) {
	C.free(unsafe.Pointer(&p[0]))
}

func FreePackRanges(p []PackRange) {
	C.free(unsafe.Pointer(&p[0]))
}

func FreeRects(p []Rect) {
	C.free(unsafe.Pointer(&p[0]))
}

func (p *PackContext) Pixels() []byte {
	len := p.width * p.height * p.stride_in_bytes
	buf := ((*[1 << 26]byte)(unsafe.Pointer(p.pixels)))[:len:len]
	return buf
}

func (p *PackContext) StrideInBytes() int {
	return int(p.stride_in_bytes)
}

func (p *PackContext) Begin(pixels []byte, width, height, stride_in_bytes, padding int) error {
	var ptr *byte
	if pixels != nil {
		ptr = &pixels[0]
	}
	rc := C.stbtt_PackBegin((*C.stbtt_pack_context)(p), (*C.uchar)(ptr), C.int(width), C.int(height), C.int(stride_in_bytes), C.int(padding), nil)
	if rc == 0 {
		return fmt.Errorf("out of memory")
	}
	return nil
}

func (p *PackContext) End() {
	C.stbtt_PackEnd((*C.stbtt_pack_context)(p))
}

func (p *PackContext) SetPixels(pixels []byte) {
	p.pixels = (*C.uchar)(&pixels[0])
}

func (p *PackContext) SetHeight(height int) {
	p.height = C.int(height)
}

func (p *PackContext) FontRangesRenderIntoRects(info *FontInfo, ranges []PackRange, rects []Rect) {
	C.stbtt_PackFontRangesRenderIntoRects((*C.stbtt_pack_context)(p), &info.info, (*C.stbtt_pack_range)(&ranges[0]), C.int(len(ranges)), (*C.struct_stbrp_rect)(&rects[0]))
}

func (p *PackContext) FontRangesPackRects(rects []Rect) {
	C.stbtt_PackFontRangesPackRects((*C.stbtt_pack_context)(p), (*C.stbrp_rect)(&rects[0]), C.int(len(rects)))
}

func (p *PackContext) SetOversampling(h_oversample, v_oversample uint) {
	C.stbtt_PackSetOversampling((*C.stbtt_pack_context)(p), C.uint(h_oversample), C.uint(v_oversample))
}

func (p *PackContext) FontRangesGatherRects(info *FontInfo, ranges []PackRange, rects []Rect) int {
	return int(C.stbtt_PackFontRangesGatherRects((*C.stbtt_pack_context)(p), &info.info, (*C.stbtt_pack_range)(&ranges[0]), C.int(len(ranges)), (*C.stbrp_rect)(&rects[0])))
}

func GetPackedQuad(p []PackedChar, pw, ph, char_index int, align_to_integer int) (xpos, ypos float64, q AlignedQuad) {
	var cxpos, cypos C.float
	var cq C.stbtt_aligned_quad
	C.stbtt_GetPackedQuad((*C.stbtt_packedchar)(&p[0]), C.int(pw), C.int(ph), C.int(char_index), &cxpos, &cypos, &cq, C.int(align_to_integer))
	xpos = float64(cxpos)
	ypos = float64(cypos)
	q = makeAlignedQuad(cq)
	return
}

func (p *PackedChar) XAdvance() float64 {
	return float64(p.xadvance)
}

func (p *PackedChar) X0() int {
	return int(p.x0)
}

func (p *PackedChar) Y0() int {
	return int(p.y0)
}

func (p *PackedChar) X1() int {
	return int(p.x1)
}

func (p *PackedChar) Y1() int {
	return int(p.y1)
}

func (p *PackRange) SetFontSize(font_size float64) {
	p.font_size = C.float(font_size)
}

func (p *PackRange) FontSize() float64 {
	return float64(p.font_size)
}

func (p *PackRange) SetFirstUnicodeCodepointInRange(range_ int) {
	p.first_unicode_codepoint_in_range = C.int(range_)
}

func (p *PackRange) FirstUnicodeCodepointInRange() int {
	return int(p.first_unicode_codepoint_in_range)
}

func (p *PackRange) SetNumChars(num_chars int) {
	p.num_chars = C.int(num_chars)
}

func (p *PackRange) NumChars() int {
	return int(p.num_chars)
}

func (p *PackRange) SetChardataForRange(chardata_for_range []PackedChar) {
	p.chardata_for_range = (*C.stbtt_packedchar)(&chardata_for_range[0])
}

func (p *PackRange) CharDataForRange() []PackedChar {
	c := ((*[1 << 26]PackedChar)(unsafe.Pointer(p.chardata_for_range)))[:p.num_chars:p.num_chars]
	return c
}

func (p *PackRange) FirstUnicodepointInRange() int {
	return int(p.first_unicode_codepoint_in_range)
}

func (f *FontInfo) Init(data []byte, offset int) error {
	f.data = C.CBytes(data)
	rc := C.stbtt_InitFont(&f.info, (*C.uchar)(f.data), C.int(offset))
	if rc == 0 {
		return fmt.Errorf("failed to load font")
	}
	return nil
}

func (f *FontInfo) Free() {
	C.free(f.data)
}

func (f *FontInfo) ScaleForPixelHeight(height float64) float64 {
	return float64(C.stbtt_ScaleForPixelHeight(&f.info, C.float(height)))
}

func (f *FontInfo) ScaleForMappingEmToPixels(pixels float64) float64 {
	return float64(C.stbtt_ScaleForMappingEmToPixels(&f.info, C.float(pixels)))
}

func (f *FontInfo) FontVMetrics() (ascent, descent, lineGap int) {
	var cascent, cdescent, clineGap C.int
	C.stbtt_GetFontVMetrics(&f.info, &cascent, &cdescent, &clineGap)
	return int(cascent), int(cdescent), int(clineGap)
}

func (f *FontInfo) CodepointHMetrics(codepoint rune) (advanceWidth, leftSideBearing int) {
	var cadvanceWidth, cleftSideBearing C.int
	C.stbtt_GetCodepointHMetrics(&f.info, C.int(codepoint), &cadvanceWidth, &cleftSideBearing)
	return int(cadvanceWidth), int(cleftSideBearing)
}

func (f *FontInfo) CodepointKernAdvance(ch1, ch2 rune) int {
	return int(C.stbtt_GetCodepointKernAdvance(&f.info, C.int(ch1), C.int(ch2)))
}

func (f *FontInfo) CodepointBox(codepoint rune) image.Rectangle {
	var x0, y0, x1, y1 C.int
	C.stbtt_GetCodepointBox(&f.info, C.int(codepoint), &x0, &y0, &x1, &y1)
	return image.Rect(int(x0), int(y0), int(x1), int(y1))
}

func (f *FontInfo) CodepointBitmap(scale_x, scale_y float64, codepoint rune) (data []byte, width, height, xoff, yoff int) {
	var cw, ch, cxoff, cyoff C.int
	cdata := C.stbtt_GetCodepointBitmap(&f.info, C.float(scale_x), C.float(scale_y), C.int(codepoint), &cw, &ch, &cxoff, &cyoff)
	defer C.free(unsafe.Pointer(cdata))
	return C.GoBytes(unsafe.Pointer(cdata), cw*ch), int(cw), int(ch), int(cxoff), int(cyoff)
}

func (f *FontInfo) CodepointBitmapBoxSubpixel(codepoint rune, scale_x, scale_y, shift_x, shift_y float64) image.Rectangle {
	var x0, y0, x1, y1 C.int
	C.stbtt_GetCodepointBitmapBoxSubpixel(&f.info, C.int(codepoint), C.float(scale_x), C.float(scale_y), C.float(shift_x), C.float(shift_y), &x0, &y0, &x1, &y1)
	return image.Rect(int(x0), int(y0), int(x1), int(y1))
}

func (f *FontInfo) MakeCodepointBitmapSubpixel(output []byte, out_w, out_h, out_stride int, scale_x, scale_y, shift_x, shift_y float64, codepoint rune) {
	C.stbtt_MakeCodepointBitmapSubpixel(&f.info, (*C.uchar)(&output[0]), C.int(out_w), C.int(out_h), C.int(out_stride), C.float(scale_x), C.float(scale_y), C.float(shift_x), C.float(shift_y), C.int(codepoint))
}

func (f *FontInfo) GlyphHMetrics(glyphIndex int) (advanceWidth, leftSideBearing int) {
	var cadvanceWidth, cleftSideBearing C.int
	C.stbtt_GetGlyphHMetrics(&f.info, C.int(glyphIndex), &cadvanceWidth, &cleftSideBearing)
	return int(cadvanceWidth), int(cleftSideBearing)
}

func (f *FontInfo) GlyphKernAdvance(glyph1, glyph2 rune) int {
	return int(C.stbtt_GetGlyphKernAdvance(&f.info, C.int(glyph1), C.int(glyph2)))
}

func (f *FontInfo) BoundingBox() image.Rectangle {
	var x0, y0, x1, y1 C.int
	C.stbtt_GetFontBoundingBox(&f.info, &x0, &y0, &x1, &y1)
	return image.Rect(int(x0), int(y0), int(x1), int(y1))
}

func GetFontOffsetForIndex(data []byte, index int) int {
	return int(C.stbtt_GetFontOffsetForIndex((*C.uchar)(&data[0]), C.int(index)))
}

func GetNumberOfFonts(data []byte) int {
	return int(C.stbtt_GetNumberOfFonts((*C.uchar)(&data[0])))
}
