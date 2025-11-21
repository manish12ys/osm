package ui

import (
	"fmt"
	"strings"

	"omarchy-monitor/core"

	"github.com/rivo/tview"
)

// MemDetailComponent shows detailed memory information
type MemDetailComponent struct {
	Flex *tview.Flex
	View *tview.TextView
}

// NewMemDetailComponent creates the memory detail view
func NewMemDetailComponent() *MemDetailComponent {
	view := tview.NewTextView().SetDynamicColors(true)
	view.SetBorder(true).SetTitle(" Memory Details ")

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(view, 0, 1, false)

	return &MemDetailComponent{
		Flex: flex,
		View: view,
	}
}

// ApplyTheme applies the current theme
func (m *MemDetailComponent) ApplyTheme() {
	m.View.SetBorderColor(CurrentTheme.Border)
	m.View.SetTitleColor(CurrentTheme.HeaderTitle)
}

// Update refreshes the memory detail view
func (m *MemDetailComponent) Update(memHist *core.History, stats *core.Stats) {
	var sb strings.Builder

	// Helper to get color tag
	cLow := fmt.Sprintf("[#%06x]", CurrentTheme.LowUsage.Hex())
	cMed := fmt.Sprintf("[#%06x]", CurrentTheme.MedUsage.Hex())
	cHigh := fmt.Sprintf("[#%06x]", CurrentTheme.HighUsage.Hex())
	cFore := fmt.Sprintf("[#%06x]", CurrentTheme.Foreground.Hex())
	cDim := "[darkgray]"
	cTitle := fmt.Sprintf("[#%06x::b]", CurrentTheme.HeaderTitle.Hex())
	cLabel := fmt.Sprintf("[#%06x]", CurrentTheme.HeaderValue.Hex())
	cVal := fmt.Sprintf("[#%06x]", CurrentTheme.HeaderValue.Hex()) // Use HeaderValue for main values

	memData := memHist.GetData()
	if len(memData) > 0 {
		current := memData[len(memData)-1]

		// Calculate stats
		var min, max, avg float64
		min = memData[0]
		max = memData[0]
		sum := 0.0

		for _, v := range memData {
			if v < min {
				min = v
			}
			if v > max {
				max = v
			}
			sum += v
		}
		avg = sum / float64(len(memData))

		sb.WriteString(fmt.Sprintf("\n %sMEMORY USAGE%s\n", cTitle, cFore))
		sb.WriteString(" " + strings.Repeat("─", 78) + "\n\n")

		sb.WriteString(fmt.Sprintf(" %sCurrent:%s %s%.1f%%%s  │  %sMin:%s %.1f%%  │  %sAvg:%s %.1f%%  │  %sMax:%s %.1f%%\n\n",
			cTitle, cFore, cVal, current, cFore, cLabel, cFore, min, cLabel, cFore, avg, cHigh, cFore, max))

		// Enhanced Sparkline with color coding
		barChars := []string{" ", "▂", "▃", "▄", "▅", "▆", "▇", "█"}

		sb.WriteString(fmt.Sprintf(" %sHistory (60s):%s\n ", cFore, cFore))
		for i, v := range memData {
			// Color based on value
			color := cLow
			if v > 85 {
				color = cHigh
			} else if v > 70 {
				color = cMed
			} else if v > 50 {
				color = cVal
			} else if v > 30 {
				color = cLabel
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
			if (i+1)%10 == 0 && i < len(memData)-1 {
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
				if percent > 85 {
					sb.WriteString(fmt.Sprintf("%s█", cHigh))
				} else if percent > 70 {
					sb.WriteString(fmt.Sprintf("%s█", cMed))
				} else if percent > 50 {
					sb.WriteString(fmt.Sprintf("%s█", cVal))
				} else if percent > 30 {
					sb.WriteString(fmt.Sprintf("%s█", cLabel))
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

	// Memory breakdown with enhanced visuals
	if stats != nil {
		totalGB := float64(stats.MemTotal) / 1024 / 1024 / 1024
		usedGB := float64(stats.MemUsed) / 1024 / 1024 / 1024
		freeGB := totalGB - usedGB
		availableGB := freeGB // Simplified - in reality would use Available field

		sb.WriteString(fmt.Sprintf(" %sMEMORY BREAKDOWN%s\n", cTitle, cFore))
		sb.WriteString(" " + strings.Repeat("─", 78) + "\n\n")

		// Detailed breakdown
		sb.WriteString(fmt.Sprintf(" Total:      %s%8.2f GB%s\n", cFore, totalGB, cFore))
		sb.WriteString(fmt.Sprintf(" %sUsed:       %s%8.2f GB%s  (%s%.1f%%%s)\n", cVal, cVal, usedGB, cFore, cVal, stats.MemUsage, cFore))
		sb.WriteString(fmt.Sprintf(" %sFree:       %s%8.2f GB%s  (%s%.1f%%%s)\n", cLow, cLow, freeGB, cFore, cLow, 100-stats.MemUsage, cFore))
		sb.WriteString(fmt.Sprintf(" %sAvailable:  %s%8.2f GB%s  (%s%.1f%%%s)\n\n", cLabel, cLabel, availableGB, cFore, cLabel, (availableGB/totalGB)*100, cFore))

		// Enhanced visual breakdown with segments
		totalWidth := 80
		usedWidth := int((stats.MemUsage / 100) * float64(totalWidth))
		freeWidth := totalWidth - usedWidth

		sb.WriteString(" ")

		// Used portion with gradient
		for i := 0; i < usedWidth; i++ {
			percent := float64(i) / float64(usedWidth) * stats.MemUsage
			if percent > 85 {
				sb.WriteString(fmt.Sprintf("%s█", cHigh))
			} else if percent > 70 {
				sb.WriteString(fmt.Sprintf("%s█", cMed))
			} else {
				sb.WriteString(fmt.Sprintf("%s█", cVal))
			}
		}

		// Free portion
		sb.WriteString(cLow)
		sb.WriteString(strings.Repeat("█", freeWidth))
		sb.WriteString(fmt.Sprintf("%s\n", cFore))

		sb.WriteString(fmt.Sprintf(" %sUsed", cVal))
		sb.WriteString(strings.Repeat(" ", totalWidth-10))
		sb.WriteString(fmt.Sprintf("%sFree%s\n\n", cLow, cFore))

		// Memory pressure indicator
		pressure := "Low"
		pressureColor := cLow

		if stats.MemUsage > 90 {
			pressure = "Critical"
			pressureColor = cHigh
		} else if stats.MemUsage > 75 {
			pressure = "High"
			pressureColor = cMed
		} else if stats.MemUsage > 50 {
			pressure = "Medium"
			pressureColor = cLabel
		}

		sb.WriteString(fmt.Sprintf(" Memory Pressure: %s%s%s\n", pressureColor, pressure, cFore))
	}

	m.View.SetText(sb.String())
	m.ApplyTheme()
}
