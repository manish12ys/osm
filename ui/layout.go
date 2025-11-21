package ui

import (
	"github.com/rivo/tview"
)

// AppLayout holds the main application layout
type AppLayout struct {
	App       *tview.Application
	Grid      *tview.Grid
	Header    *HeaderComponent
	Footer    *tview.TextView
	Pages     *tview.Pages
	ProcTable *ProcTableComponent
	DiskTable *DiskTableComponent
	NetTable  *NetTableComponent
	Temp      *TempComponent
	CPUDetail *CPUDetailComponent
	MemDetail *MemDetailComponent
}

// NewAppLayout initializes the TUI layout
func NewAppLayout() *AppLayout {
	app := tview.NewApplication()

	// Header
	headerComp := NewHeaderComponent()
	headerComp.View.SetTextAlign(tview.AlignCenter).
		SetText("Omarchy System Monitor")
	headerComp.View.SetBorder(true)

	// Footer
	footer := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetText(" q: Quit | 1: Procs | 2: Disks | 3: Net | 4: Temp | 5: CPU | 6: Mem | k: Kill | t: Theme ")
	footer.SetBorder(true)

	// Components
	procTableComp := NewProcTableComponent()
	procTableComp.Table.SetBorder(true).SetTitle(" Processes (1) ")

	diskTableComp := NewDiskTableComponent()
	diskTableComp.Table.SetBorder(true).SetTitle(" Disk Usage (2) ")

	netTableComp := NewNetTableComponent()
	netTableComp.Table.SetBorder(true).SetTitle(" Network (3) ")

	tempComp := NewTempComponent()

	cpuDetailComp := NewCPUDetailComponent()

	memDetailComp := NewMemDetailComponent()

	// Pages
	pages := tview.NewPages()
	pages.AddPage("procs", procTableComp.Table, true, true)
	pages.AddPage("disks", diskTableComp.Table, true, false)
	pages.AddPage("net", netTableComp.Table, true, false)
	pages.AddPage("temp", tempComp.Flex, true, false)
	pages.AddPage("cpu", cpuDetailComp.Flex, true, false)
	pages.AddPage("mem", memDetailComp.Flex, true, false)

	// Grid Layout
	grid := tview.NewGrid().
		SetRows(5, 0, 3). // Initial state: SearchBar hidden (height 0)
		SetColumns(0).
		SetBorders(false).
		AddItem(headerComp.View, 0, 0, 1, 1, 0, 0, false).
		AddItem(pages, 1, 0, 1, 1, 0, 0, true).
		AddItem(footer, 2, 0, 1, 1, 0, 0, false)

	l := &AppLayout{
		App:       app,
		Grid:      grid,
		Header:    headerComp,
		Footer:    footer,
		Pages:     pages,
		ProcTable: procTableComp,
		DiskTable: diskTableComp,
		NetTable:  netTableComp,
		Temp:      tempComp,
		CPUDetail: cpuDetailComp,
		MemDetail: memDetailComp,
	}
	l.ApplyTheme()
	return l
}

// ApplyTheme applies the current theme to all components
func (l *AppLayout) ApplyTheme() {
	l.Header.View.SetBorderColor(CurrentTheme.Border)
	l.Header.View.SetTitleColor(CurrentTheme.HeaderTitle)

	l.Footer.SetBorderColor(CurrentTheme.Border)
	l.Footer.SetTextColor(CurrentTheme.Foreground)

	l.ProcTable.ApplyTheme()
	l.DiskTable.ApplyTheme()
	l.NetTable.ApplyTheme()
	l.Temp.ApplyTheme()
	l.CPUDetail.ApplyTheme()
	l.MemDetail.ApplyTheme()
}

// Run starts the application
func (l *AppLayout) Run() error {
	l.App.SetRoot(l.Grid, true).SetFocus(l.ProcTable.Table)
	return l.App.Run()
}
