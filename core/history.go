package core

import "sync"

// History holds historical data points
type History struct {
	mu   sync.Mutex
	Data []float64
	Max  int
}

// NewHistory creates a new history buffer
func NewHistory(max int) *History {
	return &History{
		Data: make([]float64, 0, max),
		Max:  max,
	}
}

// Add appends a value and maintains size
func (h *History) Add(val float64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if len(h.Data) >= h.Max {
		h.Data = h.Data[1:]
	}
	h.Data = append(h.Data, val)
}

// GetData returns a copy of the data
func (h *History) GetData() []float64 {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Return a copy to avoid race conditions
	cpy := make([]float64, len(h.Data))
	copy(cpy, h.Data)
	return cpy
}
