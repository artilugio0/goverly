package goverly

import (
	"fmt"
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
	Stopped      bool   `json:"stopped"`
}

func NewCountdown(textSize, x, y int, t time.Duration) *WidgetCountdown {
	return &WidgetCountdown{
		FontFamily:   "Courier New",
		FontFill:     "white",
		DoneFontFill: "lime",
		FontSize:     textSize,
		X:            x,
		Y:            y,
		Stopped:      false,
	}
}

func (wc *WidgetCountdown) Type() string {
	return "countdown"
}

func (wc *WidgetCountdown) ApplyCustomConfig(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("countdown widgets needs exactly 1 argument for its custom config")
	}

	if args[0] == "toggle" {
		wc.Stopped = !wc.Stopped
		return nil
	}

	duration, err := time.ParseDuration(args[0])
	if err != nil {
		return err
	}

	wc.EndTime = time.Now().Add(duration).Unix()

	return nil
}
