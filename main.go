package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"omarchy-monitor/config"
	"omarchy-monitor/core"
	"omarchy-monitor/plugins"
	"omarchy-monitor/remote"
	"omarchy-monitor/ui"

	"github.com/gdamore/tcell/v2"
)

func main() {
	// Load Config
	cfg, err := config.LoadConfig()
	if err == nil {
		ui.SetTheme(cfg.Theme)
	}

	layout := ui.NewAppLayout()

	// Enable mouse support
	layout.App.EnableMouse(true)

	currentSort := core.SortCPU
	activePage := "procs"
	layout.ActivePage = activePage
	if cfg != nil && cfg.DefaultPage != "" {
		activePage = cfg.DefaultPage
		layout.ActivePage = activePage
		layout.Pages.SwitchToPage(activePage)
		// Focus logic needs to be updated if we switch page initially
		// But let's keep it simple for now or handle it below
	}

	// Start remote monitoring server if enabled
	if cfg != nil && cfg.RemoteServer {
		server := remote.NewServer(cfg.RemotePort)
		go func() {
			if err := server.Start(); err != nil {
				fmt.Printf("Remote server error: %v\n", err)
			}
		}()
	}

	// Initialize plugin manager if enabled
	var pluginManager *plugins.Manager
	if cfg != nil && cfg.PluginsEnabled {
		home, _ := os.UserHomeDir()
		pluginsDir := filepath.Join(home, ".config", "osm", "plugins")
		pluginManager = plugins.NewManager(pluginsDir)
		pluginManager.LoadPlugins()

		// Create example plugins if directory is empty
		if len(pluginManager.Plugins) == 0 {
			plugins.CreateExamplePlugins(pluginsDir)
			pluginManager.LoadPlugins()
		}
	}

	// History buffers for graphs
	cpuHistory := core.NewHistory(60)
	memHistory := core.NewHistory(60)

	// Per-core CPU tracking
	var perCoreUsage []float64

	// Cache last processes for live filtering
	var lastProcs []core.Process

	// Set up sort callback for mouse clicks on process table headers
	layout.ProcTable.SortCallback = func(sortType int) {
		currentSort = core.SortType(sortType)
		if lastProcs != nil {
			core.SortProcesses(lastProcs, currentSort)
			layout.ProcTable.Update(lastProcs)
		}
	}

	// Search Bar Handler
	layout.SearchBar.SetChangedFunc(func(text string) {
		layout.ProcTable.Filter = text
		if activePage == "procs" && lastProcs != nil {
			layout.ProcTable.Update(lastProcs)
		}
	})

	// Input handling
	layout.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// If search bar is focused, handle specific keys
		if layout.SearchBar.HasFocus() {
			if event.Key() == tcell.KeyEsc {
				layout.ToggleSearchBar(false)
				layout.ProcTable.Filter = "" // Clear filter
				if lastProcs != nil {
					layout.ProcTable.Update(lastProcs)
				}
				return nil
			}
			if event.Key() == tcell.KeyEnter {
				layout.ToggleSearchBar(false)
				return nil
			}
			return event // Let InputField handle typing
		}

		switch event.Rune() {
		case 'q':
			layout.App.Stop()
		case '/':
			if activePage == "procs" {
				layout.ToggleSearchBar(true)
				return nil
			}
		case '1':
			layout.Pages.SwitchToPage("procs")
			activePage = "procs"
			layout.ActivePage = activePage
			layout.App.SetFocus(layout.ProcTable.Table)
		case '2':
			layout.Pages.SwitchToPage("disks")
			activePage = "disks"
			layout.ActivePage = activePage
			layout.App.SetFocus(layout.DiskTable.Table)
		case '3':
			layout.Pages.SwitchToPage("net")
			activePage = "net"
			layout.ActivePage = activePage
			layout.App.SetFocus(layout.NetTable.Table)
		case '4':
			layout.Pages.SwitchToPage("temp")
			activePage = "temp"
			layout.ActivePage = activePage
			layout.App.SetFocus(layout.Temp.Flex)
		case '5':
			layout.Pages.SwitchToPage("gpu")
			activePage = "gpu"
			layout.ActivePage = activePage
			layout.App.SetFocus(layout.GPU.Table)
		case '6':
			layout.Pages.SwitchToPage("docker")
			activePage = "docker"
			layout.ActivePage = activePage
			layout.App.SetFocus(layout.DockerTable.Table)
		case '7':
			layout.Pages.SwitchToPage("cpu")
			activePage = "cpu"
			layout.ActivePage = activePage
			layout.App.SetFocus(layout.CPUDetail.Flex)
		case '8':
			layout.Pages.SwitchToPage("mem")
			activePage = "mem"
			layout.ActivePage = activePage
			layout.App.SetFocus(layout.MemDetail.Flex)
		case '9':
			layout.Pages.SwitchToPage("plugins")
			activePage = "plugins"
			layout.ActivePage = activePage
			layout.App.SetFocus(layout.Plugins.Flex)
		case 't':
			ui.CycleTheme()
			layout.ApplyTheme()
		case 'T':
			layout.ProcTable.ShowTree = !layout.ProcTable.ShowTree
			if activePage == "procs" && lastProcs != nil {
				layout.ProcTable.Update(lastProcs)
			}
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
		case 'e':
			// Trigger export
			stats, _ := core.FetchStats()
			procs, _ := core.FetchProcesses()
			filename, err := core.ExportSnapshot(stats, procs)
			if err != nil {
				layout.ShowFlashMessage(fmt.Sprintf("Export Failed: %v", err))
			} else {
				layout.ShowFlashMessage(fmt.Sprintf("Exported to %s", filename))
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

			// Adjust layout for responsive design
			layout.AdjustLayout()

			// Fetch per-core CPU for detail page
			if activePage == "cpu" {
				perCoreUsage, _ = core.FetchPerCoreCPU()
			}

			// Fetch data based on active page
			var procs []core.Process
			var disks []core.DiskStat
			var nets []core.NetStat
			var temps []core.TempStat
			var gpus []core.GPUStat
			var containers []core.ContainerStat
			var pluginResults map[string]*plugins.PluginOutput

			if activePage == "procs" {
				var err error
				procs, err = core.FetchProcesses()
				if err != nil {
					// Handle error
				}
				core.SortProcesses(procs, currentSort)
				lastProcs = procs // Update cache
			} else if activePage == "disks" {
				disks, _ = core.FetchDiskStats()
			} else if activePage == "net" {
				nets, _ = core.FetchNetStats()
			} else if activePage == "temp" {
				// Fetch temps
				temps, _ = core.FetchTemps()
			} else if activePage == "gpu" {
				gpus, _ = core.FetchGPUStats()
			} else if activePage == "docker" {
				containers, _ = core.FetchContainers()
			} else if activePage == "plugins" && pluginManager != nil {
				pluginResults = pluginManager.ExecuteAll()
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
				} else if activePage == "gpu" {
					layout.GPU.Update(gpus)
				} else if activePage == "docker" {
					layout.DockerTable.Update(containers)
				} else if activePage == "plugins" && pluginResults != nil {
					layout.Plugins.Update(pluginResults)
				} else if activePage == "cpu" {
					layout.CPUDetail.Update(cpuHistory, perCoreUsage)
				} else if activePage == "mem" {
					layout.MemDetail.Update(memHistory, stats)
				}
			})

			sleepDuration := 1000 * time.Millisecond
			if cfg != nil && cfg.RefreshRate > 0 {
				sleepDuration = time.Duration(cfg.RefreshRate) * time.Millisecond
			}
			time.Sleep(sleepDuration)
		}
	}()

	if err := layout.Run(); err != nil {
		fmt.Printf("Error running application: %v\n", err)
		os.Exit(1)
	}
}
