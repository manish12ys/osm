package core

import (
	"github.com/shirou/gopsutil/v3/disk"
)

type DiskStat struct {
	Device      string
	Mountpoint  string
	Fstype      string
	Total       uint64
	Used        uint64
	Free        uint64
	UsedPercent float64
}

func FetchDiskStats() ([]DiskStat, error) {
	parts, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	var stats []DiskStat
	for _, p := range parts {
		u, err := disk.Usage(p.Mountpoint)
		if err != nil {
			continue // Skip if permission denied or other error
		}
		stats = append(stats, DiskStat{
			Device:      p.Device,
			Mountpoint:  p.Mountpoint,
			Fstype:      p.Fstype,
			Total:       u.Total,
			Used:        u.Used,
			Free:        u.Free,
			UsedPercent: u.UsedPercent,
		})
	}
	return stats, nil
}
