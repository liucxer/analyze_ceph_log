package panels

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/liucxer/analyze_ceph_log/pkg/types"
)

// GenerateMetadataHTML generates HTML for Metadata Sync analysis
func GenerateMetadataHTML(result types.MetadataSyncAnalysisResult) string {
	html := `
    <h2 data-zh="元数据同步" data-en="Metadata Sync">元数据同步</h2>
    <div class="layout">
        <div class="left-panel">
            <div class="query-principle">
                <h3 data-zh="查询原理" data-en="Query Principle">查询原理</h3>
                <p data-zh="此分析解析 Ceph 日志中的元数据同步操作。它识别：" data-en="This analysis parses metadata sync operations from the Ceph log. It identifies:">此分析解析 Ceph 日志中的元数据同步操作。它识别：</p>
                <ul>
                    <li data-zh="<strong>元数据同步事件</strong>：包含 "_kv_sync_thread committed" 的日志行" data-en="<strong>Metadata sync events</strong>: Log lines containing "_kv_sync_thread committed"">元数据同步事件：包含 "_kv_sync_thread committed" 的日志行</li>
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
            <div class="pagination">
                <div class="pagination-info" id="metadata-table-pagination-info"></div>
                <div class="pagination-buttons" id="metadata-table-pagination-buttons"></div>
            </div>
        </div>
    </div>`

	return html
}
