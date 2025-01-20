//go:build js

package goverly

import (
	"fmt"
	"strings"
	"syscall/js"
	"time"
)

type WidgetTodoListJS struct {
	*WidgetTodoList
	appended bool     `json:"-"`
	element  js.Value `json:"-"`
}

func (wtl *WidgetTodoListJS) Update(timePassed time.Duration) []RenderAction {
	return nil
}

func (wtl *WidgetTodoListJS) UpdateConfig(newConfig Widget) []RenderAction {
	updated := false

	newCfg, ok := newConfig.(*WidgetTodoListJS)
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

func (wtl *WidgetTodoListJS) Render(svg js.Value) {
	document := js.Global().Get("document")
	wtl.element = document.Call("createElementNS", "http://www.w3.org/2000/svg", "g")
	svg.Call("appendChild", wtl.element)

	wtl.renderItems(svg)
}

func (wtl *WidgetTodoListJS) renderItems(svg js.Value) {
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

func (wtl *WidgetTodoListJS) SaveState() js.Value {
	state := map[string]interface{}{}

	return js.ValueOf(state)
}

func (wtl *WidgetTodoListJS) LoadState(state js.Value) {
}

func (wtl *WidgetTodoListJS) RemoveFromDOM() {
	wtl.element.Call("remove")
}
