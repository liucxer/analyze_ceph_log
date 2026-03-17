package html

import (
	"github.com/liucxer/analyze_ceph_log/pkg/html/panels"
	"github.com/liucxer/analyze_ceph_log/pkg/types"
)

// GenerateHTML generates HTML for all analysis types
func GenerateHTML(aioResult types.AnalysisResult, repopResult types.AnalysisResult, osdOpResult types.OSDOpAnalysisResult, transactionResult types.TransactionAnalysisResult, metadataResult types.MetadataSyncAnalysisResult, clientResult types.ClientOpAnalysisResult, dequeueResult types.DequeueAnalysisResult, osdOpResultV2 types.AnalysisResult) string {
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
        /* Pagination styles */
        .pagination {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-top: 15px;
            padding: 10px;
            background-color: #f8f9fa;
            border-radius: 8px;
        }
        .pagination-info {
            font-size: 14px;
            color: #666;
        }
        .pagination-buttons button {
            padding: 6px 12px;
            margin: 0 3px;
            background-color: #3498db;
            color: white;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            font-size: 14px;
        }
        .pagination-buttons button:hover {
            background-color: #2980b9;
        }
        .pagination-buttons button.active {
            background-color: #2c3e50;
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
            <button class="tab active" onclick="openTab(event, 'osdop')" data-zh="osd_op事件" data-en="osd_op Events">osd_op事件</button>
            <button class="tab" onclick="openTab(event, 'repop')" data-zh="osd_repop事件" data-en="osd_repop Events">osd_repop事件</button>
            <button class="tab" onclick="openTab(event, 'dequeue')" data-zh="dequeue_op事件" data-en="dequeue_op Events">dequeue_op事件</button>
            <button class="tab" onclick="openTab(event, 'aio')" data-zh="AIO 事件" data-en="AIO Events">AIO 事件</button>
            <button class="tab" onclick="openTab(event, 'metadata')" data-zh="元数据同步" data-en="Metadata Sync">元数据同步</button>
            <button class="tab" onclick="openTab(event, 'osd')" data-zh="OSD 事件" data-en="OSD Events">OSD 事件</button>
            <button class="tab" onclick="openTab(event, 'transaction')" data-zh="复制事务" data-en="Replication Transaction">复制事务</button>
            <button class="tab" onclick="openTab(event, 'client')" data-zh="客户端事件" data-en="Client Events">客户端事件</button>
            <button class="tab" onclick="openTab(event, 'client-detail')" data-zh="客户端事件详情" data-en="Client Event Detail">客户端事件详情</button>
        </div>
        
        <!-- OSD Op Events Panel -->
        <div id="osdop" class="panel active">
`

	// Add OSD Op Events analysis
	htmlContent += panels.GenerateOSDOpEventsHTML(osdOpResultV2)

	htmlContent += `
        </div>
        
        <!-- AIO Panel -->
        <div id="aio" class="panel">
`

	// Add AIO analysis
	htmlContent += panels.GenerateAIOHTML(aioResult)

	htmlContent += `
        </div>
        
        <!-- Repop Panel -->
        <div id="repop" class="panel">
`

	// Add Repop analysis
	htmlContent += panels.GenerateRepopHTML(repopResult)

	htmlContent += `
        </div>
        
        <!-- OSD Panel -->
        <div id="osd" class="panel">
`

	// Add OSD Op analysis
	htmlContent += panels.GenerateOSDHTML(osdOpResult)

	htmlContent += `
        </div>
        
        <!-- Transaction Panel -->
        <div id="transaction" class="panel">
`

	// Add Transaction analysis
	htmlContent += panels.GenerateTransactionHTML(transactionResult)

	htmlContent += `
        </div>
        
        <!-- Metadata Sync Panel -->
        <div id="metadata" class="panel">
`

	// Add Metadata Sync analysis
	htmlContent += panels.GenerateMetadataHTML(metadataResult)

	htmlContent += `
        </div>
        
        <!-- Client Operation Panel -->
        <div id="client" class="panel">
`

	// Add Client Operation analysis
	htmlContent += panels.GenerateClientOpHTML(clientResult)

	htmlContent += `
        </div>
        
        <!-- Client Operation Detail Panel -->
        <div id="client-detail" class="panel">
`

	// Add Client Operation Detail analysis
	htmlContent += panels.GenerateClientDetailHTML(clientResult)

	htmlContent += `
        </div>
        
        <!-- Dequeue Panel -->
        <div id="dequeue" class="panel">
`

	// Add Dequeue analysis
	htmlContent += panels.GenerateDequeueHTML(dequeueResult)

	htmlContent += `
        </div>
        
        <script src="js/script.js"></script>
    </div>
</body>
</html>`

	return htmlContent
}