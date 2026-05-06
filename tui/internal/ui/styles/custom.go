package styles

import "charm.land/lipgloss/v2"

// prob will extract these vars into seperate files if they grow to be bigger
var Mono = Theme{
	AccentBright: lipgloss.Color("#d0d0d0"),
	AccentMid: 	  lipgloss.Color("#888888"),
	AccentDim: 	  lipgloss.Color("#555555"),
}

var Warm = Theme{
	AccentBright: lipgloss.Color("#e8c97a"),
	AccentMid: 	  lipgloss.Color("#9a7e3a"),
	AccentDim: 	  lipgloss.Color("#3a2e18"),
}

var Cold = Theme{
	AccentBright: lipgloss.Color("#7ae8e0"),
	AccentMid: 	  lipgloss.Color("#3a8a84"),
	AccentDim: 	  lipgloss.Color("#1e3830"),
}