package main

import (
	"fmt"
	"strconv"
	"strings"
	"syscall/js"
)

type WidgetTodoList struct {
	appended bool
	element  js.Value
}

func NewTodoList(textSize, width, x, y int, items []TodoListItem) *WidgetTodoList {
	titleTextSize := textSize * 11 / 9
	itemMarginBottom := 5

	document := js.Global().Get("document")

	g := document.Call("createElementNS", "http://www.w3.org/2000/svg", "g")

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
			if curWidth+len(f)*(textSize*80/100) > width-2 {
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

			curWidth += len(f) * (textSize * 80 / 100)
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

	return &WidgetTodoList{
		element:  g,
		appended: false,
	}
}

func (wtl *WidgetTodoList) Update(svg js.Value) {
	if !wtl.appended {
		svg.Call("appendChild", wtl.element)
		wtl.appended = true
	}
}

type TodoListItem struct {
	Description string
	Done        bool
}
