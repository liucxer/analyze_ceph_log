package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// AIOEvent represents an AIO start or finish event
type AIOEvent struct {
	timestamp time.Time
	blockAddr string
	rangeStr  string
	duration  time.Duration
}

// RepopEvent represents an OSD repop event
type RepopEvent struct {
	startTime time.Time
	endTime   time.Time
	duration  time.Duration
	opID      string
}

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

// Event represents a generic event with start/end times and duration
type Event struct {
	startTime time.Time
	endTime   time.Time
	duration  time.Duration
	rangeStr  string
	opID      string
}

func main() {
	// Check if a log file is provided
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run analyze_ceph.go <log_file> <analysis_type> [output.html]")
		fmt.Println("Analysis types:")
		fmt.Println("  aio - Analyze AIO operations")
		fmt.Println("  repop - Analyze OSD repop operations")
		fmt.Println("  op - Analyze OSD operations")
		fmt.Println("  all - Analyze all operation types")
		os.Exit(1)
	}

	logFile := os.Args[1]
	analysisType := os.Args[2]

	// Determine output file
	outputFile := "analysis.html"
	if len(os.Args) > 3 {
		outputFile = os.Args[3]
	}

	// Open the log file
	file, err := os.Open(logFile)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Analyze based on type
	switch analysisType {
	case "aio":
		analyzeAIO(file, outputFile)
	case "repop":
		analyzeRepop(file, outputFile)
	case "op":
		analyzeOSDOp(file, outputFile)
	case "all":
		analyzeAll(file, outputFile)
	default:
		fmt.Printf("Unknown analysis type: %s\n", analysisType)
		os.Exit(1)
	}

	fmt.Printf("Analysis completed. Results saved to %s\n", outputFile)
}

// analyzeAIO analyzes AIO operations
func analyzeAIO(file *os.File, outputFile string) {
	// Reset file pointer
	file.Seek(0, 0)

	// Regular expressions to match start and finish events
	startRegex := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}) .* _aio_log_start (0x[0-9a-fA-F]+)~([0-9a-fA-F]+)`)
	finishRegex := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}) .* _aio_log_finish 1 (0x[0-9a-fA-F]+)~([0-9a-fA-F]+)`)

	// Map to store all events by block address
	eventsMap := make(map[string]struct {
		startTime time.Time
		endTime   time.Time
		duration  time.Duration
		rangeStr  string
	})

	// Read the file line by line
	scanner := bufio.NewScanner(file)
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
	var events []Event
	for _, event := range eventsMap {
		if event.duration >= 0 {
			events = append(events, Event{
				startTime: event.startTime,
				endTime:   event.endTime,
				duration:  event.duration,
				rangeStr:  event.rangeStr,
			})
		}
	}

	// Sort events by start time
	sort.Slice(events, func(i, j int) bool {
		return events[i].startTime.Before(events[j].startTime)
	})

	// Calculate statistics
	totalAIOs := len(events)
	totalDuration := time.Duration(0)
	maxDuration := time.Duration(0)
	minDuration := time.Duration(1000000000) // Start with a large value

	for _, event := range events {
		totalDuration += event.duration
		if event.duration > maxDuration {
			maxDuration = event.duration
		}
		if event.duration < minDuration {
			minDuration = event.duration
		}
	}

	// Generate HTML
	generateHTML(outputFile, "AIO Operations Analysis", events, totalAIOs, totalDuration, maxDuration, minDuration)
}

// analyzeRepop analyzes OSD repop operations
func analyzeRepop(file *os.File, outputFile string) {
	// Reset file pointer
	file.Seek(0, 0)

	// Regular expressions to match dequeue_op and repop_commit events
	dequeueRegex := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}) .* dequeue_op .* osd_repop\(([\s\S]+?)\)`)
	commitRegex := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}) .* repop_commit on op osd_repop\(([\s\S]+?)\)`)

	// Map to store all events by op ID
	eventsMap := make(map[string]struct {
		startTime time.Time
		endTime   time.Time
		duration  time.Duration
		opID      string
	})

	// Read the file line by line
	scanner := bufio.NewScanner(file)
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
	var events []Event
	for _, event := range eventsMap {
		if event.duration > 0 {
			events = append(events, Event{
				startTime: event.startTime,
				endTime:   event.endTime,
				duration:  event.duration,
				opID:      event.opID,
			})
		}
	}

	// Sort events by start time
	sort.Slice(events, func(i, j int) bool {
		return events[i].startTime.Before(events[j].startTime)
	})

	// Calculate statistics
	totalRepops := len(events)
	totalDuration := time.Duration(0)
	maxDuration := time.Duration(0)
	minDuration := time.Duration(1000000000) // Start with a large value

	for _, event := range events {
		totalDuration += event.duration
		if event.duration > maxDuration {
			maxDuration = event.duration
		}
		if event.duration < minDuration {
			minDuration = event.duration
		}
	}

	// Generate HTML
	generateRepopHTML(outputFile, "OSD Repop Operations Analysis", events, totalRepops, totalDuration, maxDuration, minDuration)
}

