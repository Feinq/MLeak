package monitor

// Process represents a running process with its ID and memory usage.
type Process struct {
	PID  int
	Name string
	RSS  int64 // Resident Set Size in bytes
}

// GetProcessMemory retrieves the memory usage of a specific process by its PID.
// This is a platform-independent wrapper that calls the appropriate implementation.
func GetProcessMemory(pid int) (int64, error) {
	return getProcessMemoryImpl(pid)
}

// GetProcesses retrieves a list of currently running processes and their memory usage.
func GetProcesses() ([]Process, error) {
	return getProcessesImpl()
}
