package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"syscall/js"
	"time"
)

const border int = 0
const svgWidth int = 1920 - 2*border
const svgHeight int = 1080 - 2*border
const rate time.Duration = 10 * time.Millisecond

func main() {
	config, err := getConfig()
	if err != nil {
		panic(err)
	}

	document := js.Global().Get("document")
	body := document.Get("body")

	svg := document.Call("createElementNS", "http://www.w3.org/2000/svg", "svg")
	svg.Call("setAttribute", "height", strconv.Itoa(svgHeight))
	svg.Call("setAttribute", "width", strconv.Itoa(svgWidth))

	svgStyle := svg.Get("style")
	svgStyle.Set("border", strconv.Itoa(border)+"px solid red")
	body.Call("appendChild", svg)

	appUpdateAvailableChan := make(chan bool)
	go func() {
		lastUpdate, err := getLastUpdate()
		if err != nil {
			panic(err)
		}

		for {
			time.Sleep(1 * time.Second)

			newUpdate, err := getLastUpdate()
			if err != nil {
				fmt.Printf("error while trying to get last update: %v\n", err)
				continue
			}

			if newUpdate > lastUpdate {
				lastUpdate = newUpdate
				go func() {
					appUpdateAvailableChan <- true
				}()
				return
			}
		}
	}()

	configUpdateChan := make(chan *Config)
	go func() {
		for {
			time.Sleep(1 * time.Second)

			newConfig, err := getConfig()
			if err != nil {
				fmt.Printf("error while trying to get new config: %v\n", err)
				continue
			}

			configUpdateChan <- newConfig
		}
	}()

	appUpdateAvailable := false

	widgets := config.Widgets
	loadAppState(widgets)
	for !appUpdateAvailable {
		select {
		case appUpdateAvailable = <-appUpdateAvailableChan:
			break
		case newConfig := <-configUpdateChan:
			updateWidgets(widgets, newConfig.Widgets)
		default:
			break
		}

		for _, w := range widgets {
			w.Update(svg)
		}

		time.Sleep(rate)
	}
	saveAppState(widgets)

	svg.Call("remove")

	execUpdate := js.Global().Get("runGo")
	execUpdate.Call("call")
}

func getLastUpdate() (int64, error) {
	resp, err := http.Get("/last-update")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	lu, err := strconv.ParseInt(string(bodyBytes), 10, 64)
	if err != nil {
		return 0, err
	}

	return lu, nil
}

func getConfig() (*Config, error) {
	resp, err := http.Get("/config")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	config := Config{}
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func saveAppState(widgets map[string]Widget) {
	appState := map[string]any{}

	for k, w := range widgets {
		if w, ok := w.(WidgetStateful); ok {
			wState := w.SaveState()
			appState[k] = wState
			continue
		}
	}

	js.Global().Set("appState", js.ValueOf(appState))
}

func loadAppState(widgets map[string]Widget) {
	jsAppState := js.Global().Get("appState")
	if jsAppState.IsNull() || jsAppState.IsUndefined() {
		return
	}

	keys := js.Global().Get("Object").Call("keys", jsAppState)
	for i := range keys.Length() {
		widgetKey := keys.Index(i).String()
		jsState := jsAppState.Get(widgetKey)
		if jsState.IsNull() {
			continue
		}

		w, ok := widgets[widgetKey]
		if !ok {
			continue
		}

		if w, ok := w.(WidgetStateful); ok {
			w.LoadState(jsState)
			fmt.Printf("widget '%s' loaded\n", widgetKey)
			continue
		}

		panic(fmt.Sprintf("read widget state of non stateful widget with key '%s'", widgetKey))
	}
}

func updateWidgets(widgets, mods map[string]Widget) error {
	// TODO: remove widgets present in widgets that have different type in mods

	for k, w := range widgets {
		if m, ok := mods[k]; !ok || reflect.TypeOf(w) != reflect.TypeOf(m) {
			switch v := w.(type) {
			case *WidgetCircle:
				v.element.Call("remove")
			case *WidgetText:
				v.element.Call("remove")
			case *WidgetCountdown:
				v.element.Call("remove")
			case *WidgetTodoList:
				v.element.Call("remove")
			}
			delete(widgets, k)
		}
	}

	for k, m := range mods {
		if w, ok := widgets[k]; ok {
			switch v := w.(type) {
			case *WidgetCircle:
				m, ok := mods[k].(*WidgetCircle)
				if !ok {
					return fmt.Errorf("invalid widget type received in update. Expected circle at key %s", k)
				}

				v.StrokeHue = m.StrokeHue
				v.StrokeWidth = m.StrokeWidth
				v.FillHue = m.FillHue
				v.Radius = m.Radius

			case *WidgetText:
				m, ok := mods[k].(*WidgetText)
				if !ok {
					return fmt.Errorf("invalid widget type received in update. Expected text at key %s", k)
				}

				v.Text = m.Text
				v.X = m.X
				v.Y = m.Y
				v.FontFamily = m.FontFamily
				v.FontFill = m.FontFill
				v.FontSize = m.FontSize

			case *WidgetCountdown:
				m, ok := mods[k].(*WidgetCountdown)
				if !ok {
					return fmt.Errorf("invalid widget type received in update. Expected countdown at key %s", k)
				}

				v.X = m.X
				v.Y = m.Y
				v.FontFamily = m.FontFamily
				v.FontFill = m.FontFill
				v.DoneFontFill = m.DoneFontFill
				v.FontSize = m.FontSize

			case *WidgetTodoList:
				m, ok := mods[k].(*WidgetTodoList)
				if !ok {
					return fmt.Errorf("invalid widget type received in update. Expected todolist at key %s", k)
				}

				v.X = m.X
				v.Y = m.Y
				v.FontFamily = m.FontFamily
				v.FontFill = m.FontFill
				v.Width = m.Width
				v.FontSize = m.FontSize
				v.DoneFontFill = m.DoneFontFill
				v.Items = m.Items
			}
			continue
		}

		fmt.Printf("new widget added! %+v\n", m)
		// if not present in widgets, add it
		widgets[k] = m
	}

	return nil
}
