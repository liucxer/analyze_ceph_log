package log

import (
	"bufio"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/liucxer/analyze_ceph_log/pkg/types"
)

// ParseAIOEvents parses AIO events from log file
func ParseAIOEvents(scanner *bufio.Scanner) ([]types.Event, error) {
	// Map to store all events by range string (which is unique for each AIO operation)
	eventsMap := make(map[string]struct {
		startTime time.Time
		endTime   time.Time
		duration  time.Duration
		rangeStr  string
		length    int
		blockType string
	})

	for scanner.Scan() {
		line := scanner.Text()

		// Check for AIO start event
		if strings.Contains(line, "_aio_log_start") {
			// Extract timestamp
			timestampStr := ""
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				timestampStr = parts[0] + " " + parts[1]
			}

			timestamp, err := time.Parse("2006-01-02 15:04:05.000", timestampStr)
			if err != nil {
				continue
			}

			// Extract device path from bdev()
			devicePath := ""
			bdevStart := strings.Index(line, "bdev(")
			if bdevStart != -1 {
				bdevEnd := strings.Index(line[bdevStart:], ")")
				if bdevEnd != -1 {
					bdevStr := line[bdevStart : bdevStart+bdevEnd+1]
					// Extract path from bdev string
					pathStart := strings.LastIndex(bdevStr, " ")
					if pathStart != -1 {
						devicePath = bdevStr[pathStart+1 : len(bdevStr)-1]
					}
				}
			}

			// Extract range string
			rangeStr := ""
			startIdx := strings.Index(line, "_aio_log_start")
			if startIdx != -1 {
				rangeParts := strings.Fields(line[startIdx:])
				if len(rangeParts) >= 2 {
					rangeStr = rangeParts[1]
				}
			}

			if rangeStr == "" {
				continue
			}

			// Extract length from range string
			length := 0
			rangeParts := strings.Split(rangeStr, "~")
			if len(rangeParts) == 2 {
				lengthHex := rangeParts[1]
				lengthVal, err := strconv.ParseInt(lengthHex, 16, 64)
				if err == nil {
					length = int(lengthVal)
				}
			}

			// Determine block type
			blockType := "block"
			if strings.Contains(devicePath, "block.wal") {
				blockType = "block.wal"
			} else if strings.Contains(devicePath, "block.db") {
				blockType = "block.db"
			}

			// Use rangeStr as key
			eventsMap[rangeStr] = struct {
				startTime time.Time
				endTime   time.Time
				duration  time.Duration
				rangeStr  string
				length    int
				blockType string
			}{
				startTime: timestamp,
				rangeStr:  rangeStr,
				length:    length,
				blockType: blockType,
			}

		} else if strings.Contains(line, "_aio_log_finish") {
			// Extract timestamp
			timestampStr := ""
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				timestampStr = parts[0] + " " + parts[1]
			}

			timestamp, err := time.Parse("2006-01-02 15:04:05.000", timestampStr)
			if err != nil {
				continue
			}

			// Extract range string
			rangeStr := ""
			finishIdx := strings.Index(line, "_aio_log_finish")
			if finishIdx != -1 {
				rangeParts := strings.Fields(line[finishIdx:])
				if len(rangeParts) >= 3 {
					rangeStr = rangeParts[2]
				}
			}

			if rangeStr == "" {
				continue
			}

			// Find the matching start event
			if event, ok := eventsMap[rangeStr]; ok && event.endTime.IsZero() {
				duration := timestamp.Sub(event.startTime)
				eventsMap[rangeStr] = struct {
					startTime time.Time
					endTime   time.Time
					duration  time.Duration
					rangeStr  string
					length    int
					blockType string
				}{
					startTime: event.startTime,
					endTime:   timestamp,
					duration:  duration,
					rangeStr:  event.rangeStr,
					length:    event.length,
					blockType: event.blockType,
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
				Length:    event.length,
				BlockType: event.blockType,
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
	repopReplyRegex := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}) .* osd_repop_reply\(([\s\S]+?)\) v2`)

	// Map to store OSD op events by opID
	eventsMap := make(map[string]types.OSDOpEvent)
	// Map to store repop reply times by opID
	repopReplyMap := make(map[string][]time.Time)

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

			eventsMap[opID] = event
		} else if matches := repopReplyRegex.FindStringSubmatch(line); len(matches) == 3 {
			timestampStr := matches[1]
			opInfo := matches[2]

			// Parse timestamp
			timestamp, err := time.Parse("2006-01-02 15:04:05.000", timestampStr)
			if err != nil {
				continue
			}

			// Extract opID from opInfo (first part)
			parts := strings.Fields(opInfo)
			if len(parts) > 0 {
				opID := parts[0]
				// Add timestamp to repopReplyMap
				repopReplyMap[opID] = append(repopReplyMap[opID], timestamp)
			}
		}
	}

	// Process repop reply times and calculate durations
	for opID, event := range eventsMap {
		if replyTimes, ok := repopReplyMap[opID]; ok && len(replyTimes) >= 2 {
			// Sort reply times to ensure correct order
			sort.Slice(replyTimes, func(i, j int) bool {
				return replyTimes[i].Before(replyTimes[j])
			})

			// Calculate time difference between first repop reply and osd_op start
			// and between first and second repop replies
			firstReply := math.Abs(float64(replyTimes[0].Sub(event.Timestamp).Milliseconds()))
			secondReply := math.Abs(float64(replyTimes[1].Sub(replyTimes[0]).Milliseconds()))
			event.FirstRepopReply = firstReply
			event.SecondRepopReply = secondReply
			eventsMap[opID] = event
		}
	}

	// Convert map to slice
	var events []types.OSDOpEvent
	for _, event := range eventsMap {
		events = append(events, event)
	}

	return events, nil
}

// ParseTransactionEvents parses transaction events from log file
func ParseTransactionEvents(scanner *bufio.Scanner) ([]types.TransactionEvent, error) {
	// Map to store transaction events by TID
	transactionsMap := make(map[string]struct {
		TID             string
		StartTime       time.Time
		IssueTime       time.Time
		FirstReplyTime  time.Time
		SecondReplyTime time.Time
		CompleteTime    time.Time
		OpID            string
		Object          string
		RangeStr        string
	})

	for scanner.Scan() {
		line := scanner.Text()

		// Check for new_repop event (transaction start)
		if strings.Contains(line, "new_repop") && strings.Contains(line, "rep_tid") {
			// Extract timestamp
			timestampStr := ""
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				timestampStr = parts[0] + " " + parts[1]
			}

			timestamp, err := time.Parse("2006-01-02 15:04:05.000", timestampStr)
			if err != nil {
				continue
			}

			// Extract TID
			tid := ""
			tidStart := strings.Index(line, "rep_tid")
			if tidStart != -1 {
				tidParts := strings.Fields(line[tidStart:])
				if len(tidParts) >= 2 {
					tid = tidParts[1]
				}
			}

			if tid == "" {
				continue
			}

			// Extract OpID
			opID := ""
			opStart := strings.Index(line, "osd_op(")
			if opStart != -1 {
				opEnd := strings.Index(line[opStart:], ")")
				if opEnd != -1 {
					opStr := line[opStart+7 : opStart+opEnd] // Skip "osd_op("
					opParts := strings.Fields(opStr)
					if len(opParts) >= 1 {
						opID = opParts[0]
					}
				}
			}

			// Extract object and range
			object := ""
			rangeStr := ""
			if opStart != -1 {
				opEnd := strings.Index(line[opStart:], ")")
				if opEnd != -1 {
					opStr := line[opStart+7 : opStart+opEnd] // Skip "osd_op("
					// Extract object (third field)
					opParts := strings.Fields(opStr)
					if len(opParts) >= 3 {
						object = opParts[2]
					}
					// Extract range from brackets
					bracketStart := strings.Index(opStr, "[")
					bracketEnd := strings.Index(opStr, "]")
					if bracketStart != -1 && bracketEnd != -1 {
						bracketContent := opStr[bracketStart+1 : bracketEnd]
						bracketParts := strings.Fields(bracketContent)
						if len(bracketParts) >= 2 {
							rangeStr = bracketParts[1]
						}
					}
				}
			}

			// Store transaction
			transactionsMap[tid] = struct {
				TID             string
				StartTime       time.Time
				IssueTime       time.Time
				FirstReplyTime  time.Time
				SecondReplyTime time.Time
				CompleteTime    time.Time
				OpID            string
				Object          string
				RangeStr        string
			}{
				TID:       tid,
				StartTime: timestamp,
				OpID:      opID,
				Object:    object,
				RangeStr:  rangeStr,
			}

		} else if strings.Contains(line, "issue_repop") && strings.Contains(line, "rep_tid") {
			// Extract timestamp
			timestampStr := ""
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				timestampStr = parts[0] + " " + parts[1]
			}

			timestamp, err := time.Parse("2006-01-02 15:04:05.000", timestampStr)
			if err != nil {
				continue
			}

			// Extract TID
			tid := ""
			tidStart := strings.Index(line, "rep_tid")
			if tidStart != -1 {
				tidParts := strings.Fields(line[tidStart:])
				if len(tidParts) >= 2 {
					tid = tidParts[1]
				}
			}

			if tid == "" {
				continue
			}

			// Update transaction with issue time
			if transaction, ok := transactionsMap[tid]; ok {
				transaction.IssueTime = timestamp
				transactionsMap[tid] = transaction
			}

		} else if strings.Contains(line, "do_repop_reply") && strings.Contains(line, "tid") {
			// Extract timestamp
			timestampStr := ""
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				timestampStr = parts[0] + " " + parts[1]
			}

			timestamp, err := time.Parse("2006-01-02 15:04:05.000", timestampStr)
			if err != nil {
				continue
			}

			// Extract TID
			tid := ""
			tidStart := strings.Index(line, "tid")
			if tidStart != -1 {
				tidParts := strings.Fields(line[tidStart:])
				if len(tidParts) >= 2 {
					tid = tidParts[1]
				}
			}

			if tid == "" {
				continue
			}

			// Update transaction with reply time
			if transaction, ok := transactionsMap[tid]; ok {
				if transaction.FirstReplyTime.IsZero() {
					transaction.FirstReplyTime = timestamp
				} else {
					transaction.SecondReplyTime = timestamp
				}
				transactionsMap[tid] = transaction
			}

		} else if strings.Contains(line, "repop_all_committed") && strings.Contains(line, "repop tid") {
			// Extract timestamp
			timestampStr := ""
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				timestampStr = parts[0] + " " + parts[1]
			}

			timestamp, err := time.Parse("2006-01-02 15:04:05.000", timestampStr)
			if err != nil {
				continue
			}

			// Extract TID
			tid := ""
			tidStart := strings.Index(line, "repop tid")
			if tidStart != -1 {
				tidParts := strings.Fields(line[tidStart:])
				if len(tidParts) >= 3 {
					tid = tidParts[2]
				}
			}

			if tid == "" {
				continue
			}

			// Update transaction with complete time
			if transaction, ok := transactionsMap[tid]; ok {
				transaction.CompleteTime = timestamp
				transactionsMap[tid] = transaction
			}
		}
	}

	// Convert map to slice and calculate durations
	var events []types.TransactionEvent
	for _, transaction := range transactionsMap {
		// Calculate durations
		totalDuration := time.Duration(0)
		issueDuration := time.Duration(0)
		firstReplyDuration := time.Duration(0)
		secondReplyDuration := time.Duration(0)

		if !transaction.CompleteTime.IsZero() && !transaction.StartTime.IsZero() {
			totalDuration = transaction.CompleteTime.Sub(transaction.StartTime)
		}

		if !transaction.IssueTime.IsZero() && !transaction.StartTime.IsZero() {
			issueDuration = transaction.IssueTime.Sub(transaction.StartTime)
		}

		if !transaction.FirstReplyTime.IsZero() && !transaction.IssueTime.IsZero() {
			firstReplyDuration = transaction.FirstReplyTime.Sub(transaction.IssueTime)
		}

		if !transaction.SecondReplyTime.IsZero() && !transaction.FirstReplyTime.IsZero() {
			secondReplyDuration = transaction.SecondReplyTime.Sub(transaction.FirstReplyTime)
		}

		// Create transaction event
		event := types.TransactionEvent{
			TID:                 transaction.TID,
			StartTime:           transaction.StartTime,
			IssueTime:           transaction.IssueTime,
			FirstReplyTime:      transaction.FirstReplyTime,
			SecondReplyTime:     transaction.SecondReplyTime,
			CompleteTime:        transaction.CompleteTime,
			TotalDuration:       totalDuration,
			IssueDuration:       issueDuration,
			FirstReplyDuration:  firstReplyDuration,
			SecondReplyDuration: secondReplyDuration,
			OpID:                transaction.OpID,
			Object:              transaction.Object,
			RangeStr:            transaction.RangeStr,
		}

		events = append(events, event)
	}

	return events, nil
}
