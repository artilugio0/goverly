//go:build js

package goverly

import (
	"fmt"
	"math"
	"syscall/js"
	"time"
)

var angles []int = []int{0, 5, 0, -5}

type WidgetCountdownJS struct {
	*WidgetCountdown
	// time remaining to reach 00:00
	remaining   time.Duration `json:"-"`
	angleMillis int           `json:"-"`
	angleIndex  int           `json:"-"`

	element  js.Value `json:"-"`
	textNode js.Value `json:"-"`
	timeText js.Value `json:"-"`
}

func (wc *WidgetCountdownJS) UpdateConfig(newConfig Widget) []RenderAction {
	actions := []RenderAction{}

	newCfg, ok := newConfig.(*WidgetCountdownJS)
	if !ok {
		return nil
	}

	if wc.X != newCfg.X || wc.Y != newCfg.Y {
		wc.X = newCfg.X
		wc.Y = newCfg.Y
		actions = append(actions, func(svg js.Value) {
			wc.element.Call("setAttribute", "transform", fmt.Sprintf("translate(%d, %d)", wc.X, wc.Y))
		})
	}

	if wc.FontFamily != newCfg.FontFamily {
		wc.FontFamily = newCfg.FontFamily
		actions = append(actions, func(svg js.Value) {
			wc.element.Call("setAttribute", "font-family", wc.FontFamily)
		})
	}

	if wc.FontFill != newCfg.FontFill {
		wc.FontFill = newCfg.FontFill
		actions = append(actions, func(svg js.Value) {
			wc.element.Call("setAttribute", "fill", wc.FontFill)
		})
	}

	if wc.FontSize != newCfg.FontSize {
		wc.FontSize = newCfg.FontSize
		actions = append(actions, func(svg js.Value) {
			wc.element.Call("setAttribute", "font-size", wc.FontSize)
		})
	}

	if wc.DoneFontFill != newCfg.DoneFontFill {
		wc.DoneFontFill = newCfg.DoneFontFill
		if wc.remaining < 0 {
			actions = append(actions, func(svg js.Value) {
				wc.timeText.Call("setAttribute", "fill", wc.DoneFontFill)
			})
		}
	}

	if wc.EndTime != newCfg.EndTime {
		wasDone := wc.remaining < 0
		wc.EndTime = newCfg.EndTime
		endDate := time.Unix(wc.EndTime, 0)
		wc.remaining = endDate.Sub(time.Now())

		actions = append(actions, func(svg js.Value) {
			mins := int(wc.remaining.Minutes())
			secs := int(wc.remaining.Seconds()) % 60
			text := fmt.Sprintf("%02d:%02d", mins, secs)
			wc.textNode.Set("nodeValue", text)
		})

		if wc.remaining < 0 && !wasDone {
			actions = append(actions, func(svg js.Value) {
				wc.timeText.Call("setAttribute", "fill", wc.DoneFontFill)
			})
		} else if wc.remaining >= 0 && wasDone {
			actions = append(actions, func(svg js.Value) {
				wc.timeText.Call("setAttribute", "fill", wc.FontFill)
				bbox := wc.timeText.Call("getBBox")
				timeTextWidth := bbox.Get("width").Int()
				transform := fmt.Sprintf("rotate(%d, %d, 0)", 0, timeTextWidth/2)
				wc.timeText.Call("setAttribute", "transform", transform)
			})
		}
	}

	if wc.Stopped != newCfg.Stopped {
		wc.Stopped = newCfg.Stopped
	}

	return actions
}

func (wc *WidgetCountdownJS) Update(timePassed time.Duration) []RenderAction {
	if wc.Stopped {
		return nil
	}

	actions := []RenderAction{}
	wasDone := wc.remaining < 0

	wc.remaining -= timePassed
	if wc.remaining > 0 && wc.remaining%1000 == 0 {
		actions = append(actions, func(svg js.Value) {
			mins := int(wc.remaining.Minutes())
			secs := int(wc.remaining.Seconds()) % 60
			text := fmt.Sprintf("%02d:%02d", mins, secs)
			wc.textNode.Set("nodeValue", text)
		})
		return actions
	}

	if !wasDone {
		actions = append(actions, func(svg js.Value) {
			wc.timeText.Call("setAttribute", "fill", wc.DoneFontFill)
		})
	}

	if wc.remaining%1000 == 0 {
		actions = append(actions, func(svg js.Value) {
			mins := int(math.Abs(wc.remaining.Minutes()))
			secs := int(math.Abs(wc.remaining.Seconds())) % 60
			text := fmt.Sprintf("%02d:%02d", mins, secs)
			wc.textNode.Set("nodeValue", text)
		})
	}

	bbox := wc.timeText.Call("getBBox")
	timeTextWidth := bbox.Get("width").Int()

	wc.angleMillis += int(rate.Milliseconds())
	if (wc.angleMillis/1000)%2 == 0 {
		if wc.angleMillis%50 == 0 {
			wc.angleIndex = (wc.angleIndex + 1) % len(angles)

			actions = append(actions, func(svg js.Value) {
				transform := fmt.Sprintf("rotate(%d, %d, 0)", angles[wc.angleIndex], timeTextWidth/2)
				wc.timeText.Call("setAttribute", "transform", transform)
			})
		}
		return actions
	} else if wc.angleMillis%1000 == 0 {
		actions = append(actions, func(svg js.Value) {
			transform := fmt.Sprintf("rotate(%d, %d, 0)", 0, timeTextWidth/2)
			wc.timeText.Call("setAttribute", "transform", transform)
		})
	}

	return actions
}

func (wc *WidgetCountdownJS) Render(svg js.Value) {
	endDate := time.Unix(wc.EndTime, 0)
	wc.remaining = endDate.Sub(time.Now())

	document := js.Global().Get("document")

	wc.element = document.Call("createElementNS", "http://www.w3.org/2000/svg", "g")

	wc.timeText = document.Call("createElementNS", "http://www.w3.org/2000/svg", "text")
	wc.timeText.Call("setAttribute", "x", 0)
	wc.timeText.Call("setAttribute", "y", 0)

	wc.textNode = document.Call("createTextNode", "00:00")
	wc.timeText.Call("appendChild", wc.textNode)

	wc.element.Call("setAttribute", "transform", fmt.Sprintf("translate(%d, %d)", wc.X, wc.Y))
	wc.element.Call("setAttribute", "font-family", wc.FontFamily)
	fill := wc.FontFill
	if wc.remaining <= 0 {
		fill = wc.DoneFontFill
	}
	wc.element.Call("setAttribute", "fill", fill)
	wc.element.Call("setAttribute", "font-size", wc.FontSize)

	wc.element.Call("appendChild", wc.timeText)
	svg.Call("appendChild", wc.element)
}

func (wc *WidgetCountdownJS) SaveState() js.Value {
	state := map[string]interface{}{
		"angle_millis": wc.angleMillis,
		"angle_index":  wc.angleIndex,
		"remaining":    wc.remaining.Milliseconds(),
	}

	return js.ValueOf(state)
}

func (wc *WidgetCountdownJS) LoadState(state js.Value) {
	wc.angleMillis = state.Get("angle_millis").Int()
	wc.angleIndex = state.Get("angle_index").Int()
	wc.remaining = time.Duration(state.Get("remaining").Int()) * time.Millisecond
}

func (wc *WidgetCountdownJS) RemoveFromDOM() {
	wc.element.Call("remove")
}
