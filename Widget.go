package main

import "syscall/js"

type Widget interface {
	Update(svg js.Value)
}

type WidgetStateful interface {
	Widget

	SaveState() js.Value
	LoadState(js.Value)
}
