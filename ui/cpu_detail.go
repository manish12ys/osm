package ui

import (
	"fmt"
	"strings"

	"omarchy-monitor/core"

	"github.com/rivo/tview"
)

// CPUDetailComponent shows detailed CPU information
type CPUDetailComponent struct {
	Flex *tview.Flex
	View *tview.TextView
}

// NewCPUDetailComponent creates the CPU detail view
func NewCPUDetailComponent() *CPUDetailComponent {
	view := tview.NewTextView().SetDynamicColors(true)
	view.SetBorder(true).SetTitle(" CPU Details ")

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(view, 0, 1, false)

	return &CPUDetailComponent{
		Flex: flex,
		View: view,
	}
}

// ApplyTheme applies the current theme
func (c *CPUDetailComponent) ApplyTheme() {
	c.View.SetBorderColor(CurrentTheme.Border)
	c.View.SetTitleColor(CurrentTheme.HeaderTitle)
}

// Update refreshes the CPU detail view
func (c *CPUDetailComponent) Update(cpuHist *core.History, perCoreUsage []float64) {
	var sb strings.Builder

	// Overall CPU usage with sparkline
	cpuData := cpuHist.GetData()
	if len(cpuData) > 0 {
		current := cpuData[len(cpuData)-1]

		// Calculate stats
		var min, max, avg float64
		min = cpuData[0]
		max = cpuData[0]
		sum := 0.0

		for _, v := range cpuData {
			if v < min {
				min = v
			}
			if v > max {
				max = v
			}
			sum += v
		}
		avg = sum / float64(len(cpuData))

		sb.WriteString("\n [green::b]CPU TOTAL USAGE[white]\n")
		sb.WriteString(" " + strings.Repeat("─", 78) + "\n\n")

		sb.WriteString(fmt.Sprintf(" [green::b]Current:[white] [green]%.1f%%[white]  │  [cyan]Min:[white] %.1f%%  │  [yellow]Avg:[white] %.1f%%  │  [red]Max:[white] %.1f%%\n\n",
			current, min, avg, max))

		// Enhanced Sparkline with better resolution
		barChars := []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}

		sb.WriteString(" [white]History (60s):[white]\n ")
		for i, v := range cpuData {
			// Color based on value
			color := "green"
			if v > 80 {
				color = "red"
			} else if v > 60 {
				color = "yellow"
			} else if v > 40 {
				color = "cyan"
			}

			idx := int((v / 100.0) * 7)
			if idx > 7 {
				idx = 7
			}
			if idx < 0 {
				idx = 0
			}
			sb.WriteString(fmt.Sprintf("[%s]%s", color, barChars[idx]))

			// Add spacing every 10 chars for readability
			if (i+1)%10 == 0 && i < len(cpuData)-1 {
				sb.WriteString("[white] ")
			}
		}
		sb.WriteString("[white]\n\n")

		// Large progress bar with gradient effect
		barWidth := 80
		filled := int((current / 100) * float64(barWidth))
		if filled > barWidth {
			filled = barWidth
		}

		sb.WriteString(" [white]")
		for i := 0; i < barWidth; i++ {
			if i < filled {
				// Color gradient based on position
				percent := float64(i) / float64(barWidth) * 100
				if percent > 80 {
					sb.WriteString("[red]█")
				} else if percent > 60 {
					sb.WriteString("[yellow]█")
				} else if percent > 40 {
					sb.WriteString("[cyan]█")
				} else {
					sb.WriteString("[green]█")
				}
			} else {
				sb.WriteString("[darkgray]░")
			}
		}
		sb.WriteString("[white]\n")
		sb.WriteString(" 0%                                                                           100%\n\n")
	}

	// Per-core usage with enhanced visualization
	if len(perCoreUsage) > 0 {
		sb.WriteString(" [cyan::b]PER-CORE USAGE[white]\n")
		sb.WriteString(" " + strings.Repeat("─", 78) + "\n\n")

		for i, usage := range perCoreUsage {
			color := "green"
			if usage > 80 {
				color = "red"
			} else if usage > 60 {
				color = "yellow"
			} else if usage > 40 {
				color = "cyan"
			}

			barWidth := 60
			filled := int((usage / 100) * float64(barWidth))
			if filled > barWidth {
				filled = barWidth
			}

			sb.WriteString(fmt.Sprintf(" Core %-2d [%s::b]%6.1f%%[white] │ ", i, color, usage))

			// Gradient bar
			for j := 0; j < barWidth; j++ {
				if j < filled {
					percent := float64(j) / float64(barWidth) * 100
					if percent > 80 {
						sb.WriteString("[red]█")
					} else if percent > 60 {
						sb.WriteString("[yellow]█")
					} else if percent > 40 {
						sb.WriteString("[cyan]█")
					} else {
						sb.WriteString("[green]█")
					}
				} else {
					sb.WriteString("[darkgray]░")
				}
			}
			sb.WriteString("[white]\n")
		}

		sb.WriteString("\n [darkgray]Legend: [green]Low[white] | [cyan]Medium[white] | [yellow]High[white] | [red]Critical[white]\n")
	}

	c.View.SetText(sb.String())
	c.ApplyTheme()
}
