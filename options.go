package fontcatalog

type BitmapFontOptions struct {
	Filename       string
	Charset        string
	DistanceRange  float64
	FontSize       int
	FontSpacing    []int
	FontPadding    []int
	TextureSize    []int
	TexturePadding int
	Border         float64
	FieldType      string
	SmartSize      bool
	Tolerance      float64
	IsRTL          bool
	AllowRotation  bool
	PackerMethod   FreeRectChoiceHeuristic
	Limit          int
	EdgeColoring   EdgeColoring
	AngleThreshold float64
	Seed           uint64
}

func DefaultBitmapFontOptions(filename string) BitmapFontOptions {
	return BitmapFontOptions{
		Filename:       filename,
		Charset:        "",
		DistanceRange:  4,
		FontSize:       42,
		FontSpacing:    []int{0, 0},
		FontPadding:    []int{4 >> 1, 4 >> 1},
		TextureSize:    []int{512, 512},
		Border:         0,
		FieldType:      MOD_MSDF,
		SmartSize:      false,
		Tolerance:      0,
		IsRTL:          false,
		AllowRotation:  false,
		PackerMethod:   RectBestShortSideFit,
		Limit:          15,
		EdgeColoring:   EdgeColoringInkTrap,
		AngleThreshold: 3.0,
		Seed:           6364136223846793005,
	}
}
