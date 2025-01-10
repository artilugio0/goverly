package main

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
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
		newTodoList([]TodoListItem{
			{"Migrate canvas to svg", true},
			{"Separate each widget in its own function", true},
			{"Todo list widget", true},
			{"Timer widget", false},
			{"Basic hot reload on recompilation", false},
		}),
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

func newTodoList(items []TodoListItem) Widget {
	const titleTextSize int = 22
	const textSize int = 18
	const marginLeft int = 10
	const itemMarginBottom int = 5
	const width int = 244 - marginLeft
	const x int = 1920 - width
	const y int = 1080 - 850

	id := newId()

	return func(svg js.Value) {
		document := js.Global().Get("document")
		todoList := document.Call("getElementById", id)
		if !todoList.IsNull() {
			return
		}

		g := document.Call("createElementNS", "http://www.w3.org/2000/svg", "g")
		g.Set("id", id)

		title := document.Call("createElementNS", "http://www.w3.org/2000/svg", "text")
		title.Call("setAttribute", "font-family", "Courier New")
		title.Call("setAttribute", "fill", "white")
		title.Call("setAttribute", "font-size", titleTextSize)
		title.Call("setAttribute", "x", 0)
		title.Call("setAttribute", "y", 0)
		textnode := document.Call("createTextNode", "To Do:")
		title.Call("appendChild", textnode)
		g.Call("appendChild", title)

		itemY := titleTextSize + itemMarginBottom
		for _, it := range items {
			listItemText := it.Description
			textFields := strings.Fields(listItemText)

			curWidth := 0
			lines := []string{}
			index := 0
			for i, f := range textFields {
				if curWidth+len(f)*(textSize*70/100) > width-2 {
					if curWidth == 0 {
						index = i + 1
						if len(lines) == 0 {
							lines = append(lines, "- "+textFields[i])
						} else {
							lines = append(lines, "  "+textFields[i])
						}
						continue
					}

					curWidth = 0
					line := strings.Join(textFields[index:i], " ")
					index = i

					if len(lines) == 0 {
						lines = append(lines, "- "+line)
					} else {
						lines = append(lines, line)
					}
				}

				curWidth += len(f) * textSize
			}

			if len(lines) == 0 {
				lines = append(lines, "- "+it.Description)
			} else {
				lines = append(lines, strings.Join(textFields[index:], " "))
			}

			color := "white"
			if it.Done {
				color = "lime"
			}

			for _, line := range lines {

				itemText := document.Call("createElementNS", "http://www.w3.org/2000/svg", "text")
				itemText.Call("setAttribute", "font-family", "Courier New")
				itemText.Call("setAttribute", "fill", color)
				itemText.Call("setAttribute", "font-size", strconv.Itoa(textSize))
				itemText.Call("setAttribute", "x", 0)
				itemText.Call("setAttribute", "y", itemY)

				textnode := document.Call("createTextNode", line)

				itemText.Call("appendChild", textnode)
				g.Call("appendChild", itemText)

				itemY += textSize + itemMarginBottom
			}
		}

		g.Call("setAttribute", "transform", fmt.Sprintf("translate(%d, %d)", x, y))

		svg.Call("appendChild", g)
	}
}

type TodoListItem struct {
	Description string
	Done        bool
}
