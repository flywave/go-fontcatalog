package fontcatalog

import (
	"fmt"
	math "math"
)

var NullPoint = Point{X: -math.MaxInt32, Y: -math.MaxInt32}

type Point struct {
	X, Y int
}

func PointFromFloats(x, y float32) Point {
	return Point{X: int(x), Y: int(y)}
}

func (p Point) Add(other Point) Point {
	return Point{p.X + other.X, p.Y + other.Y}
}

func (p Point) Scale(scale float32) Point {
	return Point{int(float32(p.X) * scale), int(float32(p.Y) * scale)}
}

func (p Point) EqualsTo(other Point) bool {
	return p.X == other.X && p.Y == other.Y
}

func (p Point) ToString() string {
	return fmt.Sprintf("{x:%d,y:%d}", p.X, p.Y)
}
