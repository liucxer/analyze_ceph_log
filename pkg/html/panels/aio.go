package panels

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/liucxer/analyze_ceph_log/pkg/types"
)

// GenerateAIOHTML generates HTML for AIO analysis
func GenerateAIOHTML(result types.AnalysisResult) string {
	html := `
    <h2 data-zh="AIO 事件分析" data-en="AIO Events Analysis">AIO 事件分析</h2>
    <div class="layout">
        <div class="left-panel">
            <div class="query-principle">
                <h3 data-zh="查询原理" data-en="Query Principle">查询原理</h3>
                <p data-zh="此分析解析 Ceph 日志中的 AIO (异步 I/O) 事件。它识别：" data-en="This analysis parses AIO (Asynchronous I/O) events from the Ceph log. It identifies:">此分析解析 Ceph 日志中的 AIO (异步 I/O) 事件。它识别：</p>
                <ul>
                    <li data-zh="<strong>开始事件</strong>：包含 "_aio_log_start" 的日志行" data-en="<strong>Start events</strong>: Log lines containing "_aio_log_start"">开始事件：包含 "_aio_log_start" 的日志行</li>
                    <li data-zh="<strong>结束事件</strong>：包含 "_aio_log_finish" 的日志行" data-en="<strong>Finish events</strong>: Log lines containing "_aio_log_finish"">结束事件：包含 "_aio_log_finish" 的日志行</li>
                </ul>
                <p data-zh="对于每个 AIO 事件，它提取：" data-en="For each AIO event, it extracts:">对于每个 AIO 事件，它提取：</p>
                <ul>
                    <li data-zh="时间戳" data-en="Timestamp">时间戳</li>
                    <li data-zh="块地址范围" data-en="Block address range">块地址范围</li>
                    <li data-zh="数据长度（从十六进制转换为十进制）" data-en="Data length (converted from hex to decimal)">数据长度（从十六进制转换为十进制）</li>
                    <li data-zh="块类型（block, block.wal, block.db）" data-en="Block type (block, block.wal, block.db)">块类型（block, block.wal, block.db）</li>
                </ul>
                <p data-zh="它使用块地址范围作为唯一键匹配开始和结束事件，然后计算每个事件的持续时间。" data-en="It matches start and finish events using the block address range as a unique key, then calculates the duration for each event.">它使用块地址范围作为唯一键匹配开始和结束事件，然后计算每个事件的持续时间。</p>
            </div>
            <div class="summary">
                <h3 data-zh="摘要" data-en="Summary">摘要</h3>
                <p data-zh="总 AIO 事件数：" data-en="Total AIO events: ">总 AIO 事件数：` + strconv.Itoa(result.TotalEvents) + `</p>`

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
            <div class="pagination">
                <div class="pagination-info" id="aio-table-pagination-info"></div>
                <div class="pagination-buttons" id="aio-table-pagination-buttons"></div>
            </div>
        </div>
    </div>`

	return html
}
