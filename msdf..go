package fontcatalog

// #include <stdlib.h>
// #include <string.h>
// #include <msdf.h>
// #cgo CFLAGS: -I ./  -std=c99 -O2
// #cgo linux LDFLAGS: -lm
import "C"
import (
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"reflect"
	"strings"
	"unsafe"

	"github.com/chai2010/webp"
)

func MinInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func MaxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func Clamp(x, upper, lower int) int {
	return MinInt(upper, MaxInt(x, lower))
}

func MinFloat(x, y float32) float32 {
	if x < y {
		return x
	}
	return y
}

func MaxFloat(x, y float32) float32 {
	if x > y {
		return x
	}
	return y
}

func ClampFloat(x, upper, lower float32) float32 {
	return MinFloat(upper, MaxFloat(x, lower))
}

func median(r, g, b float32) float32 {
	return float32(math.Max(math.Min(float64(r), float64(g)), math.Min(math.Max(float64(r), float64(g)), float64(b))))
}

func lerp(s, e, t float32) float32 {
	return s + (e-s)*t
}

func blerp(c00, c10, c01, c11, tx, ty float32) float32 {
	return lerp(lerp(c00, c10, tx), lerp(c01, c11, tx), ty)
}

func calc_index(x, y, size, num_channels int) int {
	x = Clamp(x, size-1, 0)
	y = Clamp(y, size-1, 0)
	return num_channels * ((y * size) + x)
}

func distVal(dist float32, pxRange *float64, midValue float32) float32 {
	if pxRange == nil {
		if dist > midValue {
			return 1
		} else {
			return 0
		}
	}
	return ClampFloat((dist-midValue)*float32(*pxRange)+.5, 1, 0)
}

type Metrics struct {
	m C.struct_ex_metrics_t
}

func msdfGlyph(finfo *FontInfo, c string, width, height, output_width, output_height int) (*Metrics, image.Image) {
	metrics := Metrics{}
	cstr := C.CString(c)
	defer C.free(unsafe.Pointer(cstr))
	raw := C.ex_msdf_glyph(&finfo.info, C.ex_utf8(cstr), C.size_t(width), C.size_t(height), &(metrics.m), 1)
	defer C.free(unsafe.Pointer(raw))

	var msdf []float32
	bufHeader := (*reflect.SliceHeader)((unsafe.Pointer(&msdf)))
	bufHeader.Cap = int(width * height * 3)
	bufHeader.Len = int(width * height * 3)
	bufHeader.Data = uintptr(unsafe.Pointer(raw))

	bitmap_sdf := image.NewRGBA(image.Rect(0, 0, output_width, output_height))

	for y := 0; y < output_height; y++ {
		for x := 0; x < output_width; x++ {
			gx := x / output_width * width
			gy := y / output_height * height
			gxi := int(gx)
			gyi := int(gy)

			c00 := msdf[calc_index(gxi, gyi, width, 3):]
			c10 := msdf[calc_index(gxi+1, gyi, width, 3):]
			c01 := msdf[calc_index(gxi, gyi+1, width, 3):]
			c11 := msdf[calc_index(gxi+1, gyi+1, width, 3):]

			r := blerp(c00[0], c10[0], c01[0], c11[0], float32(gx-gxi), float32(gy-gyi))
			g := blerp(c00[1], c10[1], c01[1], c11[1], float32(gx-gxi), float32(gy-gyi))
			b := blerp(c00[2], c10[2], c01[2], c11[2], float32(gx-gxi), float32(gy-gyi))

			sigDist := float32(math.Abs(float64(float32(output_width) * (median(r, g, b) - 0.5))))
			a := ClampFloat(sigDist+0.5, 1.0, 0.0)

			color := color.RGBA{}
			color.R = uint8(255 * a)
			color.G = uint8(255 * a)
			color.B = uint8(255 * a)
			color.A = 255

			bitmap_sdf.Set(x, y, color)
		}
	}

	return &metrics, bitmap_sdf
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
