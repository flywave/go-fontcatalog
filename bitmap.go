package fontcatalog

// #include <stb_truetype.h>
import "C"

import (
	"fmt"
	"image"
	"image/color"
	"unsafe"
)

type (
	BakedChar C.stbtt_bakedchar
)

type Bitmap struct {
	*image.RGBA
	*FontInfo
	Chardata    []BakedChar
	Firstchar   rune
	PixelHeight float64
	FG, BG      color.RGBA
}

func BakeFontBitmap(data []byte, offset int, pixel_height float64, pw, ph int, first_char rune, num_chars int) (bmp *Bitmap, numfits int, err error) {
	pixels := make([]byte, pw*ph)
	chardata := make([]BakedChar, num_chars)
	numfits = int(C.stbtt_BakeFontBitmap((*C.uchar)(unsafe.Pointer(&data[0])), C.int(offset), C.float(pixel_height), (*C.uchar)(unsafe.Pointer(&pixels[0])), C.int(pw), C.int(ph), C.int(first_char), C.int(num_chars), (*C.stbtt_bakedchar)(unsafe.Pointer(&chardata[0]))))
	img := image.NewRGBA(image.Rect(0, 0, pw, ph))
	for y := 0; y < ph; y++ {
		for x := 0; x < pw; x++ {
			if pixels[y*pw+x] != 0 {
				img.Set(x, y, color.White)
			}
		}
	}
	bmp = &Bitmap{
		RGBA:        img,
		FontInfo:    &FontInfo{},
		Chardata:    chardata,
		Firstchar:   first_char,
		PixelHeight: pixel_height,
		FG:          color.RGBA{255, 255, 255, 255},
	}
	err = bmp.Init(data, offset)
	return
}

func BakedQuad(cdata []BakedChar, pw, ph, char_index int, xpos, ypos float64, opengl_fillrule int) (float64, float64, AlignedQuad) {
	var cq C.stbtt_aligned_quad
	cxpos := C.float(xpos)
	cypos := C.float(ypos)
	C.stbtt_GetBakedQuad((*C.stbtt_bakedchar)(unsafe.Pointer(&cdata[0])), C.int(pw), C.int(ph), C.int(char_index), &cxpos, &cypos, &cq, C.int(opengl_fillrule))
	return float64(cxpos), float64(cypos), makeAlignedQuad(cq)
}

func (b *Bitmap) Print(m *image.RGBA, x, y int, args ...interface{}) {
	b.print(m, x, y, fmt.Sprint(args...))
}

func (b *Bitmap) Printf(m *image.RGBA, x, y int, format string, args ...interface{}) {
	b.print(m, x, y, fmt.Sprintf(format, args...))
}

func (b *Bitmap) print(m *image.RGBA, x, y int, s string) {
	r := b.Bounds()
	px := float64(x)
	py := float64(y)
	sx := px
	for _, c := range s {
		if c == '\n' {
			px = sx
			py += b.PixelHeight
			continue
		}
		c -= rune(b.Firstchar)

		_, _, q := BakedQuad(b.Chardata, r.Dx(), r.Dy(), int(c), px, py, 1)
		dr := image.Rect(int(q.X0), int(q.Y0), int(q.X1), int(q.Y1))
		dr = dr.Add(image.Pt(0, int(b.PixelHeight)))

		s0 := q.S0 * float64(r.Dx())
		t0 := q.T0 * float64(r.Dy())
		sp := image.Pt(int(s0), int(t0))

		for y, ty := dr.Min.Y, 0; y < dr.Max.Y; y, ty = y+1, ty+1 {
			for x, tx := dr.Min.X, 0; x < dr.Max.X; x, tx = x+1, tx+1 {
				col := b.RGBAAt(sp.X+tx, sp.Y+ty)
				if col == (color.RGBA{}) {
					if b.BG != (color.RGBA{}) {
						m.Set(x, y, b.BG)
					}
				} else {
					m.Set(x, y, b.FG)
				}
			}
		}
	}
}

func (b *Bitmap) StringSize(text string) (width, height float64) {
	w, h := 0.0, b.PixelHeight
	mw := w
	for _, c := range text {
		if c == '\n' {
			w = 0
			h += b.PixelHeight
			continue
		}

		a, _ := b.CodepointHMetrics(c)
		w += float64(a)
		if mw < w {
			mw = w
		}
	}
	w *= b.ScaleForPixelHeight(b.PixelHeight)
	return mw, h
}
