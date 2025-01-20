package goverly

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
