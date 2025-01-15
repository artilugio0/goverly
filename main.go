package main

import (
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
	document := js.Global().Get("document")
	body := document.Get("body")

	svg := document.Call("createElementNS", "http://www.w3.org/2000/svg", "svg")
	svg.Call("setAttribute", "height", strconv.Itoa(svgHeight))
	svg.Call("setAttribute", "width", strconv.Itoa(svgWidth))

	svgStyle := svg.Get("style")
	svgStyle.Set("border", strconv.Itoa(border)+"px solid red")
	body.Call("appendChild", svg)

	widgets := []Widget{
		NewCircle(10),
		/*
			NewText("Break de 3 minutos, ya vuelvo!!!", 30, 50, svgHeight-60),
			NewCountdown(40, svgWidth-180, svgHeight-850, 3*time.Minute),
		*/
		NewText("Hoy: Browser overlay para OBS usando Golang + WebAssembly", 30, 50, svgHeight-60),
		NewCountdown(40, svgWidth-180, svgHeight-850, 3*time.Second),
		NewTodoList(
			18, 230, svgWidth-235, svgHeight-800,
			[]TodoListItem{
				{"Conservar estado cuando se hace un hot reload", true},
				{"Extraer definiciones de widgets hardcodeadas", false},
			},
		),
	}

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

	loadAppState(widgets)
	for !update {
		select {
		case update = <-updateChan:
			break
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
