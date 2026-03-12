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
    </div>`

	return html
}

// generateRepopHTML generates HTML for Repop analysis
func generateRepopHTML(result types.AnalysisResult) string {
	html := `
    <h2>OSD Repop Operations Analysis</h2>
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
    </div>`

	return html
}

// generateOSDOpHTML generates HTML for OSD Op analysis
func generateOSDOpHTML(result types.OSDOpAnalysisResult) string {
	html := `
    <h2>OSD Operations Analysis</h2>
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
            <th>1st Repop (ms)</th>
            <th>2nd Repop (ms)</th>
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
            <td>%.3f</td>
            <td>%.3f</td>
        </tr>`,
		event.Timestamp.Format("2006-01-02 15:04:05.000"),
		event.OpID,
		event.PgID,
		event.Object,
		event.OpType,
		event.RangeStr,
		event.InBytes,
		event.OutBytes,
		event.Latency,
		event.FirstRepopReply,
		event.SecondRepopReply)
	}

	html += `
    </table>
    </div>`

	return html
}

// generateTransactionHTML generates HTML for Transaction analysis
func generateTransactionHTML(result types.TransactionAnalysisResult) string {
	html := `
    <h2>Transaction Analysis</h2>
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

	// Add filter form
	html += `
        <h4>Filter Options:</h4>
        <div class="filter-form">
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
    </div>`

	return html
}
