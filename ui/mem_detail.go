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

		sb.WriteString("\n [magenta::b]MEMORY USAGE[white]\n")
		sb.WriteString(" " + strings.Repeat("─", 78) + "\n\n")

		sb.WriteString(fmt.Sprintf(" [magenta::b]Current:[white] [magenta]%.1f%%[white]  │  [cyan]Min:[white] %.1f%%  │  [yellow]Avg:[white] %.1f%%  │  [red]Max:[white] %.1f%%\n\n",
			current, min, avg, max))

		// Enhanced Sparkline with color coding
		barChars := []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}

		sb.WriteString(" [white]History (60s):[white]\n ")
		for i, v := range memData {
			// Color based on value
			color := "green"
			if v > 85 {
				color = "red"
			} else if v > 70 {
				color = "yellow"
			} else if v > 50 {
				color = "magenta"
			} else if v > 30 {
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
			if (i+1)%10 == 0 && i < len(memData)-1 {
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
				if percent > 85 {
					sb.WriteString("[red]█")
				} else if percent > 70 {
					sb.WriteString("[yellow]█")
				} else if percent > 50 {
					sb.WriteString("[magenta]█")
				} else if percent > 30 {
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

	// Memory breakdown with enhanced visuals
	if stats != nil {
		totalGB := float64(stats.MemTotal) / 1024 / 1024 / 1024
		usedGB := float64(stats.MemUsed) / 1024 / 1024 / 1024
		freeGB := totalGB - usedGB
		availableGB := freeGB // Simplified - in reality would use Available field

		sb.WriteString(" [cyan::b]MEMORY BREAKDOWN[white]\n")
		sb.WriteString(" " + strings.Repeat("─", 78) + "\n\n")

		// Detailed breakdown
		sb.WriteString(fmt.Sprintf(" Total:      [white::b]%8.2f GB[white]\n", totalGB))
		sb.WriteString(fmt.Sprintf(" [magenta]Used:       [magenta::b]%8.2f GB[white]  ([magenta]%.1f%%[white])\n", usedGB, stats.MemUsage))
		sb.WriteString(fmt.Sprintf(" [green]Free:       [green::b]%8.2f GB[white]  ([green]%.1f%%[white])\n", freeGB, 100-stats.MemUsage))
		sb.WriteString(fmt.Sprintf(" [cyan]Available:  [cyan::b]%8.2f GB[white]  ([cyan]%.1f%%[white])\n\n", availableGB, (availableGB/totalGB)*100))

		// Enhanced visual breakdown with segments
		totalWidth := 80
		usedWidth := int((stats.MemUsage / 100) * float64(totalWidth))
		freeWidth := totalWidth - usedWidth

		sb.WriteString(" ")

		// Used portion with gradient
		for i := 0; i < usedWidth; i++ {
			percent := float64(i) / float64(usedWidth) * stats.MemUsage
			if percent > 85 {
				sb.WriteString("[red]█")
			} else if percent > 70 {
				sb.WriteString("[yellow]█")
			} else {
				sb.WriteString("[magenta]█")
			}
		}

		// Free portion
		sb.WriteString("[green]")
		sb.WriteString(strings.Repeat("█", freeWidth))
		sb.WriteString("[white]\n")

		sb.WriteString(" [magenta]Used")
		sb.WriteString(strings.Repeat(" ", totalWidth-10))
		sb.WriteString("[green]Free[white]\n\n")

		// Memory pressure indicator
		pressure := "Low"
		pressureColor := "green"

		if stats.MemUsage > 90 {
			pressure = "Critical"
			pressureColor = "red"
		} else if stats.MemUsage > 75 {
			pressure = "High"
			pressureColor = "yellow"
		} else if stats.MemUsage > 50 {
			pressure = "Medium"
			pressureColor = "cyan"
		}

		sb.WriteString(fmt.Sprintf(" Memory Pressure: [%s::b]%s[white]\n", pressureColor, pressure))
	}

	m.View.SetText(sb.String())
	m.ApplyTheme()
}
