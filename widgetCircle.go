package main

import (
	"fmt"
	"math/rand/v2"
	"syscall/js"
)

type WidgetCircle struct {
	appended bool
	element  js.Value
	x        int
	y        int
}

func NewCircle(radius int) *WidgetCircle {
	document := js.Global().Get("document")

	element := document.Call("createElementNS", "http://www.w3.org/2000/svg", "circle")
	element.Call("setAttribute", "r", radius)
	element.Call("setAttribute", "stroke", fmt.Sprintf("hsl(%d, 100%%, 50%%)", rand.Int()%361))
	element.Call("setAttribute", "stroke-width", "4")
	element.Call("setAttribute", "fill", fmt.Sprintf("hsl(%d, 100%%, 50%%)", rand.Int()%361))

	return &WidgetCircle{
		appended: false,
		element:  element,
		x:        0,
		y:        0,
	}
}

func (wc *WidgetCircle) Update(svg js.Value) {
	if !wc.appended {
		wc.x = svgWidth / 2
		wc.y = svgHeight / 2
		svg.Call("appendChild", wc.element)
		wc.appended = true
	}

	wc.x = (wc.x + 11) % svgWidth
	wc.y = (wc.y + 11) % svgHeight

	wc.element.Call("setAttribute", "cx", wc.x)
	wc.element.Call("setAttribute", "cy", wc.y)
}
