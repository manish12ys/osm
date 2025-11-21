package core

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

// Stats holds system statistics
type Stats struct {
	CPUUsage  float64
	MemUsage  float64
	MemTotal  uint64
	MemUsed   uint64
	Uptime    uint64
	LoadAvg1  float64
	LoadAvg5  float64
	LoadAvg15 float64
	Hostname  string
	OS        string
}

// FetchStats gathers current system stats
func FetchStats() (*Stats, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	c, err := cpu.Percent(0, false)
	if err != nil {
		return nil, err
	}

	h, err := host.Info()
	if err != nil {
		return nil, err
	}

	l, err := load.Avg()
	if err != nil {
		return nil, err
	}

	cpuVal := 0.0
	if len(c) > 0 {
		cpuVal = c[0]
	}

	return &Stats{
		CPUUsage:  cpuVal,
		MemUsage:  v.UsedPercent,
		MemTotal:  v.Total,
		MemUsed:   v.Used,
		Uptime:    h.Uptime,
		LoadAvg1:  l.Load1,
		LoadAvg5:  l.Load5,
		LoadAvg15: l.Load15,
		Hostname:  h.Hostname,
		OS:        h.Platform,
	}, nil
}

// FetchPerCoreCPU returns per-core CPU usage
func FetchPerCoreCPU() ([]float64, error) {
	percentages, err := cpu.Percent(0, true) // true for per-CPU
	if err != nil {
		return nil, err
	}
	return percentages, nil
}

// TempStat holds temperature info
type TempStat struct {
	Key         string
	Temperature float64
}

// FetchTemps returns system temperatures
func FetchTemps() ([]TempStat, error) {
	// gopsutil/host doesn't have sensors directly, usually it's in host.SensorsTemperatures
	// But wait, gopsutil v3 has it in "host" package as SensorsTemperatures
	ts, err := host.SensorsTemperatures()
	if err != nil {
		return nil, err
	}

	var stats []TempStat
	for _, t := range ts {
		stats = append(stats, TempStat{
			Key:         t.SensorKey,
			Temperature: t.Temperature,
		})
	}
	return stats, nil
}
