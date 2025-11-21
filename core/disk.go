package core

import (
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/disk"
)

type DiskStat struct {
	Device        string
	Mountpoint    string
	Fstype        string
	Total         uint64
	Used          uint64
	Free          uint64
	UsedPercent   float64
	ReadBytes     uint64
	WriteBytes    uint64
	ReadCount     uint64
	WriteCount    uint64
	ReadSpeed     float64 // Bytes per second
	WriteSpeed    float64 // Bytes per second
	IOPS          float64 // Operations per second
	InodesTotal   uint64
	InodesUsed    uint64
	InodesFree    uint64
	InodesPercent float64
}

var (
	lastDiskIO     map[string]disk.IOCountersStat
	lastDiskIOTime time.Time
	diskIOMutex    sync.Mutex
)

func init() {
	lastDiskIO = make(map[string]disk.IOCountersStat)
	lastDiskIOTime = time.Now()
}

func FetchDiskStats() ([]DiskStat, error) {
	parts, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	// Fetch I/O counters
	ioCounters, _ := disk.IOCounters()
	currentTime := time.Now()

	diskIOMutex.Lock()
	timeDelta := currentTime.Sub(lastDiskIOTime).Seconds()
	if timeDelta == 0 {
		timeDelta = 1
	}
	diskIOMutex.Unlock()

	var stats []DiskStat
	for _, p := range parts {
		u, err := disk.Usage(p.Mountpoint)
		if err != nil {
			continue // Skip if permission denied or other error
		}

		stat := DiskStat{
			Device:        p.Device,
			Mountpoint:    p.Mountpoint,
			Fstype:        p.Fstype,
			Total:         u.Total,
			Used:          u.Used,
			Free:          u.Free,
			UsedPercent:   u.UsedPercent,
			InodesTotal:   u.InodesTotal,
			InodesUsed:    u.InodesUsed,
			InodesFree:    u.InodesFree,
			InodesPercent: u.InodesUsedPercent,
		}

		// Try to match I/O stats by device name
		if ioStat, ok := ioCounters[p.Device]; ok {
			stat.ReadBytes = ioStat.ReadBytes
			stat.WriteBytes = ioStat.WriteBytes
			stat.ReadCount = ioStat.ReadCount
			stat.WriteCount = ioStat.WriteCount

			// Calculate speeds if we have previous data
			diskIOMutex.Lock()
			if lastIO, exists := lastDiskIO[p.Device]; exists {
				readDelta := float64(ioStat.ReadBytes - lastIO.ReadBytes)
				writeDelta := float64(ioStat.WriteBytes - lastIO.WriteBytes)
				readOpsDelta := float64(ioStat.ReadCount - lastIO.ReadCount)
				writeOpsDelta := float64(ioStat.WriteCount - lastIO.WriteCount)

				stat.ReadSpeed = readDelta / timeDelta
				stat.WriteSpeed = writeDelta / timeDelta
				stat.IOPS = (readOpsDelta + writeOpsDelta) / timeDelta
			}
			lastDiskIO[p.Device] = ioStat
			diskIOMutex.Unlock()
		}

		stats = append(stats, stat)
	}

	diskIOMutex.Lock()
	lastDiskIOTime = currentTime
	diskIOMutex.Unlock()

	return stats, nil
}
