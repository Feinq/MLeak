//go:build windows
// +build windows

package monitor

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

// Windows API constants
const (
	processQueryInformation = 0x0400
	memoryCommit            = 0x1000
	memoryPrivate           = 0x20000
)

// Windows DLL and procedure references
var (
	modPsapi                 = syscall.NewLazyDLL("psapi.dll")
	procGetProcessMemoryInfo = modPsapi.NewProc("GetProcessMemoryInfo")
)

// getProcessMemoryImpl retrieves memory usage for a specific process on Windows
func getProcessMemoryImpl(pid int) (int64, error) {
	handle, err := syscall.OpenProcess(processQueryInformation, false, uint32(pid))
	if err != nil {
		return 0, fmt.Errorf("failed to open process: %w", err)
	}
	defer syscall.CloseHandle(handle)

	var memCounters processMemoryCounters
	memCounters.CB = uint32(unsafe.Sizeof(memCounters))

	if err := getProcessMemoryInfo(handle, &memCounters, uint32(unsafe.Sizeof(memCounters))); err != nil {
		return 0, fmt.Errorf("failed to get process memory info: %w", err)
	}

	return int64(memCounters.WorkingSetSize), nil
}

// getProcessesImpl retrieves list of processes on Windows
func getProcessesImpl() ([]Process, error) {
	var processes []Process

	// Execute the command to get the list of processes (using tasklist on Windows)
	cmd := exec.Command("tasklist", "/FO", "CSV", "/NH")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Parse the CSV output
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Parse CSV format: "name","pid","session","session#","mem usage"
		fields := strings.Split(line, ",")
		if len(fields) < 5 {
			continue
		}

		// Clean up the fields (remove quotes)
		name := strings.Trim(fields[0], "\"")
		pidStr := strings.Trim(fields[1], "\"")
		memStr := strings.Trim(fields[4], "\" K")

		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			continue
		}

		// Memory is in KB, convert to bytes
		memKB, err := strconv.ParseInt(memStr, 10, 64)
		if err != nil {
			continue
		}
		memBytes := memKB * 1024

		processes = append(processes, Process{
			PID:  pid,
			Name: name,
			RSS:  memBytes,
		})
	}

	return processes, nil
}

// Windows specific structures and function declarations
type processMemoryCounters struct {
	CB                         uint32
	PageFaultCount             uint32
	PeakWorkingSetSize         uintptr
	WorkingSetSize             uintptr
	QuotaPeakPagedPoolUsage    uintptr
	QuotaPagedPoolUsage        uintptr
	QuotaPeakNonPagedPoolUsage uintptr
	QuotaNonPagedPoolUsage     uintptr
	PagefileUsage              uintptr
	PeakPagefileUsage          uintptr
}

func getProcessMemoryInfo(handle syscall.Handle, memCounters *processMemoryCounters, size uint32) error {
	ret, _, err := procGetProcessMemoryInfo.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(memCounters)),
		uintptr(size),
	)

	if ret == 0 {
		if err != syscall.Errno(0) {
			return err
		}
		return fmt.Errorf("GetProcessMemoryInfo call failed")
	}

	return nil
}
