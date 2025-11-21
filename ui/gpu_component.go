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
			SetTextColor(tcell.ColorYellow).
			SetSelectable(false).
			SetAlign(tview.AlignCenter))
	}

	return &GPUComponent{Table: t}
}

func (g *GPUComponent) ApplyTheme() {
	g.Table.SetBorderColor(CurrentTheme.Border)
	g.Table.SetTitleColor(CurrentTheme.HeaderTitle)
}

func (g *GPUComponent) Update(stats []core.GPUStat) {
	g.Table.Clear()

	// Re-add headers
	headers := []string{"Index", "Name", "Util%", "Mem Used", "Mem Total", "Temp"}
	for i, h := range headers {
		g.Table.SetCell(0, i, tview.NewTableCell(h).
			SetTextColor(tcell.ColorYellow).
			SetSelectable(false).
			SetAlign(tview.AlignCenter))
	}

	if len(stats) == 0 {
		g.Table.SetCell(1, 0, tview.NewTableCell("No NVIDIA GPU found or nvidia-smi not available").
			SetTextColor(tcell.ColorRed).
			SetAlign(tview.AlignCenter).
			SetExpansion(1))
		return
	}

	for i, stat := range stats {
		row := i + 1

		g.Table.SetCell(row, 0, tview.NewTableCell(fmt.Sprintf("%d", stat.Index)).SetAlign(tview.AlignCenter))
		g.Table.SetCell(row, 1, tview.NewTableCell(stat.Name).SetAlign(tview.AlignLeft))

		utilBar := makeMiniBar(float64(stat.Utilization), 10)
		g.Table.SetCell(row, 2, tview.NewTableCell(fmt.Sprintf("%3d%% %s", stat.Utilization, utilBar)).SetAlign(tview.AlignRight))

		g.Table.SetCell(row, 3, tview.NewTableCell(fmt.Sprintf("%d MiB", stat.MemoryUsed)).SetAlign(tview.AlignRight))
		g.Table.SetCell(row, 4, tview.NewTableCell(fmt.Sprintf("%d MiB", stat.MemoryTotal)).SetAlign(tview.AlignRight))
		g.Table.SetCell(row, 5, tview.NewTableCell(fmt.Sprintf("%d°C", stat.Temp)).SetAlign(tview.AlignCenter))
	}
	g.ApplyTheme()
}
