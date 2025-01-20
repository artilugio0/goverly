//go:build js

package goverly

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

func RunOverlay() {
	widgets, err := getWidgetConfigs()
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

	configUpdateChan := make(chan map[string]Widget)
	go func() {
		for {
			time.Sleep(1 * time.Second)

			widgetsCfg, err := getWidgetConfigs()
			if err != nil {
				fmt.Printf("error while trying to get new config: %v\n", err)
				continue
			}

			configUpdateChan <- widgetsCfg
		}
	}()

	appUpdateAvailable := false

	loadAppState(widgets)

	for _, w := range widgets {
		w.Render(svg)
	}

	// TODO: check sleep and breaks, maybe use tick?
APP:
	for !appUpdateAvailable {
		renderActions := []RenderAction{}

		select {
		case appUpdateAvailable = <-appUpdateAvailableChan:
			break APP
		case newWidgetsConfig := <-configUpdateChan:
			acts := updateWidgets(widgets, newWidgetsConfig)
			renderActions = append(renderActions, acts...)
		case <-time.After(rate):
		}

		for _, w := range widgets {
			acts := w.Update(rate)
			renderActions = append(renderActions, acts...)
		}

		for _, ra := range renderActions {
			ra(svg)
		}
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

func getWidgetConfigs() (map[string]Widget, error) {
	resp, err := http.Get("/config")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	config := Config{}
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return nil, err
	}

	widgets := map[string]Widget{}
	for k, w := range config.Widgets {
		switch v := w.(type) {
		case *WidgetText:
			widgets[k] = &WidgetTextJS{WidgetText: v}
		case *WidgetCircle:
			widgets[k] = &WidgetCircleJS{WidgetCircle: v}
		case *WidgetCountdown:
			widgets[k] = &WidgetCountdownJS{WidgetCountdown: v}
		case *WidgetTodoList:
			widgets[k] = &WidgetTodoListJS{WidgetTodoList: v}
		}
	}

	return widgets, nil
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

func updateWidgets(widgets, mods map[string]Widget) []RenderAction {
	renderActions := []RenderAction{}

	for k, w := range widgets {
		if m, ok := mods[k]; !ok || reflect.TypeOf(w) != reflect.TypeOf(m) {
			// TODO: check if it can be added to render actions
			w.RemoveFromDOM()
			delete(widgets, k)
		}
	}

	for k, m := range mods {
		if w, ok := widgets[k]; ok {
			acts := w.UpdateConfig(m)
			renderActions = append(renderActions, acts...)
			continue
		}

		// if not present in widgets, add it
		widgets[k] = m
		renderActions = append(renderActions, func(svg js.Value) {
			m.Render(svg)
		})
	}

	return renderActions
}
