//go:build js

package goverly

import (
	"syscall/js"
	"time"
)

type Widget interface {
	RemoveFromDOM()
	Type() string

	Render(svg js.Value)
	Update(timePassed time.Duration) []RenderAction
	UpdateConfig(newConfig Widget) []RenderAction
}

type RenderAction func(svg js.Value)

type WidgetStateful interface {
	Widget

	SaveState() js.Value
	LoadState(js.Value)
}
