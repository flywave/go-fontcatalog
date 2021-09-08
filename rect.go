package fontcatalog

import (
	"fmt"
	"image"
	"math"
)

type Alignment uint32

const (
	AlignmentHCenter Alignment = 1 << iota
	AlignmentHLeft
	AlignmentHRight
	AlignmentVCenter
	AlignmentVTop
	AlignmentVBottom
	AlignmentNone   Alignment = 0
	AlignmentCenter           = AlignmentHCenter | AlignmentVCenter
)

type FitMode int

const (
	FitModeAlign FitMode = iota
	FitModeFill
	FitModeAspectFit
	FitModeAspectFill
)

type Insets struct {
	Top, Right, Bottom, Left int
}

func HomogeneousInsets(inset int) Insets {
	return Insets{Top: inset, Right: inset, Bottom: inset, Left: inset}
}

type Rect struct {
	X, Y, W, H int
}

var EmptyRect = Rect{X: 0, Y: 0, W: 0, H: 0}

func RectFromArray(values [4]int) Rect {
	return Rect{X: values[0], Y: values[1], W: values[2], H: values[3]}
}

func RectFromFloats(x, y, w, h float32) Rect {
	return Rect{X: int(x), Y: int(y), W: int(w), H: int(h)}
}

func RectFromRectangle(r image.Rectangle) Rect {
	return Rect{X: r.Min.X, Y: r.Min.Y, W: r.Dx(), H: r.Dy()}
}

func (r Rect) ToRectangle() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{X: r.X, Y: r.Y},
		Max: image.Point{X: r.X + r.W, Y: r.Y + r.H},
	}
}

func (r Rect) Left() int    { return r.X }
func (r Rect) Top() int     { return r.Y }
func (r Rect) Right() int   { return r.X + r.W }
func (r Rect) Bottom() int  { return r.Y + r.H }
func (r Rect) CenterX() int { return r.X + r.W/2 }
func (r Rect) CenterY() int { return r.Y + r.H/2 }

func (r Rect) Center() Point {
	return Point{r.X + r.W/2, r.Y + r.H/2}
}

func (r Rect) MinPoint() Point {
	return Point{r.X, r.Y}
}

func (r Rect) MaxPoint() Point {
	return Point{r.Right(), r.Bottom()}
}

func (r Rect) Size() Size {
	return Size{r.W, r.H}
}

func (r Rect) Empty() bool {
	return r.X == 0 && r.Y == 0 && r.W == 0 && r.H == 0
}

func (r Rect) MoveTo(x, y int) Rect {
	r.X = x
	r.Y = y
	return r
}

func (r Rect) Translate(x, y int) Rect {
	r.X += x
	r.Y += y
	return r
}

func (r Rect) TranslatePoint(point Point) Rect {
	r.X += point.X
	r.Y += point.Y
	return r
}

func (r Rect) TranslateFloat(x, y float32) Rect {
	r.X += int(x)
	r.Y += int(y)
	return r
}

func (r Rect) ResizeTo(w, h int) Rect {
	r.W = w
	r.H = h
	return r
}

func (r Rect) ShrinkByInsets(i Insets) Rect {
	return Rect{
		X: r.X + i.Left,
		Y: r.Y + i.Top,
		W: r.W - i.Right - i.Left,
		H: r.H - i.Bottom - i.Top,
	}
}

func (r Rect) ShrinkByInt(i int) Rect {
	return Rect{
		X: r.X + i,
		Y: r.Y + i,
		W: r.W - i*2,
		H: r.H - i*2,
	}
}

func Round32(value float32) float32 {
	return float32(math.Round(float64(value)))
}

func (r Rect) Scale(factor float32) Rect {
	return Rect{
		X: int(Round32(float32(r.X) * factor)),
		Y: int(Round32(float32(r.Y) * factor)),
		W: int(Round32(float32(r.W) * factor)),
		H: int(Round32(float32(r.H) * factor)),
	}
}

