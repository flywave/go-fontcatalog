package fontcatalog

type BitmapFontOptions struct {
	Filename       string
	FontSpacing    []int
	TextureSize    []int
	TexturePadding []int
	FieldType      string
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
		FontSpacing:    []int{0, 0},
		TextureSize:    []int{512, 512},
		TexturePadding: []int{1, 1},
		FieldType:      MOD_SDF,
		AllowRotation:  false,
		PackerMethod:   RectBestShortSideFit,
		Limit:          100,
		EdgeColoring:   EdgeColoringInkTrap,
		AngleThreshold: 3.0,
		Seed:           6364136223846793005,
	}
}
