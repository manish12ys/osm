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
		{
			Name:        "Dracula",
			Background:  tcell.NewHexColor(0x282a36),
			Foreground:  tcell.NewHexColor(0xf8f8f2),
			Border:      tcell.NewHexColor(0xbd93f9), // Purple
			HeaderTitle: tcell.NewHexColor(0x50fa7b), // Green
			HeaderValue: tcell.NewHexColor(0xf8f8f2),
			TableHead:   tcell.NewHexColor(0x8be9fd), // Cyan
			RowAlt:      tcell.NewHexColor(0x44475a),
			HighUsage:   tcell.NewHexColor(0xff5555), // Red
			MedUsage:    tcell.NewHexColor(0xffb86c), // Orange
			LowUsage:    tcell.NewHexColor(0x50fa7b), // Green
		},
		{
			Name:        "Solarized",
			Background:  tcell.NewHexColor(0x002b36),
			Foreground:  tcell.NewHexColor(0x839496),
			Border:      tcell.NewHexColor(0x2aa198), // Cyan
			HeaderTitle: tcell.NewHexColor(0x859900), // Green
			HeaderValue: tcell.NewHexColor(0x93a1a1),
			TableHead:   tcell.NewHexColor(0xb58900), // Yellow
			RowAlt:      tcell.NewHexColor(0x073642),
			HighUsage:   tcell.NewHexColor(0xdc322f), // Red
			MedUsage:    tcell.NewHexColor(0xcb4b16), // Orange
			LowUsage:    tcell.NewHexColor(0x859900), // Green
		},
		{
			Name:        "Nord",
			Background:  tcell.NewHexColor(0x2e3440),
			Foreground:  tcell.NewHexColor(0xd8dee9),
			Border:      tcell.NewHexColor(0x88c0d0), // Frost
			HeaderTitle: tcell.NewHexColor(0xa3be8c), // Green
			HeaderValue: tcell.NewHexColor(0xeceff4),
			TableHead:   tcell.NewHexColor(0x81a1c1), // Blue
			RowAlt:      tcell.NewHexColor(0x3b4252),
			HighUsage:   tcell.NewHexColor(0xbf616a), // Red
			MedUsage:    tcell.NewHexColor(0xd08770), // Orange
			LowUsage:    tcell.NewHexColor(0xa3be8c), // Green
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

// SetTheme sets the current theme by name
func SetTheme(name string) {
	for i, t := range Themes {
		if t.Name == name {
			CurrentTheme = t
			themeIndex = i
			return
		}
	}
}
