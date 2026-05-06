package styles

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

type Theme struct {
	AccentBright color.Color
	AccentMid 	 color.Color
	AccentDim	 color.Color
}

var (
	Background  = lipgloss.Color("#0d0d0d")
	Surface	    = lipgloss.Color("#141414")
	BorderFocus = lipgloss.Color("#2a2a2a")
	TextPrimary = lipgloss.Color("#d0d0d0")
	TextMuted   = lipgloss.Color("#555555")
	TextDim		= lipgloss.Color("#333333")
)