package ui

import (
	"fmt"
	"omarchy-monitor/remote"
	"time"

	"github.com/rivo/tview"
)

type RemoteComponent struct {
	Flex       *tview.Flex
	StatsView  *tview.TextView
	StatusView *tview.TextView
	Client     *remote.Client
	LastUpdate time.Time
	Connected  bool
}

func NewRemoteComponent(remoteURL string) *RemoteComponent {
	statsView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true)
	statsView.SetBorder(true).SetTitle(" Remote System Stats ")

	statusView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	statusView.SetBorder(true).SetTitle(" Connection Status ")

	var client *remote.Client
	if remoteURL != "" {
		client = remote.NewClient(remoteURL)
	}

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(statusView, 3, 0, false).
		AddItem(statsView, 0, 1, true)

	return &RemoteComponent{
		Flex:       flex,
		StatsView:  statsView,
		StatusView: statusView,
		Client:     client,
		Connected:  false,
	}
}

func (r *RemoteComponent) ApplyTheme() {
	r.StatsView.SetBorderColor(CurrentTheme.Border)
	r.StatsView.SetTitleColor(CurrentTheme.HeaderTitle)
	r.StatsView.SetTextColor(CurrentTheme.Foreground)
	r.StatusView.SetBorderColor(CurrentTheme.Border)
	r.StatusView.SetTitleColor(CurrentTheme.HeaderTitle)
	r.StatusView.SetTextColor(CurrentTheme.Foreground)
}

