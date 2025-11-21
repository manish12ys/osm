package core

import (
	"bytes"
	"encoding/csv"
	"os/exec"
	"strconv"
	"strings"
)

type GPUStat struct {
	Index       int
	Name        string
	Utilization int // %
	MemoryUsed  int // MiB
	MemoryTotal int // MiB
	Temp        int // Celsius
}

func FetchGPUStats() ([]GPUStat, error) {
	// nvidia-smi --query-gpu=index,name,utilization.gpu,memory.used,memory.total,temperature.gpu --format=csv,noheader,nounits
	cmd := exec.Command("nvidia-smi", "--query-gpu=index,name,utilization.gpu,memory.used,memory.total,temperature.gpu", "--format=csv,noheader,nounits")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(&out)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var stats []GPUStat
	for _, record := range records {
		if len(record) < 6 {
			continue
		}

		idx, _ := strconv.Atoi(strings.TrimSpace(record[0]))
		util, _ := strconv.Atoi(strings.TrimSpace(record[2]))
		memUsed, _ := strconv.Atoi(strings.TrimSpace(record[3]))
		memTotal, _ := strconv.Atoi(strings.TrimSpace(record[4]))
		temp, _ := strconv.Atoi(strings.TrimSpace(record[5]))

		stats = append(stats, GPUStat{
			Index:       idx,
			Name:        strings.TrimSpace(record[1]),
			Utilization: util,
			MemoryUsed:  memUsed,
			MemoryTotal: memTotal,
			Temp:        temp,
		})
	}
	return stats, nil
}
