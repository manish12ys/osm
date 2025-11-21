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
	batStr := ""
	if stats.Battery != nil {
		batColor := "[green]"
		if stats.Battery.Capacity < 20 && stats.Battery.Status == "Discharging" {
			batColor = "[red]"
		} else if stats.Battery.Capacity < 50 {
			batColor = "[yellow]"
		}
		batStr = fmt.Sprintf("  %sBat:%s %s%d%% (%s)", cTitle, cValue, batColor, stats.Battery.Capacity, stats.Battery.Status)
	}

	// Simplified Header Layout
	// Line 1: Hostname | OS | Battery (if present)
	text.WriteString(fmt.Sprintf(" %sHost:%s %s (%s)%s\n",
		cTitle, cValue, stats.Hostname, stats.OS, batStr))

	// Line 2: CPU Usage Bar
	text.WriteString(fmt.Sprintf(" %sCPU:%s  %5.1f%% %s\n",
		cTitle, cValue, stats.CPUUsage, cpuBar))

	// Line 3: Memory Usage Bar + Details
	text.WriteString(fmt.Sprintf(" %sMEM:%s  %5.1f%% %s  %s(%d/%d MB)\n",
		cTitle, cValue, stats.MemUsage, memBar,
		cValue, stats.MemUsed/1024/1024, stats.MemTotal/1024/1024))

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
				bar.WriteString(fmt.Sprintf("[#%06x]█", CurrentTheme.HighUsage.Hex()))
			} else if currentPercent > 60 {
				bar.WriteString(fmt.Sprintf("[#%06x]█", CurrentTheme.MedUsage.Hex()))
			} else if currentPercent > 40 {
				bar.WriteString(fmt.Sprintf("[#%06x]█", CurrentTheme.TableHead.Hex())) // Use TableHead (often cyan/yellow) for mid-low
			} else {
				bar.WriteString(fmt.Sprintf("[#%06x]█", CurrentTheme.LowUsage.Hex()))
			}
		} else {
			bar.WriteString(fmt.Sprintf("[#%06x]·", CurrentTheme.RowAlt.Hex()))
		}
	}
	bar.WriteString("[white]")

	return bar.String()
}
