package goverly

import (
	"fmt"
	"strconv"
)

type WidgetTodoList struct {
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

func (wtl *WidgetTodoList) Type() string {
	return "todolist"
}

func (wtl *WidgetTodoList) ApplyCustomConfig(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("todolist widgets needs exactly 2 arguments for its custom config")
	}

	op := args[0]

	switch op {
	case "add":
		wtl.Items = append(wtl.Items, TodoListItem{
			Description: args[1],
			Done:        false,
		})
	case "del", "delete":
		elem, err := strconv.Atoi(args[1])
		if err != nil {
			return err
		}
		if elem >= len(wtl.Items) {
			return fmt.Errorf("element index out of bounds")
		}
		wtl.Items = append(wtl.Items[:elem], wtl.Items[elem+1:]...)
	case "toggle":
		elem, err := strconv.Atoi(args[1])
		if err != nil {
			return err
		}
		if elem >= len(wtl.Items) {
			return fmt.Errorf("element index out of bounds")
		}
		wtl.Items[elem].Done = !wtl.Items[elem].Done
	default:
		return fmt.Errorf("invalid operation '%s'", op)

	}

	return nil
}
