package main

import (
	"syscall/js"
	"time"
)

type WidgetText struct {
	element  js.Value `json:"-"`
	textNode js.Value `json:"-"`

	FontFamily string `json:"font_family"`
	FontFill   string `json:"font_fill"`
	FontSize   int    `json:"font_size"`
	Text       string `json:"text"`
	X          int    `json:"x"`
	Y          int    `json:"y"`
}

func NewText(text string, textSize, x, y int) *WidgetText {
	return &WidgetText{
		Text:       text,
		X:          x,
		Y:          y,
		FontFamily: "Courier New",
		FontFill:   "white",
		FontSize:   textSize,
	}
}

func (wt *WidgetText) Update(timePassed time.Duration) []RenderAction {
	return nil
}

func (wt *WidgetText) UpdateConfig(newConfig Widget) []RenderAction {
	actions := []RenderAction{}

	newCfg, ok := newConfig.(*WidgetText)
	if !ok {
		return nil
	}

	if wt.Text != newCfg.Text {
		wt.Text = newCfg.Text
		actions = append(actions, func(svg js.Value) {
			wt.textNode.Set("nodeValue", wt.Text)
		})
	}

	if wt.X != newCfg.X || wt.Y != newCfg.Y {
		wt.X = newCfg.X
		wt.Y = newCfg.Y
		actions = append(actions, func(svg js.Value) {
			wt.element.Call("setAttribute", "x", wt.X)
			wt.element.Call("setAttribute", "y", wt.Y)
		})
	}

	if wt.FontFamily != newCfg.FontFamily {
		wt.FontFamily = newCfg.FontFamily
		actions = append(actions, func(svg js.Value) {
			wt.element.Call("setAttribute", "font-family", wt.FontFamily)
		})
	}

	if wt.FontFill != newCfg.FontFill {
		wt.FontFill = newCfg.FontFill
		actions = append(actions, func(svg js.Value) {
			wt.element.Call("setAttribute", "fill", wt.FontFill)
		})
	}

	if wt.FontSize != newCfg.FontSize {
		wt.FontSize = newCfg.FontSize
		actions = append(actions, func(svg js.Value) {
			wt.element.Call("setAttribute", "font-size", wt.FontSize)
		})
	}

	return actions
}

func (wt *WidgetText) Render(svg js.Value) {
	document := js.Global().Get("document")
	wt.element = document.Call("createElementNS", "http://www.w3.org/2000/svg", "text")

	wt.textNode = document.Call("createTextNode", wt.Text)
	wt.element.Call("appendChild", wt.textNode)

	wt.element.Call("setAttribute", "x", wt.X)
	wt.element.Call("setAttribute", "y", wt.Y)
	wt.element.Call("setAttribute", "font-family", wt.FontFamily)
	wt.element.Call("setAttribute", "fill", wt.FontFill)
	wt.element.Call("setAttribute", "font-size", wt.FontSize)

	svg.Call("appendChild", wt.element)
}

func (wt *WidgetText) SaveState() js.Value {
	state := map[string]interface{}{}
	return js.ValueOf(state)
}

func (wt *WidgetText) LoadState(state js.Value) {
}

func (wt *WidgetText) RemoveFromDOM() {
	wt.element.Call("remove")
}

func (wt *WidgetText) Type() string {
	return "text"
}
