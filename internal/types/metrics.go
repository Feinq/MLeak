package types

type MemoryMetrics struct {
	TotalAllocated uint64 `json:"total_allocated"`
	TotalFreed     uint64 `json:"total_freed"`
	CurrentUsage   uint64 `json:"current_usage"`
	PeakUsage      uint64 `json:"peak_usage"`
	LeakRisk       string `json:"leak_risk"`
	Timestamp      string `json:"timestamp"`
}