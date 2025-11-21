package ui

import (
	"fmt"
	"strings"

	"omarchy-monitor/core"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ProcTableComponent manages the process table view
type ProcTableComponent struct {
	Table        *tview.Table
	Filter       string
	ShowTree     bool
	SortCallback func(sortType int) // Callback to trigger sorting
}

// NewProcTableComponent creates a new process table
func NewProcTableComponent() *ProcTableComponent {
	t := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false). // Select rows, not individual cells
		SetFixed(1, 0)              // Fix the header row

	// Set Headers
	headers := []string{"PID", "Name", "User", "CPU%", "Mem%", "State"}
	alignments := []int{tview.AlignRight, tview.AlignLeft, tview.AlignLeft, tview.AlignRight, tview.AlignRight, tview.AlignCenter}
	for i, h := range headers {
		t.SetCell(0, i, tview.NewTableCell(h).
			SetTextColor(tcell.ColorYellow).
			SetSelectable(false).
			SetAlign(alignments[i]))
	}

	comp := &ProcTableComponent{
		Table: t,
	}

	// Enable mouse support for header clicks
	t.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		if action == tview.MouseLeftClick {
			clickRow, clickCol := event.Position()

			// Check if clicked on header row
			x, y, _, _ := t.GetInnerRect()
			if clickRow == y && comp.SortCallback != nil {
				// Map column to sort type
				switch clickCol - x {
				case 0: // PID column area
					comp.SortCallback(2) // SortPID
				case 1: // Name column area
					comp.SortCallback(3) // SortName
				case 3: // CPU% column area
					comp.SortCallback(0) // SortCPU
				case 4: // Mem% column area
					comp.SortCallback(1) // SortMem
				}
				return action, nil
			}
		}
		return action, event
	})

	return comp
}

// ApplyTheme updates the table style
func (p *ProcTableComponent) ApplyTheme() {
	p.Table.SetBorderColor(CurrentTheme.Border)
	p.Table.SetTitleColor(CurrentTheme.HeaderTitle)

	// Update header row color
	for i := 0; i < p.Table.GetColumnCount(); i++ {
		cell := p.Table.GetCell(0, i)
		if cell != nil {
			cell.SetTextColor(CurrentTheme.TableHead)
		}
	}
}

// TreeNode represents a node in the process tree
type TreeNode struct {
	Process  core.Process
	Children []*TreeNode
	Level    int
}

