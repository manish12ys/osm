package core

import (
	"fmt"

	"github.com/distatus/battery"
)

type BatteryStats struct {
	Capacity int    // Percentage
	Status   string // Charging, Discharging, Full, Unknown
}

func FetchBatteryStats() (*BatteryStats, error) {
	batteries, err := battery.GetAll()
	if err != nil {
		return nil, err
	}

	if len(batteries) == 0 {
		return nil, fmt.Errorf("no battery found")
	}

	// Use the first battery found
	bat := batteries[0]

	// Map status
	status := bat.State.String()

	// Calculate percentage
	// distatus/battery provides Current and Full in mWh.
	capacity := 0
	if bat.Full > 0 {
		capacity = int((bat.Current / bat.Full) * 100)
	}

	return &BatteryStats{
		Capacity: capacity,
		Status:   status,
	}, nil
}
