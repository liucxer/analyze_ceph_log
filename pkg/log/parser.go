package log

import (
	"bufio"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/liucxer/analyze_ceph_log/pkg/types"
)

// ParseAIOEvents parses AIO events from log file
func ParseAIOEvents(scanner *bufio.Scanner) ([]types.Event, error) {
	startRegex := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}) .* _aio_log_start (0x[0-9a-fA-F]+)~([0-9a-fA-F]+)`)
	finishRegex := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}) .* _aio_log_finish 1 (0x[0-9a-fA-F]+)~([0-9a-fA-F]+)`)

	// Map to store all events by a unique key (using counter to ensure uniqueness)
	eventsMap := make(map[string]struct {
		startTime time.Time
		endTime   time.Time
		duration  time.Duration
		rangeStr  string
		blockAddr string
	})

	// Counter to ensure unique keys
	counter := 0

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

			// Use a unique key with counter to ensure uniqueness
			counter++
			key := blockAddr + "_" + timestampStr + "_" + string(rune(counter))
			eventsMap[key] = struct {
				startTime time.Time
				endTime   time.Time
				duration  time.Duration
				rangeStr  string
				blockAddr string
			}{
				startTime: timestamp,
				rangeStr:  rangeStr,
				blockAddr: blockAddr,
			}

		} else if finishMatches := finishRegex.FindStringSubmatch(line); len(finishMatches) == 4 {
			timestampStr := finishMatches[1]
			blockAddr := finishMatches[2]

			timestamp, err := time.Parse("2006-01-02 15:04:05.000", timestampStr)
			if err != nil {
				continue
			}

			// Find the matching start event by block address
			for key, event := range eventsMap {
				if event.blockAddr == blockAddr && event.endTime.IsZero() {
					duration := timestamp.Sub(event.startTime)
					eventsMap[key] = struct {
						startTime time.Time
						endTime   time.Time
						duration  time.Duration
						rangeStr  string
						blockAddr string
					}{
						startTime: event.startTime,
						endTime:   timestamp,
						duration:  duration,
						rangeStr:  event.rangeStr,
						blockAddr: event.blockAddr,
					}
					break
				}
			}
		}
	}

	// Convert map to slice for sorting
	var events []types.Event
	for _, event := range eventsMap {
		if !event.endTime.IsZero() && event.duration >= 0 {
			events = append(events, types.Event{
				StartTime: event.startTime,
				EndTime:   event.endTime,
				Duration:  event.duration,
				RangeStr:  event.rangeStr,
			})
		}
	}

	return events, nil
}

// ParseRepopEvents parses OSD repop events from log file
func ParseRepopEvents(scanner *bufio.Scanner) ([]types.Event, error) {
	dequeueRegex := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}) .* dequeue_op .* osd_repop\(([\s\S]+?)\)`)
	commitRegex := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}) .* repop_commit on op osd_repop\(([\s\S]+?)\)`)

	// Map to store all events by op ID
	eventsMap := make(map[string]struct {
		startTime time.Time
		endTime   time.Time
		duration  time.Duration
		opID      string
	})

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
	var events []types.Event
	for _, event := range eventsMap {
		if event.duration >= 0 {
			events = append(events, types.Event{
				StartTime: event.startTime,
				EndTime:   event.endTime,
				Duration:  event.duration,
				OpID:      event.opID,
			})
		}
	}

	return events, nil
}

// ParseOSDOpEvents parses OSD operation events from log file
func ParseOSDOpEvents(scanner *bufio.Scanner) ([]types.OSDOpEvent, error) {
	osdOpRegex := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}) .* log_op_stats osd_op\(([\s\S]+?)\) .* inb (\d+) outb (\d+) lat ([\d\.]+)`)

	var events []types.OSDOpEvent

	for scanner.Scan() {
		line := scanner.Text()

		// Check for osd_op event
		if matches := osdOpRegex.FindStringSubmatch(line); len(matches) == 6 {
			timestampStr := matches[1]
			opInfo := matches[2]
			inBytesStr := matches[3]
			outBytesStr := matches[4]
			latencyStr := matches[5]

			// Parse timestamp
			timestamp, err := time.Parse("2006-01-02 15:04:05.000", timestampStr)
			if err != nil {
				continue
			}

			// Parse input/output bytes
			inBytes, err := strconv.Atoi(inBytesStr)
			if err != nil {
				continue
			}

			outBytes, err := strconv.Atoi(outBytesStr)
			if err != nil {
				continue
			}

			// Parse latency
			latency, err := strconv.ParseFloat(latencyStr, 64)
			if err != nil {
				continue
			}
			// Convert seconds to milliseconds
			latency = latency * 1000.0

			// Extract opID, pgID, object, opType, rangeStr from opInfo
			// First, split by spaces but preserve brackets content
			var parts []string
			current := ""
			inBrackets := false
			
			for _, char := range opInfo {
				if char == '[' {
					inBrackets = true
					if current != "" {
						parts = append(parts, current)
						current = ""
					}
					current += string(char)
				} else if char == ']' {
					current += string(char)
					parts = append(parts, current)
					current = ""
					inBrackets = false
				} else if char == ' ' && !inBrackets {
					if current != "" {
						parts = append(parts, current)
						current = ""
					}
				} else {
					current += string(char)
				}
			}
			if current != "" {
				parts = append(parts, current)
			}

			var opID, pgID, object, opType, rangeStr string
			if len(parts) >= 4 {
				opID = parts[0]
				pgID = parts[1]
				object = parts[2]
				
				// Extract opType and rangeStr from the brackets
				for _, part := range parts {
					if len(part) > 2 && part[0] == '[' && part[len(part)-1] == ']' {
						bracketContent := part[1 : len(part)-1]
						bracketParts := regexp.MustCompile(`\s+`).Split(bracketContent, -1)
						if len(bracketParts) >= 2 {
							opType = bracketParts[0]
							rangeStr = bracketParts[1]
						}
						break
					}
				}
			}

			// Create event
			event := types.OSDOpEvent{
				Timestamp: timestamp,
				OpID:      opID,
				PgID:      pgID,
				Object:    object,
				OpType:    opType,
				RangeStr:  rangeStr,
				InBytes:   inBytes,
				OutBytes:  outBytes,
				Latency:   latency,
			}

			events = append(events, event)
		}
	}

	return events, nil
}
