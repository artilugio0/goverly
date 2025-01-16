package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func saveAppState(widgets []Widget) {
	appState := []interface{}{}

	for _, w := range widgets {
		if w, ok := w.(WidgetStateful); ok {
			wState := w.SaveState()
			appState = append(appState, wState)
			continue
		}

		appState = append(appState, nil)
	}

	js.Global().Set("appState", js.ValueOf(appState))
}

func loadAppState(widgets []Widget) {
	// here I am assuming that the state and the widgets array have the same length
	// TODO: remove this assumption

	jsAppState := js.Global().Get("appState")
	if jsAppState.IsNull() || jsAppState.IsUndefined() {
		return
	}

	if jsAppState.Length() != len(widgets) {
		return
	}

	for i := range jsAppState.Length() {
		jsState := jsAppState.Index(i)
		if jsState.IsNull() {
			continue
		}

		w := widgets[i]
		if w, ok := w.(WidgetStateful); ok {
			w.LoadState(jsState)
			fmt.Printf("widget %d loaded\n", i)
			continue
		}

		panic(fmt.Sprintf("read widget state of non stateful widget at index %d", i))
	}
}

func updateWidgets(widgets, mods []Widget) error {
	if len(widgets) != len(mods) {
		return fmt.Errorf("ERROR: unsupported feature, updates with different widgets")
	}

	for i, w := range widgets {
		switch v := w.(type) {
		case *WidgetCircle:
			m, ok := mods[i].(*WidgetCircle)
			if !ok {
				return fmt.Errorf("invalid widget type received in update. Expected circle at position %d", i)
			}

			v.StrokeHue = m.StrokeHue
			v.StrokeWidth = m.StrokeWidth
			v.FillHue = m.FillHue
			v.Radius = m.Radius

		case *WidgetText:
			m, ok := mods[i].(*WidgetText)
			if !ok {
				return fmt.Errorf("invalid widget type received in update. Expected text at position %d", i)
			}

			v.Text = m.Text
			v.X = m.X
			v.Y = m.Y
			v.FontFamily = m.FontFamily
			v.FontFill = m.FontFill
			v.FontSize = m.FontSize

		case *WidgetCountdown:
			m, ok := mods[i].(*WidgetCountdown)
			if !ok {
				return fmt.Errorf("invalid widget type received in update. Expected countdown at position %d", i)
			}

			v.X = m.X
			v.Y = m.Y
			v.FontFamily = m.FontFamily
			v.FontFill = m.FontFill
			v.DoneFontFill = m.DoneFontFill
			v.FontSize = m.FontSize

		case *WidgetTodoList:
			m, ok := mods[i].(*WidgetTodoList)
			if !ok {
				return fmt.Errorf("invalid widget type received in update. Expected todolist at position %d", i)
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
	}

	return nil
}
