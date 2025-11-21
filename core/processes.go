package core

import (
	"sort"
	"strings"

	"github.com/shirou/gopsutil/v3/process"
)

// Process represents a single system process
type Process struct {
	PID      int32
	PPID     int32
	Name     string
	User     string
	CPUUsage float64
	MemUsage float32
	State    string
}

// FetchProcesses returns a list of running processes
func FetchProcesses() ([]Process, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}

	var results []Process
	for _, p := range procs {
		// Basic info
		name, _ := p.Name()
		user, _ := p.Username()
		ppid, _ := p.Ppid()

		// CPU & Mem
		// Note: CPUPercent might need a duration in some versions, or 0 for "since last call"
		cpuP, _ := p.CPUPercent()
		memP, _ := p.MemoryPercent()

		// State
		state, _ := p.Status() // returns []string usually

		results = append(results, Process{
			PID:      p.Pid,
			PPID:     ppid,
			Name:     name,
			User:     user,
			CPUUsage: cpuP,
			MemUsage: memP,
			State:    strings.Join(state, ","),
		})
	}

	// Default sort by CPU usage descending
	SortProcesses(results, SortCPU)

	return results, nil
}

// KillProcess terminates a process by PID
func KillProcess(pid int32) error {
	p, err := process.NewProcess(pid)
	if err != nil {
		return err
	}
	return p.Kill()
}

// SortType defines how to sort processes
type SortType int

const (
	SortCPU SortType = iota
	SortMem
	SortPID
	SortName
)

// SortProcesses sorts the process list in place
func SortProcesses(procs []Process, sortBy SortType) {
	sort.Slice(procs, func(i, j int) bool {
		switch sortBy {
		case SortMem:
			return procs[i].MemUsage > procs[j].MemUsage
		case SortPID:
			return procs[i].PID < procs[j].PID
		case SortName:
			return strings.ToLower(procs[i].Name) < strings.ToLower(procs[j].Name)
		case SortCPU:
			fallthrough
		default:
			return procs[i].CPUUsage > procs[j].CPUUsage
		}
	})
}
