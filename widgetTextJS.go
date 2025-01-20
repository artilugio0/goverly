//go:build js

package goverly

import (
	"syscall/js"
	"time"
)

type WidgetTextJS struct {
	*WidgetText
	element  js.Value `json:"-"`
	textNode js.Value `json:"-"`
}

func (wt *WidgetTextJS) Update(timePassed time.Duration) []RenderAction {
	return nil
}

func (wt *WidgetTextJS) UpdateConfig(newConfig Widget) []RenderAction {
	actions := []RenderAction{}

	newCfg, ok := newConfig.(*WidgetTextJS)
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

func (wt *WidgetTextJS) Render(svg js.Value) {
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

func (wt *WidgetTextJS) SaveState() js.Value {
	state := map[string]interface{}{}
	return js.ValueOf(state)
}

func (wt *WidgetTextJS) LoadState(state js.Value) {
}

func (wt *WidgetTextJS) RemoveFromDOM() {
	wt.element.Call("remove")
}
