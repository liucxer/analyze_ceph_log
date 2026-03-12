package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"time"
)

// AIOEvent represents an AIO start or finish event
type AIOEvent struct {
	timestamp time.Time
	blockAddr string
	isStart   bool
}

func main() {
	// Check if a log file is provided
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run analyze_aio.go <log_file>")
		os.Exit(1)
	}

	logFile := os.Args[1]

	// Open the log file
	file, err := os.Open(logFile)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Regular expressions to match start and finish events
	startRegex := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}) .* _aio_log_start (0x[0-9a-fA-F]+)~([0-9a-fA-F]+)`)
	finishRegex := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}) .* _aio_log_finish 1 (0x[0-9a-fA-F]+)~([0-9a-fA-F]+)`)

	// Map to store start events by block address and full range
	type StartEvent struct {
		timestamp time.Time
		rangeStr  string
	}
	startEvents := make(map[string]StartEvent)

	// Variables to track statistics
	totalAIOs := 0
	totalDuration := time.Duration(0)
	maxDuration := time.Duration(0)
	minDuration := time.Duration(1000000000) // Start with a large value

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Check for start event
		if startMatches := startRegex.FindStringSubmatch(line); len(startMatches) == 4 {
			timestampStr := startMatches[1]
			blockAddr := startMatches[2]
			rangeStr := startMatches[2] + "~" + startMatches[3]

			// Parse timestamp
			timestamp, err := time.Parse("2006-01-02 15:04:05.000", timestampStr)
			if err != nil {
				fmt.Printf("Error parsing timestamp: %v\n", err)
				continue
			}

			// Store start event
			startEvents[blockAddr] = StartEvent{
				timestamp: timestamp,
				rangeStr:  rangeStr,
			}

		} else if finishMatches := finishRegex.FindStringSubmatch(line); len(finishMatches) == 4 {
			timestampStr := finishMatches[1]
			blockAddr := finishMatches[2]

			// Parse timestamp
			timestamp, err := time.Parse("2006-01-02 15:04:05.000", timestampStr)
			if err != nil {
				fmt.Printf("Error parsing timestamp: %v\n", err)
				continue
			}

			// Check if we have a corresponding start event
			if startEvent, ok := startEvents[blockAddr]; ok {
				// Calculate duration
				duration := timestamp.Sub(startEvent.timestamp)

				// Update statistics
				totalAIOs++
				totalDuration += duration

				// Update max and min durations
				if duration > maxDuration {
					maxDuration = duration
				}
				if duration < minDuration {
					minDuration = duration
				}

				// Remove the start event from the map
				delete(startEvents, blockAddr)
			} else {
				fmt.Printf("No start event found for block address: %s\n", blockAddr)
			}
		}
	}

	// Re-read the file to collect all events and print in order
	file.Seek(0, 0)
	scanner = bufio.NewScanner(file)
	
	// Map to store all events by block address
	eventsMap := make(map[string]struct {
		startTime time.Time
		endTime   time.Time
		duration  time.Duration
		rangeStr  string
	})

	// First pass to collect all events
	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if startMatches := startRegex.FindStringSubmatch(line); len(startMatches) == 4 {
			timestampStr := startMatches[1]
			blockAddr := startMatches[2]
			rangeStr := startMatches[2] + "~" + startMatches[3]

			timestamp, err := time.Parse("2006-01-02 15:04:05.000", timestampStr)
			if err != nil {
				continue
			}

			eventsMap[blockAddr] = struct {
				startTime time.Time
				endTime   time.Time
				duration  time.Duration
				rangeStr  string
			}{
				startTime: timestamp,
				rangeStr:  rangeStr,
			}

		} else if finishMatches := finishRegex.FindStringSubmatch(line); len(finishMatches) == 4 {
			timestampStr := finishMatches[1]
			blockAddr := finishMatches[2]

			timestamp, err := time.Parse("2006-01-02 15:04:05.000", timestampStr)
			if err != nil {
				continue
			}

			if event, ok := eventsMap[blockAddr]; ok {
				duration := timestamp.Sub(event.startTime)
				eventsMap[blockAddr] = struct {
					startTime time.Time
					endTime   time.Time
					duration  time.Duration
					rangeStr  string
				}{
					startTime: event.startTime,
					endTime:   timestamp,
					duration:  duration,
					rangeStr:  event.rangeStr,
				}
			}
		}
	}

	// Convert map to slice for sorting
	type Event struct {
		startTime time.Time
		endTime   time.Time
		duration  time.Duration
		rangeStr  string
	}
	var events []Event
	for _, event := range eventsMap {
		if event.duration > 0 {
			events = append(events, Event{
				startTime: event.startTime,
				endTime:   event.endTime,
				duration:  event.duration,
				rangeStr:  event.rangeStr,
			})
		}
	}

	// Sort events by start time
	for i := 0; i < len(events); i++ {
		for j := i + 1; j < len(events); j++ {
			if events[i].startTime.After(events[j].startTime) {
				events[i], events[j] = events[j], events[i]
			}
		}
	}

	// Print header
	fmt.Printf("%-25s %-25s %-15s %-20s\n", "Start Time", "End Time", "Duration (ms)", "Range")
	fmt.Println("-----------------------------------------------------------------------------------------")

	// Print events with non-zero duration in sorted order
	for _, event := range events {
		fmt.Printf("%-25s %-25s %-15.3f %-20s\n", 
			event.startTime.Format("2006-01-02 15:04:05.000"),
			event.endTime.Format("2006-01-02 15:04:05.000"),
			float64(event.duration.Microseconds())/1000.0,
			event.rangeStr)
	}

	// Check for any remaining start events
	if len(startEvents) > 0 {
		fmt.Printf("Found %d start events without corresponding finish events:\n", len(startEvents))
		for blockAddr, _ := range startEvents {
			fmt.Printf("  %s\n", blockAddr)
		}
	}

	// Print summary statistics
	fmt.Println("\nSummary:")
	fmt.Printf("Total AIO operations: %d\n", totalAIOs)
	if totalAIOs > 0 {
		averageDuration := totalDuration / time.Duration(totalAIOs)
		fmt.Printf("Average duration: %v\n", averageDuration)
		fmt.Printf("Maximum duration: %v\n", maxDuration)
		fmt.Printf("Minimum duration: %v\n", minDuration)
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}
}