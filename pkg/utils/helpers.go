package utils

import (
	"fmt"
	"runtime"
)

func GetMemoryStats() (uint64, uint64) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return memStats.Alloc, memStats.TotalAlloc
}

func FormatMemoryUsage(bytes uint64) string {
	const (
		_  = iota
		KB = 1 << (10 * iota)
		MB
		GB
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
