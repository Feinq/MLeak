package monitor

import (
	"fmt"
	"runtime"
)

type MemoryStats struct {
	Alloc      uint64 // bytes allocated and still in use
	TotalAlloc uint64 // bytes allocated (even if freed)
	Sys        uint64 // bytes obtained from the OS
	NumGC      uint32 // number of garbage collections
}

func GetMemoryUsage() MemoryStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return MemoryStats{
		Alloc:      m.Alloc,
		TotalAlloc: m.TotalAlloc,
		Sys:        m.Sys,
		NumGC:      m.NumGC,
	}
}

func PrintMemoryUsage() {
	stats := GetMemoryUsage()
	fmt.Printf("Current Memory Usage:\n")
	fmt.Printf("Allocated Memory: %v bytes\n", stats.Alloc)
	fmt.Printf("Total Allocated Memory: %v bytes\n", stats.TotalAlloc)
	fmt.Printf("System Memory: %v bytes\n", stats.Sys)
	fmt.Printf("Number of Garbage Collections: %v\n", stats.NumGC)
}
