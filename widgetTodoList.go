package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"syscall/js"
)

type WidgetTodoList struct {
	appended bool     `json:"-"`
	element  js.Value `json:"-"`

	X          int            `json:"x"`
	Y          int            `json:"y"`
	Width      int            `json:"width"`
	FontSize   int            `json:"font_size"`
	FontFamily string         `json:"font_family"`
	FontFill   string         `json:"font_fill"`
	Items      []TodoListItem `json:"items"`
}

type TodoListItem struct {
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

func NewTodoList(textSize, width, x, y int, items []TodoListItem) *WidgetTodoList {
	return &WidgetTodoList{
		Items:      items,
		FontSize:   textSize,
		FontFamily: "Courier New",
		FontFill:   "white",
		Width:      width,
		X:          x,
		Y:          y,
	}
}

func (wtl *WidgetTodoList) Update(svg js.Value) {
	if !wtl.appended {
		titleTextSize := wtl.FontSize * 11 / 9
		itemMarginBottom := 5

		document := js.Global().Get("document")

		g := document.Call("createElementNS", "http://www.w3.org/2000/svg", "g")

		title := document.Call("createElementNS", "http://www.w3.org/2000/svg", "text")
		title.Call("setAttribute", "font-family", wtl.FontFamily)
		title.Call("setAttribute", "fill", "white")
		title.Call("setAttribute", "font-size", titleTextSize)
		title.Call("setAttribute", "x", 0)
		title.Call("setAttribute", "y", 0)
		textnode := document.Call("createTextNode", "To Do:")
		title.Call("appendChild", textnode)
		g.Call("appendChild", title)

		itemY := titleTextSize + itemMarginBottom
		for _, it := range wtl.Items {
			listItemText := it.Description
			textFields := strings.Fields(listItemText)

			curWidth := 0
			lines := []string{}
			index := 0
			for i, f := range textFields {
				if curWidth+len(f)*(wtl.FontSize*80/100) > wtl.Width-2 {
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

				curWidth += len(f) * (wtl.FontSize * 80 / 100)
			}

			if len(lines) == 0 {
				lines = append(lines, "- "+it.Description)
			} else {
				lines = append(lines, strings.Join(textFields[index:], " "))
			}

			color := wtl.FontFill
			if it.Done {
				color = "lime"
			}

			for _, line := range lines {

				itemText := document.Call("createElementNS", "http://www.w3.org/2000/svg", "text")
				itemText.Call("setAttribute", "font-family", wtl.FontFamily)
				itemText.Call("setAttribute", "fill", color)
				itemText.Call("setAttribute", "font-size", wtl.FontSize)
				itemText.Call("setAttribute", "x", 0)
				itemText.Call("setAttribute", "y", itemY)

				textnode := document.Call("createTextNode", line)

				itemText.Call("appendChild", textnode)
				g.Call("appendChild", itemText)

				itemY += wtl.FontSize + itemMarginBottom
			}
		}

		g.Call("setAttribute", "transform", fmt.Sprintf("translate(%d, %d)", wtl.X, wtl.Y))
		svg.Call("appendChild", g)

		wtl.appended = true
		wtl.element = g
	}
}

func (wtl *WidgetTodoList) SaveState() js.Value {
	itemsB, err := json.Marshal(wtl.Items)
	if err != nil {
		panic(err)
	}

	state := map[string]interface{}{
		"x":           wtl.X,
		"y":           wtl.Y,
		"width":       wtl.Width,
		"font_size":   wtl.FontSize,
		"font_family": wtl.FontFamily,
		"font_fill":   wtl.FontFill,
		"items":       string(itemsB),
	}

	return js.ValueOf(state)
}

func (wtl *WidgetTodoList) LoadState(state js.Value) {
	wtl.X = state.Get("x").Int()
	wtl.Y = state.Get("y").Int()
	wtl.Width = state.Get("width").Int()
	wtl.FontSize = state.Get("font_size").Int()
	wtl.FontFamily = state.Get("font_family").String()
	wtl.FontFill = state.Get("font_fill").String()
	itemsS := state.Get("items").String()

	if err := json.Unmarshal([]byte(itemsS), &wtl.Items); err != nil {
		panic(err)
	}
}
