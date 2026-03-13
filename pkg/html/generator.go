package html

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/liucxer/analyze_ceph_log/pkg/types"
)

// GenerateHTML generates HTML for all analysis types
func GenerateHTML(aioResult types.AnalysisResult, repopResult types.AnalysisResult, osdOpResult types.OSDOpAnalysisResult, transactionResult types.TransactionAnalysisResult) string {
	htmlContent := `<!DOCTYPE html>
<html>
<head>
    <title>Ceph Log Analysis</title>
    <style>
        * {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
        }
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            margin: 0;
            padding: 20px;
            background-color: #f5f7fa;
            color: #333;
            line-height: 1.6;
        }
        h1 {
            color: #2c3e50;
            text-align: center;
            margin-bottom: 30px;
            font-size: 28px;
            font-weight: 600;
        }
        h2 {
            color: #34495e;
            margin-top: 0;
            margin-bottom: 20px;
            font-size: 22px;
            font-weight: 500;
            border-bottom: 2px solid #3498db;
            padding-bottom: 10px;
        }
        h3 {
            color: #555;
            margin-top: 20px;
            margin-bottom: 15px;
            font-size: 18px;
            font-weight: 500;
        }
        h4 {
            color: #666;
            margin-top: 15px;
            margin-bottom: 10px;
            font-size: 16px;
            font-weight: 500;
        }
        .container {
            max-width: 1400px;
            margin: 0 auto;
            background-color: white;
            border-radius: 10px;
            box-shadow: 0 0 20px rgba(0,0,0,0.1);
            overflow: hidden;
        }
        /* Tab styles */
        .tabs {
            display: flex;
            background-color: #f8f9fa;
            border-bottom: 1px solid #e9ecef;
            overflow-x: auto;
        }
        .tab {
            padding: 15px 25px;
            cursor: pointer;
            background-color: transparent;
            border: none;
            outline: none;
            transition: all 0.3s ease;
            font-size: 16px;
            font-weight: 500;
            color: #666;
            white-space: nowrap;
        }
        .tab:hover {
            background-color: #e9ecef;
            color: #3498db;
        }
        .tab.active {
            background-color: white;
            color: #3498db;
            border-bottom: 3px solid #3498db;
        }
        /* Panel styles */
        .panel {
            padding: 30px;
            display: none;
        }
        .panel.active {
            display: block;
        }
        /* Layout styles */
        .layout {
            display: flex;
            gap: 30px;
            margin-top: 20px;
        }
        .left-panel {
            flex: 0 0 350px;
        }
        .right-panel {
            flex: 1;
        }
        /* Summary styles */
        .summary {
            background-color: #f8f9fa;
            padding: 20px;
            border-radius: 8px;
            margin-bottom: 25px;
            border-left: 4px solid #3498db;
            box-shadow: 0 2px 4px rgba(0,0,0,0.05);
        }
        .summary p {
            margin-bottom: 8px;
        }
        /* Filter form styles */
        .filter-form {
            margin-bottom: 25px;
            padding: 20px;
            background-color: #f8f9fa;
            border-radius: 8px;
            display: flex;
            flex-wrap: wrap;
            gap: 15px;
            align-items: center;
            box-shadow: 0 2px 4px rgba(0,0,0,0.05);
        }
        .filter-form label {
            font-weight: 500;
            margin-right: 8px;
            color: #555;
        }
        .filter-form input,
        .filter-form select {
            padding: 8px 12px;
            border: 1px solid #ddd;
            border-radius: 4px;
            font-size: 14px;
            transition: border-color 0.3s ease;
        }
        .filter-form input:focus,
        .filter-form select:focus {
            outline: none;
            border-color: #3498db;
            box-shadow: 0 0 0 2px rgba(52, 152, 219, 0.2);
        }
        .filter-form button {
            padding: 8px 16px;
            background-color: #3498db;
            color: white;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            font-size: 14px;
            font-weight: 500;
            transition: background-color 0.3s ease;
        }
        .filter-form button:hover {
            background-color: #2980b9;
        }
        /* Table styles */
        table {
            width: 100%;
            border-collapse: collapse;
            margin: 25px 0;
            background-color: white;
            box-shadow: 0 2px 4px rgba(0,0,0,0.05);
            border-radius: 8px;
            overflow: hidden;
        }
        th, td {
            border: 1px solid #e9ecef;
            padding: 12px 15px;
            text-align: left;
        }
        th {
            background-color: #f8f9fa;
            font-weight: 600;
            color: #555;
            text-transform: uppercase;
            font-size: 14px;
            letter-spacing: 0.5px;
        }
        tr:nth-child(even) {
            background-color: #f8f9fa;
        }
        tr:hover {
            background-color: #f1f3f5;
        }
        .table-container {
            overflow-x: auto;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.05);
        }
        /* Query principle section */
        .query-principle {
            background-color: #f8f9fa;
            padding: 20px;
            border-radius: 8px;
            margin-bottom: 25px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.05);
        }
        .query-principle h3 {
            color: #34495e;
            margin-bottom: 15px;
        }
        .query-principle p {
            margin-bottom: 10px;
            color: #666;
        }
        .query-principle ul {
            margin-left: 20px;
            margin-bottom: 15px;
        }
        .query-principle li {
            margin-bottom: 5px;
            color: #666;
        }
        /* Responsive design */
        @media (max-width: 1200px) {
            .layout {
                flex-direction: column;
            }
            .left-panel {
                flex: 1;
            }
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
            <button class="tab" onclick="openTab(event, 'transaction')">Transaction Analysis</button>
        </div>
        
        <!-- AIO Panel -->
        <div id="aio" class="panel active">
`

	// Add AIO analysis
	htmlContent += generateAIOHTML(aioResult)

	htmlContent += `
        </div>
        
        <!-- Repop Panel -->
        <div id="repop" class="panel">
`

	// Add Repop analysis
	htmlContent += generateRepopHTML(repopResult)

	htmlContent += `
        </div>
        
        <!-- OSD Panel -->
        <div id="osd" class="panel">
`

	// Add OSD Op analysis
	htmlContent += generateOSDOpHTML(osdOpResult)

	htmlContent += `
        </div>
        
        <!-- Transaction Panel -->
        <div id="transaction" class="panel">
`

	// Add Transaction analysis
	htmlContent += generateTransactionHTML(transactionResult)

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
                var blockType = document.getElementById("aio-block-type").value;
                var minLength = document.getElementById("aio-min-length").value;
                var maxLength = document.getElementById("aio-max-length").value;
                var table = document.getElementById("aio-table");
                var tr = table.getElementsByTagName("tr");
                
                for (var i = 1; i < tr.length; i++) {
                    var tdStartTime = tr[i].getElementsByTagName("td")[0].textContent;
                    var tdEndTime = tr[i].getElementsByTagName("td")[1].textContent;
                    var tdDuration = parseFloat(tr[i].getElementsByTagName("td")[2].textContent);
                    var tdLength = parseInt(tr[i].getElementsByTagName("td")[4].textContent);
                    var tdBlockType = tr[i].getElementsByTagName("td")[5].textContent;
                    
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
                    
                    if (blockType) {
                        if (tdBlockType !== blockType) match = false;
                    }
                    
                    if (minLength) {
                        if (tdLength < parseInt(minLength)) match = false;
                    }
                    
                    if (maxLength) {
                        if (tdLength > parseInt(maxLength)) match = false;
                    }
                    
                    tr[i].style.display = match ? "" : "none";
                }
            }
            
            function resetAIOFilter() {
                document.getElementById("aio-start-time").value = "";
                document.getElementById("aio-end-time").value = "";
                document.getElementById("aio-min-duration").value = "";
                document.getElementById("aio-max-duration").value = "";
                document.getElementById("aio-block-type").value = "";
                document.getElementById("aio-min-length").value = "";
                document.getElementById("aio-max-length").value = "";
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
            
            // Transaction Table Filter
            function filterTransactionTable() {
                var startTime = document.getElementById("transaction-start-time").value;
                var endTime = document.getElementById("transaction-end-time").value;
                var minDuration = document.getElementById("transaction-min-duration").value;
                var maxDuration = document.getElementById("transaction-max-duration").value;
                var table = document.getElementById("transaction-table");
                var tr = table.getElementsByTagName("tr");
                
                for (var i = 1; i < tr.length; i++) {
                    var tdStartTime = tr[i].getElementsByTagName("td")[1].textContent;
                    var tdEndTime = tr[i].getElementsByTagName("td")[5].textContent;
                    var tdDuration = parseFloat(tr[i].getElementsByTagName("td")[6].textContent);
                    
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
            
            function resetTransactionFilter() {
                document.getElementById("transaction-start-time").value = "";
                document.getElementById("transaction-end-time").value = "";
                document.getElementById("transaction-min-duration").value = "";
                document.getElementById("transaction-max-duration").value = "";
                filterTransactionTable();
            }
        </script>
    </div>
</body>
</html>`

	return htmlContent
}

// generateAIOHTML generates HTML for AIO analysis
func generateAIOHTML(result types.AnalysisResult) string {
	html := `
    <h2>AIO Operations Analysis</h2>
    <div class="layout">
        <div class="left-panel">
            <div class="query-principle">
                <h3>Query Principle</h3>
                <p>This analysis parses AIO (Asynchronous I/O) operations from the Ceph log. It identifies:</p>
                <ul>
                    <li><strong>Start events</strong>: Log lines containing "_aio_log_start"</li>
                    <li><strong>Finish events</strong>: Log lines containing "_aio_log_finish"</li>
                </ul>
                <p>For each AIO operation, it extracts:</p>
                <ul>
                    <li>Timestamp</li>
                    <li>Block address range</li>
                    <li>Data length (converted from hex to decimal)</li>
                    <li>Block type (block, block.wal, block.db)</li>
                </ul>
                <p>It matches start and finish events using the block address range as a unique key, then calculates the duration for each operation.</p>
            </div>
            <div class="summary">
                <h3>Summary</h3>
                <p>Total AIO operations: ` + strconv.Itoa(result.TotalEvents) + `</p>`

	if result.TotalEvents > 0 {
		averageDuration := result.TotalDuration / time.Duration(result.TotalEvents)
		html += fmt.Sprintf(`
                <p>Average duration: %.3f ms</p>
                <p>Maximum duration: %.3f ms</p>
                <p>Minimum duration: %.3f ms</p>`,
			float64(averageDuration.Microseconds())/1000.0,
			float64(result.MaxDuration.Microseconds())/1000.0,
			float64(result.MinDuration.Microseconds())/1000.0)
	}

	// Add duration counts
	html += `
                <h4>Duration Counts:</h4>`

	// Sort durations for consistent output
	var durations []int
	for duration := range result.DurationCounts {
		durations = append(durations, duration)
	}
	sort.Ints(durations)

	for _, duration := range durations {
		html += fmt.Sprintf(`
                <p>%dms: %d requests</p>`, duration, result.DurationCounts[duration])
	}

	html += `
            </div>
        </div>
        <div class="right-panel">
            <div class="filter-form">
                <h4>Filter Options:</h4>
                <label>Start Time:</label>
                <input type="datetime-local" id="aio-start-time">
                <label>End Time:</label>
                <input type="datetime-local" id="aio-end-time">
                <label>Min Duration (ms):</label>
                <input type="number" id="aio-min-duration" min="0">
                <label>Max Duration (ms):</label>
                <input type="number" id="aio-max-duration" min="0">
                <label>Block Type:</label>
                <select id="aio-block-type">
                    <option value="">All</option>
                    <option value="block">block</option>
                    <option value="block.wal">block.wal</option>
                    <option value="block.db">block.db</option>
                </select>
                <label>Min Length (bytes):</label>
                <input type="number" id="aio-min-length" min="0">
                <label>Max Length (bytes):</label>
                <input type="number" id="aio-max-length" min="0">
                <button type="button" onclick="filterAIOTable()">Filter</button>
                <button type="button" onclick="resetAIOFilter()">Reset</button>
            </div>
            <div class="table-container">
            <table id="aio-table">
                <tr>
                    <th>Start Time</th>
                    <th>End Time</th>
                    <th>Duration (ms)</th>
                    <th>Range</th>
                    <th>Length (bytes)</th>
                    <th>Block Type</th>
                </tr>`

	for _, event := range result.Events {
		html += fmt.Sprintf(`
                <tr>
                    <td>%s</td>
                    <td>%s</td>
                    <td>%.3f</td>
                    <td>%s</td>
                    <td>%d</td>
                    <td>%s</td>
                </tr>`,
		event.StartTime.Format("2006-01-02 15:04:05.000"),
		event.EndTime.Format("2006-01-02 15:04:05.000"),
		float64(event.Duration.Microseconds())/1000.0,
		event.RangeStr,
		event.Length,
		event.BlockType)
	}

	html += `
            </table>
            </div>
        </div>
    </div>`

	return html
}

// generateRepopHTML generates HTML for Repop analysis
func generateRepopHTML(result types.AnalysisResult) string {
	html := `
    <h2>OSD Repop Operations Analysis</h2>
    <div class="layout">
        <div class="left-panel">
            <div class="query-principle">
                <h3>Query Principle</h3>
                <p>This analysis parses OSD repop (replication population) operations from the Ceph log. It identifies:</p>
                <ul>
                    <li><strong>Start events</strong>: Log lines containing "dequeue_op" with "osd_repop"</li>
                    <li><strong>Finish events</strong>: Log lines containing "repop_commit" with "osd_repop"</li>
                </ul>
                <p>For each repop operation, it extracts:</p>
                <ul>
                    <li>Timestamp</li>
                    <li>Operation ID</li>
                </ul>
                <p>It matches start and finish events using the operation ID as a unique key, then calculates the duration for each operation.</p>
            </div>
            <div class="summary">
                <h3>Summary</h3>
                <p>Total repop operations: ` + strconv.Itoa(result.TotalEvents) + `</p>`

	if result.TotalEvents > 0 {
		averageDuration := result.TotalDuration / time.Duration(result.TotalEvents)
		html += fmt.Sprintf(`
                <p>Average duration: %.3f ms</p>
                <p>Maximum duration: %.3f ms</p>
                <p>Minimum duration: %.3f ms</p>`,
			float64(averageDuration.Microseconds())/1000.0,
			float64(result.MaxDuration.Microseconds())/1000.0,
			float64(result.MinDuration.Microseconds())/1000.0)
	}

	// Add duration counts
	html += `
                <h4>Duration Counts:</h4>`

	// Sort durations for consistent output
	var durations []int
	for duration := range result.DurationCounts {
		durations = append(durations, duration)
	}
	sort.Ints(durations)

	for _, duration := range durations {
		html += fmt.Sprintf(`
                <p>%dms: %d requests</p>`, duration, result.DurationCounts[duration])
	}

	html += `
            </div>
        </div>
        <div class="right-panel">
            <div class="filter-form">
                <h4>Filter Options:</h4>
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
            <div class="table-container">
            <table id="repop-table">
                <tr>
                    <th>Start Time</th>
                    <th>End Time</th>
                    <th>Duration (ms)</th>
                    <th>OP ID</th>
                </tr>`

	for _, event := range result.Events {
		html += fmt.Sprintf(`
                <tr>
                    <td>%s</td>
                    <td>%s</td>
                    <td>%.3f</td>
                    <td>%s</td>
                </tr>`,
		event.StartTime.Format("2006-01-02 15:04:05.000"),
		event.EndTime.Format("2006-01-02 15:04:05.000"),
		float64(event.Duration.Microseconds())/1000.0,
		event.OpID)
	}

	html += `
            </table>
            </div>
        </div>
    </div>`

	return html
}

// generateOSDOpHTML generates HTML for OSD Op analysis
func generateOSDOpHTML(result types.OSDOpAnalysisResult) string {
	html := `
    <h2>OSD Operations Analysis</h2>
    <div class="layout">
        <div class="left-panel">
            <div class="query-principle">
                <h3>Query Principle</h3>
                <p>This analysis parses OSD (Object Storage Daemon) operations from the Ceph log. It identifies:</p>
                <ul>
                    <li><strong>OSD operation events</strong>: Log lines containing "log_op_stats osd_op"</li>
                </ul>
                <p>For each OSD operation, it extracts:</p>
                <ul>
                    <li>Timestamp</li>
                    <li>Operation ID</li>
                    <li>PG ID</li>
                    <li>Object name</li>
                    <li>Operation type</li>
                    <li>Range</li>
                    <li>Input bytes</li>
                    <li>Output bytes</li>
                    <li>Latency (converted from seconds to milliseconds)</li>
                </ul>
                <p>It analyzes each operation individually, calculating latency and other metrics directly from the log entries.</p>
            </div>
            <div class="summary">
                <h3>Summary</h3>
                <p>Total operations: ` + strconv.Itoa(result.TotalOps) + `</p>`

	if result.TotalOps > 0 {
		avgLatency := result.TotalLatency / float64(result.TotalOps)
		avgInBytes := result.TotalInBytes / result.TotalOps
		avgOutBytes := result.TotalOutBytes / result.TotalOps
		html += fmt.Sprintf(`
                <p>Average latency: %.6f ms</p>
                <p>Maximum latency: %.6f ms</p>
                <p>Average input: %d bytes</p>
                <p>Average output: %d bytes</p>`,
			avgLatency,
			result.MaxLatency,
			avgInBytes,
			avgOutBytes)
	}

	// Add latency distribution
	html += `
                <h4>Latency Distribution:</h4>`
	latencyRanges := []string{"0-1ms", "1-2ms", "2-3ms", "3-4ms", "4-5ms", "5-10ms", "10ms+"}
	for _, rangeStr := range latencyRanges {
		html += fmt.Sprintf(`
                <p>%s: %d requests</p>`, rangeStr, result.LatencyCounts[rangeStr])
	}

	html += `
            </div>
        </div>
        <div class="right-panel">
            <div class="filter-form">
                <h4>Filter Options:</h4>
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

	for _, event := range result.Events {
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
		event.Timestamp.Format("2006-01-02 15:04:05.000"),
		event.OpID,
		event.PgID,
		event.Object,
		event.OpType,
		event.RangeStr,
		event.InBytes,
		event.OutBytes,
		event.Latency)
	}

	html += `
            </table>
            </div>
        </div>
    </div>`

	return html
}

// generateTransactionHTML generates HTML for Transaction analysis
func generateTransactionHTML(result types.TransactionAnalysisResult) string {
	html := `
    <h2>Transaction Analysis</h2>
    <div class="layout">
        <div class="left-panel">
            <div class="query-principle">
                <h3>Query Principle</h3>
                <p>This analysis parses transaction operations from the Ceph log. It identifies different stages of a transaction:</p>
                <ul>
                    <li><strong>Start events</strong>: Log lines containing "new_repop" with "rep_tid"</li>
                    <li><strong>Issue events</strong>: Log lines containing "issue_repop" with "rep_tid"</li>
                    <li><strong>Reply events</strong>: Log lines containing "do_repop_reply" with "tid"</li>
                    <li><strong>Complete events</strong>: Log lines containing "repop_all_committed" with "repop tid"</li>
                </ul>
                <p>For each transaction, it extracts:</p>
                <ul>
                    <li>Transaction ID (TID)</li>
                    <li>Timestamps for each stage</li>
                    <li>Operation ID</li>
                    <li>Object name</li>
                    <li>Range</li>
                </ul>
                <p>It matches events using the transaction ID as a unique key, then calculates durations for each stage and the total transaction time.</p>
            </div>
            <div class="summary">
                <h3>Summary</h3>
                <p>Total transactions: ` + strconv.Itoa(result.TotalTransactions) + `</p>`

	if result.TotalTransactions > 0 {
		averageDuration := result.TotalDuration / time.Duration(result.TotalTransactions)
		html += fmt.Sprintf(`
                <p>Average total duration: %.3f ms</p>
                <p>Maximum total duration: %.3f ms</p>
                <p>Minimum total duration: %.3f ms</p>`,
			float64(averageDuration.Microseconds())/1000.0,
			float64(result.MaxDuration.Microseconds())/1000.0,
			float64(result.MinDuration.Microseconds())/1000.0)
	}

	// Add duration counts
	html += `
                <h4>Duration Counts:</h4>`

	// Sort durations for consistent output
	var durations []int
	for duration := range result.DurationCounts {
		durations = append(durations, duration)
	}
	sort.Ints(durations)

	for _, duration := range durations {
		html += fmt.Sprintf(`
                <p>%dms: %d transactions</p>`, duration, result.DurationCounts[duration])
	}

	html += `
            </div>
        </div>
        <div class="right-panel">
            <div class="filter-form">
                <h4>Filter Options:</h4>
                <label>Start Time:</label>
                <input type="datetime-local" id="transaction-start-time">
                <label>End Time:</label>
                <input type="datetime-local" id="transaction-end-time">
                <label>Min Duration (ms):</label>
                <input type="number" id="transaction-min-duration" min="0">
                <label>Max Duration (ms):</label>
                <input type="number" id="transaction-max-duration" min="0">
                <button type="button" onclick="filterTransactionTable()">Filter</button>
                <button type="button" onclick="resetTransactionFilter()">Reset</button>
            </div>
            <div class="table-container">
            <table id="transaction-table">
                <tr>
                    <th>TID</th>
                    <th>Start Time</th>
                    <th>Issue Time</th>
                    <th>1st Reply</th>
                    <th>2nd Reply</th>
                    <th>Complete Time</th>
                    <th>Total (ms)</th>
                    <th>Issue (ms)</th>
                    <th>1st Reply (ms)</th>
                    <th>2nd Reply (ms)</th>
                    <th>OP ID</th>
                    <th>Object</th>
                    <th>Range</th>
                </tr>`

	for _, event := range result.Events {
		html += fmt.Sprintf(`
                <tr>
                    <td>%s</td>
                    <td>%s</td>
                    <td>%s</td>
                    <td>%s</td>
                    <td>%s</td>
                    <td>%s</td>
                    <td>%.3f</td>
                    <td>%.3f</td>
                    <td>%.3f</td>
                    <td>%.3f</td>
                    <td>%s</td>
                    <td>%s</td>
                    <td>%s</td>
                </tr>`,
		event.TID,
		event.StartTime.Format("2006-01-02 15:04:05.000"),
		event.IssueTime.Format("2006-01-02 15:04:05.000"),
		event.FirstReplyTime.Format("2006-01-02 15:04:05.000"),
		event.SecondReplyTime.Format("2006-01-02 15:04:05.000"),
		event.CompleteTime.Format("2006-01-02 15:04:05.000"),
		float64(event.TotalDuration.Microseconds())/1000.0,
		float64(event.IssueDuration.Microseconds())/1000.0,
		float64(event.FirstReplyDuration.Microseconds())/1000.0,
		float64(event.SecondReplyDuration.Microseconds())/1000.0,
		event.OpID,
		event.Object,
		event.RangeStr)
	}

	html += `
            </table>
            </div>
        </div>
    </div>`

	return html
}
