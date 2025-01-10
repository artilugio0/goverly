package main

import (
	"math/rand/v2"
	"strconv"
	"syscall/js"
	"time"
)

const border int = 0
const svgWidth int = 1920 - 2*border
const svgHeight int = 1080 - 2*border

func main() {
	document := js.Global().Get("document")
	body := document.Get("body")

	svg := document.Call("createElementNS", "http://www.w3.org/2000/svg", "svg")
	svg.Call("setAttribute", "height", strconv.Itoa(svgHeight))
	svg.Call("setAttribute", "width", strconv.Itoa(svgWidth))

	svgStyle := svg.Get("style")
	svgStyle.Set("border", strconv.Itoa(border)+"px solid red")
	body.Call("appendChild", svg)

	widgets := []Widget{
		newCircle(),
		newText("Hoy: Browser overlay para OBS usando Golang + WebAssembly"),
	}

	for {
		for _, w := range widgets {
			w(svg)
		}

		time.Sleep(10 * time.Millisecond)
	}
}

type Widget func(svg js.Value)

func newId() string {
	id := rand.Uint64()
	return strconv.FormatUint(id, 16)
}

func newCircle() Widget {
	const cRadius int = 10

	id := newId()
	document := js.Global().Get("document")
	x, y := 0, 0

	return func(svg js.Value) {
		circle := document.Call("getElementById", id)

		if circle.IsNull() {
			circle = document.Call("createElementNS", "http://www.w3.org/2000/svg", "circle")
			circle.Set("id", id)

			circle.Call("setAttribute", "r", cRadius)
			circle.Call("setAttribute", "r", cRadius)
			circle.Call("setAttribute", "stroke", "lime")
			circle.Call("setAttribute", "stroke-width", "4")
			circle.Call("setAttribute", "fill", "yellow")
			svg.Call("appendChild", circle)

			x, y = svgWidth/2, svgHeight/2
		}

		x = (x + 11) % svgWidth
		y = (y + 11) % svgHeight

		circle.Call("setAttribute", "cx", x)
		circle.Call("setAttribute", "cy", y)
	}
}

func newText(text string) Widget {
	const textSize int = 30
	const textAreaHeight int = 109
	const marginLeft = 50

	id := newId()

	return func(svg js.Value) {
		document := js.Global().Get("document")
		svgtext := document.Call("getElementById", id)
		if !svgtext.IsNull() {
			return
		}

		svgtext = document.Call("createElementNS", "http://www.w3.org/2000/svg", "text")
		svgtext.Set("id", id)

		svgtext.Call("setAttribute", "font-family", "Courier New")
		svgtext.Call("setAttribute", "fill", "white")
		svgtext.Call("setAttribute", "font-size", strconv.Itoa(textSize))
		svgtext.Call("setAttribute", "x", marginLeft)
		svgtext.Call("setAttribute", "y", svgHeight-(textAreaHeight-textSize)/2)

		textnode := document.Call("createTextNode", text)
		svgtext.Call("appendChild", textnode)
		svg.Call("appendChild", svgtext)
	}
}
