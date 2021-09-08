package fontcatalog

import (
	"fmt"
	"image"
	"image/color"

	"github.com/flywave/imaging"
	"github.com/fogleman/gg"
)

type BitmapFontGenerater struct {
	Opt      BitmapFontOptions
	Charsets *Charsets
	holder   *FontHolder
	font     *FontGeometry
	glyphs   *GlyphGeometryList
	attr     *GeneratorAttributes
}

func NewBitmapFontGenerater(holder *FontHolder, charsets *Charsets, opt BitmapFontOptions) *BitmapFontGenerater {
	ret := &BitmapFontGenerater{Opt: opt, Charsets: charsets, holder: holder, glyphs: NewGlyphGeometryList(), attr: NewGeneratorAttributes()}
	ret.font = NewFontGeometryWithGlyphs(ret.glyphs)
	ret.font.LoadFromCharset(ret.holder, float64(ret.Opt.FontSize), ret.Charsets)
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

	fontmetric := g.holder.getFontInfo()
	fontsize := g.Opt.FontSize
	km := g.font.GetKerning()
	font.Kerning = km.GetKernings()

	charsets := make([]string, len(font.Chars))
	for i, c := range font.Chars {
		charsets[i] = c.Char
	}

	font.Info = FontInfo{
		Face:         g.font.GetName(),
		Size:         fontsize,
		Bold:         false,
		Italic:       false,
		Charset:      charsets,
		Unicode:      true,
		StretchHeigt: 100,
		Smooth:       1,
		AA:           1,
		Padding:      [4]int{g.Opt.TexturePadding, g.Opt.TexturePadding, g.Opt.TexturePadding, g.Opt.TexturePadding},
		Spacing:      [2]int{g.Opt.FontSpacing[0], g.Opt.FontSpacing[1]},
	}

	font.Common = FontCommon{
		LineHeight:   fontmetric.LineHeight,
		Base:         fontmetric.BaseLine,
		ScaleW:       g.Opt.TextureSize[0],
		ScaleH:       g.Opt.TextureSize[1],
		Pages:        p,
		Packed:       false,
		AlphaChannel: 0,
		RedChannel:   0,
		GreenChannel: 0,
		BlueChannel:  0,
	}

	font.DistanceField = DistanceField{
		FieldType:     g.Opt.FieldType,
		DistanceRange: g.Opt.DistanceRange,
	}

	return font
}

func (g *BitmapFontGenerater) mapCharsets(start, end int, chars []rune) []*CharsetImage {
	ret := []*CharsetImage{}
	for i := start; i < end; i++ {
		if chars[i] != 0 {
			cimg := generateImage(g.font, chars[i], g.Opt.FieldType, g.Opt.DistanceRange, g.Opt.Border, g.Opt.EdgeColoring, g.Opt.AngleThreshold, g.Opt.Seed, g.attr)
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
	packer := NewMaxRectsBinPacker(g.Opt.TextureSize[0], g.Opt.TextureSize[1], g.Opt.TexturePadding, g.Opt.TexturePadding, g.Opt.AllowRotation)
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
