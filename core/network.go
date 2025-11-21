package core

import (
	"github.com/shirou/gopsutil/v3/net"
)

type NetStat struct {
	Name        string
	BytesSent   uint64
	BytesRecv   uint64
	PacketsSent uint64
	PacketsRecv uint64
}

func FetchNetStats() ([]NetStat, error) {
	io, err := net.IOCounters(true) // true for per-interface
	if err != nil {
		return nil, err
	}

	var stats []NetStat
	for _, i := range io {
		stats = append(stats, NetStat{
			Name:        i.Name,
			BytesSent:   i.BytesSent,
			BytesRecv:   i.BytesRecv,
			PacketsSent: i.PacketsSent,
			PacketsRecv: i.PacketsRecv,
		})
	}
	return stats, nil
}
