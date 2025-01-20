package goverly

import (
	"fmt"
	"math/rand/v2"
	"syscall/js"
	"time"
)

type WidgetCircle struct {
	element js.Value `json:"-"`

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

func (wc *WidgetCircle) Update(timePassed time.Duration) []RenderAction {
	wc.X = (wc.X + 1) % svgWidth
	wc.Y = (wc.Y + 1) % svgHeight

	return []RenderAction{
		func(svg js.Value) { wc.element.Call("setAttribute", "cx", wc.X) },
		func(svg js.Value) { wc.element.Call("setAttribute", "cy", wc.Y) },
	}
}

func (wc *WidgetCircle) UpdateConfig(newConfig Widget) []RenderAction {
	actions := []RenderAction{}

	newCfg, ok := newConfig.(*WidgetCircle)
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

func (wc *WidgetCircle) Render(svg js.Value) {
	document := js.Global().Get("document")
	wc.element = document.Call("createElementNS", "http://www.w3.org/2000/svg", "circle")

	wc.element.Call("setAttribute", "stroke", fmt.Sprintf("hsl(%d, 100%%, 50%%)", wc.StrokeHue))
	wc.element.Call("setAttribute", "stroke-width", wc.StrokeWidth)
	wc.element.Call("setAttribute", "fill", fmt.Sprintf("hsl(%d, 100%%, 50%%)", wc.FillHue))
	wc.element.Call("setAttribute", "r", wc.Radius)

	svg.Call("appendChild", wc.element)
}

func (wc *WidgetCircle) SaveState() js.Value {
	state := map[string]interface{}{
		"x": wc.X,
		"y": wc.Y,
	}

	return js.ValueOf(state)
}

func (wc *WidgetCircle) LoadState(state js.Value) {
	wc.X = state.Get("x").Int()
	wc.Y = state.Get("y").Int()
}

func (wc *WidgetCircle) RemoveFromDOM() {
	wc.element.Call("remove")
}

func (wc *WidgetCircle) Type() string {
	return "circle"
}
