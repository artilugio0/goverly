package main

import (
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"strconv"
	"strings"
	"syscall/js"
	"time"
)

const border int = 0
const svgWidth int = 1920 - 2*border
const svgHeight int = 1080 - 2*border
const rate time.Duration = 10 * time.Millisecond

func main() {
	fmt.Println("starting")

	document := js.Global().Get("document")
	body := document.Get("body")

	svg := document.Call("createElementNS", "http://www.w3.org/2000/svg", "svg")
	svg.Call("setAttribute", "height", strconv.Itoa(svgHeight))
	svg.Call("setAttribute", "width", strconv.Itoa(svgWidth))

	svgStyle := svg.Get("style")
	svgStyle.Set("border", strconv.Itoa(border)+"px solid red")
	body.Call("appendChild", svg)

	widgets := []Widget{
		newCircle(),
		newText("Hoy: Browser overlay para OBS usando Golang + WebAssembly"),
		//newCountdown(30 * time.Minute),
		//newText("Break de 3 minutos, ya vuelvo!!!"),
		//newCountdown(3 * time.Minute),
		newTodoList([]TodoListItem{
			{"Remover IDs random de widgets", true},
			{"Fix longitud de todo list item", true},
			{"Crear widget de timer", true},
			{"Hot reload basico", true},
			{"Sacar coordenadas de widgets hardcodeadas", false},
		}),
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
	for !update {
		select {
		case update = <-updateChan:
			break
		default:
			break
		}

		for _, w := range widgets {
			w(svg)
		}

		time.Sleep(rate)
	}

	svg.Call("remove")

	execUpdate := js.Global().Get("runGo")
	execUpdate.Call("call")
	fmt.Println("exiting")
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

type Widget func(svg js.Value)

func newId() string {
	id := rand.Uint64()
	return strconv.FormatUint(id, 16)
}

func newCircle() Widget {
	const cRadius int = 10
	x, y := 0, 0

	document := js.Global().Get("document")

	circle := document.Call("createElementNS", "http://www.w3.org/2000/svg", "circle")
	circle.Call("setAttribute", "r", cRadius)
	circle.Call("setAttribute", "r", cRadius)
	circle.Call("setAttribute", "stroke", fmt.Sprintf("hsl(%d, 100%%, 50%%)", rand.Int()%361))
	circle.Call("setAttribute", "stroke-width", "4")
	circle.Call("setAttribute", "fill", fmt.Sprintf("hsl(%d, 100%%, 50%%)", rand.Int()%361))

	x, y = svgWidth/2, svgHeight/2

	appended := false
	return func(svg js.Value) {
		if !appended {
			svg.Call("appendChild", circle)
			appended = true
		}

		x = (x + 11) % svgWidth
		y = (y + 11) % svgHeight

		circle.Call("setAttribute", "cx", x)
		circle.Call("setAttribute", "cy", y)
	}
}

func newText(text string) Widget {
	const textSize int = 30
	const textAreaHeight int = 109
	const marginLeft = 50

	document := js.Global().Get("document")

	svgtext := document.Call("createElementNS", "http://www.w3.org/2000/svg", "text")
	svgtext.Call("setAttribute", "font-family", "Courier New")
	svgtext.Call("setAttribute", "fill", "white")
	svgtext.Call("setAttribute", "font-size", strconv.Itoa(textSize))
	svgtext.Call("setAttribute", "x", marginLeft)
	svgtext.Call("setAttribute", "y", svgHeight-(textAreaHeight-textSize)/2)

	textnode := document.Call("createTextNode", text)
	svgtext.Call("appendChild", textnode)
	appended := false

	return func(svg js.Value) {
		if !appended {
			svg.Call("appendChild", svgtext)
			appended = true

		}
	}
}

func newTodoList(items []TodoListItem) Widget {
	const titleTextSize int = 22
	const textSize int = 18
	const marginLeft int = 10
	const itemMarginBottom int = 5
	const width int = 244 - marginLeft
	const x int = 1920 - width
	const y int = 1080 - 850 /* cam */ + 40 /* countdown */

	document := js.Global().Get("document")

	g := document.Call("createElementNS", "http://www.w3.org/2000/svg", "g")

	title := document.Call("createElementNS", "http://www.w3.org/2000/svg", "text")
	title.Call("setAttribute", "font-family", "Courier New")
	title.Call("setAttribute", "fill", "white")
	title.Call("setAttribute", "font-size", titleTextSize)
	title.Call("setAttribute", "x", 0)
	title.Call("setAttribute", "y", 0)
	textnode := document.Call("createTextNode", "To Do:")
	title.Call("appendChild", textnode)
	g.Call("appendChild", title)

	itemY := titleTextSize + itemMarginBottom
	for _, it := range items {
		listItemText := it.Description
		textFields := strings.Fields(listItemText)

		curWidth := 0
		lines := []string{}
		index := 0
		for i, f := range textFields {
			if curWidth+len(f)*(textSize*80/100) > width-2 {
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

			curWidth += len(f) * (textSize * 80 / 100)
		}

		if len(lines) == 0 {
			lines = append(lines, "- "+it.Description)
		} else {
			lines = append(lines, strings.Join(textFields[index:], " "))
		}

		color := "white"
		if it.Done {
			color = "lime"
		}

		for _, line := range lines {

			itemText := document.Call("createElementNS", "http://www.w3.org/2000/svg", "text")
			itemText.Call("setAttribute", "font-family", "Courier New")
			itemText.Call("setAttribute", "fill", color)
			itemText.Call("setAttribute", "font-size", strconv.Itoa(textSize))
			itemText.Call("setAttribute", "x", 0)
			itemText.Call("setAttribute", "y", itemY)

			textnode := document.Call("createTextNode", line)

			itemText.Call("appendChild", textnode)
			g.Call("appendChild", itemText)

			itemY += textSize + itemMarginBottom
		}
	}

	g.Call("setAttribute", "transform", fmt.Sprintf("translate(%d, %d)", x, y))

	appended := false
	return func(svg js.Value) {
		if !appended {
			svg.Call("appendChild", g)
			appended = true
		}
	}
}

type TodoListItem struct {
	Description string
	Done        bool
}

func newCountdown(t time.Duration) Widget {
	const textSize int = 40
	const marginLeft int = 60
	const width int = 244 - marginLeft
	const x int = 1920 - width
	const y int = 1080 - 850

	document := js.Global().Get("document")

	g := document.Call("createElementNS", "http://www.w3.org/2000/svg", "g")

	timeText := document.Call("createElementNS", "http://www.w3.org/2000/svg", "text")
	timeText.Call("setAttribute", "font-family", "Courier New")
	timeText.Call("setAttribute", "fill", "white")
	timeText.Call("setAttribute", "font-size", textSize)
	timeText.Call("setAttribute", "x", 0)
	timeText.Call("setAttribute", "y", 0)
	textnode := document.Call("createTextNode", "00:00")
	timeText.Call("appendChild", textnode)
	g.Call("appendChild", timeText)
	g.Call("setAttribute", "transform", fmt.Sprintf("translate(%d, %d)", x, y))

	appended := false
	angleMillis := 0
	angles := []int{0, 5, 0, -5}
	angleIndex := 0
	passed := 0 * time.Second
	return func(svg js.Value) {
		if !appended {
			svg.Call("appendChild", g)
			appended = true
		}
		t -= rate
		if t > 0 {
			mins := int(t.Minutes())
			secs := int(t.Seconds()) % 60
			text := fmt.Sprintf("%02d:%02d", mins, secs)
			textnode.Set("nodeValue", text)
			return
		}

		// time is up

		passed += rate
		timeText.Call("setAttribute", "fill", "red")
		mins := int(passed.Minutes())
		secs := int(passed.Seconds()) % 60
		text := fmt.Sprintf("%02d:%02d", mins, secs)
		textnode.Set("nodeValue", text)

		bbox := timeText.Call("getBBox")
		timeTextWidth := bbox.Get("width").Int()

		angleMillis += int(rate.Milliseconds())
		if (angleMillis/1000)%2 == 0 {
			if angleMillis%50 == 0 {
				angleIndex = (angleIndex + 1) % len(angles)
			}

			transform := fmt.Sprintf("rotate(%d, %d, 0)", angles[angleIndex], timeTextWidth/2)
			timeText.Call("setAttribute", "transform", transform)
			return
		}

		transform := fmt.Sprintf("rotate(%d, %d, 0)", 0, timeTextWidth/2)
		timeText.Call("setAttribute", "transform", transform)
	}
}
