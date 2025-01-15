package main

import (
	"fmt"
	"math"
	"syscall/js"
	"time"
)

var angles []int = []int{0, 5, 0, -5}

type WidgetCountdown struct {
	element  js.Value `json:"-"`
	textNode js.Value `json:"-"`
	timeText js.Value `json:"-"`
	appended bool     `json:"-"`

	AngleMillis int           `json:"angle_millis"`
	AngleIndex  int           `json:"angle_index"`
	FontFamily  string        `json:"font_family"`
	FontFill    string        `json:"font_fill"`
	FontSize    int           `json:"font_size"`
	Remaining   time.Duration `json:"remaining"`
	X           int           `json:"x"`
	Y           int           `json:"y"`
}

func NewCountdown(textSize, x, y int, t time.Duration) *WidgetCountdown {
	return &WidgetCountdown{
		AngleMillis: 0,
		AngleIndex:  0,
		Remaining:   t,
		FontFamily:  "Courier New",
		FontFill:    "white",
		FontSize:    textSize,
		X:           x,
		Y:           y,
	}
}

func (wc *WidgetCountdown) Update(svg js.Value) {
	if !wc.appended {
		document := js.Global().Get("document")

		g := document.Call("createElementNS", "http://www.w3.org/2000/svg", "g")

		timeText := document.Call("createElementNS", "http://www.w3.org/2000/svg", "text")
		timeText.Call("setAttribute", "font-family", wc.FontFamily)
		timeText.Call("setAttribute", "fill", wc.FontFill)
		timeText.Call("setAttribute", "font-size", wc.FontSize)
		timeText.Call("setAttribute", "x", 0)
		timeText.Call("setAttribute", "y", 0)
		textnode := document.Call("createTextNode", "00:00")
		timeText.Call("appendChild", textnode)
		g.Call("appendChild", timeText)
		g.Call("setAttribute", "transform", fmt.Sprintf("translate(%d, %d)", wc.X, wc.Y))
		svg.Call("appendChild", g)

		wc.textNode = textnode
		wc.timeText = timeText
		wc.element = g
		wc.appended = true
	}

	wc.Remaining -= rate
	if wc.Remaining > 0 {
		mins := int(wc.Remaining.Minutes())
		secs := int(wc.Remaining.Seconds()) % 60
		text := fmt.Sprintf("%02d:%02d", mins, secs)
		wc.textNode.Set("nodeValue", text)
		return
	}

	// time is up
	wc.timeText.Call("setAttribute", "fill", "red")
	mins := int(math.Abs(wc.Remaining.Minutes()))
	secs := int(math.Abs(wc.Remaining.Seconds())) % 60
	text := fmt.Sprintf("%02d:%02d", mins, secs)
	wc.textNode.Set("nodeValue", text)

	bbox := wc.timeText.Call("getBBox")
	timeTextWidth := bbox.Get("width").Int()

	wc.AngleMillis += int(rate.Milliseconds())
	if (wc.AngleMillis/1000)%2 == 0 {
		if wc.AngleMillis%50 == 0 {
			wc.AngleIndex = (wc.AngleIndex + 1) % len(angles)
		}

		transform := fmt.Sprintf("rotate(%d, %d, 0)", angles[wc.AngleIndex], timeTextWidth/2)
		wc.timeText.Call("setAttribute", "transform", transform)
		return
	}

	transform := fmt.Sprintf("rotate(%d, %d, 0)", 0, timeTextWidth/2)
	wc.timeText.Call("setAttribute", "transform", transform)
}

func (wc *WidgetCountdown) SaveState() js.Value {
	state := map[string]interface{}{
		"angle_millis": wc.AngleMillis,
		"angle_index":  wc.AngleIndex,
		"font_family":  wc.FontFamily,
		"font_fill":    wc.FontFill,
		"font_size":    wc.FontSize,
		"remaining":    wc.Remaining.Milliseconds(),
		"x":            wc.X,
		"y":            wc.Y,
	}

	return js.ValueOf(state)
}

func (wc *WidgetCountdown) LoadState(state js.Value) {
	wc.AngleMillis = state.Get("angle_millis").Int()
	wc.AngleIndex = state.Get("angle_index").Int()
	wc.Remaining = time.Duration(state.Get("remaining").Int()) * time.Millisecond
	wc.FontFamily = state.Get("font_family").String()
	wc.FontFill = state.Get("font_fill").String()
	wc.FontSize = state.Get("font_size").Int()
	wc.X = state.Get("x").Int()
	wc.Y = state.Get("y").Int()
}
