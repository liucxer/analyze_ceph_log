package panels

import (
	"fmt"
	"strconv"

	"github.com/liucxer/analyze_ceph_log/pkg/types"
)

// GenerateOSDHTML generates HTML for OSD Op analysis
func GenerateOSDHTML(result types.OSDOpAnalysisResult) string {
	html := `
    <h2 data-zh="OSD 事件分析" data-en="OSD Events Analysis">OSD 事件分析</h2>
    <div class="layout">
        <div class="left-panel">
            <div class="query-principle">
                <h3 data-zh="查询原理" data-en="Query Principle">查询原理</h3>
                <p data-zh="此分析解析 Ceph 日志中的 OSD（对象存储守护进程）事件。它识别：" data-en="This analysis parses OSD (Object Storage Daemon) events from the Ceph log. It identifies:">此分析解析 Ceph 日志中的 OSD（对象存储守护进程）事件。它识别：</p>
                <ul>
                    <li data-zh="<strong>OSD 事件</strong>：包含 "log_op_stats osd_op" 的日志行" data-en="<strong>OSD events</strong>: Log lines containing "log_op_stats osd_op"">OSD 事件：包含 "log_op_stats osd_op" 的日志行</li>
                </ul>
                <p data-zh="对于每个 OSD 事件，它提取：" data-en="For each OSD event, it extracts:">对于每个 OSD 事件，它提取：</p>
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
                <p data-zh="它单独分析每个事件，直接从日志条目中计算延迟和其他指标。" data-en="It analyzes each event individually, calculating latency and other metrics directly from the log entries.">它单独分析每个事件，直接从日志条目中计算延迟和其他指标。</p>
            </div>
            <div class="summary">
                <h3 data-zh="摘要" data-en="Summary">摘要</h3>
                <p data-zh="总事件数：" data-en="Total events: ">总事件数：` + strconv.Itoa(result.TotalOps) + `</p>`

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
            <div class="pagination">
                <div class="pagination-info" id="osd-table-pagination-info"></div>
                <div class="pagination-buttons" id="osd-table-pagination-buttons"></div>
            </div>
        </div>
    </div>`

	return html
}
