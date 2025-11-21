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

	sb.WriteString("\n [red::b]SYSTEM TEMPERATURES[white]\n")
	sb.WriteString(" " + strings.Repeat("─", 78) + "\n\n")

	if len(temps) == 0 {
		sb.WriteString("\n\n")
		sb.WriteString("  [yellow]No temperature sensors detected[white]\n\n")
		sb.WriteString("  [darkgray]Temperature monitoring requires hardware sensor support.\n")
		sb.WriteString("  Install lm-sensors and run 'sensors-detect' to enable.[white]\n\n")
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
			color := "green"
			status := "NORMAL"

			if temp.Temperature > 85 {
				color = "red"
				status = "CRITICAL"
			} else if temp.Temperature > 75 {
				color = "yellow"
				status = "HIGH"
			} else if temp.Temperature > 60 {
				color = "cyan"
				status = "WARM"
			}

			// Create enhanced bar with gradient (max 120°C for safety margin)
			barWidth := 50
			maxScale := 120.0
			filled := int(temp.Temperature / maxScale * float64(barWidth))
			if filled > barWidth {
				filled = barWidth
			}

			sb.WriteString(fmt.Sprintf(" %-35s [%s::b]%6.1f°C[white] │ ",
				temp.Key, color, temp.Temperature))

			// Gradient bar
			for i := 0; i < barWidth; i++ {
				if i < filled {
					tempAtPos := (float64(i) / float64(barWidth)) * maxScale
					if tempAtPos > 85 {
						sb.WriteString("[red]█")
					} else if tempAtPos > 75 {
						sb.WriteString("[yellow]█")
					} else if tempAtPos > 60 {
						sb.WriteString("[cyan]█")
					} else {
						sb.WriteString("[green]█")
					}
				} else {
					sb.WriteString("[darkgray]░")
				}
			}
			sb.WriteString(fmt.Sprintf("[white] [%s]%s[white]\n", color, status))
		}

		sb.WriteString("\n " + strings.Repeat("─", 78) + "\n")
		sb.WriteString(fmt.Sprintf(" Highest Temperature: [red::b]%.1f°C[white]  │  Scale: 0°C ─────────── 120°C\n", maxTemp))
		sb.WriteString(" [darkgray]Legend: [green]Normal[white] | [cyan]Warm[white] | [yellow]High[white] | [red]Critical[white]\n")
	}

	t.TempView.SetText(sb.String())
	t.ApplyTheme()
}