// analyzeOSDOp analyzes OSD operations
func analyzeOSDOp(file *os.File, outputFile string) {
	// Reset file pointer
	file.Seek(0, 0)

	// Regular expression to match osd_op events
	osdOpRegex := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}) .* log_op_stats osd_op\(([\s\S]+?)\) .* inb (\d+) outb (\d+) lat ([\d\.]+)`)

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	var events []OSDOpEvent

	totalOps := 0
	totalLatency := 0.0
	totalInBytes := 0
	totalOutBytes := 0
	maxLatency := 0.0

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

	// Sort events by timestamp
	sort.Slice(events, func(i, j int) bool {
		return events[i].timestamp.Before(events[j].timestamp)
	})

	// Generate HTML
	generateOSDOpHTML(outputFile, "OSD Operations Analysis", events, totalOps, totalLatency, maxLatency, totalInBytes, totalOutBytes)
}

// analyzeAll analyzes all operation types
func analyzeAll(file *os.File, outputFile string) {
	// Generate HTML with tabbed interface
	htmlContent := `<!DOCTYPE html>
<html>
<head>
    <title>Ceph Log Analysis</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 20px;
            background-color: #f5f5f5;
        }
        h1 {
            color: #333;
            text-align: center;
            margin-bottom: 30px;
        }
        h2 {
            color: #555;
            margin-top: 0;
            margin-bottom: 20px;
        }
        h3 {
            color: #666;
            margin-top: 0;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
        }
        /* Tab styles */
        .tabs {
            display: flex;
            margin-bottom: 20px;
            background-color: white;
            border-radius: 8px 8px 0 0;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            overflow: hidden;
        }
        .tab {
            padding: 15px 20px;
            cursor: pointer;
            background-color: #f1f1f1;
            border: none;
            outline: none;
            transition: background-color 0.3s;
            font-size: 16px;
            font-weight: bold;
        }
        .tab:hover {
            background-color: #ddd;
        }
        .tab.active {
            background-color: #3498db;
            color: white;
        }
        /* Panel styles */
        .panel {
            background-color: white;
            padding: 20px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            border-radius: 0 8px 8px 8px;
            display: none;
        }
        .panel.active {
            display: block;
        }
        /* Summary styles */
        .summary {
            background-color: #f9f9f9;
            padding: 15px;
            border-radius: 6px;
            margin-bottom: 20px;
            border-left: 4px solid #3498db;
        }
        /* Filter form styles */
        .filter-form {
            margin-bottom: 20px;
            padding: 15px;
            background-color: #f0f0f0;
            border-radius: 6px;
            display: flex;
            flex-wrap: wrap;
            gap: 10px;
            align-items: center;
        }
        .filter-form label {
            font-weight: bold;
            margin-right: 5px;
        }
        .filter-form input {
            padding: 5px;
            border: 1px solid #ddd;
            border-radius: 4px;
        }
        .filter-form button {
            padding: 5px 10px;
            background-color: #3498db;
            color: white;
            border: none;
            border-radius: 4px;
            cursor: pointer;
        }
        .filter-form button:hover {
            background-color: #2980b9;
        }
        /* Table styles */
        table {
            width: 100%;
            border-collapse: collapse;
            margin: 20px 0;
            background-color: white;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
        }
        th, td {
            border: 1px solid #ddd;
            padding: 8px;
            text-align: left;
        }
        th {
            background-color: #f2f2f2;
            font-weight: bold;
        }
        tr:nth-child(even) {
            background-color: #f9f9f9;
        }
        .table-container {
            overflow-x: auto;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Ceph Log Analysis</h1>
        
        <!-- Tabs -->
        <div class="tabs">
            <button class="tab active" onclick="openTab(event, 'aio')">AIO Operations</button>
            <button class="tab" onclick="openTab(event, 'repop')">OSD Repop Operations</button>
            <button class="tab" onclick="openTab(event, 'osd')">OSD Operations</button>
        </div>
        
        <!-- AIO Panel -->
        <div id="aio" class="panel active">
`

	// Add AIO analysis with summary at top
	file.Seek(0, 0)
	htmlContent += analyzeAIOForAll(file)

	htmlContent += `
        </div>
        
        <!-- Repop Panel -->
        <div id="repop" class="panel">
`

	// Add Repop analysis with summary at top
	file.Seek(0, 0)
	htmlContent += analyzeRepopForAll(file)

	htmlContent += `
        </div>
        
        <!-- OSD Panel -->
        <div id="osd" class="panel">
`

	// Add OSD Op analysis with summary at top
	file.Seek(0, 0)
	htmlContent += analyzeOSDOpForAll(file)

	htmlContent += `
        </div>
        
        <script>
            function openTab(evt, tabName) {
                var i, tabcontent, tablinks;
                
                // Hide all tab content
                tabcontent = document.getElementsByClassName("panel");
                for (i = 0; i < tabcontent.length; i++) {
                    tabcontent[i].classList.remove("active");
                }
                
                // Remove active class from all tabs
                tablinks = document.getElementsByClassName("tab");
                for (i = 0; i < tablinks.length; i++) {
                    tablinks[i].classList.remove("active");
                }
                
                // Show the selected tab content and set active tab
                document.getElementById(tabName).classList.add("active");
                evt.currentTarget.classList.add("active");
            }
            
            // AIO Table Filter
            function filterAIOTable() {
                var startTime = document.getElementById("aio-start-time").value;
                var endTime = document.getElementById("aio-end-time").value;
                var minDuration = document.getElementById("aio-min-duration").value;
                var maxDuration = document.getElementById("aio-max-duration").value;
                var table = document.getElementById("aio-table");
                var tr = table.getElementsByTagName("tr");
                
                for (var i = 1; i < tr.length; i++) {
                    var tdStartTime = tr[i].getElementsByTagName("td")[0].textContent;
                    var tdEndTime = tr[i].getElementsByTagName("td")[1].textContent;
                    var tdDuration = parseFloat(tr[i].getElementsByTagName("td")[2].textContent);
                    
                    var match = true;
                    
                    if (startTime) {
                        var startDate = new Date(startTime.replace('T', ' '));
                        var rowStartDate = new Date(tdStartTime);
                        if (rowStartDate < startDate) match = false;
                    }
                    
                    if (endTime) {
                        var endDate = new Date(endTime.replace('T', ' '));
                        var rowEndDate = new Date(tdEndTime);
                        if (rowEndDate > endDate) match = false;
                    }
                    
                    if (minDuration) {
                        if (tdDuration < parseFloat(minDuration)) match = false;
                    }
                    
                    if (maxDuration) {
                        if (tdDuration > parseFloat(maxDuration)) match = false;
                    }
                    
                    tr[i].style.display = match ? "" : "none";
                }
            }
            
            function resetAIOFilter() {
                document.getElementById("aio-start-time").value = "";
                document.getElementById("aio-end-time").value = "";
                document.getElementById("aio-min-duration").value = "";
                document.getElementById("aio-max-duration").value = "";
                filterAIOTable();
            }
            
            // Repop Table Filter
            function filterRepopTable() {
                var startTime = document.getElementById("repop-start-time").value;
                var endTime = document.getElementById("repop-end-time").value;
                var minDuration = document.getElementById("repop-min-duration").value;
                var maxDuration = document.getElementById("repop-max-duration").value;
                var table = document.getElementById("repop-table");
                var tr = table.getElementsByTagName("tr");
                
                for (var i = 1; i < tr.length; i++) {
                    var tdStartTime = tr[i].getElementsByTagName("td")[0].textContent;
                    var tdEndTime = tr[i].getElementsByTagName("td")[1].textContent;
                    var tdDuration = parseFloat(tr[i].getElementsByTagName("td")[2].textContent);
                    
                    var match = true;
                    
                    if (startTime) {
                        var startDate = new Date(startTime.replace('T', ' '));
                        var rowStartDate = new Date(tdStartTime);
                        if (rowStartDate < startDate) match = false;
                    }
                    
                    if (endTime) {
                        var endDate = new Date(endTime.replace('T', ' '));
                        var rowEndDate = new Date(tdEndTime);
                        if (rowEndDate > endDate) match = false;
                    }
                    
                    if (minDuration) {
                        if (tdDuration < parseFloat(minDuration)) match = false;
                    }
                    
                    if (maxDuration) {
                        if (tdDuration > parseFloat(maxDuration)) match = false;
                    }
                    
                    tr[i].style.display = match ? "" : "none";
                }
            }
            
            function resetRepopFilter() {
                document.getElementById("repop-start-time").value = "";
                document.getElementById("repop-end-time").value = "";
                document.getElementById("repop-min-duration").value = "";
                document.getElementById("repop-max-duration").value = "";
                filterRepopTable();
            }
            
            // OSD Table Filter
            function filterOSDTable() {
                var startTime = document.getElementById("osd-start-time").value;
                var endTime = document.getElementById("osd-end-time").value;
                var minLatency = document.getElementById("osd-min-latency").value;
                var maxLatency = document.getElementById("osd-max-latency").value;
                var table = document.getElementById("osd-table");
                var tr = table.getElementsByTagName("tr");
                
                for (var i = 1; i < tr.length; i++) {
                    var tdTime = tr[i].getElementsByTagName("td")[0].textContent;
                    var tdLatency = parseFloat(tr[i].getElementsByTagName("td")[8].textContent);
                    
                    var match = true;
                    
                    if (startTime) {
                        var startDate = new Date(startTime.replace('T', ' '));
                        var rowDate = new Date(tdTime);
                        if (rowDate < startDate) match = false;
                    }
                    
                    if (endTime) {
                        var endDate = new Date(endTime.replace('T', ' '));
                        var rowDate = new Date(tdTime);
                        if (rowDate > endDate) match = false;
                    }
                    
                    if (minLatency) {
                        if (tdLatency < parseFloat(minLatency)) match = false;
                    }
                    
                    if (maxLatency) {
                        if (tdLatency > parseFloat(maxLatency)) match = false;
                    }
                    
                    tr[i].style.display = match ? "" : "none";
                }
            }
            
            function resetOSDFilter() {
                document.getElementById("osd-start-time").value = "";
                document.getElementById("osd-end-time").value = "";
                document.getElementById("osd-min-latency").value = "";
                document.getElementById("osd-max-latency").value = "";
                filterOSDTable();
            }
        </script>
    </div>
</body>
</html>`

	// Write to file
	err := os.WriteFile(outputFile, []byte(htmlContent), 0644)
	if err != nil {
		fmt.Printf("Error writing HTML file: %v\n", err)
		os.Exit(1)
	}
}

// calculateOverallStats calculates total counts of each operation type
func calculateOverallStats(file *os.File) (int, int, int) {
	file.Seek(0, 0)
	
	// Regular expressions for each operation type
	aioStartRegex := regexp.MustCompile(`_aio_log_start`)
	repopRegex := regexp.MustCompile(`osd_repop`)
	osdOpRegex := regexp.MustCompile(`log_op_stats osd_op`)
	
	// Counters
	aioCount := 0
	repopCount := 0
	osdOpCount := 0
	
	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		
		if aioStartRegex.MatchString(line) {
			aioCount++
		} else if repopRegex.MatchString(line) {
			repopCount++
		} else if osdOpRegex.MatchString(line) {
			osdOpCount++
		}
	}
	
	return aioCount, repopCount, osdOpCount
}

// analyzeAIOForAll analyzes AIO operations for combined report
func analyzeAIOForAll(file *os.File) string {
	// Same as analyzeAIO but returns HTML string
	// Regular expressions to match start and finish events
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

	// Read the file line by line
	scanner := bufio.NewScanner(file)
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
			key := fmt.Sprintf("%s_%s_%d", blockAddr, timestampStr, counter)
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
			// We need to iterate through the map to find the matching start event
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
	var events []Event
	for _, event := range eventsMap {
		if !event.endTime.IsZero() && event.duration >= 0 {
			events = append(events, Event{
				startTime: event.startTime,
				endTime:   event.endTime,
				duration:  event.duration,
				rangeStr:  event.rangeStr,
			})
		}
	}

	// Sort events by start time
	sort.Slice(events, func(i, j int) bool {
		return events[i].startTime.Before(events[j].startTime)
	})

	// Calculate statistics
	totalAIOs := len(events)
	totalDuration := time.Duration(0)
	maxDuration := time.Duration(0)
	minDuration := time.Duration(1000000000) // Start with a large value
	
	// Count requests by all durations present in data
	durationCounts := make(map[int]int)

	for _, event := range events {
		totalDuration += event.duration
		if event.duration > maxDuration {
			maxDuration = event.duration
		}
		if event.duration < minDuration {
			minDuration = event.duration
		}
		
		// Categorize by duration
		durationMs := int(float64(event.duration.Microseconds()) / 1000.0)
		durationCounts[durationMs]++
	}

	// Generate HTML with summary at top
	html := `
    <h2>AIO Operations Analysis</h2>
    <div class="summary">
        <h3>Summary</h3>
        <p>Total AIO operations: ` + strconv.Itoa(totalAIOs) + `</p>`

	if totalAIOs > 0 {
		averageDuration := totalDuration / time.Duration(totalAIOs)
		html += fmt.Sprintf(`
        <p>Average duration: %.3f ms</p>
        <p>Maximum duration: %.3f ms</p>
        <p>Minimum duration: %.3f ms</p>`,
			float64(averageDuration.Microseconds())/1000.0,
			float64(maxDuration.Microseconds())/1000.0,
			float64(minDuration.Microseconds())/1000.0)
	}

	// Add duration counts
	html += `
        <h4>Duration Counts:</h4>`
	
	// Sort durations for consistent output
	var durations []int
	for duration := range durationCounts {
		durations = append(durations, duration)
	}
	sort.Ints(durations)
	
	for _, duration := range durations {
		html += fmt.Sprintf(`
        <p>%dms: %d requests</p>`, duration, durationCounts[duration])
	}

	// Add filter form
	html += `
        <h4>Filter Options:</h4>
        <div class="filter-form">
            <label>Start Time:</label>
            <input type="datetime-local" id="aio-start-time">
            <label>End Time:</label>
            <input type="datetime-local" id="aio-end-time">
            <label>Min Duration (ms):</label>
            <input type="number" id="aio-min-duration" min="0">
            <label>Max Duration (ms):</label>
            <input type="number" id="aio-max-duration" min="0">
            <button type="button" onclick="filterAIOTable()">Filter</button>
            <button type="button" onclick="resetAIOFilter()">Reset</button>
        </div>
    </div>
    <div class="table-container">
    <table id="aio-table">
        <tr>
            <th>Start Time</th>
            <th>End Time</th>
            <th>Duration (ms)</th>
            <th>Range</th>
        </tr>`

	for _, event := range events {
		html += fmt.Sprintf(`
        <tr>
            <td>%s</td>
            <td>%s</td>
            <td>%.3f</td>
            <td>%s</td>
        </tr>`,
			event.startTime.Format("2006-01-02 15:04:05.000"),
			event.endTime.Format("2006-01-02 15:04:05.000"),
			float64(event.duration.Microseconds())/1000.0,
			event.rangeStr)
	}

	html += `
    </table>
    </div>`

	return html
}

// analyzeRepopForAll analyzes Repop operations for combined report
func analyzeRepopForAll(file *os.File) string {
	// Same as analyzeRepop but returns HTML string
	// Regular expressions to match dequeue_op and repop_commit events
	dequeueRegex := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}) .* dequeue_op .* osd_repop\(([\s\S]+?)\)`)
	commitRegex := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}) .* repop_commit on op osd_repop\(([\s\S]+?)\)`)

	// Map to store all events by op ID
	eventsMap := make(map[string]struct {
		startTime time.Time
		endTime   time.Time
		duration  time.Duration
		opID      string
	})

	// Read the file line by line
	scanner := bufio.NewScanner(file)
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
	var events []Event
	for _, event := range eventsMap {
		if event.duration >= 0 {
			events = append(events, Event{
				startTime: event.startTime,
				endTime:   event.endTime,
				duration:  event.duration,
				opID:      event.opID,
			})
		}
	}

	// Sort events by start time
	sort.Slice(events, func(i, j int) bool {
		return events[i].startTime.Before(events[j].startTime)
	})

	// Calculate statistics
	totalRepops := len(events)
	totalDuration := time.Duration(0)
	maxDuration := time.Duration(0)
	minDuration := time.Duration(1000000000) // Start with a large value
	
	// Count requests by all durations present in data
	durationCounts := make(map[int]int)

	for _, event := range events {
		totalDuration += event.duration
		if event.duration > maxDuration {
			maxDuration = event.duration
		}
		if event.duration < minDuration {
			minDuration = event.duration
		}
		
		// Categorize by duration
		durationMs := int(float64(event.duration.Microseconds()) / 1000.0)
		durationCounts[durationMs]++
	}

	// Generate HTML with summary at top
	html := `
    <h2>OSD Repop Operations Analysis</h2>
    <div class="summary">
        <h3>Summary</h3>
        <p>Total repop operations: ` + strconv.Itoa(totalRepops) + `</p>`

	if totalRepops > 0 {
		averageDuration := totalDuration / time.Duration(totalRepops)
		html += fmt.Sprintf(`
        <p>Average duration: %.3f ms</p>
        <p>Maximum duration: %.3f ms</p>
        <p>Minimum duration: %.3f ms</p>`,
			float64(averageDuration.Microseconds())/1000.0,
			float64(maxDuration.Microseconds())/1000.0,
			float64(minDuration.Microseconds())/1000.0)
	}

	// Add duration counts
	html += `
        <h4>Duration Counts:</h4>`
	
	// Sort durations for consistent output
	var durations []int
	for duration := range durationCounts {
		durations = append(durations, duration)
	}
	sort.Ints(durations)
	
	for _, duration := range durations {
		html += fmt.Sprintf(`
        <p>%dms: %d requests</p>`, duration, durationCounts[duration])
	}

	// Add filter form
	html += `
        <h4>Filter Options:</h4>
        <div class="filter-form">
            <label>Start Time:</label>
            <input type="datetime-local" id="repop-start-time">
            <label>End Time:</label>
            <input type="datetime-local" id="repop-end-time">
            <label>Min Duration (ms):</label>
            <input type="number" id="repop-min-duration" min="0">
            <label>Max Duration (ms):</label>
            <input type="number" id="repop-max-duration" min="0">
            <button type="button" onclick="filterRepopTable()">Filter</button>
            <button type="button" onclick="resetRepopFilter()">Reset</button>
        </div>
    </div>
    <div class="table-container">
    <table id="repop-table">
        <tr>
            <th>Start Time</th>
            <th>End Time</th>
            <th>Duration (ms)</th>
            <th>OP ID</th>
        </tr>`

	for _, event := range events {
		html += fmt.Sprintf(`
        <tr>
            <td>%s</td>
            <td>%s</td>
            <td>%.3f</td>
            <td>%s</td>
        </tr>`,
			event.startTime.Format("2006-01-02 15:04:05.000"),
			event.endTime.Format("2006-01-02 15:04:05.000"),
			float64(event.duration.Microseconds())/1000.0,
			event.opID)
	}

	html += `
    </table>
    </div>`

	return html
}

// analyzeOSDOpForAll analyzes OSD operations for combined report
func analyzeOSDOpForAll(file *os.File) string {
	// Same as analyzeOSDOp but returns HTML string
	// Regular expression to match osd_op events
	osdOpRegex := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}) .* log_op_stats osd_op\(([\s\S]+?)\) .* inb (\d+) outb (\d+) lat ([\d\.]+)`)

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	var events []OSDOpEvent

	totalOps := 0
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
			
			// Categorize by latency
			switch {
			case latency < 1:
				latencyCounts["0-1ms"]++
			case latency < 2:
				latencyCounts["1-2ms"]++
			case latency < 3:
				latencyCounts["2-3ms"]++
			case latency < 4:
				latencyCounts["3-4ms"]++
			case latency < 5:
				latencyCounts["4-5ms"]++
			case latency < 10:
				latencyCounts["5-10ms"]++
			default:
				latencyCounts["10ms+"]++
			}
		}
	}

	// Sort events by timestamp
	sort.Slice(events, func(i, j int) bool {
		return events[i].timestamp.Before(events[j].timestamp)
	})

	// Generate HTML with summary at top
	html := `
    <h2>OSD Operations Analysis</h2>
    <div class="summary">
        <h3>Summary</h3>
        <p>Total operations: ` + strconv.Itoa(totalOps) + `</p>`

	if totalOps > 0 {
		avgLatency := totalLatency / float64(totalOps)
		avgInBytes := totalInBytes / totalOps
		avgOutBytes := totalOutBytes / totalOps
		html += fmt.Sprintf(`
        <p>Average latency: %.6f ms</p>
        <p>Maximum latency: %.6f ms</p>
        <p>Average input: %d bytes</p>
        <p>Average output: %d bytes</p>`,
			avgLatency,
			maxLatency,
			avgInBytes,
			avgOutBytes)
	}

	// Add latency distribution
	html += `
        <h4>Latency Distribution:</h4>`
	for _, rangeStr := range latencyRanges {
		html += fmt.Sprintf(`
        <p>%s: %d requests</p>`, rangeStr, latencyCounts[rangeStr])
	}

	// Add filter form
	html += `
        <h4>Filter Options:</h4>
        <div class="filter-form">
            <label>Start Time:</label>
            <input type="datetime-local" id="osd-start-time">
            <label>End Time:</label>
            <input type="datetime-local" id="osd-end-time">
            <label>Min Latency (ms):</label>
            <input type="number" id="osd-min-latency" min="0">
            <label>Max Latency (ms):</label>
            <input type="number" id="osd-max-latency" min="0">
            <button type="button" onclick="filterOSDTable()">Filter</button>
            <button type="button" onclick="resetOSDFilter()">Reset</button>
        </div>
    </div>
    <div class="table-container">
    <table id="osd-table">
        <tr>
            <th>Timestamp</th>
            <th>OP ID</th>
            <th>PG ID</th>
            <th>Object</th>
            <th>Op Type</th>
            <th>Range</th>
            <th>In (bytes)</th>
            <th>Out (bytes)</th>
            <th>Latency (ms)</th>
        </tr>`

	for _, event := range events {
		html += fmt.Sprintf(`
        <tr>
            <td>%s</td>
            <td>%s</td>
            <td>%s</td>
            <td>%s</td>
            <td>%s</td>
            <td>%s</td>
            <td>%d</td>
            <td>%d</td>
            <td>%.6f</td>
        </tr>`,
			event.timestamp.Format("2006-01-02 15:04:05.000"),
			event.opID,
			event.pgID,
			event.object,
			event.opType,
			event.rangeStr,
			event.inBytes,
			event.outBytes,
			event.latency)
	}

	html += `
    </table>
    </div>`

	return html
}

