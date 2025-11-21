package ui

import "github.com/gdamore/tcell/v2"

// Theme defines the color palette
type Theme struct {
	Name        string
	Background  tcell.Color
	Foreground  tcell.Color
	Border      tcell.Color
	HeaderTitle tcell.Color
	HeaderValue tcell.Color
	TableHead   tcell.Color
	RowAlt      tcell.Color
	HighUsage   tcell.Color
	MedUsage    tcell.Color
	LowUsage    tcell.Color
}

var (
	Themes = []Theme{
		{
			Name:        "Omarchy Default",
			Background:  tcell.ColorBlack,
			Foreground:  tcell.ColorWhite,
			Border:      tcell.ColorBlue,
			HeaderTitle: tcell.ColorGreen,
			HeaderValue: tcell.ColorWhite,
			TableHead:   tcell.ColorYellow,
			RowAlt:      tcell.ColorBlack,
			HighUsage:   tcell.ColorRed,
			MedUsage:    tcell.ColorOrange,
			LowUsage:    tcell.ColorGreen,
		},
		{
			Name:        "Cyberpunk",
			Background:  tcell.ColorBlack,
			Foreground:  tcell.NewHexColor(0xFF00FF), // Neon Pink
			Border:      tcell.NewHexColor(0x00FFFF), // Neon Blue/Cyan
			HeaderTitle: tcell.NewHexColor(0x00FFFF),
			HeaderValue: tcell.NewHexColor(0xFF00FF),
			TableHead:   tcell.NewHexColor(0x00FFFF),
			RowAlt:      tcell.Color16,
			HighUsage:   tcell.ColorRed,
			MedUsage:    tcell.ColorYellow,
			LowUsage:    tcell.ColorGreen,
		},
		{
			Name:        "Retro CRT",
			Background:  tcell.ColorBlack,
			Foreground:  tcell.ColorGreen,
			Border:      tcell.ColorGreen,
			HeaderTitle: tcell.ColorGreen,
			HeaderValue: tcell.ColorDarkGreen,
			TableHead:   tcell.ColorGreen,
			RowAlt:      tcell.ColorBlack,
			HighUsage:   tcell.ColorGreen, // Monochromatic feel
			MedUsage:    tcell.ColorDarkGreen,
			LowUsage:    tcell.ColorDarkGreen,
		},
	}
	CurrentTheme = Themes[0]
	themeIndex   = 0
)

// CycleTheme switches to the next available theme
func CycleTheme() Theme {
	themeIndex = (themeIndex + 1) % len(Themes)
	CurrentTheme = Themes[themeIndex]
	return CurrentTheme
}
