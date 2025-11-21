package ui

import (
	"fmt"

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
}

func (n *NetTableComponent) Update(stats []core.NetStat) {
	n.Table.Clear()

	headers := []string{"Interface", "Sent", "Received", "Pkts Out", "Pkts In", "Status"}
	alignments := []int{tview.AlignLeft, tview.AlignRight, tview.AlignRight, tview.AlignRight, tview.AlignRight, tview.AlignCenter}
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

		// Format bytes sent (auto-scale to MB/GB)
		sentMB := float64(s.BytesSent) / 1024 / 1024
		sentStr := ""
		if sentMB > 1024 {
			sentStr = fmt.Sprintf("%.2f GB", sentMB/1024)
		} else {
			sentStr = fmt.Sprintf("%.1f MB", sentMB)
		}
		n.setCellAligned(row, 1, sentStr, CurrentTheme.Foreground, tview.AlignRight)

		// Format bytes received (auto-scale to MB/GB)
		recvMB := float64(s.BytesRecv) / 1024 / 1024
		recvStr := ""
		if recvMB > 1024 {
			recvStr = fmt.Sprintf("%.2f GB", recvMB/1024)
		} else {
			recvStr = fmt.Sprintf("%.1f MB", recvMB)
		}
		n.setCellAligned(row, 2, recvStr, CurrentTheme.Foreground, tview.AlignRight)

		// Packets with formatting
		n.setCellAligned(row, 3, formatNumber(s.PacketsSent), CurrentTheme.HeaderValue, tview.AlignRight)
		n.setCellAligned(row, 4, formatNumber(s.PacketsRecv), CurrentTheme.HeaderValue, tview.AlignRight)

		// Status indicator
		status := "Active"
		statusColor := CurrentTheme.LowUsage
		if s.BytesSent == 0 && s.BytesRecv == 0 {
			status = "Idle"
			statusColor = tcell.ColorGray
		}
		n.setCellAligned(row, 5, status, statusColor, tview.AlignCenter)
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
