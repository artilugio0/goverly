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
	resp, err := http.Get("/config")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	config := Config{}
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
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

	updateChan := make(chan bool)
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
					updateChan <- true
				}()
				return
			}
		}
	}()

	update := false

	loadAppState(config.Widgets)
	for !update {
		select {
		case update = <-updateChan:
			break
		default:
			break
		}

		for _, w := range config.Widgets {
			w.Update(svg)
		}

		time.Sleep(rate)
	}
	saveAppState(config.Widgets)

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
