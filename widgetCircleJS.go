//go:build js

package goverly

import (
	"fmt"
	"syscall/js"
	"time"
)

type WidgetCircleJS struct {
	*WidgetCircle
	element js.Value `json:"-"`
}

func (wc *WidgetCircleJS) Update(timePassed time.Duration) []RenderAction {
	wc.X = (wc.X + 1) % svgWidth
	wc.Y = (wc.Y + 1) % svgHeight

	return []RenderAction{
		func(svg js.Value) { wc.element.Call("setAttribute", "cx", wc.X) },
		func(svg js.Value) { wc.element.Call("setAttribute", "cy", wc.Y) },
	}
}

func (wc *WidgetCircleJS) UpdateConfig(newConfig Widget) []RenderAction {
	actions := []RenderAction{}

	newCfg, ok := newConfig.(*WidgetCircleJS)
	if !ok {
		return nil
	}

	if wc.StrokeHue != newCfg.StrokeHue {
		wc.StrokeHue = newCfg.StrokeHue
		actions = append(actions, func(svg js.Value) {
			wc.element.Call("setAttribute", "stroke", fmt.Sprintf("hsl(%d, 100%%, 50%%)", wc.StrokeHue))
		})
	}

	if wc.StrokeWidth != newCfg.StrokeWidth {
		wc.StrokeWidth = newCfg.StrokeWidth
		actions = append(actions, func(svg js.Value) {
			wc.element.Call("setAttribute", "stroke-width", wc.StrokeWidth)
		})
	}

	if wc.FillHue != newCfg.FillHue {
		wc.FillHue = newCfg.FillHue
		actions = append(actions, func(svg js.Value) {
			wc.element.Call("setAttribute", "fill", fmt.Sprintf("hsl(%d, 100%%, 50%%)", wc.FillHue))
		})
	}

	if wc.Radius != newCfg.Radius {
		wc.Radius = newCfg.Radius
		actions = append(actions, func(svg js.Value) {
			wc.element.Call("setAttribute", "r", wc.Radius)
		})
	}

	return actions
}

func (wc *WidgetCircleJS) Render(svg js.Value) {
	document := js.Global().Get("document")
	wc.element = document.Call("createElementNS", "http://www.w3.org/2000/svg", "circle")

	wc.element.Call("setAttribute", "stroke", fmt.Sprintf("hsl(%d, 100%%, 50%%)", wc.StrokeHue))
	wc.element.Call("setAttribute", "stroke-width", wc.StrokeWidth)
	wc.element.Call("setAttribute", "fill", fmt.Sprintf("hsl(%d, 100%%, 50%%)", wc.FillHue))
	wc.element.Call("setAttribute", "r", wc.Radius)

	svg.Call("appendChild", wc.element)
}

func (wc *WidgetCircleJS) SaveState() js.Value {
	state := map[string]interface{}{
		"x": wc.X,
		"y": wc.Y,
	}

	return js.ValueOf(state)
}

func (wc *WidgetCircleJS) LoadState(state js.Value) {
	wc.X = state.Get("x").Int()
	wc.Y = state.Get("y").Int()
}

func (wc *WidgetCircleJS) RemoveFromDOM() {
	wc.element.Call("remove")
}
