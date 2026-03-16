package analyzer

import (
	"sort"
	"time"

	"github.com/liucxer/analyze_ceph_log/pkg/types"
)

// AnalyzeAIOEvents analyzes AIO events and returns statistics
func AnalyzeAIOEvents(events []types.Event) types.AnalysisResult {
	totalEvents := len(events)
	totalDuration := time.Duration(0)
	maxDuration := time.Duration(0)
	minDuration := time.Duration(1000000000) // Start with a large value

	// Count requests by all durations present in data
	durationCounts := make(map[int]int)

	for _, event := range events {
		totalDuration += event.Duration
		if event.Duration > maxDuration {
			maxDuration = event.Duration
		}
		if event.Duration < minDuration {
			minDuration = event.Duration
		}

		// Categorize by duration
		durationMs := int(float64(event.Duration.Microseconds()) / 1000.0)
		durationCounts[durationMs]++
	}

	// Sort events by start time
	sort.Slice(events, func(i, j int) bool {
		return events[i].StartTime.Before(events[j].StartTime)
	})

	return types.AnalysisResult{
		Events:         events,
		TotalEvents:    totalEvents,
		TotalDuration:  totalDuration,
		MaxDuration:    maxDuration,
		MinDuration:    minDuration,
		DurationCounts: durationCounts,
	}
}

// AnalyzeRepopEvents analyzes OSD repop events and returns statistics
func AnalyzeRepopEvents(events []types.Event) types.AnalysisResult {
	totalEvents := len(events)
	totalDuration := time.Duration(0)
	maxDuration := time.Duration(0)
	minDuration := time.Duration(1000000000) // Start with a large value

	// Count requests by all durations present in data
	durationCounts := make(map[int]int)

	for _, event := range events {
		totalDuration += event.Duration
		if event.Duration > maxDuration {
			maxDuration = event.Duration
		}
		if event.Duration < minDuration {
			minDuration = event.Duration
		}

		// Categorize by duration
		durationMs := int(float64(event.Duration.Microseconds()) / 1000.0)
		durationCounts[durationMs]++
	}

	// Sort events by start time
	sort.Slice(events, func(i, j int) bool {
		return events[i].StartTime.Before(events[j].StartTime)
	})

	return types.AnalysisResult{
		Events:         events,
		TotalEvents:    totalEvents,
		TotalDuration:  totalDuration,
		MaxDuration:    maxDuration,
		MinDuration:    minDuration,
		DurationCounts: durationCounts,
	}
}

// AnalyzeOSDOpEvents analyzes OSD operation events and returns statistics
func AnalyzeOSDOpEvents(events []types.OSDOpEvent) types.OSDOpAnalysisResult {
	totalOps := len(events)
	totalLatency := 0.0
	totalInBytes := 0
	totalOutBytes := 0
	maxLatency := 0.0

	// Count requests by latency ranges
	latencyCounts := make(map[string]int)
	latencyRanges := []string{"0-1ms", "1-2ms", "2-3ms", "3-4ms", "4-5ms", "5-10ms", "10ms+"}

	for _, rangeStr := range latencyRanges {
		latencyCounts[rangeStr] = 0
	}

	for _, event := range events {
		totalLatency += event.Latency
		totalInBytes += event.InBytes
		totalOutBytes += event.OutBytes

		// Update max latency
		if event.Latency > maxLatency {
			maxLatency = event.Latency
		}

		// Categorize by latency
		switch {
		case event.Latency < 1:
			latencyCounts["0-1ms"]++
		case event.Latency < 2:
			latencyCounts["1-2ms"]++
		case event.Latency < 3:
			latencyCounts["2-3ms"]++
		case event.Latency < 4:
			latencyCounts["3-4ms"]++
		case event.Latency < 5:
			latencyCounts["4-5ms"]++
		case event.Latency < 10:
			latencyCounts["5-10ms"]++
		default:
			latencyCounts["10ms+"]++
		}
	}

	// Sort events by timestamp
	sort.Slice(events, func(i, j int) bool {
		return events[i].Timestamp.Before(events[j].Timestamp)
	})

	return types.OSDOpAnalysisResult{
		Events:        events,
		TotalOps:      totalOps,
		TotalLatency:  totalLatency,
		MaxLatency:    maxLatency,
		TotalInBytes:  totalInBytes,
		TotalOutBytes: totalOutBytes,
		LatencyCounts: latencyCounts,
	}
}

// AnalyzeTransactionEvents analyzes transaction events and returns statistics
func AnalyzeTransactionEvents(events []types.TransactionEvent) types.TransactionAnalysisResult {
	totalTransactions := len(events)
	totalDuration := time.Duration(0)
	maxDuration := time.Duration(0)
	minDuration := time.Duration(1000000000) // Start with a large value

	// Count requests by all durations present in data
	durationCounts := make(map[int]int)

	for _, event := range events {
		totalDuration += event.TotalDuration
		if event.TotalDuration > maxDuration {
			maxDuration = event.TotalDuration
		}
		if event.TotalDuration < minDuration {
			minDuration = event.TotalDuration
		}

		// Categorize by duration
		durationMs := int(float64(event.TotalDuration.Microseconds()) / 1000.0)
		durationCounts[durationMs]++
	}

	// Sort events by start time
	sort.Slice(events, func(i, j int) bool {
		return events[i].StartTime.Before(events[j].StartTime)
	})

	return types.TransactionAnalysisResult{
		Events:            events,
		TotalTransactions: totalTransactions,
		TotalDuration:     totalDuration,
		MaxDuration:       maxDuration,
		MinDuration:       minDuration,
		DurationCounts:    durationCounts,
	}
}

// AnalyzeMetadataSyncEvents analyzes metadata sync events and returns statistics
func AnalyzeMetadataSyncEvents(events []types.MetadataSyncEvent) types.MetadataSyncAnalysisResult {
	totalEvents := len(events)
	totalDuration := time.Duration(0)
	maxDuration := time.Duration(0)
	minDuration := time.Duration(1000000000) // Start with a large value

	// Count requests by all durations present in data
	durationCounts := make(map[int]int)

	for _, event := range events {
		totalDuration += event.Duration
		if event.Duration > maxDuration {
			maxDuration = event.Duration
		}
		if event.Duration < minDuration {
			minDuration = event.Duration
		}

		// Categorize by duration
		durationMs := int(float64(event.Duration.Microseconds()) / 1000.0)
		durationCounts[durationMs]++
	}

	// Sort events by timestamp
	sort.Slice(events, func(i, j int) bool {
		return events[i].Timestamp.Before(events[j].Timestamp)
	})

	return types.MetadataSyncAnalysisResult{
		Events:         events,
		TotalEvents:    totalEvents,
		TotalDuration:  totalDuration,
		MaxDuration:    maxDuration,
		MinDuration:    minDuration,
		DurationCounts: durationCounts,
	}
}