// Update refreshes the process list
func (p *ProcTableComponent) Update(procs []core.Process) {
	// We keep the header (row 0), so start clearing from row 1
	rowCount := p.Table.GetRowCount()

	var displayProcs []core.Process
	var displayNames []string // To store indented names

	if p.ShowTree {
		// Build Tree
		procMap := make(map[int32]*TreeNode)
		var roots []*TreeNode

		// 1. Create all nodes
		for _, proc := range procs {
			procMap[proc.PID] = &TreeNode{Process: proc}
		}

		// 2. Build hierarchy
		for _, proc := range procs {
			node := procMap[proc.PID]
			if parent, ok := procMap[proc.PPID]; ok && proc.PPID != 0 && proc.PPID != proc.PID {
				parent.Children = append(parent.Children, node)
			} else {
				roots = append(roots, node)
			}
		}

		// 3. Flatten tree to list
		var flatten func(node *TreeNode, level int)
		flatten = func(node *TreeNode, level int) {
			displayProcs = append(displayProcs, node.Process)
			indent := strings.Repeat("  ", level)
			prefix := ""
			if level > 0 {
				prefix = "└─ "
			}
			displayNames = append(displayNames, indent+prefix+node.Process.Name)

			for _, child := range node.Children {
				flatten(child, level+1)
			}
		}

		for _, root := range roots {
			flatten(root, 0)
		}

	} else {
		// Flat View (Filter logic applies here usually, but for now let's keep it simple or re-apply filter)
		// Filter processes
		if p.Filter != "" {
			filterLower := strings.ToLower(p.Filter)
			for _, proc := range procs {
				if strings.Contains(strings.ToLower(proc.Name), filterLower) ||
					strings.Contains(fmt.Sprintf("%d", proc.PID), filterLower) {
					displayProcs = append(displayProcs, proc)
					displayNames = append(displayNames, proc.Name)
				}
			}
		} else {
			displayProcs = procs
			for _, proc := range procs {
				displayNames = append(displayNames, proc.Name)
			}
		}
	}

	for i, proc := range displayProcs {
		row := i + 1 // Skip header

		// Color logic
		color := CurrentTheme.Foreground
		if proc.CPUUsage > 50 {
			color = CurrentTheme.HighUsage
		} else if proc.CPUUsage > 10 {
			color = CurrentTheme.MedUsage
		}

		p.setCellAligned(row, 0, fmt.Sprintf("%7d", proc.PID), color, tview.AlignRight)
		p.setCellAligned(row, 1, displayNames[i], CurrentTheme.Foreground, tview.AlignLeft)
		p.setCellAligned(row, 2, proc.User, CurrentTheme.HeaderValue, tview.AlignLeft)

		// CPU with mini bar
		cpuBar := makeMiniBar(proc.CPUUsage, 8)
		p.setCellAligned(row, 3, fmt.Sprintf("%5.1f%% %s", proc.CPUUsage, cpuBar), color, tview.AlignRight)

		// Memory with mini bar
		memColor := CurrentTheme.Foreground
		if proc.MemUsage > 10 {
			memColor = CurrentTheme.MedUsage
		}
		if proc.MemUsage > 50 {
			memColor = CurrentTheme.HighUsage
		}
		memBar := makeMiniBar(float64(proc.MemUsage), 8)
		p.setCellAligned(row, 4, fmt.Sprintf("%5.1f%% %s", proc.MemUsage, memBar), memColor, tview.AlignRight)

		// State with simple indicator
		stateIcon := "S"
		stateColor := tcell.ColorGray
		if strings.Contains(proc.State, "R") {
			stateIcon = "R"
			stateColor = CurrentTheme.LowUsage
		} else if strings.Contains(proc.State, "Z") {
			stateIcon = "Z"
			stateColor = CurrentTheme.HighUsage
		} else if strings.Contains(proc.State, "T") {
			stateIcon = "T"
			stateColor = CurrentTheme.MedUsage
		}
		p.setCellAligned(row, 5, stateIcon, stateColor, tview.AlignCenter)
	}

	// Remove stale rows if the new list is shorter
	if len(displayProcs)+1 < rowCount {
		for i := rowCount - 1; i > len(displayProcs); i-- {
			p.Table.RemoveRow(i)
		}
	}

	p.ApplyTheme() // Re-apply theme to ensure borders/headers are correct
}

func makeMiniBar(percent float64, width int) string {
	filled := int((percent / 100.0) * float64(width))
	if filled > width {
		filled = width
	}

	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "▪"
		} else {
			bar += "·"
		}
	}
	return bar
}

func (p *ProcTableComponent) setCell(row, col int, text string, color tcell.Color) {
	cell := p.Table.GetCell(row, col)
	if cell == nil {
		cell = tview.NewTableCell(text)
	} else {
		cell.SetText(text)
	}
	cell.SetTextColor(color)
	p.Table.SetCell(row, col, cell)
}

func (p *ProcTableComponent) setCellAligned(row, col int, text string, color tcell.Color, align int) {
	cell := p.Table.GetCell(row, col)
	if cell == nil {
		cell = tview.NewTableCell(text)
	} else {
		cell.SetText(text)
	}
	cell.SetTextColor(color).SetAlign(align)
	p.Table.SetCell(row, col, cell)
}

// GetSelectedPID returns the PID of the currently selected row
func (p *ProcTableComponent) GetSelectedPID() (int32, error) {
	row, _ := p.Table.GetSelection()
	if row <= 0 { // Header or nothing selected
		return 0, fmt.Errorf("no process selected")
	}

	// PID is in column 0
	cell := p.Table.GetCell(row, 0)
	if cell == nil {
		return 0, fmt.Errorf("empty cell")
	}

	var pid int32
	_, err := fmt.Sscanf(cell.Text, "%d", &pid)
	return pid, err
}