// generateHTML generates HTML for AIO analysis
func generateHTML(outputFile, title string, events []Event, totalAIOs int, totalDuration, maxDuration, minDuration time.Duration) {
	htmlContent := `<!DOCTYPE html>
<html>
<head>
    <title>` + title + `</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 20px;
            background-color: #f5f5f5;
        }
        h1 {
            color: #333;
            text-align: center;
        }
        h2 {
            color: #555;
            border-bottom: 2px solid #ddd;
            padding-bottom: 5px;
            margin-top: 0;
        }
        h3 {
            color: #666;
            margin-top: 0;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
        }
        .summary {
            background-color: #f9f9f9;
            padding: 15px;
            margin-bottom: 20px;
            border-radius: 6px;
            border-left: 4px solid #3498db;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            margin: 20px 0;
            background-color: white;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
        }
        th, td {
            border: 1px solid #ddd;
            padding: 8px;
            text-align: left;
        }
        th {
            background-color: #f2f2f2;
            font-weight: bold;
        }
        tr:nth-child(even) {
            background-color: #f9f9f9;
        }
        .table-container {
            overflow-x: auto;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>` + title + `</h1>
        <div class="summary">
            <h2>Summary</h2>
            <p>Total AIO operations: ` + strconv.Itoa(totalAIOs) + `</p>`

	if totalAIOs > 0 {
		averageDuration := totalDuration / time.Duration(totalAIOs)
		htmlContent += fmt.Sprintf(`
            <p>Average duration: %.3f ms</p>
            <p>Maximum duration: %.3f ms</p>
            <p>Minimum duration: %.3f ms</p>`,
			float64(averageDuration.Microseconds())/1000.0,
			float64(maxDuration.Microseconds())/1000.0,
			float64(minDuration.Microseconds())/1000.0)
	}

	htmlContent += `
        </div>
        <div class="table-container">
        <table>
            <tr>
                <th>Start Time</th>
                <th>End Time</th>
                <th>Duration (ms)</th>
                <th>Range</th>
            </tr>`

	for _, event := range events {
		htmlContent += fmt.Sprintf(`
            <tr>
                <td>%s</td>
                <td>%s</td>
                <td>%.3f</td>
                <td>%s</td>
            </tr>`,
			event.startTime.Format("2006-01-02 15:04:05.000"),
			event.endTime.Format("2006-01-02 15:04:05.000"),
			float64(event.duration.Microseconds())/1000.0,
			event.rangeStr)
	}

	htmlContent += `
        </table>
        </div>
    </div>
</body>
</html>`

	// Write to file
	err := os.WriteFile(outputFile, []byte(htmlContent), 0644)
	if err != nil {
		fmt.Printf("Error writing HTML file: %v\n", err)
		os.Exit(1)
	}
}

