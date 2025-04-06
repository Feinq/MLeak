package monitor

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// getProcessMemoryImpl retrieves memory usage for a specific process on Linux
func getProcessMemoryImpl(pid int) (int64, error) {
	// Read the memory usage from /proc/[pid]/statm
	statmPath := "/proc/" + strconv.Itoa(pid) + "/statm"
	data, err := os.ReadFile(statmPath)
	if err != nil {
		return 0, err
	}

	fields := strings.Fields(string(data))
	if len(fields) < 2 {
		return 0, nil
	}

	rssPages, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return 0, err
	}

	// Convert pages to bytes (assuming a page size of 4KB)
	return rssPages * 4096, nil
}

// getProcessesImpl retrieves list of processes on Linux
func getProcessesImpl() ([]Process, error) {
	var processes []Process

	// Execute the command to get the list of processes
	cmd := exec.Command("ps", "-eo", "pid,comm,rss")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines[1:] {
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		pid, err := strconv.Atoi(fields[0])
		if err != nil {
			continue
		}

		rss, err := strconv.ParseInt(fields[2], 10, 64)
		if err != nil {
			continue
		}

		// RSS from ps is in KB, convert to bytes
		rssBytes := rss * 1024

		processes = append(processes, Process{
			PID:  pid,
			Name: fields[1],
			RSS:  rssBytes,
		})
	}

	return processes, nil
}
