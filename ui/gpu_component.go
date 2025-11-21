package ui

import (
	"fmt"

	"omarchy-monitor/core"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type GPUComponent struct {
	Table *tview.Table
}

func NewGPUComponent() *GPUComponent {
	t := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false).
		SetFixed(1, 0)

	headers := []string{"Index", "Name", "Util%", "Mem Used", "Mem Total", "Temp"}
	for i, h := range headers {
		t.SetCell(0, i, tview.NewTableCell(h).
			SetTextColor(CurrentTheme.TableHead).
			SetSelectable(false).
			SetAlign(tview.AlignCenter))
	}

	return &GPUComponent{Table: t}
}

func (g *GPUComponent) ApplyTheme() {
	g.Table.SetBorderColor(CurrentTheme.Border)
	g.Table.SetTitleColor(CurrentTheme.HeaderTitle)
	g.Table.SetSelectedStyle(tcell.StyleDefault.Foreground(CurrentTheme.SelectedFg).Background(CurrentTheme.SelectedBg))
}

func (g *GPUComponent) Update(stats []core.GPUStat) {
	g.Table.Clear()

	// Re-add headers
	headers := []string{"#", "Name", "GPU Util", "Memory", "Temp", "Power", "Fan", "Clocks"}
	alignments := []int{
		tview.AlignCenter, tview.AlignLeft, tview.AlignRight,
		tview.AlignRight, tview.AlignCenter, tview.AlignRight,
		tview.AlignCenter, tview.AlignRight,
	}
	for i, h := range headers {
		g.Table.SetCell(0, i, tview.NewTableCell(h).
			SetTextColor(CurrentTheme.TableHead).
			SetSelectable(false).
			SetAlign(alignments[i]))
	}

	if len(stats) == 0 {
		g.Table.SetCell(1, 0, tview.NewTableCell("No NVIDIA GPU found or nvidia-smi not available").
			SetTextColor(CurrentTheme.HighUsage).
			SetAlign(tview.AlignCenter).
			SetExpansion(1))
		return
	}

	for i, stat := range stats {
		row := i + 1

		// Index
		g.Table.SetCell(row, 0, tview.NewTableCell(fmt.Sprintf("%d", stat.Index)).
			SetTextColor(CurrentTheme.Foreground).
			SetAlign(tview.AlignCenter))

		// Name
		g.Table.SetCell(row, 1, tview.NewTableCell(stat.Name).
			SetTextColor(CurrentTheme.Foreground).
			SetAlign(tview.AlignLeft))

		// GPU Utilization with bar
		utilBar := makeMiniBar(float64(stat.Utilization), 12)
		utilColor := CurrentTheme.LowUsage
		if stat.Utilization > 80 {
			utilColor = CurrentTheme.HighUsage
		} else if stat.Utilization > 50 {
			utilColor = CurrentTheme.MedUsage
		}
		g.Table.SetCell(row, 2, tview.NewTableCell(fmt.Sprintf("%3d%% %s", stat.Utilization, utilBar)).
			SetTextColor(utilColor).
			SetAlign(tview.AlignRight))

		// Memory with usage bar
		memBar := makeMiniBar(stat.MemoryPercent, 12)
		memColor := CurrentTheme.LowUsage
		if stat.MemoryPercent > 80 {
			memColor = CurrentTheme.HighUsage
		} else if stat.MemoryPercent > 60 {
			memColor = CurrentTheme.MedUsage
		}
		memText := fmt.Sprintf("%d/%d MiB %s", stat.MemoryUsed, stat.MemoryTotal, memBar)
		g.Table.SetCell(row, 3, tview.NewTableCell(memText).
			SetTextColor(memColor).
			SetAlign(tview.AlignRight))

		// Temperature with color coding
		tempColor := CurrentTheme.LowUsage
		if stat.Temp > 80 {
			tempColor = CurrentTheme.HighUsage
		} else if stat.Temp > 70 {
			tempColor = CurrentTheme.MedUsage
		}
		g.Table.SetCell(row, 4, tview.NewTableCell(fmt.Sprintf("%d°C", stat.Temp)).
			SetTextColor(tempColor).
			SetAlign(tview.AlignCenter))

		// Power consumption
		powerText := "-"
		if stat.PowerDraw > 0 {
			powerPercent := 0.0
			if stat.PowerLimit > 0 {
				powerPercent = (stat.PowerDraw / stat.PowerLimit) * 100
			}
			powerText = fmt.Sprintf("%.0fW", stat.PowerDraw)
			if powerPercent > 0 {
				powerText = fmt.Sprintf("%.0fW (%.0f%%)", stat.PowerDraw, powerPercent)
			}
		}
		g.Table.SetCell(row, 5, tview.NewTableCell(powerText).
			SetTextColor(CurrentTheme.Foreground).
			SetAlign(tview.AlignRight))

		// Fan speed
		fanText := "-"
		fanColor := CurrentTheme.Foreground
		if stat.FanSpeed > 0 {
			fanText = fmt.Sprintf("%d%%", stat.FanSpeed)
			if stat.FanSpeed > 80 {
				fanColor = CurrentTheme.HighUsage
			} else if stat.FanSpeed > 50 {
				fanColor = CurrentTheme.MedUsage
			}
		}
		g.Table.SetCell(row, 6, tview.NewTableCell(fanText).
			SetTextColor(fanColor).
			SetAlign(tview.AlignCenter))

		// Clock speeds
		clockText := "-"
		if stat.ClockCore > 0 || stat.ClockMemory > 0 {
			clockText = fmt.Sprintf("C:%d M:%d MHz", stat.ClockCore, stat.ClockMemory)
		}
		g.Table.SetCell(row, 7, tview.NewTableCell(clockText).
			SetTextColor(CurrentTheme.HeaderValue).
			SetAlign(tview.AlignRight))
	}
	g.ApplyTheme()
}
