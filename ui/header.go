package ui

import (
	"fmt"
	"strings"

	"omarchy-monitor/core"

	"github.com/rivo/tview"
)

// HeaderComponent manages the header view
type HeaderComponent struct {
	View *tview.TextView
}

// NewHeaderComponent creates a new header
func NewHeaderComponent() *HeaderComponent {
	h := &HeaderComponent{
		View: tview.NewTextView().SetDynamicColors(true),
	}
	return h
}

// Update refreshes the header with new stats
func (h *HeaderComponent) Update(stats *core.Stats) {
	// Enhanced ASCII bars with gradient
	cpuBar := makeEnhancedBar(stats.CPUUsage, 30)
	memBar := makeEnhancedBar(stats.MemUsage, 30)

	// Helper to colorize text
	cTitle := fmt.Sprintf("[#%06x::b]", CurrentTheme.HeaderTitle.Hex())
	cValue := fmt.Sprintf("[#%06x]", CurrentTheme.HeaderValue.Hex())

	var text strings.Builder
	text.WriteString("\n")

	// Host and system info line
	text.WriteString(fmt.Sprintf(" %sHost:%s %s (%s)  %sUptime:%s %dh %dm\n",
		cTitle, cValue, stats.Hostname, stats.OS,
		cTitle, cValue, stats.Uptime/3600, (stats.Uptime%3600)/60))

	// CPU line with bar and load average
	text.WriteString(fmt.Sprintf(" %sCPU:%s  %5.1f%% %s  %sLoad:%s %.2f %.2f %.2f\n",
		cTitle, cValue, stats.CPUUsage, cpuBar,
		cTitle, cValue, stats.LoadAvg1, stats.LoadAvg5, stats.LoadAvg15))

	// Memory line with bar and usage
	text.WriteString(fmt.Sprintf(" %sMEM:%s  %5.1f%% %s  %sUsed:%s %d / %d MB\n",
		cTitle, cValue, stats.MemUsage, memBar,
		cTitle, cValue, stats.MemUsed/1024/1024, stats.MemTotal/1024/1024))

	h.View.SetText(text.String())
	h.View.SetBorderColor(CurrentTheme.Border)
	h.View.SetTitleColor(CurrentTheme.HeaderTitle)
}

func makeEnhancedBar(percent float64, width int) string {
	filled := int(percent / 100.0 * float64(width))
	if filled > width {
		filled = width
	}

	var bar strings.Builder
	bar.WriteString("[white]")

	for i := 0; i < width; i++ {
		if i < filled {
			// Color based on percentage
			currentPercent := (float64(i) / float64(width)) * 100
			if currentPercent > 80 {
				bar.WriteString("[red]█")
			} else if currentPercent > 60 {
				bar.WriteString("[yellow]█")
			} else if currentPercent > 40 {
				bar.WriteString("[cyan]█")
			} else {
				bar.WriteString("[green]█")
			}
		} else {
			bar.WriteString("[darkgray]░")
		}
	}
	bar.WriteString("[white]")

	return bar.String()
}
