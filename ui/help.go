package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// HelpComponent displays keybindings and usage information
type HelpComponent struct {
	Flex     *tview.Flex
	Table    *tview.Table
	DoneFunc func()
}

// NewHelpComponent creates a new help component
func NewHelpComponent() *HelpComponent {
	table := tview.NewTable().
		SetBorders(false).
		SetSelectable(false, false)

	// Define keybindings content
	sections := []struct {
		Title string
		Items []struct {
			Key  string
			Desc string
		}
	}{
		{
			Title: "Navigation",
			Items: []struct{ Key, Desc string }{
				{"1-9, 0", "Switch Views"},
				{"/", "Search Processes"},
				{"Esc", "Close Help / Clear Search"},
			},
		},
		{
			Title: "General Actions",
			Items: []struct{ Key, Desc string }{
				{"q", "Quit Application"},
				{"t", "Cycle Themes"},
				{"e", "Export Snapshot"},
				{"?", "Toggle Help"},
			},
		},
		{
			Title: "Process View",
			Items: []struct{ Key, Desc string }{
				{"k", "Kill Selected Process"},
				{"T", "Toggle Tree View"},
				{"c", "Sort by CPU"},
				{"m", "Sort by Memory"},
				{"p", "Sort by PID"},
				{"n", "Sort by Name"},
			},
		},
	}

	row := 0
	for _, section := range sections {
		// Section Header
		table.SetCell(row, 0, tview.NewTableCell(section.Title).
			SetTextColor(tcell.ColorYellow).
			SetAttributes(tcell.AttrBold).
			SetSelectable(false).
			SetAlign(tview.AlignLeft).
			SetExpansion(1))
		table.SetCell(row, 1, tview.NewTableCell("").SetSelectable(false))
		row++

		// Items
		for _, item := range section.Items {
			table.SetCell(row, 0, tview.NewTableCell("  "+item.Key).
				SetTextColor(tcell.ColorWhite).
				SetAlign(tview.AlignLeft))

			table.SetCell(row, 1, tview.NewTableCell(item.Desc).
				SetTextColor(tcell.ColorGray).
				SetAlign(tview.AlignLeft))
			row++
		}
		// Spacer
		row++
	}

	// Wrap table in a centered modal-like box
	// We use a Flex with spacers to center the table
	// And we give the table a border

	table.SetBorder(true).SetTitle(" Help & Keybindings ")

	// Calculate approximate height needed
	height := row + 4
	width := 60

	flex := tview.NewFlex().
		AddItem(nil, 0, 1, false). // Top spacer
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).       // Left spacer
			AddItem(table, height, 1, true). // The Table
			AddItem(nil, 0, 1, false),       // Right spacer
						width, 1, true).
		AddItem(nil, 0, 1, false) // Bottom spacer

	h := &HelpComponent{
		Flex:  flex,
		Table: table,
	}

	// Handle input to close
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc || event.Key() == tcell.KeyEnter || event.Rune() == 'q' || event.Rune() == '?' {
			if h.DoneFunc != nil {
				h.DoneFunc()
			}
			return nil
		}
		return event
	})

	return h
}

// SetDoneFunc sets the callback for when the help is closed
func (h *HelpComponent) SetDoneFunc(handler func()) {
	h.DoneFunc = handler
}

// ApplyTheme applies the current theme to the help component
func (h *HelpComponent) ApplyTheme() {
	h.Table.SetBackgroundColor(CurrentTheme.Background)
	h.Table.SetBorderColor(CurrentTheme.Border)
	h.Table.SetTitleColor(CurrentTheme.HeaderTitle)

	// Re-color headers and items
	rowCount := h.Table.GetRowCount()
	for r := 0; r < rowCount; r++ {
		cellKey := h.Table.GetCell(r, 0)
		cellDesc := h.Table.GetCell(r, 1)

		if cellKey != nil {
			// Heuristic: if it has bold attribute, it's a header
			if cellKey.Attributes&tcell.AttrBold != 0 {
				cellKey.SetTextColor(CurrentTheme.HeaderTitle)
			} else {
				cellKey.SetTextColor(CurrentTheme.Foreground)
			}
		}
		if cellDesc != nil {
			cellDesc.SetTextColor(CurrentTheme.HeaderValue)
		}
	}
}
