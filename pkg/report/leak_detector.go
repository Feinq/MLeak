package report

import (
	"math"
	"time"
)

type MemoryTrend struct {
	Samples       []float64     // Memory samples over time
	Timestamps    []time.Time   // Corresponding timestamps
	MaxSamples    int           // Maximum samples to keep
	GrowthRate    float64       // Current growth rate (MB/s)
	IsStable      bool          // Whether memory usage is stable
	StableFor     time.Duration // How long memory has been stable
	LastStableRSS float64       // Last known stable memory usage
}

func NewMemoryTrend(maxSamples int) *MemoryTrend {
	return &MemoryTrend{
		Samples:    make([]float64, 0, maxSamples),
		Timestamps: make([]time.Time, 0, maxSamples),
		MaxSamples: maxSamples,
		IsStable:   true,
	}
}

func (mt *MemoryTrend) AddSample(rss float64, timestamp time.Time) {
	if len(mt.Samples) >= mt.MaxSamples {
		mt.Samples = mt.Samples[1:]
		mt.Timestamps = mt.Timestamps[1:]
	}

	mt.Samples = append(mt.Samples, rss)
	mt.Timestamps = append(mt.Timestamps, timestamp)

	mt.updateMetrics()
}

func (mt *MemoryTrend) updateMetrics() {
	if len(mt.Samples) < 2 {
		mt.GrowthRate = 0
		return
	}

	recent := mt.Samples[len(mt.Samples)-1]
	previous := mt.Samples[len(mt.Samples)-2]
	timeDiff := mt.Timestamps[len(mt.Timestamps)-1].Sub(mt.Timestamps[len(mt.Timestamps)-2]).Seconds()

	if timeDiff > 0 {
		mt.GrowthRate = ((recent - previous) / 1024 / 1024) / timeDiff
	}

	growthTolerance := 0.1 // MB/s
	isCurrentlyStable := math.Abs(mt.GrowthRate) < growthTolerance

	if isCurrentlyStable {
		if !mt.IsStable {
			mt.IsStable = true
			mt.StableFor = 0
		} else {
			mt.StableFor += time.Duration(timeDiff) * time.Second
		}
		mt.LastStableRSS = recent
	} else {
		mt.IsStable = false
		mt.StableFor = 0
	}
}

func (mt *MemoryTrend) GetTrendMetrics() (avgGrowthRate float64, consistentGrowth bool, stabilizationTime time.Duration) {
	if len(mt.Samples) < 3 {
		return 0, false, 0
	}

	firstSample := mt.Samples[0]
	lastSample := mt.Samples[len(mt.Samples)-1]
	totalTime := mt.Timestamps[len(mt.Timestamps)-1].Sub(mt.Timestamps[0]).Seconds()

	if totalTime > 0 {
		avgGrowthRate = ((lastSample - firstSample) / 1024 / 1024) / totalTime
	}

	growthCount := 0
	shrinkCount := 0

	for i := 1; i < len(mt.Samples); i++ {
		if mt.Samples[i] > mt.Samples[i-1] {
			growthCount++
		} else if mt.Samples[i] < mt.Samples[i-1] {
			shrinkCount++
		}
	}

	consistentGrowth = growthCount > (len(mt.Samples)-1)*2/3

	if mt.IsStable {
		stabilizationTime = mt.StableFor
	} else {
		stabilizationTime = 0
	}

	return
}

func AnalyzeLeakRisk(trend *MemoryTrend) string {
	avgGrowthRate, consistentGrowth, stabilizationTime := trend.GetTrendMetrics()

	if len(trend.Samples) < 3 {
		return "Low"
	}

	// Check for memory stability
	if trend.IsStable && trend.StableFor > 5*time.Second {
		return "Low"
	}

	if consistentGrowth {
		if avgGrowthRate > 1.0 { // More than 1 MB/s
			return "High"
		} else if avgGrowthRate > 0.2 { // More than 0.2 MB/s
			return "Medium"
		}
	}

	if !trend.IsStable && stabilizationTime < 10*time.Second && avgGrowthRate > 0.5 {
		return "Medium"
	}

	if avgGrowthRate > 0 {
		return "Low"
	}

	return "Low"
}
