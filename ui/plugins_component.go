package ui

import (
	"fmt"
	"omarchy-monitor/plugins"

	"github.com/rivo/tview"
)

type PluginsComponent struct {
	Flex     *tview.Flex
	TextView *tview.TextView
}

func NewPluginsComponent() *PluginsComponent {
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true)
	textView.SetBorder(true).SetTitle(" Custom Plugins (9) ")

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(textView, 0, 1, true)

	return &PluginsComponent{
		Flex:     flex,
		TextView: textView,
	}
}

func (p *PluginsComponent) ApplyTheme() {
	p.TextView.SetBorderColor(CurrentTheme.Border)
	p.TextView.SetTitleColor(CurrentTheme.HeaderTitle)
	p.TextView.SetTextColor(CurrentTheme.Foreground)
}

func (p *PluginsComponent) Update(results map[string]*plugins.PluginOutput) {
	p.TextView.Clear()

	if len(results) == 0 {
		fmt.Fprintf(p.TextView, "[yellow]No plugins loaded[-]\n\n")
		fmt.Fprintf(p.TextView, "Create plugins in: [cyan]~/.config/osm/plugins/[-]\n\n")
		fmt.Fprintf(p.TextView, "[white]Example plugin (uptime.sh):[-]\n")
		fmt.Fprintf(p.TextView, "[gray]#!/bin/bash\n")
		fmt.Fprintf(p.TextView, "cat <<EOF\n")
		fmt.Fprintf(p.TextView, "{\n")
		fmt.Fprintf(p.TextView, "  \"title\": \"System Uptime\",\n")
		fmt.Fprintf(p.TextView, "  \"value\": \"$(uptime -p)\"\n")
		fmt.Fprintf(p.TextView, "}\n")
		fmt.Fprintf(p.TextView, "EOF[-]\n")
		return
	}

	for _, output := range results {
		if output.Error != "" {
			fmt.Fprintf(p.TextView, "[red]✗ %s[-]\n", output.Title)
			fmt.Fprintf(p.TextView, "  [gray]Error: %s[-]\n\n", output.Error)
			continue
		}

		fmt.Fprintf(p.TextView, "[green]✓ %s[-]\n", output.Title)

		if output.Value != "" {
			fmt.Fprintf(p.TextView, "  [white]%s[-]\n", output.Value)
		}

		if len(output.Metrics) > 0 {
			for key, value := range output.Metrics {
				fmt.Fprintf(p.TextView, "  [cyan]%s:[-] [white]%s[-]\n", key, value)
			}
		}

		fmt.Fprintf(p.TextView, "\n")
	}
}
