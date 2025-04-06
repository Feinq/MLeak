package report

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const (
	topLeft     = "┌"
	topRight    = "┐"
	bottomLeft  = "└"
	bottomRight = "┘"
	horizontal  = "─"
	vertical    = "│"
	teeDown     = "┬"
	teeUp       = "┴"
	teeRight    = "├"
	teeLeft     = "┤"
	cross       = "┼"
)

type TableMemoryReport struct {
	Reports     []MemoryReport
	ProcessID   int
	ProcessName string
	Interval    time.Duration
	MaxEntries  int
}

func NewTableMemoryReport(pid int, processName string, interval time.Duration, maxEntries int) *TableMemoryReport {
	return &TableMemoryReport{
		Reports:     make([]MemoryReport, 0, maxEntries),
		ProcessID:   pid,
		ProcessName: processName,
		Interval:    interval,
		MaxEntries:  maxEntries,
	}
}

func (t *TableMemoryReport) AddReport(report MemoryReport) {
	if len(t.Reports) >= t.MaxEntries {
		t.Reports = t.Reports[1:]
	}
	t.Reports = append(t.Reports, report)
}

func ClearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (t *TableMemoryReport) DisplayTable() {
	ClearScreen()

	timeWidth := 12 // timestamp column
	rssWidth := 11  // RSS column
	riskWidth := 11 // risk column

	processInfo := fmt.Sprintf("PID %d", t.ProcessID)
	if t.ProcessName != "" {
		processInfo = fmt.Sprintf("PID %d (%s)", t.ProcessID, t.ProcessName)
	}

	intervalStr := fmt.Sprintf("%.0fs", t.Interval.Seconds())
	fmt.Printf("[MLeak] Monitoring %s | Interval: %s\n", processInfo, intervalStr)

	drawHorizontalBorder := func(left, middle, right string, width ...int) string {
		var result strings.Builder
		result.WriteString(left)

		for i, w := range width {
			result.WriteString(strings.Repeat(horizontal, w))
			if i < len(width)-1 {
				result.WriteString(middle)
			}
		}

		result.WriteString(right)
		return result.String()
	}

	topBorder := drawHorizontalBorder(topLeft, teeDown, topRight, timeWidth, rssWidth, riskWidth)
	fmt.Println(topBorder)

	fmt.Printf("%s %-*s %s %-*s %s %-*s %s\n",
		vertical, timeWidth-2, "Timestamp",
		vertical, rssWidth-2, "RSS",
		vertical, riskWidth-2, "Leak Risk",
		vertical)

	middleBorder := drawHorizontalBorder(teeRight, cross, teeLeft, timeWidth, rssWidth, riskWidth)
	fmt.Println(middleBorder)

	for i, report := range t.Reports {
		ts, _ := time.Parse(time.RFC3339, report.Timestamp)
		timeStr := ts.Format("15:04:05")

		rssMB := fmt.Sprintf("%.2f MB", report.RSS/1024/1024)

		growthIndicator := ""
		if i > 0 {
			prevRSS := t.Reports[i-1].RSS
			growth := report.RSS - prevRSS
			if growth > 0 {
				growthMB := growth / 1024 / 1024
				timeDiff := t.Interval.Seconds()
				warningEmoji := ""

				if report.LeakRisk == "Medium" || report.LeakRisk == "High" {
					warningEmoji = " !"
				}

				if growthMB > 0 {
					growthIndicator = fmt.Sprintf("%s +%dMB/%.0fs", warningEmoji, int(growthMB), timeDiff)
				}
			}
		}

		fmt.Printf("%s %-*s %s %-*s %s %-*s %s%s\n",
			vertical, timeWidth-2, timeStr,
			vertical, rssWidth-2, rssMB,
			vertical, riskWidth-2, report.LeakRisk,
			vertical, growthIndicator)
	}

	bottomBorder := drawHorizontalBorder(bottomLeft, teeUp, bottomRight, timeWidth, rssWidth, riskWidth)
	fmt.Println(bottomBorder)
}
