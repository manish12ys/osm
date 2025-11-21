package core

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// ExportData represents the snapshot of data to be exported
type ExportData struct {
	Timestamp string    `json:"timestamp"`
	Stats     *Stats    `json:"stats"`
	Processes []Process `json:"processes"`
}

// ExportSnapshot saves the current system state to a JSON file
func ExportSnapshot(stats *Stats, procs []Process) (string, error) {
	data := ExportData{
		Timestamp: time.Now().Format(time.RFC3339),
		Stats:     stats,
		Processes: procs,
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}

	filename := fmt.Sprintf("osm_snapshot_%s.json", time.Now().Format("20060102_150405"))
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return "", err
	}

	return filename, nil
}
