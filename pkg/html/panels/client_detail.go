package panels

import (
	"fmt"
	"strconv"

	"github.com/liucxer/analyze_ceph_log/pkg/types"
)

// GenerateClientDetailHTML generates HTML for Client Operation Detail analysis
func GenerateClientDetailHTML(result types.ClientOpAnalysisResult) string {
	html := `
    <h2 data-zh="客户端事件详情分析" data-en="Client Event Detail Analysis">客户端事件详情分析</h2>
    <div class="layout">
        <div class="left-panel">
            <div class="query-principle">
                <h3 data-zh="查询原理" data-en="Query Principle">查询原理</h3>
                <p data-zh="此分析解析 Ceph 日志中的单个客户端事件详情。它识别：" data-en="This analysis parses detailed client event information from the Ceph log. It identifies:">此分析解析 Ceph 日志中的单个客户端事件详情。它识别：</p>
                <ul>
                    <li data-zh="<strong>客户端事件入队</strong>：osd_op 日志行" data-en="<strong>Client event enqueue</strong>: osd_op log lines">客户端事件入队：osd_op 日志行</li>
                    <li data-zh="<strong>客户端事件出队</strong>：dequeue_op 日志行" data-en="<strong>Client event dequeue</strong>: dequeue_op log lines">客户端事件出队：dequeue_op 日志行</li>
                    <li data-zh="<strong>事件处理</strong>：do_op 日志行" data-en="<strong>Event processing</strong>: do_op log lines">事件处理：do_op 日志行</li>
                    <li data-zh="<strong>复制事务创建</strong>：new_repop 日志行" data-en="<strong>Replication transaction creation</strong>: new_repop log lines">复制事务创建：new_repop 日志行</li>
                    <li data-zh="<strong>日志追加</strong>：append_log 日志行" data-en="<strong>Log append</strong>: append_log log lines">日志追加：append_log 日志行</li>
                    <li data-zh="<strong>复制回复</strong>：osd_repop_reply 日志行" data-en="<strong>Replication reply</strong>: osd_repop_reply log lines">复制回复：osd_repop_reply 日志行</li>
                    <li data-zh="<strong>事件统计</strong>：log_op_stats 日志行" data-en="<strong>Event statistics</strong>: log_op_stats log lines">事件统计：log_op_stats 日志行</li>
                </ul>
                <p data-zh="它分析单个客户端事件的完整执行流程，包括各阶段的时间点和延迟。" data-en="It analyzes the complete execution flow of a single client event, including timestamps and latencies at each stage.">它分析单个客户端事件的完整执行流程，包括各阶段的时间点和延迟。</p>
            </div>
            <div class="summary">
                <h3 data-zh="摘要" data-en="Summary">摘要</h3>
                <p data-zh="总客户端事件数：" data-en="Total client events: ">总客户端事件数：` + strconv.Itoa(result.TotalOps) + `</p>`

	if result.TotalOps > 0 {
		avgLatency := result.TotalLatency / float64(result.TotalOps)
		avgDataSize := result.TotalDataSize / result.TotalOps
		html += fmt.Sprintf(`
                <p data-zh="平均延迟：" data-en="Average latency: ">平均延迟：%.6f ms</p>
                <p data-zh="最大延迟：" data-en="Maximum latency: ">最大延迟：%.6f ms</p>
                <p data-zh="平均数据大小：" data-en="Average data size: ">平均数据大小：%d bytes</p>
                <p data-zh="总数据量：" data-en="Total data size: ">总数据量：%d bytes</p>`,
			avgLatency,
			result.MaxLatency,
			avgDataSize,
			result.TotalDataSize)
	}

	// Add latency distribution
	html += `
                <h4 data-zh="延迟分布：" data-en="Latency Distribution:">延迟分布：</h4>`
	latencyRanges := []string{"0-1ms", "1-2ms", "2-3ms", "3-4ms", "4-5ms", "5-10ms", "10-50ms", "50-100ms", "100-200ms", "200-400ms", "400ms+"}
	for _, rangeStr := range latencyRanges {
		html += fmt.Sprintf(`
                <p>%s: %d operations</p>`, rangeStr, result.LatencyCounts[rangeStr])
	}

	html += `
            </div>
        </div>
        <div class="right-panel">
            <div class="filter-form">
                <h4 data-zh="筛选选项：" data-en="Filter Options:">筛选选项：</h4>
                <label data-zh="开始时间：" data-en="Start Time:">开始时间：</label>
                <input type="datetime-local" id="client-detail-start-time">
                <label data-zh="结束时间：" data-en="End Time:">结束时间：</label>
                <input type="datetime-local" id="client-detail-end-time">
                <label data-zh="最小延迟 (ms)：" data-en="Min Latency (ms):">最小延迟 (ms)：</label>
                <input type="number" id="client-detail-min-latency" min="0">
                <label data-zh="最大延迟 (ms)：" data-en="Max Latency (ms):">最大延迟 (ms)：</label>
                <input type="number" id="client-detail-max-latency" min="0">
                <button type="button" onclick="filterClientDetailTable()" data-zh="筛选" data-en="Filter">筛选</button>
                <button type="button" onclick="resetClientDetailFilter()" data-zh="重置" data-en="Reset">重置</button>
            </div>
            <div class="table-container">
            <table id="client-detail-table">
                <tr>
                    <th data-zh="时间戳" data-en="Timestamp">时间戳</th>
                    <th data-zh="客户端 ID" data-en="Client ID">客户端 ID</th>
                    <th data-zh="操作 ID" data-en="Op ID">操作 ID</th>
                    <th data-zh="PG ID" data-en="PG ID">PG ID</th>
                    <th data-zh="对象" data-en="Object">对象</th>
                    <th data-zh="数据大小 (字节)" data-en="Data Size (bytes)">数据大小 (字节)</th>
                    <th data-zh="总延迟 (ms)" data-en="Total Latency (ms)">总延迟 (ms)</th>
                    <th data-zh="第一次回复 (ms)" data-en="First Reply (ms)">第一次回复 (ms)</th>
                    <th data-zh="第二次回复 (ms)" data-en="Second Reply (ms)">第二次回复 (ms)</th>
                    <th data-zh="本地处理 (ms)" data-en="Local Processing (ms)">本地处理 (ms)</th>
                </tr>`

	for _, event := range result.Events {
		html += fmt.Sprintf(`
                <tr>
                    <td>%s</td>
                    <td>%s</td>
                    <td>%s</td>
                    <td>%s</td>
                    <td>%s</td>
                    <td>%d</td>
                    <td>%.6f</td>
                    <td>%.6f</td>
                    <td>%.6f</td>
                    <td>%.6f</td>
                </tr>`,
			event.Timestamp.Format("2006-01-02 15:04:05.000"),
			event.ClientID,
			event.OpID,
			event.PGID,
			event.Object,
			event.DataSize,
			event.TotalLatency,
			event.FirstReplyLatency,
			event.SecondReplyLatency,
			event.LocalProcessingTime)
	}

	html += `
            </table>
            </div>
            <div class="pagination">
                <div class="pagination-info" id="client-detail-table-pagination-info"></div>
                <div class="pagination-buttons" id="client-detail-table-pagination-buttons"></div>
            </div>
        </div>
    </div>`

	return html
}