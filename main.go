package main

import (
	"strconv"
	"syscall/js"
	"time"
)

const cRadius int = 10
const textSize int = 30
const border int = 0
const canvasWidth int = 1920 - 2*border
const canvasHeight int = 1080 - 2*border
const text string = "Hoy: PoC overlay de OBS para el stream - [Golang + WebAssembly]"
const marginBottom = 10
const marginLeft = 50
const textAreaHeight int = 109

func main() {
	document := js.Global().Get("document")

	canvas := document.Call("createElement", "canvas")
	canvas.Set("id", "game-canvas")

	canvas.Set("width", strconv.Itoa(canvasWidth))
	canvas.Set("height", strconv.Itoa(canvasHeight))

	canvasStyle := canvas.Get("style")
	canvasStyle.Set("border", strconv.Itoa(border)+"px solid red")

	document.Get("body").Call("appendChild", canvas)

	ctx := canvas.Call("getContext", "2d")

	x, y := canvasWidth/2, canvasHeight/2
	onClick := func(this js.Value, args []js.Value) interface{} {
		e := args[0]
		x = e.Get("clientX").Int()
		y = e.Get("clientY").Int()
		return nil
	}

	canvas.Call("addEventListener", "click", js.FuncOf(onClick))

	for {
		x = (x + 11) % canvasWidth
		y = (y + 11) % canvasHeight

		draw(ctx, x, y)

		time.Sleep(10 * time.Millisecond)
	}
}

func draw(ctx js.Value, x int, y int) {
	width := js.Global().Get("innerWidth").Int()
	height := js.Global().Get("innerHeight").Int()

	ctx.Set("font", strconv.Itoa(textSize)+"px serif")
	metrics := ctx.Call("measureText", text)
	fontHeight := metrics.Get("actualBoundingBoxAscent").Int() + metrics.Get("actualBoundingBoxDescent").Int()

	ctx.Call("clearRect", 0, 0, width, height)

	ctx.Call("beginPath")
	ctx.Set("fillStyle", "lime")
	ctx.Call("arc", x, y, cRadius, 0, js.Global().Get("Math").Get("PI").Float()*2)
	ctx.Call("fill")

	ctx.Set("fillStyle", "white")
	ctx.Set("font", strconv.Itoa(textSize)+"px Courier New")
	ctx.Call("fillText", text, marginLeft, canvasHeight-(textAreaHeight-fontHeight)/2)
}
