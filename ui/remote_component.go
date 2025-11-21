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
	// Helper to get color tag
	cLow := fmt.Sprintf("[#%06x]", CurrentTheme.LowUsage.Hex())
	cMed := fmt.Sprintf("[#%06x]", CurrentTheme.MedUsage.Hex())
	cHigh := fmt.Sprintf("[#%06x]", CurrentTheme.HighUsage.Hex())
	cFore := fmt.Sprintf("[#%06x]", CurrentTheme.Foreground.Hex())
	cDim := "[darkgray]"
	cLabel := fmt.Sprintf("[#%06x]", CurrentTheme.HeaderValue.Hex())

	if r.Client == nil {
		r.StatusView.Clear()
		fmt.Fprintf(r.StatusView, "%s✗ Not Configured[-]", cHigh)
		r.StatsView.Clear()
		fmt.Fprintf(r.StatsView, "%sRemote monitoring not configured[-]\n\n", cMed)
		fmt.Fprintf(r.StatsView, "To enable remote monitoring, add to config.yaml:\n\n")
		fmt.Fprintf(r.StatsView, "%sremote_mode: true[-]\n", cLabel)
		fmt.Fprintf(r.StatsView, "%sremote_url: \"http://192.168.1.100:8080\"[-]\n", cLabel)
		return
	}

	stats, err := r.Client.FetchRemoteStats()
	if err != nil {
		r.Connected = false
		r.StatusView.Clear()
		fmt.Fprintf(r.StatusView, "%s✗ Disconnected[-] - %v", cHigh, err)
		r.StatsView.Clear()
		fmt.Fprintf(r.StatsView, "%sFailed to connect to remote server[-]\n\n", cHigh)
		fmt.Fprintf(r.StatsView, "%sURL: %s[-]\n", cDim, r.Client.BaseURL)
		fmt.Fprintf(r.StatsView, "%sError: %v[-]\n", cDim, err)
		return
	}

	r.Connected = true
	r.LastUpdate = time.Now()

	// Update status
	r.StatusView.Clear()
	fmt.Fprintf(r.StatusView, "%s✓ Connected[-] to %s%s[-] | Last Update: %s%s[-]",
		cLow, cLabel, stats.Hostname, cFore, r.LastUpdate.Format("15:04:05"))

	// Update stats view
	r.StatsView.Clear()

	// System Overview
	fmt.Fprintf(r.StatsView, "%s═══ REMOTE SYSTEM: %s ═══[-]\n\n", cMed, stats.Hostname)

	if stats.Stats != nil {
		fmt.Fprintf(r.StatsView, "%sCPU Usage:[-]     %s%.1f%%[-]\n", cLabel, cFore, stats.Stats.CPUUsage)
		fmt.Fprintf(r.StatsView, "%sMemory Usage:[-]  %s%.1f%%[-] (%s%s[-] / %s%s[-])\n",
			cLabel, cFore, stats.Stats.MemUsage, cFore, formatBytes(stats.Stats.MemUsed), cFore, formatBytes(stats.Stats.MemTotal))
		fmt.Fprintf(r.StatsView, "%sUptime:[-]        %s%s[-]\n", cLabel, cFore, formatUptime(stats.Stats.Uptime))
		fmt.Fprintf(r.StatsView, "%sLoad Average:[-]  %s%.2f, %.2f, %.2f[-]\n\n",
			cLabel, cFore, stats.Stats.LoadAvg1, stats.Stats.LoadAvg5, stats.Stats.LoadAvg15)
	}

	// Processes
	if len(stats.Processes) > 0 {
		fmt.Fprintf(r.StatsView, "%s═══ TOP PROCESSES ═══[-]\n", cMed)
		fmt.Fprintf(r.StatsView, "%s%-8s %-6s %-6s %-20s[-]\n", cDim, "PID", "CPU%", "MEM%", "NAME")
		count := 0
		for _, proc := range stats.Processes {
			if count >= 10 {
				break
			}
			fmt.Fprintf(r.StatsView, "%s%-8d %-6.1f %-6.1f %-20s[-]\n",
				cFore, proc.PID, proc.CPUUsage, float64(proc.MemUsage), truncate(proc.Name, 20))
			count++
		}
		fmt.Fprintf(r.StatsView, "\n")
	}

	// Disks
	if len(stats.Disks) > 0 {
		fmt.Fprintf(r.StatsView, "%s═══ DISK USAGE ═══[-]\n", cMed)
		fmt.Fprintf(r.StatsView, "%s%-20s %-10s %-10s %-6s[-]\n", cDim, "DEVICE", "USED", "TOTAL", "USE%")
		for _, disk := range stats.Disks {
			fmt.Fprintf(r.StatsView, "%s%-20s %-10s %-10s %-6.1f%%[-]\n",
				cFore, truncate(disk.Device, 20), formatBytes(disk.Used), formatBytes(disk.Total), disk.UsedPercent)
		}
		fmt.Fprintf(r.StatsView, "\n")
	}

	// Network
	if len(stats.Nets) > 0 {
		fmt.Fprintf(r.StatsView, "%s═══ NETWORK ═══[-]\n", cMed)
		fmt.Fprintf(r.StatsView, "%s%-15s %-12s %-12s[-]\n", cDim, "INTERFACE", "RX", "TX")
		for _, net := range stats.Nets {
			fmt.Fprintf(r.StatsView, "%s%-15s %-12s %-12s[-]\n",
				cFore, truncate(net.Name, 15), formatBytes(net.BytesRecv), formatBytes(net.BytesSent))
		}
		fmt.Fprintf(r.StatsView, "\n")
	}

	// Temperatures
	if len(stats.Temps) > 0 {
		fmt.Fprintf(r.StatsView, "%s═══ TEMPERATURES ═══[-]\n", cMed)
		for _, temp := range stats.Temps {
			color := cFore
			if temp.Temperature > 80 {
				color = cHigh
			} else if temp.Temperature > 60 {
				color = cMed
			}
			fmt.Fprintf(r.StatsView, "%s%-20s[-] [%s]%.1f°C[-]\n", cLabel, temp.Key, color, temp.Temperature)
		}
		fmt.Fprintf(r.StatsView, "\n")
	}

	// GPUs
	if len(stats.GPUs) > 0 {
		fmt.Fprintf(r.StatsView, "%s═══ GPU ═══[-]\n", cMed)
		for _, gpu := range stats.GPUs {
			fmt.Fprintf(r.StatsView, "%s%s[-]\n", cLabel, gpu.Name)
			fmt.Fprintf(r.StatsView, "  Usage: %s%d%%[-] | Temp: %s%d°C[-] | Memory: %s%d MiB[-] / %s%d MiB[-]\n",
				cFore, gpu.Utilization, cFore, gpu.Temp, cFore, gpu.MemoryUsed, cFore, gpu.MemoryTotal)
		}
		fmt.Fprintf(r.StatsView, "\n")
	}

	// Docker Containers
	if len(stats.Containers) > 0 {
		fmt.Fprintf(r.StatsView, "%s═══ DOCKER CONTAINERS ═══[-]\n", cMed)
		fmt.Fprintf(r.StatsView, "%s%-20s %-15s[-]\n", cDim, "NAME", "IMAGE")
		for _, container := range stats.Containers {
			statusColor := cLow
			if container.Status != "running" {
				statusColor = cDim
			}
			fmt.Fprintf(r.StatsView, "[%s]%-20s[-] %s%-15s[-]\n",
				statusColor, truncate(container.Name, 20), cFore, truncate(container.Image, 15))
		}
	}
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
