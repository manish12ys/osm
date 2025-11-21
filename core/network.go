package core

import (
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/net"
)

type NetStat struct {
	Name        string
	BytesSent   uint64
	BytesRecv   uint64
	PacketsSent uint64
	PacketsRecv uint64
	Errin       uint64
	Errout      uint64
	Dropin      uint64
	Dropout     uint64
	SendSpeed   float64 // Bytes per second
	RecvSpeed   float64 // Bytes per second
	SendPPS     float64 // Packets per second (send)
	RecvPPS     float64 // Packets per second (recv)
}

var (
	lastNetIO     map[string]net.IOCountersStat
	lastNetIOTime time.Time
	netIOMutex    sync.Mutex
)

func init() {
	lastNetIO = make(map[string]net.IOCountersStat)
	lastNetIOTime = time.Now()
}

func FetchNetStats() ([]NetStat, error) {
	io, err := net.IOCounters(true) // true for per-interface
	if err != nil {
		return nil, err
	}

	currentTime := time.Now()

	netIOMutex.Lock()
	timeDelta := currentTime.Sub(lastNetIOTime).Seconds()
	if timeDelta == 0 {
		timeDelta = 1
	}
	netIOMutex.Unlock()

	var stats []NetStat
	for _, i := range io {
		stat := NetStat{
			Name:        i.Name,
			BytesSent:   i.BytesSent,
			BytesRecv:   i.BytesRecv,
			PacketsSent: i.PacketsSent,
			PacketsRecv: i.PacketsRecv,
			Errin:       i.Errin,
			Errout:      i.Errout,
			Dropin:      i.Dropin,
			Dropout:     i.Dropout,
		}

		// Calculate speeds if we have previous data
		netIOMutex.Lock()
		if lastIO, exists := lastNetIO[i.Name]; exists {
			sendDelta := float64(i.BytesSent - lastIO.BytesSent)
			recvDelta := float64(i.BytesRecv - lastIO.BytesRecv)
			sendPktDelta := float64(i.PacketsSent - lastIO.PacketsSent)
			recvPktDelta := float64(i.PacketsRecv - lastIO.PacketsRecv)

			stat.SendSpeed = sendDelta / timeDelta
			stat.RecvSpeed = recvDelta / timeDelta
			stat.SendPPS = sendPktDelta / timeDelta
			stat.RecvPPS = recvPktDelta / timeDelta
		}
		lastNetIO[i.Name] = i
		netIOMutex.Unlock()

		stats = append(stats, stat)
	}

	netIOMutex.Lock()
	lastNetIOTime = currentTime
	netIOMutex.Unlock()

	return stats, nil
}
