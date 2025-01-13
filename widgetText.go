package main

import (
	"strconv"
	"syscall/js"
)

type WidgetText struct {
	appended bool
	element  js.Value
}

func NewText(text string, textSize, x, y int) *WidgetText {
	document := js.Global().Get("document")

	element := document.Call("createElementNS", "http://www.w3.org/2000/svg", "text")
	element.Call("setAttribute", "font-family", "Courier New")
	element.Call("setAttribute", "fill", "white")
	element.Call("setAttribute", "font-size", strconv.Itoa(textSize))
	element.Call("setAttribute", "x", x)
	element.Call("setAttribute", "y", y)

	textnode := document.Call("createTextNode", text)
	element.Call("appendChild", textnode)

	return &WidgetText{
		element: element,
	}
}

func (wt *WidgetText) Update(svg js.Value) {
	if !wt.appended {
		svg.Call("appendChild", wt.element)
		wt.appended = true
	}
}
