package shandler

import (
	"fmt"
	"testing"

	"github.com/lucasb-eyer/go-colorful"
)

func TestTheme(t *testing.T) {
	pink, _ := colorful.Hex("#c65f9f")
	blue, _ := colorful.Hex("#2f81ff")
	orange, _ := colorful.Hex("#e16c00")
	red, _ := colorful.Hex("#ff000a")
	lightEmpty := colorful.Color{R: 1, G: 1, B: 1}
	darkEmpty := colorful.Color{}

	tests := []struct {
		lightF    colorful.Color
		darkF     colorful.Color
		lightB    colorful.Color
		darkB     colorful.Color
		message   string
		bold      bool
		faint     bool
		italic    bool
		underline bool
		overline  bool
		blink     bool
		reverse   bool
		crossOut  bool
	}{
		{
			pink, blue, lightEmpty, darkEmpty,
			"Foreground",
			false, false, false, false,
			false, false, false, false,
		},
		{
			pink, blue, lightEmpty, darkEmpty,
			"Foreground Bold",
			true, false, false, false,
			false, false, false, false,
		},
		{
			orange, red, lightEmpty, darkEmpty,
			"Foreground Bold Italic Underline",
			true, false, true, true,
			false, false, false, false,
		},
		{
			orange, red, lightEmpty, darkEmpty,
			"Foreground Bold Italic",
			true, false, true, false,
			false, false, false, false,
		},
		{
			orange, red, lightEmpty, darkEmpty,
			"Foreground Bold Italic Underline Overline",
			true, false, true, true,
			true, false, false, false,
		},
	}

	for _, s := range tests {
		theme := NewTheme().Foreground(s.lightF, s.darkF) // .Background(s.lightB, s.darkB)
		theme.Bold(s.bold).Faint(s.faint).Italic(s.italic).Underline(s.underline)
		theme.Overline(s.overline).Blink(s.blink).Reverse(s.reverse).Reverse(s.reverse)
		message := theme.Format().Render(s.message)
		println(message, fmt.Sprintf("%q", message))
	}
}
