package core

import (
	"bytes"
	"encoding/csv"
	"os/exec"
	"strconv"
	"strings"
)

type GPUStat struct {
	Index         int
	Name          string
	Utilization   int     // %
	MemoryUsed    int     // MiB
	MemoryTotal   int     // MiB
	MemoryPercent float64 // %
	Temp          int     // Celsius
	PowerDraw     float64 // Watts
	PowerLimit    float64 // Watts
	FanSpeed      int     // %
	ClockCore     int     // MHz
	ClockMemory   int     // MHz
}

func FetchGPUStats() ([]GPUStat, error) {
	// nvidia-smi with extended query
	cmd := exec.Command("nvidia-smi",
		"--query-gpu=index,name,utilization.gpu,memory.used,memory.total,temperature.gpu,power.draw,power.limit,fan.speed,clocks.gr,clocks.mem",
		"--format=csv,noheader,nounits")
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

		stat := GPUStat{
			Index:       idx,
			Name:        strings.TrimSpace(record[1]),
			Utilization: util,
			MemoryUsed:  memUsed,
			MemoryTotal: memTotal,
			Temp:        temp,
		}

		// Calculate memory percentage
		if memTotal > 0 {
			stat.MemoryPercent = float64(memUsed) / float64(memTotal) * 100.0
		}

		// Parse optional fields (may not be available on all GPUs)
		if len(record) > 6 {
			powerDraw, _ := strconv.ParseFloat(strings.TrimSpace(record[6]), 64)
			stat.PowerDraw = powerDraw
		}
		if len(record) > 7 {
			powerLimit, _ := strconv.ParseFloat(strings.TrimSpace(record[7]), 64)
			stat.PowerLimit = powerLimit
		}
		if len(record) > 8 {
			fanSpeed, _ := strconv.Atoi(strings.TrimSpace(record[8]))
			stat.FanSpeed = fanSpeed
		}
		if len(record) > 9 {
			clockCore, _ := strconv.Atoi(strings.TrimSpace(record[9]))
			stat.ClockCore = clockCore
		}
		if len(record) > 10 {
			clockMemory, _ := strconv.Atoi(strings.TrimSpace(record[10]))
			stat.ClockMemory = clockMemory
		}

		stats = append(stats, stat)
	}
	return stats, nil
}