func (r Rect) CenterIn(o Rect) Rect {
	hW := (o.W - r.W) / 2
	hH := (o.H - r.H) / 2
	return Rect{
		X: r.X + hW,
		Y: r.Y + hH,
		W: r.W + hW*2,
		H: r.H + hH*2,
	}
}

func (r Rect) AlignIn(b Rect, alignment Alignment) Rect {
	newRect := r
	if alignment&AlignmentHLeft != 0 {
		newRect.X = b.X
	} else if alignment&AlignmentHCenter != 0 {
		newRect.X = b.X + (b.W-r.W)/2
	} else if alignment&AlignmentHRight != 0 {
		newRect.X = b.Right() - r.W
	}
	if alignment&AlignmentVTop != 0 {
		newRect.Y = b.Y
	} else if alignment&AlignmentVCenter != 0 {
		newRect.Y = b.Y + (b.H-r.H)/2
	} else if alignment&AlignmentVBottom != 0 {
		newRect.Y = b.Bottom() - r.H
	}
	return newRect
}

func (r Rect) FitIn(b Rect, mode FitMode, alignment Alignment) Rect {
	switch mode {
	case FitModeFill:
		return b
	case FitModeAspectFit:
		if r.W > r.H {
			f := float32(b.W) / float32(r.W)
			r.W = int(f * float32(r.W))
			r.H = int(f * float32(r.H))
		} else {
			f := float32(b.H) / float32(r.H)
			r.W = int(f * float32(r.W))
			r.H = int(f * float32(r.H))
		}
	case FitModeAspectFill:
		if r.W > r.H {
			f := float32(b.H) / float32(r.H)
			r.W = int(f * float32(r.W))
			r.H = int(f * float32(r.H))
		} else {
			f := float32(b.W) / float32(r.W)
			r.W = int(f * float32(r.W))
			r.H = int(f * float32(r.H))
		}
	}
	return r.AlignIn(b, alignment)
}

func (r Rect) UnionWith(other Rect) Rect {
	x1 := Min(r.X, other.X)
	y1 := Min(r.Y, other.Y)
	x2 := Max(r.Right(), other.Right())
	y2 := Max(r.Bottom(), other.Bottom())
	return Rect{X: x1, Y: y1, W: x2 - x1, H: y2 - y1}
}

func (r Rect) ContainsPoint(pointX, pointY int) bool {
	pointX -= r.X
	pointY -= r.Y
	return pointX >= 0 && pointY >= 0 && pointX < r.W && pointY < r.H
}

func (r Rect) IsContainedIn(b Rect) bool {
	return r.X >= b.X && r.Y >= b.Y && r.Right() <= b.Right() && r.Bottom() <= b.Bottom()
}

func (r Rect) Intersect(r2 Rect) bool {
	if r.X >= r2.X+r2.W || r2.X >= r.X+r.W {
		return false
	}
	if r.Y >= r2.Y+r2.H || r2.Y >= r.Y+r.H {
		return false
	}
	return true
}

func (r Rect) IntersectWithCircle(circleX, circleY, circleRadius int) bool {
	dX := circleX - Max(r.X, Min(circleX, r.X+r.W))
	dY := circleY - Max(r.Y, Min(circleY, r.Y+r.H))
	return (dX*dX + dY*dY) < (circleRadius * circleRadius)
}

func (r Rect) Intersection(s Rect) Rect {
	x2 := Min(r.Right(), s.Right())
	y2 := Min(r.Bottom(), s.Bottom())

	if r.X < s.X {
		r.X = s.X
	}
	if r.Y < s.Y {
		r.Y = s.Y
	}

	r.W = x2 - r.X
	r.H = y2 - r.Y

	if r.W > 0 && r.H > 0 {
		return r
	}
	return Rect{}
}

func (r Rect) EqualsTo(other Rect) bool {
	return r.X == other.X && r.Y == other.Y && r.W == other.W && r.H == other.H
}

func (r Rect) ToString() string {
	return fmt.Sprintf("{x:%d,y:%d,w:%d,h:%d}", r.X, r.Y, r.W, r.H)
}
