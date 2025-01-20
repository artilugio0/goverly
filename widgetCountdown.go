package goverly

import (
	"time"
)

type WidgetCountdown struct {
	FontFamily   string `json:"font_family"`
	FontFill     string `json:"font_fill"`
	DoneFontFill string `json:"done_font_fill"`
	FontSize     int    `json:"font_size"`
	EndTime      int64  `json:"end_time"`
	X            int    `json:"x"`
	Y            int    `json:"y"`
}

func NewCountdown(textSize, x, y int, t time.Duration) *WidgetCountdown {
	return &WidgetCountdown{
		FontFamily:   "Courier New",
		FontFill:     "white",
		DoneFontFill: "lime",
		FontSize:     textSize,
		X:            x,
		Y:            y,
	}
}
func (wc *WidgetCountdown) Type() string {
	return "countdown"
}
