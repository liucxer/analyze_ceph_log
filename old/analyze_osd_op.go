package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

// OSDOpEvent represents an OSD operation event
type OSDOpEvent struct {
	timestamp time.Time
	opID      string
	pgID      string
	object    string
	opType    string
	rangeStr  string
	inBytes   int
	outBytes  int
	latency   float64 // in milliseconds
}

func main() {
	// Check if a log file is provided
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run analyze_osd_op.go <log_file>")
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

	// Regular expression to match osd_op events
	osdOpRegex := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}) .* log_op_stats osd_op\(([^\)]+)\) .* inb (\d+) outb (\d+) lat ([\d\.]+)`)

	// Variables to track statistics
	totalOps := 0
	totalLatency := 0.0
	totalInBytes := 0
	totalOutBytes := 0
	maxLatency := 0.0

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	var events []OSDOpEvent

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
				fmt.Printf("Error parsing timestamp: %v\n", err)
				continue
			}

			// Parse input/output bytes
			inBytes, err := strconv.Atoi(inBytesStr)
			if err != nil {
				fmt.Printf("Error parsing inBytes: %v\n", err)
				continue
			}

			outBytes, err := strconv.Atoi(outBytesStr)
			if err != nil {
				fmt.Printf("Error parsing outBytes: %v\n", err)
				continue
			}

			// Parse latency
			latency, err := strconv.ParseFloat(latencyStr, 64)
			if err != nil {
				fmt.Printf("Error parsing latency: %v\n", err)
				continue
			}
			// Convert seconds to milliseconds
			latency = latency * 1000.0

			// Extract opID, pgID, object, opType, rangeStr from opInfo
			// Format: client.74318.0:1 1.2b 1:d4682f9e:::liuc999:head [writefull 0~169200] snapc 0=[] ondisk+write+known_if_redirected e84

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
			event := OSDOpEvent{
				timestamp: timestamp,
				opID:      opID,
				pgID:      pgID,
				object:    object,
				opType:    opType,
				rangeStr:  rangeStr,
				inBytes:   inBytes,
				outBytes:  outBytes,
				latency:   latency,
			}

			events = append(events, event)

			// Update statistics
			totalOps++
			totalLatency += latency
			totalInBytes += inBytes
			totalOutBytes += outBytes
			// Update max latency
			if latency > maxLatency {
				maxLatency = latency
			}
		}
	}

	// Print header
	fmt.Printf("%-25s %-20s %-10s %-30s %-12s %-15s %-10s %-10s\n",
		"Timestamp", "OP ID", "PG ID", "Object", "Op Type", "Range", "In (bytes)", "Out (bytes)")
	fmt.Println("--------------------------------------------------------------------------------------------------------------------------------------------")

	// Print events
	for _, event := range events {
		fmt.Printf("%-25s %-20s %-10s %-30s %-12s %-15s %-10d %-10d\n",
			event.timestamp.Format("2006-01-02 15:04:05.000"),
			event.opID,
			event.pgID,
			event.object,
			event.opType,
			event.rangeStr,
			event.inBytes,
			event.outBytes)
		// Print latency
		fmt.Printf("%-130s Latency: %.6f ms\n", "", event.latency)
		fmt.Println("--------------------------------------------------------------------------------------------------------------------------------------------")
	}

	// Print summary statistics
	fmt.Println("\nSummary:")
	fmt.Printf("Total operations: %d\n", totalOps)
	if totalOps > 0 {
		avgLatency := totalLatency / float64(totalOps)
		avgInBytes := totalInBytes / totalOps
		avgOutBytes := totalOutBytes / totalOps
		fmt.Printf("Average latency: %.6f ms\n", avgLatency)
		fmt.Printf("Maximum latency: %.6f ms\n", maxLatency)
		fmt.Printf("Average input: %d bytes\n", avgInBytes)
		fmt.Printf("Average output: %d bytes\n", avgOutBytes)
	}
	fmt.Printf("Total input: %d bytes\n", totalInBytes)
	fmt.Printf("Total output: %d bytes\n", totalOutBytes)

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}
}
