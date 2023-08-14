package shandler

import (
	"fmt"
	"github.com/k0kubun/pp/v3"
	"github.com/lucasb-eyer/go-colorful"
	"testing"
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
			"Foreground Bold Italic",
			true, false, true, false,
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
			"Foreground Bold Italic Underline Overline",
			true, false, true, true,
			true, false, false, false,
		},
	}

	for _, s := range tests {
		theme := NewTheme().Foreground(s.lightF, s.darkF) // .Background(s.lightB, s.darkB)
		if s.bold {
			theme.Bold()
		}
		if s.faint {
			theme.Faint()
		}
		if s.italic {
			theme.Italic()
		}
		if s.underline {
			theme.Underline()
		}
		if s.overline {
			theme.Overline()
		}
		if s.blink {
			theme.Blink()
		}
		if s.reverse {
			theme.Reverse()
		}
		if s.crossOut {
			theme.Reverse()
		}
		message := theme.Format().Render(s.message)
		println(message, fmt.Sprintf("%q", message))
	}
}

func TestColoredStruct(t *testing.T) {
	txt := pp.Sprint(struct {
		Name  string
		value string
	}{
		"should be name",
		"is a value",
	})
	println(txt)
	pp.Println("just a string")
}
