package main

import (
	"fmt"
	"syscall/js"
	"time"
)

type WidgetCountdown struct {
	element  js.Value
	textNode js.Value
	timeText js.Value

	appended bool

	angleMillis int
	angles      []int
	angleIndex  int
	passed      time.Duration

	remaining time.Duration
}

func NewCountdown(textSize, x, y int, t time.Duration) *WidgetCountdown {
	document := js.Global().Get("document")

	g := document.Call("createElementNS", "http://www.w3.org/2000/svg", "g")

	timeText := document.Call("createElementNS", "http://www.w3.org/2000/svg", "text")
	timeText.Call("setAttribute", "font-family", "Courier New")
	timeText.Call("setAttribute", "fill", "white")
	timeText.Call("setAttribute", "font-size", textSize)
	timeText.Call("setAttribute", "x", 0)
	timeText.Call("setAttribute", "y", 0)
	textnode := document.Call("createTextNode", "00:00")
	timeText.Call("appendChild", textnode)
	g.Call("appendChild", timeText)
	g.Call("setAttribute", "transform", fmt.Sprintf("translate(%d, %d)", x, y))

	return &WidgetCountdown{
		appended:    false,
		angleMillis: 0,
		angles:      []int{0, 5, 0, -5},
		angleIndex:  0,
		passed:      0 * time.Second,
		remaining:   t,
		textNode:    textnode,
		timeText:    timeText,
		element:     g,
	}
}

func (wc *WidgetCountdown) Update(svg js.Value) {
	if !wc.appended {
		svg.Call("appendChild", wc.element)
		wc.appended = true
	}
	wc.remaining -= rate
	if wc.remaining > 0 {
		mins := int(wc.remaining.Minutes())
		secs := int(wc.remaining.Seconds()) % 60
		text := fmt.Sprintf("%02d:%02d", mins, secs)
		wc.textNode.Set("nodeValue", text)
		return
	}

	// time is up

	// TODO: replace with wc.remaning
	wc.passed += rate
	wc.timeText.Call("setAttribute", "fill", "red")
	mins := int(wc.passed.Minutes())
	secs := int(wc.passed.Seconds()) % 60
	text := fmt.Sprintf("%02d:%02d", mins, secs)
	wc.textNode.Set("nodeValue", text)

	bbox := wc.timeText.Call("getBBox")
	timeTextWidth := bbox.Get("width").Int()

	wc.angleMillis += int(rate.Milliseconds())
	if (wc.angleMillis/1000)%2 == 0 {
		if wc.angleMillis%50 == 0 {
			wc.angleIndex = (wc.angleIndex + 1) % len(wc.angles)
		}

		transform := fmt.Sprintf("rotate(%d, %d, 0)", wc.angles[wc.angleIndex], timeTextWidth/2)
		wc.timeText.Call("setAttribute", "transform", transform)
		return
	}

	transform := fmt.Sprintf("rotate(%d, %d, 0)", 0, timeTextWidth/2)
	wc.timeText.Call("setAttribute", "transform", transform)
}
