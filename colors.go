package shandler

import (
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/exp/rand"
	"time"
)

var (
	boldStyle = lipgloss.NewStyle().Bold(true)
	timeStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "#6085b9",
		Dark:  "#7d467c",
	})
	debugStyle = boldStyle.Copy().Foreground(lipgloss.AdaptiveColor{
		Light: "#4746ff",
		Dark:  "#2f81ff",
	})
	infoStyle = boldStyle.Copy().Foreground(lipgloss.AdaptiveColor{
		Light: "#009adc",
		Dark:  "#00FFD5",
	})
	warnStyle = boldStyle.Copy().Foreground(lipgloss.AdaptiveColor{
		Light: "#e16c00",
		Dark:  "#ff9c01",
	})
	errorStyle = boldStyle.Copy().Foreground(lipgloss.AdaptiveColor{
		Light: "#ff000a",
		Dark:  "#FF4F86",
	})
	prefixStyle = boldStyle.Copy().Foreground(lipgloss.AdaptiveColor{
		Light: "#579159",
		Dark:  "#008708",
	})
	callerStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "#9d86b9",
		Dark:  "#2f6982",
	})
	keyStyles = []lipgloss.Style{
		lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
			Light: "#7F7F7F",
			Dark:  "#7F7F7F",
		}),
		lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
			Light: "#f2be45",
			Dark:  "#7F7F7F",
		}),
		lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
			Light: "#622a1d",
			Dark:  "#7F7F7F",
		}),
		lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
			Light: "#88ada6",
			Dark:  "#7F7F7F",
		}),
		lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
			Light: "#30d7eb",
			Dark:  "#7F7F7F",
		}),
		lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
			Light: "#4b5cc4",
			Dark:  "#7F7F7F",
		}),
	}
)

func randomKeyColor() lipgloss.Style {
	rand.Seed(uint64(time.Now().UnixNano()))
	return keyStyles[rand.Intn(len(keyStyles))]
}
