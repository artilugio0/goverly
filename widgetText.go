package goverly

type WidgetText struct {
	FontFamily string `json:"font_family"`
	FontFill   string `json:"font_fill"`
	FontSize   int    `json:"font_size"`
	Text       string `json:"text"`
	X          int    `json:"x"`
	Y          int    `json:"y"`
}

func NewText(text string, textSize, x, y int) *WidgetText {
	return &WidgetText{
		Text:       text,
		X:          x,
		Y:          y,
		FontFamily: "Courier New",
		FontFill:   "white",
		FontSize:   textSize,
	}
}

func (wt *WidgetText) Type() string {
	return "text"
}