// generateRepopHTML generates HTML for Repop analysis
func generateRepopHTML(outputFile, title string, events []Event, totalRepops int, totalDuration, maxDuration, minDuration time.Duration) {
	htmlContent := `<!DOCTYPE html>
<html>
<head>
    <title>` + title + `</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 20px;
            background-color: #f5f5f5;
        }
        h1 {
            color: #333;
            text-align: center;
        }
        h2 {
            color: #555;
            border-bottom: 2px solid #ddd;
            padding-bottom: 5px;
            margin-top: 0;
        }
        h3 {
            color: #666;
            margin-top: 0;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
        }
        .summary {
            background-color: #f9f9f9;
            padding: 15px;
            margin-bottom: 20px;
            border-radius: 6px;
            border-left: 4px solid #3498db;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            margin: 20px 0;
            background-color: white;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
        }
        th, td {
            border: 1px solid #ddd;
            padding: 8px;
            text-align: left;
        }
        th {
            background-color: #f2f2f2;
            font-weight: bold;
        }
        tr:nth-child(even) {
            background-color: #f9f9f9;
        }
        .table-container {
            overflow-x: auto;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>` + title + `</h1>
        <div class="summary">
            <h2>Summary</h2>
            <p>Total repop operations: ` + strconv.Itoa(totalRepops) + `</p>`

	if totalRepops > 0 {
		averageDuration := totalDuration / time.Duration(totalRepops)
		htmlContent += fmt.Sprintf(`
            <p>Average duration: %.3f ms</p>
            <p>Maximum duration: %.3f ms</p>
            <p>Minimum duration: %.3f ms</p>`,
			float64(averageDuration.Microseconds())/1000.0,
			float64(maxDuration.Microseconds())/1000.0,
			float64(minDuration.Microseconds())/1000.0)
	}

	htmlContent += `
        </div>
        <div class="table-container">
        <table>
            <tr>
                <th>Start Time</th>
                <th>End Time</th>
                <th>Duration (ms)</th>
                <th>OP ID</th>
            </tr>`

	for _, event := range events {
		htmlContent += fmt.Sprintf(`
            <tr>
                <td>%s</td>
                <td>%s</td>
                <td>%.3f</td>
                <td>%s</td>
            </tr>`,
			event.startTime.Format("2006-01-02 15:04:05.000"),
			event.endTime.Format("2006-01-02 15:04:05.000"),
			float64(event.duration.Microseconds())/1000.0,
			event.opID)
	}

	htmlContent += `
        </table>
        </div>
    </div>
</body>
</html>`

	// Write to file
	err := os.WriteFile(outputFile, []byte(htmlContent), 0644)
	if err != nil {
		fmt.Printf("Error writing HTML file: %v\n", err)
		os.Exit(1)
	}
}

