package goverly

import (
	"fmt"
	"strings"
	"syscall/js"
	"time"
)

type WidgetTodoList struct {
	appended bool     `json:"-"`
	element  js.Value `json:"-"`

	X            int            `json:"x"`
	Y            int            `json:"y"`
	Width        int            `json:"width"`
	FontSize     int            `json:"font_size"`
	FontFamily   string         `json:"font_family"`
	FontFill     string         `json:"font_fill"`
	DoneFontFill string         `json:"done_font_fill"`
	Items        []TodoListItem `json:"items"`
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

func (wtl *WidgetTodoList) Update(timePassed time.Duration) []RenderAction {
	return nil
}

func (wtl *WidgetTodoList) UpdateConfig(newConfig Widget) []RenderAction {
	updated := false

	newCfg, ok := newConfig.(*WidgetTodoList)
	if !ok {
		return nil
	}

	if wtl.X != newCfg.X || wtl.Y != newCfg.Y {
		wtl.X = newCfg.X
		wtl.Y = newCfg.Y
		updated = true
	}

	if wtl.Width != newCfg.Width {
		wtl.Width = newCfg.Width
		updated = true
	}

	if wtl.FontFamily != newCfg.FontFamily {
		wtl.FontFamily = newCfg.FontFamily
		updated = true
	}

	if wtl.FontFill != newCfg.FontFill {
		wtl.FontFill = newCfg.FontFill
		updated = true
	}

	if wtl.DoneFontFill != newCfg.DoneFontFill {
		wtl.DoneFontFill = newCfg.DoneFontFill
		updated = true
	}

	if wtl.FontSize != newCfg.FontSize {
		wtl.FontSize = newCfg.FontSize
		updated = true
	}

	if len(wtl.Items) != len(newCfg.Items) {
		updated = true
		wtl.Items = newCfg.Items
	} else {
		for i, item := range wtl.Items {
			if item.Description != newCfg.Items[i].Description || item.Done != newCfg.Items[i].Done {
				updated = true
				wtl.Items = newCfg.Items
				break
			}
		}
	}

	if updated {
		return []RenderAction{wtl.renderItems}
	}

	return nil
}

func (wtl *WidgetTodoList) Render(svg js.Value) {
	document := js.Global().Get("document")
	wtl.element = document.Call("createElementNS", "http://www.w3.org/2000/svg", "g")
	svg.Call("appendChild", wtl.element)

	wtl.renderItems(svg)
}

func (wtl *WidgetTodoList) renderItems(svg js.Value) {
	document := js.Global().Get("document")

	wtl.element.Set("innerHTML", "")

	titleTextSize := wtl.FontSize * 11 / 9
	itemMarginBottom := 5

	title := document.Call("createElementNS", "http://www.w3.org/2000/svg", "text")
	title.Call("setAttribute", "font-family", wtl.FontFamily)
	title.Call("setAttribute", "fill", "white")
	title.Call("setAttribute", "font-size", titleTextSize)
	title.Call("setAttribute", "x", 0)
	title.Call("setAttribute", "y", 0)
	textnode := document.Call("createTextNode", "To Do:")

	title.Call("appendChild", textnode)
	wtl.element.Call("appendChild", title)

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
			color = wtl.DoneFontFill
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
			wtl.element.Call("appendChild", itemText)

			itemY += wtl.FontSize + itemMarginBottom
		}
	}

	wtl.element.Call("setAttribute", "transform", fmt.Sprintf("translate(%d, %d)", wtl.X, wtl.Y))
}

func (wtl *WidgetTodoList) SaveState() js.Value {
	state := map[string]interface{}{}

	return js.ValueOf(state)
}

func (wtl *WidgetTodoList) LoadState(state js.Value) {
}

func (wtl *WidgetTodoList) RemoveFromDOM() {
	wtl.element.Call("remove")
}

func (wtl *WidgetTodoList) Type() string {
	return "todolist"
}
