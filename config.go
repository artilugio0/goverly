package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Config struct {
	Widgets map[string]Widget `json:"widgets"`
}

func readConfig(configPath string) (*Config, error) {
	fileBytes, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	config := Config{}
	if err := json.Unmarshal(fileBytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) UnmarshalJSON(b []byte) error {
	body := map[string]any{}
	if err := json.Unmarshal(b, &body); err != nil {
		return err
	}

	widgetDefs, ok := body["widgets"].(map[string]any)
	if !ok {
		return nil
	}

	c.Widgets = map[string]Widget{}

	for k, d := range widgetDefs {
		widgetDef := d.(map[string]any)
		wType, ok := widgetDef["type"].(string)
		if !ok {
			return fmt.Errorf("invalid value for 'type' key in widget. Found: '%+v'", widgetDef["type"])
		}

		wObj, ok := widgetDef["widget"].(map[string]any)
		if !ok {
			return fmt.Errorf("invalid value for 'widget' key in widget. Found: '%+v'", widgetDef["widget"])
		}

		var widget Widget
		switch wType {
		case "circle":
			circle := WidgetCircle{}
			x, ok := wObj["x"].(float64)
			if !ok {
				return fmt.Errorf("invalid value for 'x' key in circle. Found: '%+v'", wObj["x"])
			}
			circle.X = int(x)

			y, ok := wObj["y"].(float64)
			if !ok {
				return fmt.Errorf("invalid value for 'y' key in circle. Found: '%+v'", wObj["y"])
			}
			circle.Y = int(y)

			strokeHue, ok := wObj["stroke_hue"].(float64)
			if !ok {
				return fmt.Errorf("invalid value for 'stroke_hue' key in circle. Found: '%+v'", wObj["stroke_hue"])
			}
			circle.StrokeHue = int(strokeHue)

			strokeWidth, ok := wObj["stroke_width"].(float64)
			if !ok {
				return fmt.Errorf("invalid value for 'stroke_width' key in circle. Found: '%+v'", wObj["stroke_width"])
			}
			circle.StrokeWidth = int(strokeWidth)

			fillHue, ok := wObj["fill_hue"].(float64)
			if !ok {
				return fmt.Errorf("invalid value for 'fill_hue' key in circle. Found: '%+v'", wObj["fill_hue"])
			}
			circle.FillHue = int(fillHue)

			radius, ok := wObj["radius"].(float64)
			if !ok {
				return fmt.Errorf("invalid value for 'radius' key in circle. Found: '%+v'", wObj["radius"])
			}
			circle.Radius = int(radius)
			widget = &circle
		case "text":
			w := WidgetText{}
			x, ok := wObj["x"].(float64)
			if !ok {
				return fmt.Errorf("invalid value for 'x' key in text. Found: '%+v'", wObj["x"])
			}
			w.X = int(x)

			y, ok := wObj["y"].(float64)
			if !ok {
				return fmt.Errorf("invalid value for 'y' key in text. Found: '%+v'", wObj["y"])
			}
			w.Y = int(y)

			fontFamily, ok := wObj["font_family"].(string)
			if !ok {
				return fmt.Errorf("invalid value for 'font_family' key in text. Found: '%+v'", wObj["font_family"])
			}
			w.FontFamily = fontFamily

			fontFill, ok := wObj["font_fill"].(string)
			if !ok {
				return fmt.Errorf("invalid value for 'font_fill' key in text. Found: '%+v'", wObj["font_fill"])
			}
			w.FontFill = fontFill

			fontSize, ok := wObj["font_size"].(float64)
			if !ok {
				return fmt.Errorf("invalid value for 'font_size' key in text. Found: '%+v'", wObj["font_size"])
			}
			w.FontSize = int(fontSize)

			text, ok := wObj["text"].(string)
			if !ok {
				return fmt.Errorf("invalid value for 'text' key in text. Found: '%+v'", wObj["text"])
			}
			w.Text = text

			widget = &w
		case "countdown":
			w := WidgetCountdown{}
			x, ok := wObj["x"].(float64)
			if !ok {
				return fmt.Errorf("invalid value for 'x' key in countdown. Found: '%+v'", wObj["x"])
			}
			w.X = int(x)

			y, ok := wObj["y"].(float64)
			if !ok {
				return fmt.Errorf("invalid value for 'y' key in countdown. Found: '%+v'", wObj["y"])
			}
			w.Y = int(y)

			fontFamily, ok := wObj["font_family"].(string)
			if !ok {
				return fmt.Errorf("invalid value for 'font_family' key in countdown. Found: '%+v'", wObj["font_family"])
			}
			w.FontFamily = fontFamily

			fontFill, ok := wObj["font_fill"].(string)
			if !ok {
				return fmt.Errorf("invalid value for 'font_fill' key in countdown. Found: '%+v'", wObj["font_fill"])
			}
			w.FontFill = fontFill

			fontSize, ok := wObj["font_size"].(float64)
			if !ok {
				return fmt.Errorf("invalid value for 'font_size' key in countdown. Found: '%+v'", wObj["font_size"])
			}
			w.FontSize = int(fontSize)

			doneFontFill, ok := wObj["done_font_fill"].(string)
			if !ok {
				return fmt.Errorf("invalid value for 'done_font_fill' key in countdown. Found: '%+v'", wObj["done_font_fill"])
			}
			w.DoneFontFill = doneFontFill

			remaining, ok := wObj["remaining"].(float64)
			if !ok {
				return fmt.Errorf("invalid value for 'remaining' key in countdown. Found: '%+v'", wObj["remaining"])
			}
			w.Remaining = time.Duration(int(remaining)) * time.Second

			widget = &w
		case "todolist":
			w := WidgetTodoList{}
			x, ok := wObj["x"].(float64)
			if !ok {
				return fmt.Errorf("invalid value for 'x' key in todolist. Found: '%+v'", wObj["x"])
			}
			w.X = int(x)

			y, ok := wObj["y"].(float64)
			if !ok {
				return fmt.Errorf("invalid value for 'y' key in todolist. Found: '%+v'", wObj["y"])
			}
			w.Y = int(y)

			fontFamily, ok := wObj["font_family"].(string)
			if !ok {
				return fmt.Errorf("invalid value for 'font_family' key in todolist. Found: '%+v'", wObj["font_family"])
			}
			w.FontFamily = fontFamily

			fontFill, ok := wObj["font_fill"].(string)
			if !ok {
				return fmt.Errorf("invalid value for 'font_fill' key in todolist. Found: '%+v'", wObj["font_fill"])
			}
			w.FontFill = fontFill

			doneFontFill, ok := wObj["done_font_fill"].(string)
			if !ok {
				return fmt.Errorf("invalid value for 'done_font_fill' key in todolist. Found: '%+v'", wObj["done_font_fill"])
			}
			w.DoneFontFill = doneFontFill

			fontSize, ok := wObj["font_size"].(float64)
			if !ok {
				return fmt.Errorf("invalid value for 'font_size' key in todolist. Found: '%+v'", wObj["font_size"])
			}
			w.FontSize = int(fontSize)

			width, ok := wObj["width"].(float64)
			if !ok {
				return fmt.Errorf("invalid value for 'width' key in todolist. Found: '%+v'", wObj["width"])
			}
			w.Width = int(width)

			itemsArr, ok := wObj["items"].([]any)
			if !ok {
				return fmt.Errorf("invalid value for 'items' key in todolist. Found: '%+v'", wObj["items"])
			}

			items := []TodoListItem{}
			for _, it := range itemsArr {
				itemObj, ok := it.(map[string]any)
				if !ok {
					return fmt.Errorf("invalid element in todolist.items. Found: '%+v'", it)
				}

				item := TodoListItem{}
				desc, ok := itemObj["description"].(string)
				if !ok {
					return fmt.Errorf("invalid value for 'description' key in todolist.item. Found: '%+v'", itemObj["description"])
				}
				item.Description = desc

				done, ok := itemObj["done"].(bool)
				if !ok {
					return fmt.Errorf("invalid value for 'done' key in todolist.item. Found: '%+v'", itemObj["done"])
				}
				item.Done = done

				items = append(items, item)
			}
			w.Items = items

			widget = &w
		default:
			return fmt.Errorf("invalid widget type found: '%s'", wType)
		}

		c.Widgets[k] = widget
	}

	return nil
}

type configWidget struct {
	Type   string `json:"type"`
	Widget Widget `json:"widget"`
}

func (c Config) MarshalJSON() ([]byte, error) {
	configObj := map[string]any{}
	configWidgets := map[string]configWidget{}

	for k, w := range c.Widgets {
		cw := configWidget{
			Widget: w,
		}

		switch w.(type) {
		case *WidgetCircle:
			cw.Type = "circle"
		case *WidgetText:
			cw.Type = "text"
		case *WidgetTodoList:
			cw.Type = "todolist"
		case *WidgetCountdown:
			cw.Type = "countdown"
		}

		configWidgets[k] = cw
	}

	configObj["widgets"] = configWidgets

	return json.Marshal(configObj)
}