func (r *RemoteComponent) Update() {
	if r.Client == nil {
		r.StatusView.Clear()
		fmt.Fprintf(r.StatusView, "[red]✗ Not Configured[-]")
		r.StatsView.Clear()
		fmt.Fprintf(r.StatsView, "[yellow]Remote monitoring not configured[-]\n\n")
		fmt.Fprintf(r.StatsView, "To enable remote monitoring, add to config.yaml:\n\n")
		fmt.Fprintf(r.StatsView, "[cyan]remote_mode: true[-]\n")
		fmt.Fprintf(r.StatsView, "[cyan]remote_url: \"http://192.168.1.100:8080\"[-]\n")
		return
	}

	stats, err := r.Client.FetchRemoteStats()
	if err != nil {
		r.Connected = false
		r.StatusView.Clear()
		fmt.Fprintf(r.StatusView, "[red]✗ Disconnected[-] - %v", err)
		r.StatsView.Clear()
		fmt.Fprintf(r.StatsView, "[red]Failed to connect to remote server[-]\n\n")
		fmt.Fprintf(r.StatsView, "[gray]URL: %s[-]\n", r.Client.BaseURL)
		fmt.Fprintf(r.StatsView, "[gray]Error: %v[-]\n", err)
		return
	}

	r.Connected = true
	r.LastUpdate = time.Now()

	// Update status
	r.StatusView.Clear()
	fmt.Fprintf(r.StatusView, "[green]✓ Connected[-] to [cyan]%s[-] | Last Update: [white]%s[-]",
		stats.Hostname, r.LastUpdate.Format("15:04:05"))

	// Update stats view
	r.StatsView.Clear()

	// System Overview
	fmt.Fprintf(r.StatsView, "[yellow]═══ REMOTE SYSTEM: %s ═══[-]\n\n", stats.Hostname)

	if stats.Stats != nil {
		fmt.Fprintf(r.StatsView, "[cyan]CPU Usage:[-]     [white]%.1f%%[-]\n", stats.Stats.CPUUsage)
		fmt.Fprintf(r.StatsView, "[cyan]Memory Usage:[-]  [white]%.1f%%[-] ([white]%s[-] / [white]%s[-])\n",
			stats.Stats.MemUsage, formatBytes(stats.Stats.MemUsed), formatBytes(stats.Stats.MemTotal))
		fmt.Fprintf(r.StatsView, "[cyan]Uptime:[-]        [white]%s[-]\n", formatUptime(stats.Stats.Uptime))
		fmt.Fprintf(r.StatsView, "[cyan]Load Average:[-]  [white]%.2f, %.2f, %.2f[-]\n\n",
			stats.Stats.LoadAvg1, stats.Stats.LoadAvg5, stats.Stats.LoadAvg15)
	}

	// Processes
	if len(stats.Processes) > 0 {
		fmt.Fprintf(r.StatsView, "[yellow]═══ TOP PROCESSES ═══[-]\n")
		fmt.Fprintf(r.StatsView, "[gray]%-8s %-6s %-6s %-20s[-]\n", "PID", "CPU%", "MEM%", "NAME")
		count := 0
		for _, proc := range stats.Processes {
			if count >= 10 {
				break
			}
			fmt.Fprintf(r.StatsView, "[white]%-8d %-6.1f %-6.1f %-20s[-]\n",
				proc.PID, proc.CPUUsage, float64(proc.MemUsage), truncate(proc.Name, 20))
			count++
		}
		fmt.Fprintf(r.StatsView, "\n")
	}

	// Disks
	if len(stats.Disks) > 0 {
		fmt.Fprintf(r.StatsView, "[yellow]═══ DISK USAGE ═══[-]\n")
		fmt.Fprintf(r.StatsView, "[gray]%-20s %-10s %-10s %-6s[-]\n", "DEVICE", "USED", "TOTAL", "USE%")
		for _, disk := range stats.Disks {
			fmt.Fprintf(r.StatsView, "[white]%-20s %-10s %-10s %-6.1f%%[-]\n",
				truncate(disk.Device, 20), formatBytes(disk.Used), formatBytes(disk.Total), disk.UsedPercent)
		}
		fmt.Fprintf(r.StatsView, "\n")
	}

	// Network
	if len(stats.Nets) > 0 {
		fmt.Fprintf(r.StatsView, "[yellow]═══ NETWORK ═══[-]\n")
		fmt.Fprintf(r.StatsView, "[gray]%-15s %-12s %-12s[-]\n", "INTERFACE", "RX", "TX")
		for _, net := range stats.Nets {
			fmt.Fprintf(r.StatsView, "[white]%-15s %-12s %-12s[-]\n",
				truncate(net.Name, 15), formatBytes(net.BytesRecv), formatBytes(net.BytesSent))
		}
		fmt.Fprintf(r.StatsView, "\n")
	}

	// Temperatures
	if len(stats.Temps) > 0 {
		fmt.Fprintf(r.StatsView, "[yellow]═══ TEMPERATURES ═══[-]\n")
		for _, temp := range stats.Temps {
			color := "white"
			if temp.Temperature > 80 {
				color = "red"
			} else if temp.Temperature > 60 {
				color = "yellow"
			}
			fmt.Fprintf(r.StatsView, "[cyan]%-20s[-] [%s]%.1f°C[-]\n", temp.Key, color, temp.Temperature)
		}
		fmt.Fprintf(r.StatsView, "\n")
	}

	// GPUs
	if len(stats.GPUs) > 0 {
		fmt.Fprintf(r.StatsView, "[yellow]═══ GPU ═══[-]\n")
		for _, gpu := range stats.GPUs {
			fmt.Fprintf(r.StatsView, "[cyan]%s[-]\n", gpu.Name)
			fmt.Fprintf(r.StatsView, "  Usage: [white]%d%%[-] | Temp: [white]%d°C[-] | Memory: [white]%d MiB[-] / [white]%d MiB[-]\n",
				gpu.Utilization, gpu.Temp, gpu.MemoryUsed, gpu.MemoryTotal)
		}
		fmt.Fprintf(r.StatsView, "\n")
	}

	// Docker Containers
	if len(stats.Containers) > 0 {
		fmt.Fprintf(r.StatsView, "[yellow]═══ DOCKER CONTAINERS ═══[-]\n")
		fmt.Fprintf(r.StatsView, "[gray]%-20s %-15s[-]\n", "NAME", "IMAGE")
		for _, container := range stats.Containers {
			statusColor := "green"
			if container.Status != "running" {
				statusColor = "gray"
			}
			fmt.Fprintf(r.StatsView, "[%s]%-20s[-] [white]%-15s[-]\n",
				statusColor, truncate(container.Name, 20), truncate(container.Image, 15))
		}
	}
}

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func formatUptime(seconds uint64) string {
	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	minutes := (seconds % 3600) / 60
	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
