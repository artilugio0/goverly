package main

import (
	"fmt"
	"math/rand/v2"
	"syscall/js"
)

type WidgetCircle struct {
	appended bool     `json:"-"`
	element  js.Value `json:"-"`

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

func (wc *WidgetCircle) Update(svg js.Value) {
	if !wc.appended {
		document := js.Global().Get("document")
		element := document.Call("createElementNS", "http://www.w3.org/2000/svg", "circle")
		element.Call("setAttribute", "r", wc.Radius)
		element.Call("setAttribute", "stroke-width", wc.StrokeWidth)
		element.Call("setAttribute", "stroke", fmt.Sprintf("hsl(%d, 100%%, 50%%)", wc.StrokeHue))
		element.Call("setAttribute", "fill", fmt.Sprintf("hsl(%d, 100%%, 50%%)", wc.FillHue))

		svg.Call("appendChild", element)

		wc.appended = true
		wc.element = element
	}

	wc.X = (wc.X + 1) % svgWidth
	wc.Y = (wc.Y + 1) % svgHeight

	wc.element.Call("setAttribute", "cx", wc.X)
	wc.element.Call("setAttribute", "cy", wc.Y)
}

func (wc *WidgetCircle) SaveState() js.Value {
	state := map[string]interface{}{
		"x":            wc.X,
		"y":            wc.Y,
		"stroke_hue":   wc.StrokeHue,
		"fill_hue":     wc.FillHue,
		"radius":       wc.Radius,
		"stroke_width": wc.StrokeWidth,
	}

	return js.ValueOf(state)
}

func (wc *WidgetCircle) LoadState(state js.Value) {
	wc.X = state.Get("x").Int()
	wc.Y = state.Get("y").Int()
	wc.StrokeHue = state.Get("stroke_hue").Int()
	wc.FillHue = state.Get("fill_hue").Int()
	wc.Radius = state.Get("radius").Int()
	wc.StrokeWidth = state.Get("stroke_width").Int()
}
