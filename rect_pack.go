package fontcatalog

// #define STBRP_LARGE_RECTS 1
// #define STB_RECT_PACK_IMPLEMENTATION
// #include "stb_rect_pack.h"
import "C"

type (
	Rect C.stbrp_rect
)

func (r *Rect) ID() int {
	return int(r.id)
}

func (r *Rect) X() int {
	return int(r.x)
}

func (r *Rect) Y() int {
	return int(r.y)
}

func (r *Rect) W() int {
	return int(r.w)
}

func (r *Rect) H() int {
	return int(r.h)
}

func (r *Rect) SetX(x int) {
	r.x = C.stbrp_coord(x)
}

func (r *Rect) SetY(y int) {
	r.y = C.stbrp_coord(y)
}

func (r *Rect) SetW(w int) {
	r.w = C.stbrp_coord(w)
}

func (r *Rect) SetH(h int) {
	r.h = C.stbrp_coord(h)
}

func (r *Rect) WasPacked() int {
	return int(r.was_packed)
}
