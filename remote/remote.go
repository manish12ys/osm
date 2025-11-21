package remote

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"omarchy-monitor/core"
)

// RemoteStats holds remote system statistics
type RemoteStats struct {
	Hostname   string
	Stats      *core.Stats
	Processes  []core.Process
	Disks      []core.DiskStat
	Nets       []core.NetStat
	Temps      []core.TempStat
	GPUs       []core.GPUStat
	Containers []core.ContainerStat
	Timestamp  time.Time
}

// Client handles remote monitoring connections
type Client struct {
	BaseURL string
	client  *http.Client
}

// NewClient creates a new remote monitoring client
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// FetchRemoteStats fetches statistics from a remote server
func (c *Client) FetchRemoteStats() (*RemoteStats, error) {
	resp, err := c.client.Get(c.BaseURL + "/api/stats")
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var stats RemoteStats
	if err := json.Unmarshal(body, &stats); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &stats, nil
}

// Server provides remote monitoring API
type Server struct {
	Port string
}

// NewServer creates a new remote monitoring server
func NewServer(port string) *Server {
	return &Server{Port: port}
}

// Start starts the remote monitoring server
func (s *Server) Start() error {
	http.HandleFunc("/api/stats", s.handleStats)
	http.HandleFunc("/health", s.handleHealth)

	fmt.Printf("Remote monitoring server starting on port %s\n", s.Port)
	return http.ListenAndServe(":"+s.Port, nil)
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Fetch all stats
	stats, _ := core.FetchStats()
	procs, _ := core.FetchProcesses()
	disks, _ := core.FetchDiskStats()
	nets, _ := core.FetchNetStats()
	temps, _ := core.FetchTemps()
	gpus, _ := core.FetchGPUStats()
	containers, _ := core.FetchContainers()

	remoteStats := RemoteStats{
		Hostname:   stats.Hostname,
		Stats:      stats,
		Processes:  procs,
		Disks:      disks,
		Nets:       nets,
		Temps:      temps,
		GPUs:       gpus,
		Containers: containers,
		Timestamp:  time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(remoteStats)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
