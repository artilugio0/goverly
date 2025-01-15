package main

import (
	"syscall/js"
)

type WidgetText struct {
	appended bool     `json:"-"`
	element  js.Value `json:"-"`
	textNode js.Value `json:"-"`

	Text string `json:"text"`
	X    int    `json:"x"`
	Y    int    `json:"y"`

	FontFamily string `json:"font_family"`
	FontFill   string `json:"font_fill"`
	FontSize   int    `json:"font_size"`
}

func NewText(text string, textSize, x, y int) *WidgetText {
	return &WidgetText{
		Text:       text,
		X:          x,
		Y:          y,
		FontFamily: "Courier New",
		FontFill:   "White",
		FontSize:   textSize,
	}
}

func (wt *WidgetText) Update(svg js.Value) {
	if !wt.appended {
		document := js.Global().Get("document")
		element := document.Call("createElementNS", "http://www.w3.org/2000/svg", "text")
		wt.textNode = document.Call("createTextNode", wt.Text)
		element.Call("appendChild", wt.textNode)

		svg.Call("appendChild", element)

		wt.appended = true
		wt.element = element
	}

	wt.textNode.Set("nodeValue", wt.Text)

	wt.element.Call("setAttribute", "font-family", wt.FontFamily)
	wt.element.Call("setAttribute", "fill", wt.FontFill)
	wt.element.Call("setAttribute", "font-size", wt.FontSize)
	wt.element.Call("setAttribute", "x", wt.X)
	wt.element.Call("setAttribute", "y", wt.Y)

}

func (wt *WidgetText) SaveState() js.Value {
	state := map[string]interface{}{}
	return js.ValueOf(state)
}

func (wt *WidgetText) LoadState(state js.Value) {
}
