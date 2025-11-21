package ui

import (
	"fmt"
	"strings"

	"omarchy-monitor/core"

	"github.com/rivo/tview"
)

// TempComponent displays system temperatures
type TempComponent struct {
	Flex     *tview.Flex
	TempView *tview.TextView
}

// NewTempComponent creates the temperature view
func NewTempComponent() *TempComponent {
	temp := tview.NewTextView().SetDynamicColors(true)
	temp.SetBorder(true).SetTitle(" System Temperatures ")

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(temp, 0, 1, false)

	return &TempComponent{
		Flex:     flex,
		TempView: temp,
	}
}

// ApplyTheme applies the current theme
func (t *TempComponent) ApplyTheme() {
	t.TempView.SetBorderColor(CurrentTheme.Border)
	t.TempView.SetTitleColor(CurrentTheme.HeaderTitle)
}

// Update refreshes the temperature display
func (t *TempComponent) Update(temps []core.TempStat) {
	var sb strings.Builder

	// Helper to get color tag
	cLow := fmt.Sprintf("[#%06x]", CurrentTheme.LowUsage.Hex())
	cMed := fmt.Sprintf("[#%06x]", CurrentTheme.MedUsage.Hex())
	cHigh := fmt.Sprintf("[#%06x]", CurrentTheme.HighUsage.Hex())
	cFore := fmt.Sprintf("[#%06x]", CurrentTheme.Foreground.Hex())
	cDim := "[darkgray]"
	cTitle := fmt.Sprintf("[#%06x::b]", CurrentTheme.HeaderTitle.Hex())
	cLabel := fmt.Sprintf("[#%06x]", CurrentTheme.HeaderValue.Hex())

	sb.WriteString(fmt.Sprintf("\n %sSYSTEM TEMPERATURES%s\n", cTitle, cFore))
	sb.WriteString(" " + strings.Repeat("─", 78) + "\n\n")

	if len(temps) == 0 {
		sb.WriteString("\n\n")
		sb.WriteString(fmt.Sprintf("  %sNo temperature sensors detected%s\n\n", cMed, cFore))
		sb.WriteString(fmt.Sprintf("  %sTemperature monitoring requires hardware sensor support.\n", cDim))
		sb.WriteString(fmt.Sprintf("  Install lm-sensors and run 'sensors-detect' to enable.%s\n\n", cFore))
	} else {
		// Find max temp for reference
		maxTemp := 0.0
		for _, temp := range temps {
			if temp.Temperature > maxTemp {
				maxTemp = temp.Temperature
			}
		}

		for _, temp := range temps {
			// Determine color and status
			color := cLow
			status := "NORMAL"

			if temp.Temperature > 85 {
				color = cHigh
				status = "CRITICAL"
			} else if temp.Temperature > 75 {
				color = cMed
				status = "HIGH"
			} else if temp.Temperature > 60 {
				color = cLabel // Use Label color (often cyan) for Warm
				status = "WARM"
			}

			// Create enhanced bar with gradient (max 120°C for safety margin)
			barWidth := 50
			maxScale := 120.0
			filled := int(temp.Temperature / maxScale * float64(barWidth))
			if filled > barWidth {
				filled = barWidth
			}

			sb.WriteString(fmt.Sprintf(" %-35s [%s::b]%6.1f°C%s │ ",
				temp.Key, color, temp.Temperature, cFore))

			// Gradient bar
			for i := 0; i < barWidth; i++ {
				if i < filled {
					tempAtPos := (float64(i) / float64(barWidth)) * maxScale
					if tempAtPos > 85 {
						sb.WriteString(fmt.Sprintf("%s█", cHigh))
					} else if tempAtPos > 75 {
						sb.WriteString(fmt.Sprintf("%s█", cMed))
					} else if tempAtPos > 60 {
						sb.WriteString(fmt.Sprintf("%s█", cLabel))
					} else {
						sb.WriteString(fmt.Sprintf("%s█", cLow))
					}
				} else {
					sb.WriteString(fmt.Sprintf("%s░", cDim))
				}
			}
			sb.WriteString(fmt.Sprintf("%s %s%s%s\n", cFore, color, status, cFore))
		}

		sb.WriteString("\n " + strings.Repeat("─", 78) + "\n")
		sb.WriteString(fmt.Sprintf(" Highest Temperature: %s%.1f°C%s  │  Scale: 0°C ─────────── 120°C\n", cHigh, maxTemp, cFore))
		sb.WriteString(fmt.Sprintf(" %sLegend: %sNormal%s | %sWarm%s | %sHigh%s | %sCritical%s\n",
			cDim, cLow, cFore, cLabel, cFore, cMed, cFore, cHigh, cFore))
	}

	t.TempView.SetText(sb.String())
	t.ApplyTheme()
}
