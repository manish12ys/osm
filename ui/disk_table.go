package ui

import (
	"fmt"
	"strings"

	"omarchy-monitor/core"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type DiskTableComponent struct {
	Table *tview.Table
}

func NewDiskTableComponent() *DiskTableComponent {
	t := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false).
		SetFixed(1, 0)

	headers := []string{"Device", "Mount", "FS", "Total", "Used", "Free", "Use%"}
	for i, h := range headers {
		t.SetCell(0, i, tview.NewTableCell(h).
			SetTextColor(tcell.ColorYellow).
			SetSelectable(false).
			SetAlign(tview.AlignLeft))
	}

	return &DiskTableComponent{Table: t}
}

// ApplyTheme updates the table style
func (d *DiskTableComponent) ApplyTheme() {
	d.Table.SetBorderColor(CurrentTheme.Border)
	d.Table.SetTitleColor(CurrentTheme.HeaderTitle)
	d.Table.SetSelectedStyle(tcell.StyleDefault.Foreground(CurrentTheme.SelectedFg).Background(CurrentTheme.SelectedBg))
}

func (d *DiskTableComponent) Update(stats []core.DiskStat) {
	d.Table.Clear()

	// Re-add headers with alignment
	headers := []string{"Device", "Mount", "FS", "Size", "Used", "Free", "Use%", "Usage", "Read", "Write", "IOPS", "Inodes"}
	alignments := []int{
		tview.AlignLeft, tview.AlignLeft, tview.AlignCenter,
		tview.AlignRight, tview.AlignRight, tview.AlignRight,
		tview.AlignRight, tview.AlignLeft, tview.AlignRight,
		tview.AlignRight, tview.AlignRight, tview.AlignRight,
	}
	for i, h := range headers {
		d.Table.SetCell(0, i, tview.NewTableCell(h).
			SetTextColor(CurrentTheme.TableHead).
			SetSelectable(false).
			SetAlign(alignments[i]))
	}

	for i, s := range stats {
		row := i + 1
		d.setCellAligned(row, 0, s.Device, CurrentTheme.Foreground, tview.AlignLeft)
		d.setCellAligned(row, 1, s.Mountpoint, CurrentTheme.Foreground, tview.AlignLeft)
		d.setCellAligned(row, 2, s.Fstype, CurrentTheme.HeaderValue, tview.AlignCenter)
		d.setCellAligned(row, 3, formatBytes(s.Total), CurrentTheme.Foreground, tview.AlignRight)
		d.setCellAligned(row, 4, formatBytes(s.Used), CurrentTheme.Foreground, tview.AlignRight)
		d.setCellAligned(row, 5, formatBytes(s.Free), CurrentTheme.Foreground, tview.AlignRight)

		color := CurrentTheme.LowUsage
		if s.UsedPercent > 90 {
			color = CurrentTheme.HighUsage
		} else if s.UsedPercent > 70 {
			color = CurrentTheme.MedUsage
		}
		d.setCellAligned(row, 6, fmt.Sprintf("%5.1f%%", s.UsedPercent), color, tview.AlignRight)

		// Add visual bar with gradient
		barWidth := 20
		filled := int((s.UsedPercent / 100.0) * float64(barWidth))
		if filled > barWidth {
			filled = barWidth
		}

		var bar strings.Builder
		for j := 0; j < barWidth; j++ {
			if j < filled {
				pct := (float64(j) / float64(barWidth)) * 100
				if pct > 90 {
					bar.WriteString("█")
				} else if pct > 70 {
					bar.WriteString("█")
				} else {
					bar.WriteString("█")
				}
			} else {
				bar.WriteString("·")
			}
		}
		d.setCellAligned(row, 7, bar.String(), color, tview.AlignLeft)

		// I/O Statistics
		readSpeed := formatSpeed(s.ReadSpeed)
		writeSpeed := formatSpeed(s.WriteSpeed)
		iops := fmt.Sprintf("%.0f", s.IOPS)

		readColor := CurrentTheme.Foreground
		if s.ReadSpeed > 100*1024*1024 { // > 100 MB/s
			readColor = CurrentTheme.LowUsage
		}
		writeColor := CurrentTheme.Foreground
		if s.WriteSpeed > 100*1024*1024 { // > 100 MB/s
			writeColor = CurrentTheme.LowUsage
		}

		d.setCellAligned(row, 8, readSpeed, readColor, tview.AlignRight)
		d.setCellAligned(row, 9, writeSpeed, writeColor, tview.AlignRight)
		d.setCellAligned(row, 10, iops, CurrentTheme.HeaderValue, tview.AlignRight)

		// Inode usage
		inodeStr := "-"
		if s.InodesTotal > 0 {
			inodeStr = fmt.Sprintf("%.1f%%", s.InodesPercent)
		}
		inodeColor := CurrentTheme.Foreground
		if s.InodesPercent > 90 {
			inodeColor = CurrentTheme.HighUsage
		} else if s.InodesPercent > 70 {
			inodeColor = CurrentTheme.MedUsage
		}
		d.setCellAligned(row, 11, inodeStr, inodeColor, tview.AlignRight)
	}
	d.ApplyTheme()
}

func (d *DiskTableComponent) setCell(row, col int, text string, color tcell.Color) {
	cell := d.Table.GetCell(row, col)
	if cell == nil {
		cell = tview.NewTableCell(text)
	} else {
		cell.SetText(text)
	}
	cell.SetTextColor(color)
	d.Table.SetCell(row, col, cell)
}

func (d *DiskTableComponent) setCellAligned(row, col int, text string, color tcell.Color, align int) {
	cell := d.Table.GetCell(row, col)
	if cell == nil {
		cell = tview.NewTableCell(text)
	} else {
		cell.SetText(text)
	}
	cell.SetTextColor(color).SetAlign(align)
	d.Table.SetCell(row, col, cell)
}