// generateOSDOpHTML generates HTML for OSD Op analysis
func generateOSDOpHTML(outputFile, title string, events []OSDOpEvent, totalOps int, totalLatency, maxLatency float64, totalInBytes, totalOutBytes int) {
	htmlContent := `<!DOCTYPE html>
<html>
<head>
    <title>` + title + `</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 20px;
            background-color: #f5f5f5;
        }
        h1 {
            color: #333;
            text-align: center;
        }
        h2 {
            color: #555;
            border-bottom: 2px solid #ddd;
            padding-bottom: 5px;
            margin-top: 0;
        }
        h3 {
            color: #666;
            margin-top: 0;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
        }
        .summary {
            background-color: #f9f9f9;
            padding: 15px;
            margin-bottom: 20px;
            border-radius: 6px;
            border-left: 4px solid #3498db;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            margin: 20px 0;
            background-color: white;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
        }
        th, td {
            border: 1px solid #ddd;
            padding: 8px;
            text-align: left;
        }
        th {
            background-color: #f2f2f2;
            font-weight: bold;
        }
        tr:nth-child(even) {
            background-color: #f9f9f9;
        }
        .table-container {
            overflow-x: auto;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>` + title + `</h1>
        <div class="summary">
            <h2>Summary</h2>
            <p>Total operations: ` + strconv.Itoa(totalOps) + `</p>`

	if totalOps > 0 {
		avgLatency := totalLatency / float64(totalOps)
		avgInBytes := totalInBytes / totalOps
		avgOutBytes := totalOutBytes / totalOps
		htmlContent += fmt.Sprintf(`
            <p>Average latency: %.6f ms</p>
            <p>Maximum latency: %.6f ms</p>
            <p>Average input: %d bytes</p>
            <p>Average output: %d bytes</p>`,
			avgLatency,
			maxLatency,
			avgInBytes,
			avgOutBytes)
	}

	htmlContent += fmt.Sprintf(`
            <p>Total input: %d bytes</p>
            <p>Total output: %d bytes</p>
        </div>
        <div class="table-container">
        <table>
            <tr>
                <th>Timestamp</th>
                <th>OP ID</th>
                <th>PG ID</th>
                <th>Object</th>
                <th>Op Type</th>
                <th>Range</th>
                <th>In (bytes)</th>
                <th>Out (bytes)</th>
                <th>Latency (ms)</th>
            </tr>`,
		totalInBytes,
		totalOutBytes)

	for _, event := range events {
		htmlContent += fmt.Sprintf(`
            <tr>
                <td>%s</td>
                <td>%s</td>
                <td>%s</td>
                <td>%s</td>
                <td>%s</td>
                <td>%s</td>
                <td>%d</td>
                <td>%d</td>
                <td>%.6f</td>
            </tr>`,
			event.timestamp.Format("2006-01-02 15:04:05.000"),
			event.opID,
			event.pgID,
			event.object,
			event.opType,
			event.rangeStr,
			event.inBytes,
			event.outBytes,
			event.latency)
	}

	htmlContent += `
        </table>
        </div>
    </div>
</body>
</html>`

	// Write to file
	err := os.WriteFile(outputFile, []byte(htmlContent), 0644)
	if err != nil {
		fmt.Printf("Error writing HTML file: %v\n", err)
		os.Exit(1)
	}
}