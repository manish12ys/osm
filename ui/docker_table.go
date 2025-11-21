package ui

import (
	"omarchy-monitor/core"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type DockerTableComponent struct {
	Table *tview.Table
}

func NewDockerTableComponent() *DockerTableComponent {
	t := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false).
		SetFixed(1, 0)

	headers := []string{"ID", "Name", "Image", "State", "Status"}
	alignments := []int{tview.AlignLeft, tview.AlignLeft, tview.AlignLeft, tview.AlignCenter, tview.AlignLeft}
	for i, h := range headers {
		t.SetCell(0, i, tview.NewTableCell(h).
			SetTextColor(tcell.ColorYellow).
			SetSelectable(false).
			SetAlign(alignments[i]))
	}

	return &DockerTableComponent{Table: t}
}

func (d *DockerTableComponent) ApplyTheme() {
	d.Table.SetBorderColor(CurrentTheme.Border)
	d.Table.SetTitleColor(CurrentTheme.HeaderTitle)
	d.Table.SetSelectedStyle(tcell.StyleDefault.Foreground(CurrentTheme.SelectedFg).Background(CurrentTheme.SelectedBg))

	for i := 0; i < d.Table.GetColumnCount(); i++ {
		cell := d.Table.GetCell(0, i)
		if cell != nil {
			cell.SetTextColor(CurrentTheme.TableHead)
		}
	}
}

func (d *DockerTableComponent) Update(containers []core.ContainerStat) {
	// Clear existing rows (except header)
	rowCount := d.Table.GetRowCount()
	for i := rowCount - 1; i > 0; i-- {
		d.Table.RemoveRow(i)
	}

	if len(containers) == 0 {
		d.Table.SetCell(1, 0, tview.NewTableCell("No containers found").
			SetTextColor(CurrentTheme.Foreground).
			SetAlign(tview.AlignCenter).
			SetExpansion(1))
		return
	}

	for i, c := range containers {
		row := i + 1

		// State color
		stateColor := CurrentTheme.Foreground
		if c.State == "running" {
			stateColor = CurrentTheme.LowUsage
		} else if c.State == "exited" {
			stateColor = CurrentTheme.MedUsage
		} else if c.State == "dead" {
			stateColor = CurrentTheme.HighUsage
		}

		d.Table.SetCell(row, 0, tview.NewTableCell(c.ID).
			SetTextColor(CurrentTheme.Foreground).
			SetAlign(tview.AlignLeft))
		d.Table.SetCell(row, 1, tview.NewTableCell(c.Name).
			SetTextColor(CurrentTheme.Foreground).
			SetAlign(tview.AlignLeft))
		d.Table.SetCell(row, 2, tview.NewTableCell(c.Image).
			SetTextColor(CurrentTheme.HeaderValue).
			SetAlign(tview.AlignLeft))
		d.Table.SetCell(row, 3, tview.NewTableCell(c.State).
			SetTextColor(stateColor).
			SetAlign(tview.AlignCenter))
		d.Table.SetCell(row, 4, tview.NewTableCell(c.Status).
			SetTextColor(CurrentTheme.Foreground).
			SetAlign(tview.AlignLeft))
	}

	d.ApplyTheme()
}
