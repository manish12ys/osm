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

	// Helper to get color tag
	cLow := fmt.Sprintf("[#%06x]", CurrentTheme.LowUsage.Hex())
	cMed := fmt.Sprintf("[#%06x]", CurrentTheme.MedUsage.Hex())
	cHigh := fmt.Sprintf("[#%06x]", CurrentTheme.HighUsage.Hex())
	cFore := fmt.Sprintf("[#%06x]", CurrentTheme.Foreground.Hex())
	cDim := "[darkgray]"
	cTitle := fmt.Sprintf("[#%06x]", CurrentTheme.HeaderTitle.Hex())
	cLabel := fmt.Sprintf("[#%06x]", CurrentTheme.HeaderValue.Hex())

	if len(results) == 0 {
		fmt.Fprintf(p.TextView, "%sNo plugins loaded[-]\n\n", cMed)
		fmt.Fprintf(p.TextView, "Create plugins in: %s~/.config/osm/plugins/[-]\n\n", cTitle)
		fmt.Fprintf(p.TextView, "%sExample plugin (uptime.sh):[-]\n", cFore)
		fmt.Fprintf(p.TextView, "%s#!/bin/bash\n", cDim)
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
			fmt.Fprintf(p.TextView, "%s✗ %s[-]\n", cHigh, output.Title)
			fmt.Fprintf(p.TextView, "  %sError: %s[-]\n\n", cDim, output.Error)
			continue
		}

		fmt.Fprintf(p.TextView, "%s✓ %s[-]\n", cLow, output.Title)

		if output.Value != "" {
			// Translate ANSI codes to tview tags for colorful output (e.g. neofetch)
			translated := tview.TranslateANSI(output.Value)
			fmt.Fprintf(p.TextView, "  %s%s[-]\n", cFore, translated)
		}

		if len(output.Metrics) > 0 {
			for key, value := range output.Metrics {
				fmt.Fprintf(p.TextView, "  %s%s:[-] %s%s[-]\n", cLabel, key, cFore, value)
			}
		}

		fmt.Fprintf(p.TextView, "\n")
	}
}
