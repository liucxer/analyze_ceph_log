package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

// RepopEvent represents an OSD repop event
type RepopEvent struct {
	startTime time.Time
	endTime   time.Time
	duration  time.Duration
	opID      string
}

func main() {
	// Check if a log file is provided
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run analyze_osd_repop.go <log_file>")
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

	// Regular expressions to match dequeue_op and repop_commit events
	dequeueRegex := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}) .* dequeue_op .* osd_repop\(([^\)]+)\)`)
	commitRegex := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}) .* repop_commit on op osd_repop\(([^\)]+)\)`)

	// Map to store start events by op ID
	startEvents := make(map[string]time.Time)

	// Variables to track statistics
	totalRepops := 0
	totalDuration := time.Duration(0)
	maxDuration := time.Duration(0)
	minDuration := time.Duration(1000000000) // Start with a large value

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Check for dequeue_op event (start)
		if dequeueMatches := dequeueRegex.FindStringSubmatch(line); len(dequeueMatches) == 3 {
			timestampStr := dequeueMatches[1]
			opInfo := dequeueMatches[2]

			// Extract op ID from osd_repop info
			// Format: client.177801796.0:14337467 113.e7f e35626/35623
			// Use first three parts as op ID
			parts := make([]string, 0)
			for i, part := range strings.Fields(opInfo) {
				if i < 3 {
					parts = append(parts, part)
				} else {
					break
				}
			}
			opID := strings.Join(parts, " ")

			// Parse timestamp
			timestamp, err := time.Parse("2006-01-02 15:04:05.000", timestampStr)
			if err != nil {
				fmt.Printf("Error parsing timestamp: %v\n", err)
				continue
			}

			// Store start event
			startEvents[opID] = timestamp

			// Check for repop_commit event (end)
		} else if commitMatches := commitRegex.FindStringSubmatch(line); len(commitMatches) == 3 {
			timestampStr := commitMatches[1]
			opInfo := commitMatches[2]

			// Extract op ID from osd_repop info
			// Format: client.177801796.0:14337467 113.e7f e35626/35623 113:fe73d84f:::...
			// Use first three parts as op ID
			parts := make([]string, 0)
			for i, part := range strings.Fields(opInfo) {
				if i < 3 {
					parts = append(parts, part)
				} else {
					break
				}
			}
			opID := strings.Join(parts, " ")

			// Parse timestamp
			timestamp, err := time.Parse("2006-01-02 15:04:05.000", timestampStr)
			if err != nil {
				fmt.Printf("Error parsing timestamp: %v\n", err)
				continue
			}

			// Check if we have a corresponding start event
			if startTime, ok := startEvents[opID]; ok {
				// Calculate duration
				duration := timestamp.Sub(startTime)

				// Update statistics
				totalRepops++
				totalDuration += duration

				// Update max and min durations
				if duration > maxDuration {
					maxDuration = duration
				}
				if duration < minDuration {
					minDuration = duration
				}

				// Remove the start event from the map
				delete(startEvents, opID)
			}
		}
	}

	// Re-read the file to collect all events and print in order
	file.Seek(0, 0)
	scanner = bufio.NewScanner(file)

	// Map to store all events by op ID
	eventsMap := make(map[string]struct {
		startTime time.Time
		endTime   time.Time
		duration  time.Duration
		opID      string
	})

	// First pass to collect all events
	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if dequeueMatches := dequeueRegex.FindStringSubmatch(line); len(dequeueMatches) == 3 {
			timestampStr := dequeueMatches[1]
			opInfo := dequeueMatches[2]

			// Use first three parts as op ID
			parts := make([]string, 0)
			for i, part := range strings.Fields(opInfo) {
				if i < 3 {
					parts = append(parts, part)
				} else {
					break
				}
			}
			opID := strings.Join(parts, " ")

			timestamp, err := time.Parse("2006-01-02 15:04:05.000", timestampStr)
			if err != nil {
				continue
			}

			eventsMap[opID] = struct {
				startTime time.Time
				endTime   time.Time
				duration  time.Duration
				opID      string
			}{
				startTime: timestamp,
				opID:      opID,
			}

		} else if commitMatches := commitRegex.FindStringSubmatch(line); len(commitMatches) == 3 {
			timestampStr := commitMatches[1]
			opInfo := commitMatches[2]

			// Use first three parts as op ID
			parts := make([]string, 0)
			for i, part := range strings.Fields(opInfo) {
				if i < 3 {
					parts = append(parts, part)
				} else {
					break
				}
			}
			opID := strings.Join(parts, " ")

			timestamp, err := time.Parse("2006-01-02 15:04:05.000", timestampStr)
			if err != nil {
				continue
			}

			if event, ok := eventsMap[opID]; ok {
				duration := timestamp.Sub(event.startTime)
				eventsMap[opID] = struct {
					startTime time.Time
					endTime   time.Time
					duration  time.Duration
					opID      string
				}{
					startTime: event.startTime,
					endTime:   timestamp,
					duration:  duration,
					opID:      event.opID,
				}
			}
		}
	}

	// Convert map to slice for sorting
	type Event struct {
		startTime time.Time
		endTime   time.Time
		duration  time.Duration
		opID      string
	}
	var events []Event
	for _, event := range eventsMap {
		events = append(events, Event{
			startTime: event.startTime,
			endTime:   event.endTime,
			duration:  event.duration,
			opID:      event.opID,
		})
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
	fmt.Printf("%-25s %-25s %-15s %-40s\n", "Start Time", "End Time", "Duration (ms)", "OP ID")
	fmt.Println("----------------------------------------------------------------------------------------------------")

	// Print events
	for _, event := range events {
		if event.duration > 0 {
			fmt.Printf("%-25s %-25s %-15.3f %-40s\n",
				event.startTime.Format("2006-01-02 15:04:05.000"),
				event.endTime.Format("2006-01-02 15:04:05.000"),
				float64(event.duration.Microseconds())/1000.0,
				event.opID)
		}
	}

	// Print summary statistics
	fmt.Println("\nSummary:")
	fmt.Printf("Total repop operations: %d\n", totalRepops)
	if totalRepops > 0 {
		averageDuration := totalDuration / time.Duration(totalRepops)
		fmt.Printf("Average duration: %v\n", averageDuration)
		fmt.Printf("Maximum duration: %v\n", maxDuration)
		fmt.Printf("Minimum duration: %v\n", minDuration)
	}

	// Check for any remaining start events
	if len(startEvents) > 0 {
		fmt.Printf("Found %d start events without corresponding finish events:\n", len(startEvents))
		for opID, _ := range startEvents {
			fmt.Printf("  %s\n", opID)
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}
}
