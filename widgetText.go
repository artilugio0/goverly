package main

import (
	"syscall/js"
)

type WidgetText struct {
	appended bool     `json:"-"`
	element  js.Value `json:"-"`

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
		element.Call("setAttribute", "font-family", wt.FontFamily)
		element.Call("setAttribute", "fill", wt.FontFill)
		element.Call("setAttribute", "font-size", wt.FontSize)
		element.Call("setAttribute", "x", wt.X)
		element.Call("setAttribute", "y", wt.Y)

		textnode := document.Call("createTextNode", wt.Text)
		element.Call("appendChild", textnode)

		svg.Call("appendChild", element)

		wt.appended = true
		wt.element = element
	}
}

func (wt *WidgetText) SaveState() js.Value {
	state := map[string]interface{}{
		"text":        wt.Text,
		"x":           wt.X,
		"y":           wt.Y,
		"font_family": wt.FontFamily,
		"font_fill":   wt.FontFill,
		"font_size":   wt.FontSize,
	}

	return js.ValueOf(state)
}

func (wt *WidgetText) LoadState(state js.Value) {
	wt.Text = state.Get("text").String()
	wt.X = state.Get("x").Int()
	wt.Y = state.Get("y").Int()
	wt.FontFamily = state.Get("font_family").String()
	wt.FontFill = state.Get("font_fill").String()
	wt.FontSize = state.Get("font_size").Int()
}
