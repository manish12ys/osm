package ui

import (
	"fmt"
	"time"

	"github.com/rivo/tview"
)

// AppLayout holds the main application layout
type AppLayout struct {
	App         *tview.Application
	Flex        *tview.Flex
	MiddleFlex  *tview.Flex
	Header      *HeaderComponent
	SearchBar   *tview.InputField
	Footer      *tview.TextView
	Pages       *tview.Pages
	ProcTable   *ProcTableComponent
	DiskTable   *DiskTableComponent
	NetTable    *NetTableComponent
	Temp        *TempComponent
	GPU         *GPUComponent
	DockerTable *DockerTableComponent
	Plugins     *PluginsComponent
	Remote      *RemoteComponent
	CPUDetail   *CPUDetailComponent
	MemDetail   *MemDetailComponent
	Help        *HelpComponent
	ActivePage  string
}

// AdjustLayout adjusts the layout based on terminal size
func (l *AppLayout) AdjustLayout() {
	// Get screen dimensions from the root primitive
	_, _, width, height := l.Flex.GetRect()

	// If not yet drawn, use default values
	if width == 0 || height == 0 {
		return
	}

	// Adjust header height based on terminal height
	headerHeight := 5
	footerHeight := 3

	if height < 20 {
		// Very small terminal - minimize header and footer
		headerHeight = 3
		footerHeight = 2
	} else if height < 30 {
		// Small terminal - compact header
		headerHeight = 4
		footerHeight = 2
	}

	// Rebuild flex with adjusted heights
	l.Flex.Clear()
	l.Flex.AddItem(l.Header.View, headerHeight, 0, false).
		AddItem(l.MiddleFlex, 0, 1, true).
		AddItem(l.Footer, footerHeight, 0, false)

	// Adjust footer text based on width
	if width < 80 {
		// Compact footer for narrow terminals
		l.Footer.SetText(" q:Quit | 1-9,0:Views | t:Theme | /:Search | ?:Help ")
	} else if width < 120 {
		// Medium footer
		l.Footer.SetText(" q:Quit | 1:Procs 2:Disks 3:Net 4:Temp 5:GPU 6:Docker 7:CPU 8:Mem 9:Plugins 0:Remote | t:Theme | ?:Help ")
	} else {
		// Full footer
		l.Footer.SetText(" q: Quit | 1: Procs | 2: Disks | 3: Net | 4: Temp | 5: GPU | 6: Docker | 7: CPU | 8: Mem | 9: Plugins | 0: Remote | k: Kill | t: Theme | /: Search | e: Export | ?: Help ")
	}
}

// NewAppLayout initializes the TUI layout
func NewAppLayout(remoteURL string) *AppLayout {
	app := tview.NewApplication()

	// Header
	headerComp := NewHeaderComponent()
	headerComp.View.SetTextAlign(tview.AlignCenter).
		SetText("Omarchy System Monitor")
	headerComp.View.SetBorder(true)

	// Search Bar
	searchBar := tview.NewInputField().
		SetLabel("Search: ").
		SetFieldWidth(0).
		SetFieldBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
	searchBar.SetBorder(true)

	// Footer
	footer := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetText(" q: Quit | 1: Procs | 2: Disks | 3: Net | 4: Temp | 5: GPU | 6: Docker | 7: CPU | 8: Mem | 9: Plugins | 0: Remote | k: Kill | t: Theme | /: Search | e: Export | ?: Help ")
	footer.SetBorder(true)

	// Components
	procTableComp := NewProcTableComponent()
	procTableComp.Table.SetBorder(true).SetTitle(" Processes (1) ")

	diskTableComp := NewDiskTableComponent()
	diskTableComp.Table.SetBorder(true).SetTitle(" Disk Usage (2) ")

	netTableComp := NewNetTableComponent()
	netTableComp.Table.SetBorder(true).SetTitle(" Network (3) ")

	tempComp := NewTempComponent()

	gpuComp := NewGPUComponent()
	gpuComp.Table.SetBorder(true).SetTitle(" GPU Monitor (5) ")

	dockerTableComp := NewDockerTableComponent()
	dockerTableComp.Table.SetBorder(true).SetTitle(" Docker Containers (6) ")

	pluginsComp := NewPluginsComponent()

	remoteComp := NewRemoteComponent(remoteURL)

	cpuDetailComp := NewCPUDetailComponent()

	memDetailComp := NewMemDetailComponent()

	helpComp := NewHelpComponent()

	// Pages
	pages := tview.NewPages()
	pages.AddPage("procs", procTableComp.Table, true, true)
	pages.AddPage("disks", diskTableComp.Table, true, false)
	pages.AddPage("net", netTableComp.Table, true, false)
	pages.AddPage("temp", tempComp.Flex, true, false)
	pages.AddPage("gpu", gpuComp.Table, true, false)
	pages.AddPage("docker", dockerTableComp.Table, true, false)
	pages.AddPage("plugins", pluginsComp.Flex, true, false)
	pages.AddPage("remote", remoteComp.Flex, true, false)
	pages.AddPage("cpu", cpuDetailComp.Flex, true, false)
	pages.AddPage("mem", memDetailComp.Flex, true, false)
	pages.AddPage("help", helpComp.Flex, true, false)

	// Middle Flex (Pages + SearchBar)
	middleFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(pages, 0, 1, true)
	// SearchBar is NOT added initially

	// Main Flex Layout
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(headerComp.View, 5, 0, false).
		AddItem(middleFlex, 0, 1, true).
		AddItem(footer, 3, 0, false)

	l := &AppLayout{
		App:         app,
		Flex:        flex,
		MiddleFlex:  middleFlex,
		Header:      headerComp,
		SearchBar:   searchBar,
		Footer:      footer,
		Pages:       pages,
		ProcTable:   procTableComp,
		DiskTable:   diskTableComp,
		NetTable:    netTableComp,
		Temp:        tempComp,
		GPU:         gpuComp,
		DockerTable: dockerTableComp,
		Plugins:     pluginsComp,
		Remote:      remoteComp,
		CPUDetail:   cpuDetailComp,
		MemDetail:   memDetailComp,
		Help:        helpComp,
	}

	// Set help done func to close help
	// Set help done func to close help
	helpComp.SetDoneFunc(func() {
		l.ToggleHelp(false)
	})

	l.ApplyTheme()
	return l
}

