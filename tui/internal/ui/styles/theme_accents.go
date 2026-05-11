package styles

import "charm.land/lipgloss/v2"

// prob will extract these vars into seperate files if they grow to be bigger
var Mono = Theme{
	AccentBright:   lipgloss.Color("#d0d0d0"),
	AccentMid:      lipgloss.Color("#888888"),
	AccentDim:      lipgloss.Color("#555555"),
	BorderFocused:  lipgloss.Color("#3a3a3a"),
	GradientDim:    lipgloss.Color("#4a6fa5"),
	GradientMid:    lipgloss.Color("#8ab0d0"),
	GradientBright: lipgloss.Color("#d0e8f5"),
}

var Warm = Theme{
	AccentBright:   lipgloss.Color("#e8c97a"),
	AccentMid:      lipgloss.Color("#9a7e3a"),
	AccentDim:      lipgloss.Color("#3a2e18"),
	BorderFocused:  lipgloss.Color("#5a4828"),
	GradientDim:    lipgloss.Color("#c47c2a"),
	GradientMid:    lipgloss.Color("#e8a84a"),
	GradientBright: lipgloss.Color("#f5dfa0"),
}

var Cold = Theme{
	AccentBright:   lipgloss.Color("#7ae8e0"),
	AccentMid:      lipgloss.Color("#3a8a84"),
	AccentDim:      lipgloss.Color("#1e3830"),
	BorderFocused:  lipgloss.Color("#1e5450"),
	GradientDim:    lipgloss.Color("#2a7a8a"),
	GradientMid:    lipgloss.Color("#4ab8c8"),
	GradientBright: lipgloss.Color("#a0eae8"),
}
