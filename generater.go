package fontcatalog

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/flywave/imaging"
	"github.com/fogleman/gg"
)

type BitmapFontGenerater struct {
	Opt           BitmapFontOptions
	Charsets      *Charsets
	holder        *FontHolder
	font          *FontGeometry
	glyphs        *GlyphGeometryList
	attr          *GeneratorAttributes
	fontSize      int
	distanceRange float64
}

func NewBitmapFontGenerater(holder *FontHolder, charsets *Charsets, fontSize int, distanceRange float64, opt BitmapFontOptions) *BitmapFontGenerater {
	ret := &BitmapFontGenerater{Opt: opt, Charsets: charsets, holder: holder, glyphs: NewGlyphGeometryList(), attr: NewGeneratorAttributes(), fontSize: fontSize, distanceRange: distanceRange}
	ret.font = NewFontGeometryWithGlyphs(ret.glyphs)
	ret.font.LoadFromCharset(ret.holder, float64(fontSize), ret.Charsets)
	return ret
}

func (g *BitmapFontGenerater) Generate() *BitmapFont {
	font := &BitmapFont{pagesMap: make(map[int]Page), pageSheets: make(map[int]image.Image)}
	start := 0
	done := true
	p := 0

	chars := g.Charsets.GetRunes()

	for done {
		limit := g.Opt.Limit
		if start+g.Opt.Limit > g.Charsets.Size() {
			limit = g.Charsets.Size() - start
			done = false
		}
		images := g.mapCharsets(start, start+limit, chars)
		image, chrs := g.packeCharsets(images, p)
		if image != nil && chrs != nil {
			font.Chars = append(font.Chars, chrs...)
			font.pageSheets[p] = image
			var page string
			if p > 1 {
				page = fmt.Sprintf("%s.%d", g.Opt.Filename, p)
			} else {
				page = g.Opt.Filename
			}
			font.Pages = append(font.Pages, page)
		}
		start += limit
		p++
	}

	if len(font.pageSheets) == 0 {
		return nil
	}

	fontmetric := g.font.GetFontMetrics()
	km := g.font.GetKerning()
	font.Kerning = km.GetKernings()

	charsets := make([]string, len(font.Chars))
	for i, c := range font.Chars {
		charsets[i] = c.Char
	}

	pad := int(0.5 * g.distanceRange)

	font.Info = FontInfo{
		Face:         g.font.GetName(),
		Size:         g.fontSize,
		Bold:         false,
		Italic:       false,
		Charset:      charsets,
		Unicode:      true,
		StretchHeigt: 100,
		Smooth:       1,
		AA:           1,
		Padding:      [4]int{pad, pad, pad, pad},
		Spacing:      [2]int{g.Opt.FontSpacing[0], g.Opt.FontSpacing[1]},
	}

	rect := font.pageSheets[0].Bounds()
	baseline := fontmetric.AscenderY*(float64(g.fontSize)/fontmetric.EmSize) + (0.5 * g.distanceRange)

	font.Common = FontCommon{
		LineHeight:   int(math.Round(fontmetric.LineHeight)),
		Base:         int(math.Round(baseline)),
		ScaleW:       rect.Dx(),
		ScaleH:       rect.Dy(),
		Pages:        p,
		Packed:       0,
		AlphaChannel: 0,
		RedChannel:   0,
		GreenChannel: 0,
		BlueChannel:  0,
	}

	font.DistanceField = DistanceField{
		FieldType:     g.Opt.FieldType,
		DistanceRange: g.distanceRange,
	}

	return font
}

func (g *BitmapFontGenerater) mapCharsets(start, end int, chars []rune) []*CharsetImage {
	ret := []*CharsetImage{}
	for i := start; i < end; i++ {
		if chars[i] != 0 {
			cimg := generateImage(g.font, chars[i], g.Opt.FieldType, g.distanceRange, g.Opt.EdgeColoring, g.Opt.AngleThreshold, g.Opt.Seed, g.attr)
			if cimg != nil {
				ret = append(ret, cimg)
			}
		}
	}
	return ret
}

func (g *BitmapFontGenerater) packeCharsets(images []*CharsetImage, page int) (image.Image, []Charset) {
	if len(images) == 0 {
		return nil, nil
	}
	packer := NewMaxRectsBinPacker(g.Opt.TextureSize[0], g.Opt.TextureSize[1], g.Opt.TexturePadding[0], g.Opt.TexturePadding[1], g.Opt.AllowRotation)
	rects := make([]RectNode, len(images))

	for i := range rects {
		rects[i] = *images[i].glyph.Rect()
	}
	res := packer.Pack(rects, g.Opt.PackerMethod)

	image := image.NewRGBA(image.Rect(0, 0, res.Width, res.Height))

	dbg := gg.NewContextForImage(image)

	if g.Opt.FieldType == MOD_MTSDF || g.Opt.FieldType == MOD_MSDF {
		dbg.SetColor(color.Black)
		dbg.Clear()
	}

	maps := make(map[int]*CharsetImage)
	chars := []Charset{}

	for i := range images {
		maps[images[i].glyph.GetIndex()] = images[i]
	}

	for _, node := range res.PlacedRects {
		img := maps[node.Index]
		if node.Rotated {
			img.image = imaging.Rotate90(img.image)
		}
		fnt := img.font
		fnt.X = node.X
		fnt.Y = node.Y
		fnt.Page = page
		chars = append(chars, fnt)
		dbg.DrawImage(img.image, node.X, node.Y)
	}

	return dbg.Image(), chars
}