// ToggleSearchBar shows or hides the search bar
func (l *AppLayout) ToggleSearchBar(show bool) {
	if show {
		// Ensure it's not added twice by removing first (safe op)
		l.MiddleFlex.RemoveItem(l.SearchBar)
		l.MiddleFlex.AddItem(l.SearchBar, 3, 0, false)
		l.App.SetFocus(l.SearchBar)
	} else {
		l.MiddleFlex.RemoveItem(l.SearchBar)
		l.App.SetFocus(l.ProcTable.Table)
		l.SearchBar.SetText("") // Clear text when hiding
	}
}

// ToggleHelp shows or hides the help modal
func (l *AppLayout) ToggleHelp(show bool) {
	if show {
		l.Pages.ShowPage("help")
		l.App.SetFocus(l.Help.Table)
	} else {
		l.Pages.HidePage("help")
		// Restore focus to active page
		switch l.ActivePage {
		case "procs":
			l.App.SetFocus(l.ProcTable.Table)
		case "disks":
			l.App.SetFocus(l.DiskTable.Table)
		case "net":
			l.App.SetFocus(l.NetTable.Table)
		case "temp":
			l.App.SetFocus(l.Temp.Flex)
		case "gpu":
			l.App.SetFocus(l.GPU.Table)
		case "docker":
			l.App.SetFocus(l.DockerTable.Table)
		case "cpu":
			l.App.SetFocus(l.CPUDetail.Flex)
		case "mem":
			l.App.SetFocus(l.MemDetail.Flex)
		case "plugins":
			l.App.SetFocus(l.Plugins.Flex)
		case "remote":
			l.App.SetFocus(l.Remote.Flex)
		default:
			l.App.SetFocus(l.ProcTable.Table)
		}
	}
}

// ShowFlashMessage displays a temporary message in the footer
func (l *AppLayout) ShowFlashMessage(msg string) {
	originalText := l.Footer.GetText(true)
	l.Footer.SetText(fmt.Sprintf(" [yellow]%s[-]", msg))

	go func() {
		time.Sleep(3 * time.Second)
		l.App.QueueUpdateDraw(func() {
			l.Footer.SetText(originalText)
		})
	}()
}

// ApplyTheme applies the current theme to all components
func (l *AppLayout) ApplyTheme() {
	l.Header.View.SetBorderColor(CurrentTheme.Border)
	l.Header.View.SetTitleColor(CurrentTheme.HeaderTitle)

	l.SearchBar.SetBorderColor(CurrentTheme.Border)
	l.SearchBar.SetLabelColor(CurrentTheme.HeaderTitle)
	l.SearchBar.SetFieldTextColor(CurrentTheme.Foreground)

	l.Footer.SetBorderColor(CurrentTheme.Border)
	l.Footer.SetTextColor(CurrentTheme.Foreground)

	l.ProcTable.ApplyTheme()
	l.DiskTable.ApplyTheme()
	l.NetTable.ApplyTheme()
	l.Temp.ApplyTheme()
	l.GPU.ApplyTheme()
	l.DockerTable.ApplyTheme()
	l.Plugins.ApplyTheme()
	l.Remote.ApplyTheme()
	l.CPUDetail.ApplyTheme()
	l.MemDetail.ApplyTheme()
	l.Help.ApplyTheme()
}

// Run starts the application
func (l *AppLayout) Run() error {
	l.App.SetRoot(l.Flex, true).SetFocus(l.ProcTable.Table)
	return l.App.Run()
}
