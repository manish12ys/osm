package main

import (
	"fmt"
	"os"
	"time"

	"omarchy-monitor/core"
	"omarchy-monitor/ui"

	"github.com/gdamore/tcell/v2"
)

func main() {
	layout := ui.NewAppLayout()

	currentSort := core.SortCPU
	activePage := "procs"

	// History buffers for graphs
	cpuHistory := core.NewHistory(60)
	memHistory := core.NewHistory(60)

	// Per-core CPU tracking
	var perCoreUsage []float64

	// Input handling
	layout.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'q':
			layout.App.Stop()
		case '1':
			layout.Pages.SwitchToPage("procs")
			activePage = "procs"
			layout.App.SetFocus(layout.ProcTable.Table)
		case '2':
			layout.Pages.SwitchToPage("disks")
			activePage = "disks"
			layout.App.SetFocus(layout.DiskTable.Table)
		case '3':
			layout.Pages.SwitchToPage("net")
			activePage = "net"
			layout.App.SetFocus(layout.NetTable.Table)
		case '4':
			layout.Pages.SwitchToPage("temp")
			activePage = "temp"
			layout.App.SetFocus(layout.Temp.Flex)
		case '5':
			layout.Pages.SwitchToPage("cpu")
			activePage = "cpu"
			layout.App.SetFocus(layout.CPUDetail.Flex)
		case '6':
			layout.Pages.SwitchToPage("mem")
			activePage = "mem"
			layout.App.SetFocus(layout.MemDetail.Flex)
		case 't':
			ui.CycleTheme()
			layout.ApplyTheme()
		case 'c':
			currentSort = core.SortCPU
		case 'm':
			currentSort = core.SortMem
		case 'p':
			currentSort = core.SortPID
		case 'n':
			currentSort = core.SortName
		case 'k':
			if activePage == "procs" {
				pid, err := layout.ProcTable.GetSelectedPID()
				if err == nil {
					core.KillProcess(pid)
				}
			}
		}
		return event
	})

	// Stats loop
	go func() {
		for {
			stats, err := core.FetchStats()
			if err != nil {
				continue
			}

			// Add to history
			cpuHistory.Add(stats.CPUUsage)
			memHistory.Add(stats.MemUsage)

			// Fetch per-core CPU for detail page
			if activePage == "cpu" {
				perCoreUsage, _ = core.FetchPerCoreCPU()
			}

			// Fetch data based on active page
			var procs []core.Process
			var disks []core.DiskStat
			var nets []core.NetStat
			var temps []core.TempStat

			if activePage == "procs" {
				var err error
				procs, err = core.FetchProcesses()
				if err != nil {
					// Handle error
				}
				core.SortProcesses(procs, currentSort)
			} else if activePage == "disks" {
				disks, _ = core.FetchDiskStats()
			} else if activePage == "net" {
				nets, _ = core.FetchNetStats()
			} else if activePage == "temp" {
				// Fetch temps
				temps, _ = core.FetchTemps()
			}

			layout.App.QueueUpdateDraw(func() {
				layout.Header.Update(stats)

				if activePage == "procs" && procs != nil {
					layout.ProcTable.Update(procs)
				} else if activePage == "disks" && disks != nil {
					layout.DiskTable.Update(disks)
				} else if activePage == "net" && nets != nil {
					layout.NetTable.Update(nets)
				} else if activePage == "temp" {
					layout.Temp.Update(temps)
				} else if activePage == "cpu" {
					layout.CPUDetail.Update(cpuHistory, perCoreUsage)
				} else if activePage == "mem" {
					layout.MemDetail.Update(memHistory, stats)
				}
			})

			time.Sleep(1 * time.Second)
		}
	}()

	if err := layout.Run(); err != nil {
		fmt.Printf("Error running application: %v\n", err)
		os.Exit(1)
	}
}
