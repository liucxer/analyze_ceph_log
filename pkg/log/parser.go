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
	enqueueRegex := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}) .* enqueue_op .* osd_repop\(([\s\S]+?)\)(?: v\d+)?`)
	commitRegex := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}) .* osd_repop\(([\s\S]+?)\)(?: v\d+)?, sending commit to`)

	// Map to store all events by op ID
	eventsMap := make(map[string]struct {
		startTime time.Time
		endTime   time.Time
		duration  time.Duration
		opID      string
	})

	for scanner.Scan() {
		line := scanner.Text()

		if enqueueMatches := enqueueRegex.FindStringSubmatch(line); len(enqueueMatches) == 3 {
			timestampStr := enqueueMatches[1]
			opInfo := enqueueMatches[2]

			// Extract client ID and op ID (first part)
			parts := strings.Fields(opInfo)
			if len(parts) == 0 {
				continue
			}
			opID := parts[0]

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

			// Extract client ID and op ID (first part)
			parts := strings.Fields(opInfo)
			if len(parts) == 0 {
				continue
			}
			opID := parts[0]

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

// ParseOSDOpEvents parses OSD op events from log file
func ParseOSDOpEventsV2(scanner *bufio.Scanner) ([]types.Event, error) {
	enqueueRegex := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}) .* enqueue_op .* osd_op\(([\s\S]+?)\)(?: v\d+)?`)
	replyRegex := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}) .* sending reply on osd_op\(([\s\S]+?)\)(?: v\d+)?`)

	// Map to store all events by op ID
	eventsMap := make(map[string]struct {
		startTime time.Time
		endTime   time.Time
		duration  time.Duration
		opID      string
	})

	for scanner.Scan() {
		line := scanner.Text()

		if enqueueMatches := enqueueRegex.FindStringSubmatch(line); len(enqueueMatches) == 3 {
			timestampStr := enqueueMatches[1]
			opInfo := enqueueMatches[2]

			// Extract client ID and op ID for osd_op
			clientOp := ""
			parts := strings.Fields(opInfo)
			if len(parts) > 0 {
				clientOp = parts[0]
			}

			timestamp, err := time.Parse("2006-01-02 15:04:05.000", timestampStr)
			if err != nil {
				continue
			}

			eventsMap[clientOp] = struct {
				startTime time.Time
				endTime   time.Time
				duration  time.Duration
				opID      string
			}{
				startTime: timestamp,
				opID:      clientOp,
			}
		} else if replyMatches := replyRegex.FindStringSubmatch(line); len(replyMatches) == 3 {
			timestampStr := replyMatches[1]
			opInfo := replyMatches[2]

			// Extract client ID and op ID for osd_op
			clientOp := ""
			parts := strings.Fields(opInfo)
			if len(parts) > 0 {
				clientOp = parts[0]
			}

			timestamp, err := time.Parse("2006-01-02 15:04:05.000", timestampStr)
			if err != nil {
				continue
			}

			if event, ok := eventsMap[clientOp]; ok {
				duration := timestamp.Sub(event.startTime)
				eventsMap[clientOp] = struct {
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

// ParseMetadataSyncEvents parses metadata sync events from log file
func ParseMetadataSyncEvents(scanner *bufio.Scanner) ([]types.MetadataSyncEvent, error) {
	var events []types.MetadataSyncEvent

	for scanner.Scan() {
		line := scanner.Text()

		// Check for metadata sync event
		if strings.Contains(line, "_kv_sync_thread committed") {
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

			// Extract committed and cleaned values
			committed := 0
			cleaned := 0

			for i, part := range parts {
				if part == "committed" && i+1 < len(parts) {
					committed, _ = strconv.Atoi(parts[i+1])
				} else if part == "cleaned" && i+1 < len(parts) {
					cleaned, _ = strconv.Atoi(parts[i+1])
				}
			}

			// Extract duration
			duration := time.Duration(0)
			flushTime := time.Duration(0)
			kvCommitTime := time.Duration(0)

			// Find duration in the line
			inIdx := -1
			for i, part := range parts {
				if part == "in" && i+1 < len(parts) {
					inIdx = i
					break
				}
			}

			if inIdx != -1 && inIdx+1 < len(parts) {
				// Extract duration string (e.g., "0.000328454s")
				durationStr := parts[inIdx+1]
				durationStr = strings.TrimSuffix(durationStr, "s")
				durationSec, err := strconv.ParseFloat(durationStr, 64)
				if err == nil {
					duration = time.Duration(durationSec * float64(time.Second))
				}

				// Extract flush and kv commit times from parentheses
				if inIdx+2 < len(parts) && strings.Contains(parts[inIdx+2], "(") {
					// Collect all parts until we find the closing parenthesis
					var parenthesesParts []string
					for i := inIdx + 2; i < len(parts); i++ {
						parenthesesParts = append(parenthesesParts, parts[i])
						if strings.Contains(parts[i], ")") {
							break
						}
					}

					// Join the parts to get the full parentheses content
					parenthesesContent := strings.Join(parenthesesParts, " ")
					parenthesesContent = strings.TrimPrefix(parenthesesContent, "(")
					parenthesesContent = strings.TrimSuffix(parenthesesContent, ")")

					// Split into flush and kv commit parts
					parts := strings.Split(parenthesesContent, " + ")
					if len(parts) == 2 {
						// Parse flush time
						flushStr := parts[0]
						// First remove " flush", then remove "s"
						flushStr = strings.TrimSuffix(flushStr, " flush")
						flushStr = strings.TrimSuffix(flushStr, "s")
						// Remove any trailing whitespace
						flushStr = strings.TrimSpace(flushStr)
						flushSec, err := strconv.ParseFloat(flushStr, 64)
						if err == nil {
							flushTime = time.Duration(flushSec * float64(time.Second))
						}

						// Parse kv commit time
						kvCommitStr := parts[1]
						// First remove " kv commit", then remove "s"
						kvCommitStr = strings.TrimSuffix(kvCommitStr, " kv commit")
						kvCommitStr = strings.TrimSuffix(kvCommitStr, "s")
						// Remove any trailing whitespace
						kvCommitStr = strings.TrimSpace(kvCommitStr)
						kvCommitSec, err := strconv.ParseFloat(kvCommitStr, 64)
						if err == nil {
							kvCommitTime = time.Duration(kvCommitSec * float64(time.Second))
						}
					}
				}
			}

			// Create metadata sync event
			event := types.MetadataSyncEvent{
				Timestamp:    timestamp,
				Committed:    committed,
				Cleaned:      cleaned,
				Duration:     duration,
				FlushTime:    flushTime,
				KVCommitTime: kvCommitTime,
			}

			events = append(events, event)
		}
	}

	return events, nil
}

// ParseClientOpEvents parses client operation events from log file
func ParseClientOpEvents(scanner *bufio.Scanner) ([]types.ClientOpEvent, error) {
	// Map to store client operations by op ID
	eventsMap := make(map[string]types.ClientOpEvent)
	// Map to store repop reply times by op ID
	repopReplyMap := make(map[string][]time.Time)

	for scanner.Scan() {
		line := scanner.Text()

		// Check for osd_op event (client operation)
		if strings.Contains(line, "osd_op(client.") {
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

			// Extract client ID and op ID
			clientID := ""
			opID := ""
			clientOp := ""
			pgID := ""
			object := ""
			dataSize := 0

			// Find osd_op section
			opStart := strings.Index(line, "osd_op(")
			if opStart != -1 {
				opEnd := strings.Index(line[opStart:], ")")
				if opEnd != -1 {
					opStr := line[opStart+7 : opStart+opEnd] // Skip "osd_op("
					opParts := strings.Fields(opStr)
					if len(opParts) >= 3 {
						// Extract client ID and op ID
						clientOp = opParts[0]
						clientParts := strings.Split(clientOp, ":")
						if len(clientParts) == 2 {
							clientID = clientParts[0]
							opID = clientParts[1]
						}
						pgID = opParts[1]
						object = opParts[2]

						// Extract data size from range string
						for _, part := range opParts {
							if strings.Contains(part, "~") {
								rangeParts := strings.Split(part, "~")
								if len(rangeParts) == 2 {
									sizeStr := rangeParts[1]
									size, err := strconv.Atoi(sizeStr)
									if err == nil {
										dataSize = size
									}
								}
								break
							}
						}
					}
				}
			}

			if clientID != "" && opID != "" && clientOp != "" {
				eventsMap[clientOp] = types.ClientOpEvent{
					Timestamp: timestamp,
					ClientID:  clientID,
					OpID:      opID,
					PGID:      pgID,
					Object:    object,
					DataSize:  dataSize,
				}
			}

		} else if strings.Contains(line, "osd_repop_reply(client.") {
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

			// Extract client ID and op ID
			clientOp := ""
			replyStart := strings.Index(line, "osd_repop_reply(")
			if replyStart != -1 {
				replyEnd := strings.Index(line[replyStart:], ")")
				if replyEnd != -1 {
					replyStr := line[replyStart+17 : replyStart+replyEnd] // Skip "osd_repop_reply("
					replyParts := strings.Fields(replyStr)
					if len(replyParts) >= 1 {
						clientOp = replyParts[0]
					}
				}
			}

			if clientOp != "" {
				repopReplyMap[clientOp] = append(repopReplyMap[clientOp], timestamp)
			}

		} else if strings.Contains(line, "log_op_stats osd_op(client.") {
			// Extract client ID and op ID
			clientOp := ""
			latency := 0.0

			// Find osd_op section
			opStart := strings.Index(line, "osd_op(")
			if opStart != -1 {
				opEnd := strings.Index(line[opStart:], ")")
				if opEnd != -1 {
					opStr := line[opStart+7 : opStart+opEnd] // Skip "osd_op("
					opParts := strings.Fields(opStr)
					if len(opParts) >= 1 {
						clientOp = opParts[0]
					}
				}
			}

			// Extract latency
			parts := strings.Fields(line)
			for i, part := range parts {
				if part == "lat" && i+1 < len(parts) {
					latencyStr := parts[i+1]
					latency, _ = strconv.ParseFloat(latencyStr, 64)
					// Convert to milliseconds
					latency *= 1000.0
					break
				}
			}

			if clientOp != "" {
				if event, ok := eventsMap[clientOp]; ok {
					event.TotalLatency = latency
					eventsMap[clientOp] = event
				}
			}
		}
	}

	// Process repop reply times
	for clientOp, event := range eventsMap {
		if replyTimes, ok := repopReplyMap[clientOp]; ok {
			// Sort reply times
			sort.Slice(replyTimes, func(i, j int) bool {
				return replyTimes[i].Before(replyTimes[j])
			})

			// Calculate reply latencies
			if len(replyTimes) >= 1 {
				firstReplyLatency := math.Abs(float64(replyTimes[0].Sub(event.Timestamp).Milliseconds()))
				event.FirstReplyLatency = firstReplyLatency
			}

			if len(replyTimes) >= 2 {
				secondReplyLatency := math.Abs(float64(replyTimes[1].Sub(replyTimes[0]).Milliseconds()))
				event.SecondReplyLatency = secondReplyLatency
			}

			// Calculate local processing time
			if event.TotalLatency > 0 && event.FirstReplyLatency > 0 && event.SecondReplyLatency > 0 {
				event.LocalProcessingTime = event.TotalLatency - (event.FirstReplyLatency + event.SecondReplyLatency)
				if event.LocalProcessingTime < 0 {
					event.LocalProcessingTime = 0
				}
			}

			eventsMap[clientOp] = event
		}
	}

	// Convert map to slice
	var events []types.ClientOpEvent
	for _, event := range eventsMap {
		events = append(events, event)
	}

	return events, nil
}

// ParseDequeueEvents parses dequeue operation events from log file
func ParseDequeueEvents(scanner *bufio.Scanner) ([]types.DequeueEvent, error) {
	var events []types.DequeueEvent

	for scanner.Scan() {
		line := scanner.Text()

		// Check for dequeue operation event
		if strings.Contains(line, "latency") && strings.Contains(line, "dequeue_op") {
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

			// Extract op type
			opType := ""
			if strings.Contains(line, "osd_op(") {
				opType = "osd_op"
			} else if strings.Contains(line, "osd_repop(") {
				opType = "osd_repop"
			} else if strings.Contains(line, "osd_repop_reply(") {
				opType = "osd_repop_reply"
			} else if strings.Contains(line, "osd_subop(") {
				opType = "osd_subop"
			} else if strings.Contains(line, "pg_update_log_missing(") {
				opType = "pg_update_log_missing"
			}

			// Extract client ID and op ID
			clientID := ""
			opID := ""
			pgID := ""

			// Extract client ID and op ID from osd_op
			if strings.Contains(line, "osd_op(") {
				opStart := strings.Index(line, "osd_op(")
				opEnd := strings.Index(line[opStart:], ")")
				if opEnd != -1 {
					opStr := line[opStart+7 : opStart+opEnd] // Skip "osd_op("
					opParts := strings.Fields(opStr)
					if len(opParts) >= 2 {
						clientOp := opParts[0]
						// Remove any leading parentheses
						clientOp = strings.TrimPrefix(clientOp, "(")
						// Remove any trailing parentheses
						clientOp = strings.TrimSuffix(clientOp, ")")
						clientParts := strings.Split(clientOp, ":")
						if len(clientParts) == 2 {
							clientID = clientParts[0]
							opID = clientParts[1]
						}
						pgID = opParts[1]
					}
				}
			} else if strings.Contains(line, "osd_repop(") {
				// Extract client ID, op ID and pgID from osd_repop
				opStart := strings.Index(line, "osd_repop(")
				opEnd := strings.Index(line[opStart:], ")")
				if opEnd != -1 {
					opStr := line[opStart+9 : opStart+opEnd] // Skip "osd_repop("
					opParts := strings.Fields(opStr)
					if len(opParts) >= 2 {
						clientOp := opParts[0]
						// Remove any leading parentheses
						clientOp = strings.TrimPrefix(clientOp, "(")
						// Remove any trailing parentheses
						clientOp = strings.TrimSuffix(clientOp, ")")
						clientParts := strings.Split(clientOp, ":")
						if len(clientParts) == 2 {
							clientID = clientParts[0]
							opID = clientParts[1]
						}
						pgID = opParts[1]
					}
				}
			} else if strings.Contains(line, "osd_repop_reply(") {
				// Extract client ID and op ID from osd_repop_reply
				opStart := strings.Index(line, "osd_repop_reply(")
				opEnd := strings.Index(line[opStart:], ")")
				if opEnd != -1 {
					opStr := line[opStart+16 : opStart+opEnd] // Skip "osd_repop_reply("
					opParts := strings.Fields(opStr)
					if len(opParts) >= 2 {
						clientOp := opParts[0]
						// Remove any leading parentheses
						clientOp = strings.TrimPrefix(clientOp, "(")
						// Remove any trailing parentheses
						clientOp = strings.TrimSuffix(clientOp, ")")
						clientParts := strings.Split(clientOp, ":")
						if len(clientParts) == 2 {
							clientID = clientParts[0]
							opID = clientParts[1]
						}
						pgID = opParts[1]
					}
				}
			} else if strings.Contains(line, "pg_update_log_missing(") {
				// Extract client ID and op ID from pg_update_log_missing
				opStart := strings.Index(line, "by client.")
				if opStart != -1 {
					clientPart := line[opStart+3:] // Skip "by "
					clientEnd := strings.Index(clientPart, " ")
					if clientEnd != -1 {
						clientOp := clientPart[:clientEnd]
						clientParts := strings.Split(clientOp, ":")
						if len(clientParts) == 2 {
							clientID = clientParts[0]
							opID = clientParts[1]
						}
					}
				}
				// Extract PG ID from pg_update_log_missing
				pgStart := strings.Index(line, "pg_update_log_missing(")
				if pgStart != -1 {
					pgPart := line[pgStart:]
					pgEnd := strings.Index(pgPart, " epoch")
					if pgEnd != -1 {
						pgStr := pgPart[22:pgEnd] // Skip "pg_update_log_missing("
						pgID = pgStr
					}
				}
			}

			// Extract dequeue latency (in microseconds)
			dequeueLatency := 0
			for i, part := range parts {
				if part == "latency" && i+1 < len(parts) {
					latencyStr := parts[i+1]
					// Parse float value (in seconds)
					latencySec, err := strconv.ParseFloat(latencyStr, 64)
					if err == nil {
						// Convert to microseconds
						dequeueLatency = int(latencySec * 1000000)
					}
					break
				}
			}

			// Extract priority
			priority := 0
			for i, part := range parts {
				if part == "prio" && i+1 < len(parts) {
					prioStr := parts[i+1]
					prio, err := strconv.Atoi(prioStr)
					if err == nil {
						priority = prio
					}
					break
				}
			}

			// Extract cost
			cost := 0
			for i, part := range parts {
				if part == "cost" && i+1 < len(parts) {
					costStr := parts[i+1]
					costVal, err := strconv.Atoi(costStr)
					if err == nil {
						cost = costVal
					}
					break
				}
			}

			// Create dequeue event
			event := types.DequeueEvent{
				Timestamp:      timestamp,
				OpType:         opType,
				ClientID:       clientID,
				OpID:           opID,
				PGID:           pgID,
				DequeueLatency: dequeueLatency,
				Priority:       priority,
				Cost:           cost,
			}

			events = append(events, event)
		}
	}

	return events, nil
}
