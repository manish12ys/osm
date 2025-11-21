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

	// Helper to get color tag
	cLow := fmt.Sprintf("[#%06x]", CurrentTheme.LowUsage.Hex())
	cMed := fmt.Sprintf("[#%06x]", CurrentTheme.MedUsage.Hex())
	cHigh := fmt.Sprintf("[#%06x]", CurrentTheme.HighUsage.Hex())
	cCrit := fmt.Sprintf("[#%06x]", CurrentTheme.HighUsage.Hex()) // Use High for Critical for now, or add Critical to theme
	cFore := fmt.Sprintf("[#%06x]", CurrentTheme.Foreground.Hex())
	cDim := "[darkgray]" // Keep darkgray for background bars or use a specific theme color if available
	cTitle := fmt.Sprintf("[#%06x::b]", CurrentTheme.HeaderTitle.Hex())
	cLabel := fmt.Sprintf("[#%06x]", CurrentTheme.HeaderValue.Hex())

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

		sb.WriteString(fmt.Sprintf("\n %sCPU TOTAL USAGE%s\n", cTitle, cFore))
		sb.WriteString(" " + strings.Repeat("─", 78) + "\n\n")

		sb.WriteString(fmt.Sprintf(" %sCurrent:%s %s%.1f%%%s  │  %sMin:%s %.1f%%  │  %sAvg:%s %.1f%%  │  %sMax:%s %.1f%%\n\n",
			cTitle, cFore, cLow, current, cFore, cLabel, cFore, min, cLabel, cFore, avg, cCrit, cFore, max))

		// Enhanced Sparkline with better resolution
		barChars := []string{" ", "▂", "▃", "▄", "▅", "▆", "▇", "█"}

		sb.WriteString(fmt.Sprintf(" %sHistory (60s):%s\n ", cFore, cFore))
		for i, v := range cpuData {
			// Color based on value
			color := cLow
			if v > 80 {
				color = cCrit
			} else if v > 60 {
				color = cHigh // Use High for > 60
			} else if v > 40 {
				color = cMed
			}

			idx := int((v / 100.0) * 7)
			if idx > 7 {
				idx = 7
			}
			if idx < 0 {
				idx = 0
			}
			sb.WriteString(fmt.Sprintf("%s%s", color, barChars[idx]))

			// Add spacing every 10 chars for readability
			if (i+1)%10 == 0 && i < len(cpuData)-1 {
				sb.WriteString(" ")
			}
		}
		sb.WriteString(fmt.Sprintf("%s\n\n", cFore))

		// Large progress bar with gradient effect
		barWidth := 80
		filled := int((current / 100) * float64(barWidth))
		if filled > barWidth {
			filled = barWidth
		}

		sb.WriteString(fmt.Sprintf(" %s", cFore))
		for i := 0; i < barWidth; i++ {
			if i < filled {
				// Color gradient based on position
				percent := float64(i) / float64(barWidth) * 100
				if percent > 80 {
					sb.WriteString(fmt.Sprintf("%s█", cCrit))
				} else if percent > 60 {
					sb.WriteString(fmt.Sprintf("%s█", cHigh))
				} else if percent > 40 {
					sb.WriteString(fmt.Sprintf("%s█", cMed))
				} else {
					sb.WriteString(fmt.Sprintf("%s█", cLow))
				}
			} else {
				sb.WriteString(fmt.Sprintf("%s·", cDim))
			}
		}
		sb.WriteString(fmt.Sprintf("%s\n", cFore))
		sb.WriteString(" 0%                                                                           100%\n\n")
	}

	// Per-core usage with enhanced visualization
	if len(perCoreUsage) > 0 {
		sb.WriteString(fmt.Sprintf(" %sPER-CORE USAGE%s\n", cTitle, cFore))
		sb.WriteString(" " + strings.Repeat("─", 78) + "\n\n")

		for i, usage := range perCoreUsage {
			color := cLow
			if usage > 80 {
				color = cCrit
			} else if usage > 60 {
				color = cHigh
			} else if usage > 40 {
				color = cMed
			}

			barWidth := 60
			filled := int((usage / 100) * float64(barWidth))
			if filled > barWidth {
				filled = barWidth
			}

			sb.WriteString(fmt.Sprintf(" Core %-2d %s%6.1f%%%s │ ", i, color, usage, cFore))

			// Gradient bar
			for j := 0; j < barWidth; j++ {
				if j < filled {
					percent := float64(j) / float64(barWidth) * 100
					if percent > 80 {
						sb.WriteString(fmt.Sprintf("%s█", cCrit))
					} else if percent > 60 {
						sb.WriteString(fmt.Sprintf("%s█", cHigh))
					} else if percent > 40 {
						sb.WriteString(fmt.Sprintf("%s█", cMed))
					} else {
						sb.WriteString(fmt.Sprintf("%s█", cLow))
					}
				} else {
					sb.WriteString(fmt.Sprintf("%s░", cDim))
				}
			}
			sb.WriteString(fmt.Sprintf("%s\n", cFore))
		}

		sb.WriteString(fmt.Sprintf("\n %sLegend: %sLow%s | %sMedium%s | %sHigh%s | %sCritical%s\n",
			cDim, cLow, cFore, cMed, cFore, cHigh, cFore, cCrit, cFore))
	}

	c.View.SetText(sb.String())
	c.ApplyTheme()
}
