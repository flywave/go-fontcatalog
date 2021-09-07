package fontcatalog

// #include <stdlib.h>
// #include "fontcatalog_lib.h"
// #cgo CFLAGS: -I ./lib
// #cgo linux CXXFLAGS: -I ./lib -std=c++14
// #cgo darwin CXXFLAGS: -I ./lib  -std=gnu++14
// #cgo darwin,arm CXXFLAGS: -I ./lib  -std=gnu++14
import "C"
import (
	"image"
	"image/color"
	"reflect"
	"runtime"
	"unsafe"
)

type BitmapChannel uint32

const (
	GRAY BitmapChannel = 1
	RGB  BitmapChannel = 3
	RGBA BitmapChannel = 4
)

type Bitmap struct {
	m *C.struct__fc_bitmap_t
}

func (h *Bitmap) free() {
	C.fc_bitmap_free(h.m)
}

func NewBitmap(channel BitmapChannel) *Bitmap {
	ret := &Bitmap{m: C.fc_new_bitmap(C.int(channel))}
	runtime.SetFinalizer(ret, (*Bitmap).free)
	return ret
}

func NewBitmapAlloc(channel BitmapChannel, size [2]int) *Bitmap {
	ret := &Bitmap{m: C.fc_new_bitmap_alloc(C.int(channel), C.int(size[0]), C.int(size[1]))}
	runtime.SetFinalizer(ret, (*Bitmap).free)
	return ret
}

func (h *Bitmap) GetWidth() int {
	return int(C.fc_bitmap_width(h.m))
}

func (h *Bitmap) GetHeight() int {
	return int(C.fc_bitmap_height(h.m))
}

func (b *Bitmap) GetChannels() BitmapChannel {
	return BitmapChannel(C.fc_bitmap_channels(b.m))
}

func (b *Bitmap) GetData() []float32 {
	si := b.GetWidth() * b.GetHeight() * int(b.GetChannels())
	cdata := C.fc_bitmap_data(b.m)
	var dSlice []float32
	dHeader := (*reflect.SliceHeader)((unsafe.Pointer(&dSlice)))
	dHeader.Cap = int(si)
	dHeader.Len = int(si)
	dHeader.Data = uintptr(unsafe.Pointer(cdata))
	return dSlice
}

func (b *Bitmap) getBlitData(pix []uint8) {
	si := len(pix)
	cdata := C.fc_bitmap_blit_data(b.m)

	defer C.free(unsafe.Pointer(cdata))

	var dSlice []uint8
	dHeader := (*reflect.SliceHeader)((unsafe.Pointer(&dSlice)))
	dHeader.Cap = int(si)
	dHeader.Len = int(si)
	dHeader.Data = uintptr(unsafe.Pointer(cdata))
	copy(pix, dSlice)
}

func (b *Bitmap) GetBlitData() []uint8 {
	si := b.GetWidth() * b.GetHeight() * int(b.GetChannels())
	ret := make([]uint8, int(si))
	b.getBlitData(ret)
	return ret
}

func (b *Bitmap) GetImage() image.Image {
	switch b.GetChannels() {
	case GRAY:
		img := image.NewGray(image.Rect(0, 0, b.GetWidth(), b.GetHeight()))
		b.getBlitData(img.Pix)
		return img
	case RGB:
		rgbimg := b.GetBlitData()
		img := image.NewRGBA(image.Rect(0, 0, b.GetWidth(), b.GetHeight()))
		for y := 0; y < b.GetHeight(); y++ {
			for x := 0; x < b.GetWidth(); x++ {
				rgb := rgbimg[(y*b.GetWidth()*3)+(x*3):]
				img.SetRGBA(x, y, color.RGBA{R: rgb[0], G: rgb[1], B: rgb[2], A: 255})
			}
		}
		return img
	case RGBA:
		img := image.NewRGBA(image.Rect(0, 0, b.GetWidth(), b.GetHeight()))
		b.getBlitData(img.Pix)
		return img
	}
	return nil
}

type BitmapRef struct {
	m    *C.struct__fc_bitmap_ref_t
	data []float32
}

func (h *BitmapRef) free() {
	C.fc_bitmap_ref_free(h.m)
}

func NewBitmapRefAlloc(data []float32, channel BitmapChannel, size [2]int) *BitmapRef {
	ret := &BitmapRef{m: C.fc_new_bitmap_ref((*C.float)(&data[0]), C.int(channel), C.int(size[0]), C.int(size[1]))}
	runtime.SetFinalizer(ret, (*BitmapRef).free)
	return ret
}

func (h *BitmapRef) GetWidth() int {
	return int(C.fc_bitmap_ref_width(h.m))
}

func (h *BitmapRef) GetHeight() int {
	return int(C.fc_bitmap_ref_height(h.m))
}

func (b *BitmapRef) GetChannels() BitmapChannel {
	return BitmapChannel(C.fc_bitmap_ref_channels(b.m))
}

func (b *BitmapRef) GetData() []float32 {
	si := b.GetWidth() * b.GetHeight() * int(b.GetChannels())
	cdata := C.fc_bitmap_ref_data(b.m)
	var dSlice []float32
	dHeader := (*reflect.SliceHeader)((unsafe.Pointer(&dSlice)))
	dHeader.Cap = int(si)
	dHeader.Len = int(si)
	dHeader.Data = uintptr(unsafe.Pointer(cdata))
	return dSlice
}

func (b *BitmapRef) getBlitData(pix []uint8) {
	si := len(pix)
	cdata := C.fc_bitmap_ref_blit_data(b.m)

	defer C.free(unsafe.Pointer(cdata))

	var dSlice []uint8
	dHeader := (*reflect.SliceHeader)((unsafe.Pointer(&dSlice)))
	dHeader.Cap = int(si)
	dHeader.Len = int(si)
	dHeader.Data = uintptr(unsafe.Pointer(cdata))
	copy(pix, dSlice)
}

func (b *BitmapRef) GetBlitData() []uint8 {
	si := b.GetWidth() * b.GetHeight() * int(b.GetChannels())
	ret := make([]uint8, int(si))
	b.getBlitData(ret)
	return ret
}

func (b *BitmapRef) GetImage() image.Image {
	switch b.GetChannels() {
	case GRAY:
		img := image.NewGray(image.Rect(0, 0, b.GetWidth(), b.GetHeight()))
		b.getBlitData(img.Pix)
		return img
	case RGB:
		rgbimg := b.GetBlitData()
		img := image.NewRGBA(image.Rect(0, 0, b.GetWidth(), b.GetHeight()))
		for y := 0; y < b.GetHeight(); y++ {
			for x := 0; x < b.GetWidth(); x++ {
				rgb := rgbimg[(y*b.GetWidth()*3)+(x*3):]
				img.SetRGBA(x, y, color.RGBA{R: rgb[0], G: rgb[1], B: rgb[2], A: 255})
			}
		}
		return img
	case RGBA:
		img := image.NewRGBA(image.Rect(0, 0, b.GetWidth(), b.GetHeight()))
		b.getBlitData(img.Pix)
		return img
	}
	return nil
}
