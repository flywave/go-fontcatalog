package fontcatalog

import (
	"fmt"
	"image"

	"github.com/flywave/imaging"
	"github.com/fogleman/gg"
)

type BitmapFontGenerater struct {
	Opt      BitmapFontOptions
	Charsets *Charsets
	mapchan  chan []*CharsetImage
	holder   *FontHolder
	font     *FontGeometry
	glyphs   *GlyphGeometryList
	attr     *GeneratorAttributes
}

func NewBitmapFontGenerater(holder *FontHolder, charsets *Charsets, opt BitmapFontOptions) *BitmapFontGenerater {
	glyphs := NewGlyphGeometryList()
	return &BitmapFontGenerater{Opt: opt, Charsets: charsets, mapchan: make(chan []*CharsetImage), holder: holder, glyphs: glyphs, font: NewFontGeometryWithGlyphs(glyphs), attr: NewGeneratorAttributes()}
}

func (g *BitmapFontGenerater) Generate() *BitmapFont {
	font := &BitmapFont{pagesMap: make(map[int]Page), pageSheets: make(map[int]image.Image)}
	start := 0
	done := false
	pageCount := 0
	chars := g.Charsets.GetRunes()

	for done {
		limit := g.Opt.Limit
		if start+g.Opt.Limit > g.Charsets.Size() {
			limit = g.Charsets.Size() - start
			done = true
		}
		go g.mapCharsets(start, start+limit, chars)
		start += limit
		pageCount++
	}

	for p := 0; p < pageCount; p++ {
		images := <-g.mapchan
		image, chrs := g.packeCharsets(images, p)
		font.Chars = append(font.Chars, chrs...)
		font.pageSheets[p] = image
		var page string
		if pageCount > 1 {
			page = fmt.Sprintf("%s.%d", g.Opt.Filename, p)
		} else {
			page = g.Opt.Filename
		}
		font.Pages = append(font.Pages, page)
	}

	fontmetric := g.holder.getFontInfo()
	fontsize := g.Opt.FontSize
	km := g.font.GetKerning()
	font.Kerning = km.GetKernings()

	charsets := make([]string, g.Charsets.Size())
	for i := range chars {
		charsets[i] = string(chars[i])
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
		Pages:        pageCount,
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

func (g *BitmapFontGenerater) mapCharsets(start, end int, chars []rune) {
	ret := make([]*CharsetImage, g.Opt.Limit)
	for i := start; i < end; i++ {
		ret[i] = generateImage(g.font, chars[i], g.Opt.FieldType, g.Opt.DistanceRange, g.Opt.Border, g.Opt.EdgeColoring, g.Opt.AngleThreshold, g.Opt.Seed, g.attr)
	}
	g.mapchan <- ret
}

func (g *BitmapFontGenerater) packeCharsets(images []*CharsetImage, page int) (image.Image, []Charset) {
	packer := NewMaxRectsBinPacker(g.Opt.TextureSize[0], g.Opt.TextureSize[1], g.Opt.TexturePadding, g.Opt.TexturePadding, g.Opt.AllowRotation)
	rects := make([]RectNode, len(images))

	for i := range rects {
		rects[i] = *images[i].glyph.Rect()
	}
	res := packer.Pack(rects, g.Opt.PackerMethod)

	image := image.NewRGBA(image.Rect(0, 0, res.Width, res.Height))

	dbg := gg.NewContextForImage(image)

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
