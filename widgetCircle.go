package goverly

import (
	"math/rand/v2"
)

type WidgetCircle struct {
	X           int `json:"x"`
	Y           int `json:"y"`
	StrokeHue   int `json:"stroke_hue"`
	FillHue     int `json:"fill_hue"`
	Radius      int `json:"radius"`
	StrokeWidth int `json:"stroke_width"`
}

func NewCircle(radius int) *WidgetCircle {
	return &WidgetCircle{
		X:           0,
		Y:           0,
		StrokeHue:   rand.Int() % 361,
		FillHue:     rand.Int() % 361,
		Radius:      radius,
		StrokeWidth: 4,
	}
}
func (wc *WidgetCircle) Type() string {
	return "circle"
}
