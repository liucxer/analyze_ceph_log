package html

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/liucxer/analyze_ceph_log/pkg/types"
)

// GenerateHTML generates HTML for all analysis types
func GenerateHTML(aioResult types.AnalysisResult, repopResult types.AnalysisResult, osdOpResult types.OSDOpAnalysisResult, transactionResult types.TransactionAnalysisResult, metadataResult types.MetadataSyncAnalysisResult) string {
	htmlContent := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
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
            max-width: 1800px;
            margin: 0 auto;
            background-color: white;
            border-radius: 10px;
            box-shadow: 0 0 20px rgba(0,0,0,0.1);
            overflow: hidden;
        }
        /* Language button styles */
        .language-btn {
            padding: 8px 16px;
            margin-left: 10px;
            background-color: #f1f1f1;
            color: #333;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            font-size: 14px;
            font-weight: 500;
            transition: all 0.3s ease;
        }
        .language-btn:hover {
            background-color: #ddd;
        }
        .language-btn.active {
            background-color: #3498db;
            color: white;
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
        <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px;">
            <h1 id="page-title">Ceph 日志分析</h1>
            <div>
                <button onclick="switchLanguage('zh')" id="zh-btn" class="language-btn active">中文</button>
                <button onclick="switchLanguage('en')" id="en-btn" class="language-btn">English</button>
            </div>
        </div>
        
        <!-- Tabs -->
        <div class="tabs">
            <button class="tab" onclick="openTab(event, 'aio')" data-zh="AIO 操作" data-en="AIO Operations">AIO 操作</button>
            <button class="tab" onclick="openTab(event, 'repop')" data-zh="OSD Repop 操作" data-en="OSD Repop Operations">OSD Repop 操作</button>
            <button class="tab" onclick="openTab(event, 'osd')" data-zh="OSD 操作" data-en="OSD Operations">OSD 操作</button>
            <button class="tab active" onclick="openTab(event, 'transaction')" data-zh="复制事务" data-en="Replication Transaction">复制事务</button>
            <button class="tab" onclick="openTab(event, 'metadata')" data-zh="元数据同步" data-en="Metadata Sync">元数据同步</button>
        </div>
        
        <!-- AIO Panel -->
        <div id="aio" class="panel">
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
        <div id="transaction" class="panel active">
`

	// Add Transaction analysis
	htmlContent += generateTransactionHTML(transactionResult)

	htmlContent += `
        </div>
        
        <!-- Metadata Sync Panel -->
        <div id="metadata" class="panel">
`

	// Add Metadata Sync analysis
	htmlContent += generateMetadataSyncHTML(metadataResult)

	htmlContent += `
        </div>
        
        <script>
            // Language switching functionality
            function switchLanguage(lang) {
                try {
                    // Update language buttons
                    var zhBtn = document.getElementById('zh-btn');
                    var enBtn = document.getElementById('en-btn');
                    if (zhBtn) zhBtn.classList.remove('active');
                    if (enBtn) enBtn.classList.remove('active');
                    var langBtn = document.getElementById(lang + '-btn');
                    if (langBtn) langBtn.classList.add('active');
                    
                    // Update page title
                    const pageTitle = document.getElementById('page-title');
                    if (pageTitle) {
                        if (lang === 'zh') {
                            pageTitle.textContent = 'Ceph 日志分析';
                        } else {
                            pageTitle.textContent = 'Ceph Log Analysis';
                        }
                    }
                    
                    // Update tab labels
                    const tabs = document.getElementsByClassName('tab');
                    for (let i = 0; i < tabs.length; i++) {
                        tabs[i].textContent = tabs[i].getAttribute('data-' + lang);
                    }
                    
                    // Update AIO panel
                    updateAIOPanel(lang);
                    
                    // Update Repop panel
                    updateRepopPanel(lang);
                    
                    // Update OSD panel
                    updateOSDPanel(lang);
                    
                    // Update Transaction panel
                    updateTransactionPanel(lang);
                    
                    // Update Metadata panel
                    updateMetadataPanel(lang);
                } catch (e) {
                    console.error("Error in switchLanguage:", e);
                }
            }
            
            function updateAIOPanel(lang) {
                try {
                    const panel = document.getElementById('aio');
                    if (!panel) return;
                    const elements = panel.querySelectorAll('[data-zh]');
                    elements.forEach(el => {
                        el.textContent = el.getAttribute('data-' + lang);
                    });
                } catch (e) {
                    console.error("Error in updateAIOPanel:", e);
                }
            }
            
            function updateRepopPanel(lang) {
                try {
                    const panel = document.getElementById('repop');
                    if (!panel) return;
                    const elements = panel.querySelectorAll('[data-zh]');
                    elements.forEach(el => {
                        el.textContent = el.getAttribute('data-' + lang);
                    });
                } catch (e) {
                    console.error("Error in updateRepopPanel:", e);
                }
            }
            
            function updateOSDPanel(lang) {
                try {
                    const panel = document.getElementById('osd');
                    if (!panel) return;
                    const elements = panel.querySelectorAll('[data-zh]');
                    elements.forEach(el => {
                        el.textContent = el.getAttribute('data-' + lang);
                    });
                } catch (e) {
                    console.error("Error in updateOSDPanel:", e);
                }
            }
            
            function updateTransactionPanel(lang) {
                try {
                    const panel = document.getElementById('transaction');
                    if (!panel) return;
                    const elements = panel.querySelectorAll('[data-zh]');
                    elements.forEach(el => {
                        el.textContent = el.getAttribute('data-' + lang);
                    });
                } catch (e) {
                    console.error("Error in updateTransactionPanel:", e);
                }
            }
            
            function updateMetadataPanel(lang) {
                try {
                    const panel = document.getElementById('metadata');
                    if (!panel) return;
                    const elements = panel.querySelectorAll('[data-zh]');
                    elements.forEach(el => {
                        el.textContent = el.getAttribute('data-' + lang);
                    });
                } catch (e) {
                    console.error("Error in updateMetadataPanel:", e);
                }
            }
            
            function openTab(evt, tabName) {
                try {
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
                    var tabContent = document.getElementById(tabName);
                    if (tabContent) tabContent.classList.add("active");
                    if (evt && evt.currentTarget) evt.currentTarget.classList.add("active");
                } catch (e) {
                    console.error("Error in openTab:", e);
                }
            }
            
            // AIO Table Filter
            function filterAIOTable() {
                try {
                    var startTime = document.getElementById("aio-start-time").value;
                    var endTime = document.getElementById("aio-end-time").value;
                    var minDuration = document.getElementById("aio-min-duration").value;
                    var maxDuration = document.getElementById("aio-max-duration").value;
                    var blockType = document.getElementById("aio-block-type").value;
                    var minLength = document.getElementById("aio-min-length").value;
                    var maxLength = document.getElementById("aio-max-length").value;
                    var table = document.getElementById("aio-table");
                    if (!table) return;
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
                } catch (e) {
                    console.error("Error in filterAIOTable:", e);
                }
            }
            
            function resetAIOFilter() {
                try {
                    var startTimeInput = document.getElementById("aio-start-time");
                    var endTimeInput = document.getElementById("aio-end-time");
                    var minDurationInput = document.getElementById("aio-min-duration");
                    var maxDurationInput = document.getElementById("aio-max-duration");
                    var blockTypeInput = document.getElementById("aio-block-type");
                    var minLengthInput = document.getElementById("aio-min-length");
                    var maxLengthInput = document.getElementById("aio-max-length");
                    
                    if (startTimeInput) startTimeInput.value = "";
                    if (endTimeInput) endTimeInput.value = "";
                    if (minDurationInput) minDurationInput.value = "";
                    if (maxDurationInput) maxDurationInput.value = "";
                    if (blockTypeInput) blockTypeInput.value = "";
                    if (minLengthInput) minLengthInput.value = "";
                    if (maxLengthInput) maxLengthInput.value = "";
                    
                    filterAIOTable();
                } catch (e) {
                    console.error("Error in resetAIOFilter:", e);
                }
            }
            
            // Repop Table Filter
            function filterRepopTable() {
                try {
                    var startTime = document.getElementById("repop-start-time").value;
                    var endTime = document.getElementById("repop-end-time").value;
                    var minDuration = document.getElementById("repop-min-duration").value;
                    var maxDuration = document.getElementById("repop-max-duration").value;
                    var table = document.getElementById("repop-table");
                    if (!table) return;
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
                } catch (e) {
                    console.error("Error in filterRepopTable:", e);
                }
            }
            
            function resetRepopFilter() {
                try {
                    var startTimeInput = document.getElementById("repop-start-time");
                    var endTimeInput = document.getElementById("repop-end-time");
                    var minDurationInput = document.getElementById("repop-min-duration");
                    var maxDurationInput = document.getElementById("repop-max-duration");
                    
                    if (startTimeInput) startTimeInput.value = "";
                    if (endTimeInput) endTimeInput.value = "";
                    if (minDurationInput) minDurationInput.value = "";
                    if (maxDurationInput) maxDurationInput.value = "";
                    
                    filterRepopTable();
                } catch (e) {
                    console.error("Error in resetRepopFilter:", e);
                }
            }
            
            // OSD Table Filter
            function filterOSDTable() {
                try {
                    var startTime = document.getElementById("osd-start-time").value;
                    var endTime = document.getElementById("osd-end-time").value;
                    var minLatency = document.getElementById("osd-min-latency").value;
                    var maxLatency = document.getElementById("osd-max-latency").value;
                    var table = document.getElementById("osd-table");
                    if (!table) return;
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
                } catch (e) {
                    console.error("Error in filterOSDTable:", e);
                }
            }
            
            function resetOSDFilter() {
                try {
                    var startTimeInput = document.getElementById("osd-start-time");
                    var endTimeInput = document.getElementById("osd-end-time");
                    var minLatencyInput = document.getElementById("osd-min-latency");
                    var maxLatencyInput = document.getElementById("osd-max-latency");
                    
                    if (startTimeInput) startTimeInput.value = "";
                    if (endTimeInput) endTimeInput.value = "";
                    if (minLatencyInput) minLatencyInput.value = "";
                    if (maxLatencyInput) maxLatencyInput.value = "";
                    
                    filterOSDTable();
                } catch (e) {
                    console.error("Error in resetOSDFilter:", e);
                }
            }
            
            // Transaction Table Filter
            function filterTransactionTable() {
                try {
                    var startTime = document.getElementById("transaction-start-time").value;
                    var endTime = document.getElementById("transaction-end-time").value;
                    var minDuration = document.getElementById("transaction-min-duration").value;
                    var maxDuration = document.getElementById("transaction-max-duration").value;
                    var table = document.getElementById("transaction-table");
                    if (!table) return;
                    var tr = table.getElementsByTagName("tr");
                    
                    for (var i = 1; i < tr.length; i++) {
                        var tds = tr[i].getElementsByTagName("td");
                        if (tds.length < 7) continue; // Skip rows with insufficient cells
                        
                        var tdStartTime = tds[1].textContent;
                        var tdEndTime = tds[5].textContent;
                        var tdDuration = parseFloat(tds[6].textContent);
                        
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
                } catch (e) {
                    console.error("Error in filterTransactionTable:", e);
                }
            }
            
            function resetTransactionFilter() {
                try {
                    var startTimeInput = document.getElementById("transaction-start-time");
                    var endTimeInput = document.getElementById("transaction-end-time");
                    var minDurationInput = document.getElementById("transaction-min-duration");
                    var maxDurationInput = document.getElementById("transaction-max-duration");
                    
                    if (startTimeInput) startTimeInput.value = "";
                    if (endTimeInput) endTimeInput.value = "";
                    if (minDurationInput) minDurationInput.value = "";
                    if (maxDurationInput) maxDurationInput.value = "";
                    
                    filterTransactionTable();
                } catch (e) {
                    console.error("Error in resetTransactionFilter:", e);
                }
            }
            
            // Metadata Sync Table Filter
            function filterMetadataTable() {
                try {
                    var startTime = document.getElementById("metadata-start-time").value;
                    var endTime = document.getElementById("metadata-end-time").value;
                    var minDuration = document.getElementById("metadata-min-duration").value;
                    var maxDuration = document.getElementById("metadata-max-duration").value;
                    var table = document.getElementById("metadata-table");
                    if (!table) return;
                    var tr = table.getElementsByTagName("tr");
                    
                    for (var i = 1; i < tr.length; i++) {
                        var tdTime = tr[i].getElementsByTagName("td")[0].textContent;
                        var tdDuration = parseFloat(tr[i].getElementsByTagName("td")[3].textContent);
                        
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
                        
                        if (minDuration) {
                            if (tdDuration < parseFloat(minDuration)) match = false;
                        }
                        
                        if (maxDuration) {
                            if (tdDuration > parseFloat(maxDuration)) match = false;
                        }
                        
                        tr[i].style.display = match ? "" : "none";
                    }
                } catch (e) {
                    console.error("Error in filterMetadataTable:", e);
                }
            }
            
            function resetMetadataFilter() {
                try {
                    var startTimeInput = document.getElementById("metadata-start-time");
                    var endTimeInput = document.getElementById("metadata-end-time");
                    var minDurationInput = document.getElementById("metadata-min-duration");
                    var maxDurationInput = document.getElementById("metadata-max-duration");
                    
                    if (startTimeInput) startTimeInput.value = "";
                    if (endTimeInput) endTimeInput.value = "";
                    if (minDurationInput) minDurationInput.value = "";
                    if (maxDurationInput) maxDurationInput.value = "";
                    
                    filterMetadataTable();
                } catch (e) {
                    console.error("Error in resetMetadataFilter:", e);
                }
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
    <h2 data-zh="AIO 操作分析" data-en="AIO Operations Analysis">AIO 操作分析</h2>
    <div class="layout">
        <div class="left-panel">
            <div class="query-principle">
                <h3 data-zh="查询原理" data-en="Query Principle">查询原理</h3>
                <p data-zh="此分析解析 Ceph 日志中的 AIO (异步 I/O) 操作。它识别：" data-en="This analysis parses AIO (Asynchronous I/O) operations from the Ceph log. It identifies:">此分析解析 Ceph 日志中的 AIO (异步 I/O) 操作。它识别：</p>
                <ul>
                    <li data-zh="<strong>开始事件</strong>：包含 "_aio_log_start" 的日志行" data-en="<strong>Start events</strong>: Log lines containing "_aio_log_start""><strong>开始事件</strong>：包含 "_aio_log_start" 的日志行</li>
                    <li data-zh="<strong>结束事件</strong>：包含 "_aio_log_finish" 的日志行" data-en="<strong>Finish events</strong>: Log lines containing "_aio_log_finish""><strong>结束事件</strong>：包含 "_aio_log_finish" 的日志行</li>
                </ul>
                <p data-zh="对于每个 AIO 操作，它提取：" data-en="For each AIO operation, it extracts:">对于每个 AIO 操作，它提取：</p>
                <ul>
                    <li data-zh="时间戳" data-en="Timestamp">时间戳</li>
                    <li data-zh="块地址范围" data-en="Block address range">块地址范围</li>
                    <li data-zh="数据长度（从十六进制转换为十进制）" data-en="Data length (converted from hex to decimal)">数据长度（从十六进制转换为十进制）</li>
                    <li data-zh="块类型（block, block.wal, block.db）" data-en="Block type (block, block.wal, block.db)">块类型（block, block.wal, block.db）</li>
                </ul>
                <p data-zh="它使用块地址范围作为唯一键匹配开始和结束事件，然后计算每个操作的持续时间。" data-en="It matches start and finish events using the block address range as a unique key, then calculates the duration for each operation.">它使用块地址范围作为唯一键匹配开始和结束事件，然后计算每个操作的持续时间。</p>
            </div>
            <div class="summary">
                <h3 data-zh="摘要" data-en="Summary">摘要</h3>
                <p data-zh="总 AIO 操作数：" data-en="Total AIO operations: ">总 AIO 操作数：` + strconv.Itoa(result.TotalEvents) + `</p>`

	if result.TotalEvents > 0 {
		averageDuration := result.TotalDuration / time.Duration(result.TotalEvents)
		html += fmt.Sprintf(`
                <p data-zh="平均持续时间：" data-en="Average duration: ">平均持续时间：%.3f ms</p>
                <p data-zh="最大持续时间：" data-en="Maximum duration: ">最大持续时间：%.3f ms</p>
                <p data-zh="最小持续时间：" data-en="Minimum duration: ">最小持续时间：%.3f ms</p>`,
			float64(averageDuration.Microseconds())/1000.0,
			float64(result.MaxDuration.Microseconds())/1000.0,
			float64(result.MinDuration.Microseconds())/1000.0)
	}

	// Add duration counts
	html += `
                <h4 data-zh="持续时间计数：" data-en="Duration Counts:">持续时间计数：</h4>`

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
                <h4 data-zh="筛选选项：" data-en="Filter Options:">筛选选项：</h4>
                <label data-zh="开始时间：" data-en="Start Time:">开始时间：</label>
                <input type="datetime-local" id="aio-start-time">
                <label data-zh="结束时间：" data-en="End Time:">结束时间：</label>
                <input type="datetime-local" id="aio-end-time">
                <label data-zh="最小持续时间 (ms)：" data-en="Min Duration (ms):">最小持续时间 (ms)：</label>
                <input type="number" id="aio-min-duration" min="0">
                <label data-zh="最大持续时间 (ms)：" data-en="Max Duration (ms):">最大持续时间 (ms)：</label>
                <input type="number" id="aio-max-duration" min="0">
                <label data-zh="块类型：" data-en="Block Type:">块类型：</label>
                <select id="aio-block-type">
                    <option value="" data-zh="全部" data-en="All">全部</option>
                    <option value="block" data-zh="block" data-en="block">block</option>
                    <option value="block.wal" data-zh="block.wal" data-en="block.wal">block.wal</option>
                    <option value="block.db" data-zh="block.db" data-en="block.db">block.db</option>
                </select>
                <label data-zh="最小长度 (字节)：" data-en="Min Length (bytes):">最小长度 (字节)：</label>
                <input type="number" id="aio-min-length" min="0">
                <label data-zh="最大长度 (字节)：" data-en="Max Length (bytes):">最大长度 (字节)：</label>
                <input type="number" id="aio-max-length" min="0">
                <button type="button" onclick="filterAIOTable()" data-zh="筛选" data-en="Filter">筛选</button>
                <button type="button" onclick="resetAIOFilter()" data-zh="重置" data-en="Reset">重置</button>
            </div>
            <div class="table-container">
            <table id="aio-table">
                <tr>
                    <th data-zh="开始时间" data-en="Start Time">开始时间</th>
                    <th data-zh="结束时间" data-en="End Time">结束时间</th>
                    <th data-zh="持续时间 (ms)" data-en="Duration (ms)">持续时间 (ms)</th>
                    <th data-zh="范围" data-en="Range">范围</th>
                    <th data-zh="长度 (字节)" data-en="Length (bytes)">长度 (字节)</th>
                    <th data-zh="块类型" data-en="Block Type">块类型</th>
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
    <h2 data-zh="OSD Repop 操作分析" data-en="OSD Repop Operations Analysis">OSD Repop 操作分析</h2>
    <div class="layout">
        <div class="left-panel">
            <div class="query-principle">
                <h3 data-zh="查询原理" data-en="Query Principle">查询原理</h3>
                <p data-zh="此分析解析 Ceph 日志中的 OSD repop（复制填充）操作。它识别：" data-en="This analysis parses OSD repop (replication population) operations from the Ceph log. It identifies:">此分析解析 Ceph 日志中的 OSD repop（复制填充）操作。它识别：</p>
                <ul>
                    <li data-zh="<strong>开始事件</strong>：包含 "dequeue_op" 和 "osd_repop" 的日志行" data-en="<strong>Start events</strong>: Log lines containing "dequeue_op" with "osd_repop""><strong>开始事件</strong>：包含 "dequeue_op" 和 "osd_repop" 的日志行</li>
                    <li data-zh="<strong>结束事件</strong>：包含 "repop_commit" 和 "osd_repop" 的日志行" data-en="<strong>Finish events</strong>: Log lines containing "repop_commit" with "osd_repop""><strong>结束事件</strong>：包含 "repop_commit" 和 "osd_repop" 的日志行</li>
                </ul>
                <p data-zh="对于每个 repop 操作，它提取：" data-en="For each repop operation, it extracts:">对于每个 repop 操作，它提取：</p>
                <ul>
                    <li data-zh="时间戳" data-en="Timestamp">时间戳</li>
                    <li data-zh="操作 ID" data-en="Operation ID">操作 ID</li>
                </ul>
                <p data-zh="它使用操作 ID 作为唯一键匹配开始和结束事件，然后计算每个操作的持续时间。" data-en="It matches start and finish events using the operation ID as a unique key, then calculates the duration for each operation.">它使用操作 ID 作为唯一键匹配开始和结束事件，然后计算每个操作的持续时间。</p>
            </div>
            <div class="summary">
                <h3 data-zh="摘要" data-en="Summary">摘要</h3>
                <p data-zh="总 repop 操作数：" data-en="Total repop operations: ">总 repop 操作数：` + strconv.Itoa(result.TotalEvents) + `</p>`

	if result.TotalEvents > 0 {
		averageDuration := result.TotalDuration / time.Duration(result.TotalEvents)
		html += fmt.Sprintf(`
                <p data-zh="平均持续时间：" data-en="Average duration: ">平均持续时间：%.3f ms</p>
                <p data-zh="最大持续时间：" data-en="Maximum duration: ">最大持续时间：%.3f ms</p>
                <p data-zh="最小持续时间：" data-en="Minimum duration: ">最小持续时间：%.3f ms</p>`,
			float64(averageDuration.Microseconds())/1000.0,
			float64(result.MaxDuration.Microseconds())/1000.0,
			float64(result.MinDuration.Microseconds())/1000.0)
	}

	// Add duration counts
	html += `
                <h4 data-zh="持续时间计数：" data-en="Duration Counts:">持续时间计数：</h4>`

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
                <h4 data-zh="筛选选项：" data-en="Filter Options:">筛选选项：</h4>
                <label data-zh="开始时间：" data-en="Start Time:">开始时间：</label>
                <input type="datetime-local" id="repop-start-time">
                <label data-zh="结束时间：" data-en="End Time:">结束时间：</label>
                <input type="datetime-local" id="repop-end-time">
                <label data-zh="最小持续时间 (ms)：" data-en="Min Duration (ms):">最小持续时间 (ms)：</label>
                <input type="number" id="repop-min-duration" min="0">
                <label data-zh="最大持续时间 (ms)：" data-en="Max Duration (ms):">最大持续时间 (ms)：</label>
                <input type="number" id="repop-max-duration" min="0">
                <button type="button" onclick="filterRepopTable()" data-zh="筛选" data-en="Filter">筛选</button>
                <button type="button" onclick="resetRepopFilter()" data-zh="重置" data-en="Reset">重置</button>
            </div>
            <div class="table-container">
            <table id="repop-table">
                <tr>
                    <th data-zh="开始时间" data-en="Start Time">开始时间</th>
                    <th data-zh="结束时间" data-en="End Time">结束时间</th>
                    <th data-zh="持续时间 (ms)" data-en="Duration (ms)">持续时间 (ms)</th>
                    <th data-zh="操作 ID" data-en="OP ID">操作 ID</th>
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
    <h2 data-zh="OSD 操作分析" data-en="OSD Operations Analysis">OSD 操作分析</h2>
    <div class="layout">
        <div class="left-panel">
            <div class="query-principle">
                <h3 data-zh="查询原理" data-en="Query Principle">查询原理</h3>
                <p data-zh="此分析解析 Ceph 日志中的 OSD（对象存储守护进程）操作。它识别：" data-en="This analysis parses OSD (Object Storage Daemon) operations from the Ceph log. It identifies:">此分析解析 Ceph 日志中的 OSD（对象存储守护进程）操作。它识别：</p>
                <ul>
                    <li data-zh="<strong>OSD 操作事件</strong>：包含 "log_op_stats osd_op" 的日志行" data-en="<strong>OSD operation events</strong>: Log lines containing "log_op_stats osd_op""><strong>OSD 操作事件</strong>：包含 "log_op_stats osd_op" 的日志行</li>
                </ul>
                <p data-zh="对于每个 OSD 操作，它提取：" data-en="For each OSD operation, it extracts:">对于每个 OSD 操作，它提取：</p>
                <ul>
                    <li data-zh="时间戳" data-en="Timestamp">时间戳</li>
                    <li data-zh="操作 ID" data-en="Operation ID">操作 ID</li>
                    <li data-zh="PG ID" data-en="PG ID">PG ID</li>
                    <li data-zh="对象名称" data-en="Object name">对象名称</li>
                    <li data-zh="操作类型" data-en="Operation type">操作类型</li>
                    <li data-zh="范围" data-en="Range">范围</li>
                    <li data-zh="输入字节" data-en="Input bytes">输入字节</li>
                    <li data-zh="输出字节" data-en="Output bytes">输出字节</li>
                    <li data-zh="延迟（从秒转换为毫秒）" data-en="Latency (converted from seconds to milliseconds)">延迟（从秒转换为毫秒）</li>
                </ul>
                <p data-zh="它单独分析每个操作，直接从日志条目中计算延迟和其他指标。" data-en="It analyzes each operation individually, calculating latency and other metrics directly from the log entries.">它单独分析每个操作，直接从日志条目中计算延迟和其他指标。</p>
            </div>
            <div class="summary">
                <h3 data-zh="摘要" data-en="Summary">摘要</h3>
                <p data-zh="总操作数：" data-en="Total operations: ">总操作数：` + strconv.Itoa(result.TotalOps) + `</p>`

	if result.TotalOps > 0 {
		avgLatency := result.TotalLatency / float64(result.TotalOps)
		avgInBytes := result.TotalInBytes / result.TotalOps
		avgOutBytes := result.TotalOutBytes / result.TotalOps
		html += fmt.Sprintf(`
                <p data-zh="平均延迟：" data-en="Average latency: ">平均延迟：%.6f ms</p>
                <p data-zh="最大延迟：" data-en="Maximum latency: ">最大延迟：%.6f ms</p>
                <p data-zh="平均输入：" data-en="Average input: ">平均输入：%d bytes</p>
                <p data-zh="平均输出：" data-en="Average output: ">平均输出：%d bytes</p>`,
			avgLatency,
			result.MaxLatency,
			avgInBytes,
			avgOutBytes)
	}

	// Add latency distribution
	html += `
                <h4 data-zh="延迟分布：" data-en="Latency Distribution:">延迟分布：</h4>`
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
                <h4 data-zh="筛选选项：" data-en="Filter Options:">筛选选项：</h4>
                <label data-zh="开始时间：" data-en="Start Time:">开始时间：</label>
                <input type="datetime-local" id="osd-start-time">
                <label data-zh="结束时间：" data-en="End Time:">结束时间：</label>
                <input type="datetime-local" id="osd-end-time">
                <label data-zh="最小延迟 (ms)：" data-en="Min Latency (ms):">最小延迟 (ms)：</label>
                <input type="number" id="osd-min-latency" min="0">
                <label data-zh="最大延迟 (ms)：" data-en="Max Latency (ms):">最大延迟 (ms)：</label>
                <input type="number" id="osd-max-latency" min="0">
                <button type="button" onclick="filterOSDTable()" data-zh="筛选" data-en="Filter">筛选</button>
                <button type="button" onclick="resetOSDFilter()" data-zh="重置" data-en="Reset">重置</button>
            </div>
            <div class="table-container">
            <table id="osd-table">
                <tr>
                    <th data-zh="时间戳" data-en="Timestamp">时间戳</th>
                    <th data-zh="操作 ID" data-en="OP ID">操作 ID</th>
                    <th data-zh="PG ID" data-en="PG ID">PG ID</th>
                    <th data-zh="对象" data-en="Object">对象</th>
                    <th data-zh="操作类型" data-en="Op Type">操作类型</th>
                    <th data-zh="范围" data-en="Range">范围</th>
                    <th data-zh="输入 (字节)" data-en="In (bytes)">输入 (字节)</th>
                    <th data-zh="输出 (字节)" data-en="Out (bytes)">输出 (字节)</th>
                    <th data-zh="延迟 (ms)" data-en="Latency (ms)">延迟 (ms)</th>
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
    <h2 data-zh="复制事务" data-en="Replication Transaction">复制事务</h2>
    <div class="layout">
        <div class="left-panel">
            <div class="query-principle">
                <h3 data-zh="查询原理" data-en="Query Principle">查询原理</h3>
                <p data-zh="此分析解析 Ceph 日志中的复制事务操作。它识别复制事务的不同阶段：" data-en="This analysis parses replication transaction operations from the Ceph log. It identifies different stages of a replication transaction:">此分析解析 Ceph 日志中的复制事务操作。它识别复制事务的不同阶段：</p>
                <ul>
                    <li data-zh="<strong>开始事件</strong>：包含 "new_repop" 和 "rep_tid" 的日志行" data-en="<strong>Start events</strong>: Log lines containing "new_repop" with "rep_tid""</li>
                    <li data-zh="<strong>发出事件</strong>：包含 "issue_repop" 和 "rep_tid" 的日志行" data-en="<strong>Issue events</strong>: Log lines containing "issue_repop" with "rep_tid""</li>
                    <li data-zh="<strong>回复事件</strong>：包含 "do_repop_reply" 和 "tid" 的日志行" data-en="<strong>Reply events</strong>: Log lines containing "do_repop_reply" with "tid""</li>
                    <li data-zh="<strong>完成事件</strong>：包含 "repop_all_committed" 和 "repop tid" 的日志行" data-en="<strong>Complete events</strong>: Log lines containing "repop_all_committed" with "repop tid""</li>
                </ul>
                <p data-zh="对于每个复制事务，它提取：" data-en="For each replication transaction, it extracts:">对于每个复制事务，它提取：</p>
                <ul>
                    <li data-zh="事务 ID (TID)" data-en="Transaction ID (TID)">事务 ID (TID)</li>
                    <li data-zh="每个阶段的时间戳" data-en="Timestamps for each stage">每个阶段的时间戳</li>
                    <li data-zh="操作 ID" data-en="Operation ID">操作 ID</li>
                    <li data-zh="对象名称" data-en="Object name">对象名称</li>
                    <li data-zh="范围" data-en="Range">范围</li>
                </ul>
                <p data-zh="它使用复制事务 ID 作为唯一键匹配事件，然后计算每个阶段的持续时间和总复制事务时间。" data-en="It matches events using the replication transaction ID as a unique key, then calculates durations for each stage and the total replication transaction time.">它使用复制事务 ID 作为唯一键匹配事件，然后计算每个阶段的持续时间和总复制事务时间。</p>
            </div>
            <div class="summary">
                <h3 data-zh="摘要" data-en="Summary">摘要</h3>
                <p data-zh="总复制事务数：" data-en="Total replication transactions: ">总复制事务数：` + strconv.Itoa(result.TotalTransactions) + `</p>`

	if result.TotalTransactions > 0 {
		averageDuration := result.TotalDuration / time.Duration(result.TotalTransactions)
		html += fmt.Sprintf(`
                <p data-zh="平均复制事务持续时间：" data-en="Average replication transaction duration: ">平均复制事务持续时间：%.3f ms</p>
                <p data-zh="最大复制事务持续时间：" data-en="Maximum replication transaction duration: ">最大复制事务持续时间：%.3f ms</p>
                <p data-zh="最小复制事务持续时间：" data-en="Minimum replication transaction duration: ">最小复制事务持续时间：%.3f ms</p>`,
			float64(averageDuration.Microseconds())/1000.0,
			float64(result.MaxDuration.Microseconds())/1000.0,
			float64(result.MinDuration.Microseconds())/1000.0)
	}

	// Add duration counts
	html += `
                <h4 data-zh="复制事务持续时间计数：" data-en="Replication Transaction Duration Counts:">复制事务持续时间计数：</h4>`

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
                <h4 data-zh="筛选选项：" data-en="Filter Options:">筛选选项：</h4>
                <label data-zh="开始时间：" data-en="Start Time:">开始时间：</label>
                <input type="datetime-local" id="transaction-start-time">
                <label data-zh="结束时间：" data-en="End Time:">结束时间：</label>
                <input type="datetime-local" id="transaction-end-time">
                <label data-zh="最小持续时间 (ms)：" data-en="Min Duration (ms):">最小持续时间 (ms)：</label>
                <input type="number" id="transaction-min-duration" min="0">
                <label data-zh="最大持续时间 (ms)：" data-en="Max Duration (ms):">最大持续时间 (ms)：</label>
                <input type="number" id="transaction-max-duration" min="0">
                <button type="button" onclick="filterTransactionTable()" data-zh="筛选" data-en="Filter">筛选</button>
                <button type="button" onclick="resetTransactionFilter()" data-zh="重置" data-en="Reset">重置</button>
            </div>
            <div class="table-container">
            <table id="transaction-table">
                <tr>
                    <th data-zh="TID" data-en="TID">TID</th>
                    <th data-zh="开始时间" data-en="Start Time">开始时间</th>
                    <th data-zh="发出时间" data-en="Issue Time">发出时间</th>
                    <th data-zh="第一次回复" data-en="1st Reply">第一次回复</th>
                    <th data-zh="第二次回复" data-en="2nd Reply">第二次回复</th>
                    <th data-zh="完成时间" data-en="Complete Time">完成时间</th>
                    <th data-zh="总耗时 (ms)" data-en="Total (ms)">总耗时 (ms)</th>
                    <th data-zh="发出耗时 (ms)" data-en="Issue (ms)">发出耗时 (ms)</th>
                    <th data-zh="第一次回复耗时 (ms)" data-en="1st Reply (ms)">第一次回复耗时 (ms)</th>
                    <th data-zh="第二次回复耗时 (ms)" data-en="2nd Reply (ms)">第二次回复耗时 (ms)</th>
                    <th data-zh="操作 ID" data-en="OP ID">操作 ID</th>
                    <th data-zh="对象" data-en="Object">对象</th>
                    <th data-zh="范围" data-en="Range">范围</th>
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

// generateMetadataSyncHTML generates HTML for Metadata Sync analysis
func generateMetadataSyncHTML(result types.MetadataSyncAnalysisResult) string {
	html := `
    <h2 data-zh="元数据同步" data-en="Metadata Sync">元数据同步</h2>
    <div class="layout">
        <div class="left-panel">
            <div class="query-principle">
                <h3 data-zh="查询原理" data-en="Query Principle">查询原理</h3>
                <p data-zh="此分析解析 Ceph 日志中的元数据同步操作。它识别：" data-en="This analysis parses metadata sync operations from the Ceph log. It identifies:">此分析解析 Ceph 日志中的元数据同步操作。它识别：</p>
                <ul>
                    <li data-zh="<strong>元数据同步事件</strong>：包含 "_kv_sync_thread committed" 的日志行" data-en="<strong>Metadata sync events</strong>: Log lines containing "_kv_sync_thread committed""</li>
                </ul>
                <p data-zh="对于每个元数据同步操作，它提取：" data-en="For each metadata sync operation, it extracts:">对于每个元数据同步操作，它提取：</p>
                <ul>
                    <li data-zh="时间戳" data-en="Timestamp">时间戳</li>
                    <li data-zh="提交的操作数" data-en="Number of committed operations">提交的操作数</li>
                    <li data-zh="清理的操作数" data-en="Number of cleaned operations">清理的操作数</li>
                    <li data-zh="总持续时间" data-en="Total duration">总持续时间</li>
                    <li data-zh="刷新时间" data-en="Flush time">刷新时间</li>
                    <li data-zh="KV 提交时间" data-en="KV commit time">KV 提交时间</li>
                </ul>
                <p data-zh="它分析元数据同步线程的执行情况，包括同步操作的频率、耗时和详细分解。" data-en="It analyzes the execution of metadata sync threads, including the frequency, duration, and detailed breakdown of sync operations.">它分析元数据同步线程的执行情况，包括同步操作的频率、耗时和详细分解。</p>
            </div>
            <div class="summary">
                <h3 data-zh="摘要" data-en="Summary">摘要</h3>
                <p data-zh="总元数据同步操作数：" data-en="Total metadata sync operations: ">总元数据同步操作数：` + strconv.Itoa(result.TotalEvents) + `</p>`

	if result.TotalEvents > 0 {
		averageDuration := result.TotalDuration / time.Duration(result.TotalEvents)
		html += fmt.Sprintf(`
                <p data-zh="平均持续时间：" data-en="Average duration: ">平均持续时间：%.3f ms</p>
                <p data-zh="最大持续时间：" data-en="Maximum duration: ">最大持续时间：%.3f ms</p>
                <p data-zh="最小持续时间：" data-en="Minimum duration: ">最小持续时间：%.3f ms</p>`,
			float64(averageDuration.Microseconds())/1000.0,
			float64(result.MaxDuration.Microseconds())/1000.0,
			float64(result.MinDuration.Microseconds())/1000.0)
	}

	// Add duration counts
	html += `
                <h4 data-zh="持续时间计数：" data-en="Duration Counts:">持续时间计数：</h4>`

	// Sort durations for consistent output
	var durations []int
	for duration := range result.DurationCounts {
		durations = append(durations, duration)
	}
	sort.Ints(durations)

	for _, duration := range durations {
		html += fmt.Sprintf(`
                <p>%dms: %d operations</p>`, duration, result.DurationCounts[duration])
	}

	html += `
            </div>
        </div>
        <div class="right-panel">
            <div class="filter-form">
                <h4 data-zh="筛选选项：" data-en="Filter Options:">筛选选项：</h4>
                <label data-zh="开始时间：" data-en="Start Time:">开始时间：</label>
                <input type="datetime-local" id="metadata-start-time">
                <label data-zh="结束时间：" data-en="End Time:">结束时间：</label>
                <input type="datetime-local" id="metadata-end-time">
                <label data-zh="最小持续时间 (ms)：" data-en="Min Duration (ms):">最小持续时间 (ms)：</label>
                <input type="number" id="metadata-min-duration" min="0">
                <label data-zh="最大持续时间 (ms)：" data-en="Max Duration (ms):">最大持续时间 (ms)：</label>
                <input type="number" id="metadata-max-duration" min="0">
                <button type="button" onclick="filterMetadataTable()" data-zh="筛选" data-en="Filter">筛选</button>
                <button type="button" onclick="resetMetadataFilter()" data-zh="重置" data-en="Reset">重置</button>
            </div>
            <div class="table-container">
            <table id="metadata-table">
                <tr>
                    <th data-zh="时间戳" data-en="Timestamp">时间戳</th>
                    <th data-zh="提交数" data-en="Committed">提交数</th>
                    <th data-zh="清理数" data-en="Cleaned">清理数</th>
                    <th data-zh="总耗时 (ms)" data-en="Total (ms)">总耗时 (ms)</th>
                    <th data-zh="刷新耗时 (ms)" data-en="Flush (ms)">刷新耗时 (ms)</th>
                    <th data-zh="KV 提交耗时 (ms)" data-en="KV Commit (ms)">KV 提交耗时 (ms)</th>
                </tr>`

	for _, event := range result.Events {
		html += fmt.Sprintf(`
                <tr>
                    <td>%s</td>
                    <td>%d</td>
                    <td>%d</td>
                    <td>%.3f</td>
                    <td>%.6f</td>
                    <td>%.3f</td>
                </tr>`,
			event.Timestamp.Format("2006-01-02 15:04:05.000"),
			event.Committed,
			event.Cleaned,
			float64(event.Duration.Microseconds())/1000.0,
			float64(event.FlushTime.Microseconds())/1000.0,
			float64(event.KVCommitTime.Microseconds())/1000.0)
	}

	html += `
            </table>
            </div>
        </div>
    </div>`

	return html
}
