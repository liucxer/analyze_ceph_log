package panels

import (
	"fmt"
	"strconv"

	"github.com/liucxer/analyze_ceph_log/pkg/types"
)

// GenerateDequeueHTML generates HTML for Dequeue analysis
func GenerateDequeueHTML(result types.DequeueAnalysisResult) string {
	html := `
    <h2 data-zh="dequeue_op事件分析" data-en="dequeue_op Events Analysis">dequeue_op事件分析</h2>
    <div class="layout">
        <div class="left-panel">
            <div class="query-principle">
                <h3 data-zh="查询原理" data-en="Query Principle">查询原理</h3>
                <p data-zh="此分析解析 Ceph 日志中的出队事件。它识别：" data-en="This analysis parses dequeue events from the Ceph log. It identifies:">此分析解析 Ceph 日志中的出队事件。它识别：</p>
                <ul>
                    <li data-zh="<strong>出队事件</strong>：包含 "dequeue_op" 的日志行" data-en="<strong>Dequeue events</strong>: Log lines containing "dequeue_op"">出队事件：包含 "dequeue_op" 的日志行</li>
                </ul>
                <p data-zh="对于每个出队事件，它提取：" data-en="For each dequeue event, it extracts:">对于每个出队事件，它提取：</p>
                <ul>
                    <li data-zh="时间戳" data-en="Timestamp">时间戳</li>
                    <li data-zh="操作类型" data-en="Operation type">操作类型</li>
                    <li data-zh="客户端 ID" data-en="Client ID">客户端 ID</li>
                    <li data-zh="操作 ID" data-en="Operation ID">操作 ID</li>
                    <li data-zh="PG ID" data-en="PG ID">PG ID</li>
                    <li data-zh="出队延迟" data-en="Dequeue latency">出队延迟</li>
                    <li data-zh="优先级" data-en="Priority">优先级</li>
                    <li data-zh="成本" data-en="Cost">成本</li>
                </ul>
                <p data-zh="它分析事件从入队到出队的等待时间，帮助识别队列积压和性能瓶颈。" data-en="It analyzes the waiting time from enqueue to dequeue, helping to identify queue backlogs and performance bottlenecks.">它分析事件从入队到出队的等待时间，帮助识别队列积压和性能瓶颈。</p>
            </div>
            <div class="summary">
                <h3 data-zh="摘要" data-en="Summary">摘要</h3>
                <p data-zh="总出队事件数：" data-en="Total dequeue events: ">总出队事件数：` + strconv.Itoa(result.TotalOps) + `</p>`

	if result.TotalOps > 0 {
		avgLatency := float64(result.TotalLatency) / float64(result.TotalOps)
		html += fmt.Sprintf(`
                <p data-zh="平均出队延迟：" data-en="Average dequeue latency: ">平均出队延迟：%.6f ms</p>
                <p data-zh="最大出队延迟：" data-en="Maximum dequeue latency: ">最大出队延迟：%.6f ms</p>
                <p data-zh="最小出队延迟：" data-en="Minimum dequeue latency: ">最小出队延迟：%.6f ms</p>`,
			avgLatency / 1000, // 转换为毫秒
			float64(result.MaxLatency) / 1000, // 转换为毫秒
			float64(result.MinLatency) / 1000) // 转换为毫秒
	}

	// Add latency distribution
	html += `
                <h4 data-zh="延迟分布：" data-en="Latency Distribution:">延迟分布：</h4>`
	latencyRanges := []string{"0-10μs", "10-20μs", "20-30μs", "30-40μs", "40-50μs", "50-60μs", "60-70μs", "70-80μs", "80-90μs", "90-100μs", "100μs+"}
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
                <input type="datetime-local" id="dequeue-start-time">
                <label data-zh="结束时间：" data-en="End Time:">结束时间：</label>
                <input type="datetime-local" id="dequeue-end-time">
                <label data-zh="最小延迟 (ms)：" data-en="Min Latency (ms):">最小延迟 (ms)：</label>
                <input type="number" id="dequeue-min-latency" min="0" step="0.001">
                <label data-zh="最大延迟 (ms)：" data-en="Max Latency (ms):">最大延迟 (ms)：</label>
                <input type="number" id="dequeue-max-latency" min="0" step="0.001">
                <label data-zh="操作类型：" data-en="Operation Type:">操作类型：</label>
                <select id="dequeue-op-type">
                    <option value="" data-zh="全部" data-en="All">全部</option>
                    <option value="osd_op" data-zh="客户端操作" data-en="Client operation">客户端操作</option>
                    <option value="osd_repop" data-zh="复制事务" data-en="Replication transaction">复制事务</option>
                    <option value="osd_repop_reply" data-zh="复制回复" data-en="Replication reply">复制回复</option>
                    <option value="pg_update_log_missing" data-zh="PG更新日志缺失" data-en="PG update log missing">PG更新日志缺失</option>
                </select>
                <button type="button" onclick="window.filterDequeueTable()" data-zh="筛选" data-en="Filter">筛选</button>
                <button type="button" onclick="window.resetDequeueFilter()" data-zh="重置" data-en="Reset">重置</button>
            </div>
            <div class="table-container">
            <table id="dequeue-table">
                <tr>
                    <th data-zh="时间戳" data-en="Timestamp">时间戳</th>
                    <th data-zh="操作类型" data-en="Operation Type">操作类型</th>
                    <th data-zh="客户端 ID" data-en="Client ID">客户端 ID</th>
                    <th data-zh="操作 ID" data-en="Op ID">操作 ID</th>
                    <th data-zh="PG ID" data-en="PG ID">PG ID</th>
                    <th data-zh="出队延迟 (ms)" data-en="Dequeue Latency (ms)">出队延迟 (ms)</th>
                    <th data-zh="优先级" data-en="Priority">优先级</th>
                    <th data-zh="成本" data-en="Cost">成本</th>
                </tr>`

	for _, event := range result.Events {
		html += fmt.Sprintf(`
                <tr>
                    <td>%s</td>
                    <td>%s</td>
                    <td>%s</td>
                    <td>%s</td>
                    <td>%s</td>
                    <td>%.6f</td>
                    <td>%d</td>
                    <td>%d</td>
                </tr>`,
			event.Timestamp.Format("2006-01-02 15:04:05.000"),
			event.OpType,
			event.ClientID,
			event.OpID,
			event.PGID,
			float64(event.DequeueLatency) / 1000, // 转换为毫秒
			event.Priority,
			event.Cost)
	}

	html += `
            </table>
            </div>
            <div class="pagination">
                <div class="pagination-info" id="dequeue-table-pagination-info"></div>
                <div class="pagination-buttons" id="dequeue-table-pagination-buttons"></div>
            </div>
        </div>
    </div>`

	return html
}