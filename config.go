package main

import (
	"encoding/json"
	"fmt"
	"os"
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
	type Tmp struct {
		Widgets map[string]struct {
			Type   string          `json:"type"`
			Widget json.RawMessage `json:"widget"`
		}
	}

	tmp := Tmp{}
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}

	c.Widgets = map[string]Widget{}

	for k, d := range tmp.Widgets {
		var widget Widget
		switch d.Type {
		case "circle":
			w := WidgetCircle{}
			if err := json.Unmarshal(d.Widget, &w); err != nil {
				return err
			}
			widget = &w
		case "text":
			w := WidgetText{}
			if err := json.Unmarshal(d.Widget, &w); err != nil {
				return err
			}
			widget = &w
		case "countdown":
			w := WidgetCountdown{}
			if err := json.Unmarshal(d.Widget, &w); err != nil {
				return err
			}
			widget = &w
		case "todolist":
			w := WidgetTodoList{}
			if err := json.Unmarshal(d.Widget, &w); err != nil {
				return err
			}
			widget = &w
		default:
			return fmt.Errorf("invalid widget type found: '%s'", d.Type)
		}

		c.Widgets[k] = widget
	}

	return nil
}

func (c Config) MarshalJSON() ([]byte, error) {
	type configWidget struct {
		Type   string `json:"type"`
		Widget Widget `json:"widget"`
	}

	configObj := map[string]any{}
	configWidgets := map[string]configWidget{}

	for k, w := range c.Widgets {
		cw := configWidget{
			Widget: w,
			Type:   w.Type(),
		}

		configWidgets[k] = cw
	}

	configObj["widgets"] = configWidgets

	return json.Marshal(configObj)
}
