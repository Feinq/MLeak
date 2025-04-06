package report

import (
	"encoding/json"
	"fmt"
	"os"
)

type MemoryReport struct {
	Timestamp string  `json:"timestamp"`
	RSS       float64 `json:"rss"`
	LeakRisk  string  `json:"leak_risk"`
	Growth    float64 `json:"growth,omitempty"` // Memory growth since last report
}

func (r MemoryReport) OutputJSON() error {
	data, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("error marshaling report to JSON: %w", err)
	}
	_, err = os.Stdout.Write(data)
	if err != nil {
		return fmt.Errorf("error writing report to stdout: %w", err)
	}
	return nil
}

func (r MemoryReport) OutputPlainText() {
	fmt.Printf("Timestamp: %s\nRSS: %.2f MB\nLeak Risk: %s\n",
		r.Timestamp, r.RSS/1024/1024, r.LeakRisk)
}

func DetermineLeakRisk(current, previous float64, timeInterval float64) string {
	if previous == 0 {
		return "Low"
	}

	growthRate := (current - previous) / previous
	growthPerSecond := growthRate / timeInterval

	if growthPerSecond > 0.05 {
		return "High"
	} else if growthPerSecond > 0.01 {
		return "Medium"
	}
	return "Low"
}

func DetermineLeakRiskWithTrend(trend *MemoryTrend) string {
	return AnalyzeLeakRisk(trend)
}
