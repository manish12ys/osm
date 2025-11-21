package ui

import (
	"fmt"
	"strings"

	"omarchy-monitor/core"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type NetTableComponent struct {
	Table *tview.Table
}

func NewNetTableComponent() *NetTableComponent {
	t := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false).
		SetFixed(1, 0)

	headers := []string{"Interface", "Sent (MB)", "Recv (MB)", "Pkts Sent", "Pkts Recv"}
	for i, h := range headers {
		t.SetCell(0, i, tview.NewTableCell(h).
			SetTextColor(tcell.ColorYellow).
			SetSelectable(false).
			SetAlign(tview.AlignLeft))
	}

	return &NetTableComponent{Table: t}
}

// ApplyTheme updates the table style
func (n *NetTableComponent) ApplyTheme() {
	n.Table.SetBorderColor(CurrentTheme.Border)
	n.Table.SetTitleColor(CurrentTheme.HeaderTitle)
	n.Table.SetSelectedStyle(tcell.StyleDefault.Foreground(CurrentTheme.SelectedFg).Background(CurrentTheme.SelectedBg))
}

func (n *NetTableComponent) Update(stats []core.NetStat) {
	n.Table.Clear()

	headers := []string{"Interface", "Total Sent", "Total Recv", "↑ Speed", "↓ Speed", "Pkts/s", "Errors", "Drops", "Status"}
	alignments := []int{
		tview.AlignLeft, tview.AlignRight, tview.AlignRight,
		tview.AlignRight, tview.AlignRight, tview.AlignRight,
		tview.AlignRight, tview.AlignRight, tview.AlignCenter,
	}
	for i, h := range headers {
		n.Table.SetCell(0, i, tview.NewTableCell(h).
			SetTextColor(CurrentTheme.TableHead).
			SetSelectable(false).
			SetAlign(alignments[i]))
	}

	for i, s := range stats {
		row := i + 1

		// Interface name
		ifaceColor := CurrentTheme.Foreground
		n.setCellAligned(row, 0, s.Name, ifaceColor, tview.AlignLeft)

		// Total bytes sent/received
		n.setCellAligned(row, 1, formatBytes(s.BytesSent), CurrentTheme.Foreground, tview.AlignRight)
		n.setCellAligned(row, 2, formatBytes(s.BytesRecv), CurrentTheme.Foreground, tview.AlignRight)

		// Real-time speeds
		sendSpeedStr := formatSpeed(s.SendSpeed)
		recvSpeedStr := formatSpeed(s.RecvSpeed)

		sendColor := CurrentTheme.Foreground
		if s.SendSpeed > 10*1024*1024 { // > 10 MB/s
			sendColor = CurrentTheme.LowUsage
		} else if s.SendSpeed > 1*1024*1024 { // > 1 MB/s
			sendColor = CurrentTheme.HeaderValue
		}

		recvColor := CurrentTheme.Foreground
		if s.RecvSpeed > 10*1024*1024 { // > 10 MB/s
			recvColor = CurrentTheme.LowUsage
		} else if s.RecvSpeed > 1*1024*1024 { // > 1 MB/s
			recvColor = CurrentTheme.HeaderValue
		}

		n.setCellAligned(row, 3, sendSpeedStr, sendColor, tview.AlignRight)
		n.setCellAligned(row, 4, recvSpeedStr, recvColor, tview.AlignRight)

		// Packets per second (combined)
		totalPPS := s.SendPPS + s.RecvPPS
		ppsStr := "-"
		if totalPPS > 0 {
			ppsStr = fmt.Sprintf("%.0f", totalPPS)
		}
		n.setCellAligned(row, 5, ppsStr, CurrentTheme.HeaderValue, tview.AlignRight)

		// Errors (combined in/out)
		totalErrors := s.Errin + s.Errout
		errStr := "-"
		errColor := CurrentTheme.Foreground
		if totalErrors > 0 {
			errStr = fmt.Sprintf("%d", totalErrors)
			errColor = CurrentTheme.HighUsage
		}
		n.setCellAligned(row, 6, errStr, errColor, tview.AlignRight)

		// Drops (combined in/out)
		totalDrops := s.Dropin + s.Dropout
		dropStr := "-"
		dropColor := CurrentTheme.Foreground
		if totalDrops > 0 {
			dropStr = fmt.Sprintf("%d", totalDrops)
			dropColor = CurrentTheme.MedUsage
		}
		n.setCellAligned(row, 7, dropStr, dropColor, tview.AlignRight)

		// Status indicator with activity bar
		status := "Idle"
		statusColor := tcell.ColorGray
		activityBar := ""

		if s.SendSpeed > 0 || s.RecvSpeed > 0 {
			status = "Active"
			statusColor = CurrentTheme.LowUsage

			// Create mini activity bar
			maxSpeed := s.SendSpeed
			if s.RecvSpeed > maxSpeed {
				maxSpeed = s.RecvSpeed
			}

			barWidth := 8
			if maxSpeed > 0 {
				var bar strings.Builder
				bar.WriteString(" [")
				filled := 0
				if maxSpeed > 100*1024*1024 { // > 100 MB/s
					filled = barWidth
				} else if maxSpeed > 10*1024*1024 { // > 10 MB/s
					filled = barWidth * 3 / 4
				} else if maxSpeed > 1*1024*1024 { // > 1 MB/s
					filled = barWidth / 2
				} else {
					filled = barWidth / 4
				}

				for j := 0; j < barWidth; j++ {
					if j < filled {
						bar.WriteString("▮")
					} else {
						bar.WriteString("·")
					}
				}
				bar.WriteString("]")
				activityBar = bar.String()
			}
		}

		n.setCellAligned(row, 8, status+activityBar, statusColor, tview.AlignCenter)
	}
	n.ApplyTheme()
}

func (n *NetTableComponent) setCell(row, col int, text string, color tcell.Color) {
	cell := n.Table.GetCell(row, col)
	if cell == nil {
		cell = tview.NewTableCell(text)
	} else {
		cell.SetText(text)
	}
	cell.SetTextColor(color)
	n.Table.SetCell(row, col, cell)
}

func (n *NetTableComponent) setCellAligned(row, col int, text string, color tcell.Color, align int) {
	cell := n.Table.GetCell(row, col)
	if cell == nil {
		cell = tview.NewTableCell(text)
	} else {
		cell.SetText(text)
	}
	cell.SetTextColor(color).SetAlign(align)
	n.Table.SetCell(row, col, cell)
}

func formatNumber(n uint64) string {
	if n > 1000000000 {
		return fmt.Sprintf("%.1fB", float64(n)/1000000000)
	} else if n > 1000000 {
		return fmt.Sprintf("%.1fM", float64(n)/1000000)
	} else if n > 1000 {
		return fmt.Sprintf("%.1fK", float64(n)/1000)
	}
	return fmt.Sprintf("%d", n)
}
