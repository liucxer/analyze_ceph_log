package panels

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/liucxer/analyze_ceph_log/pkg/types"
)

// GenerateTransactionHTML generates HTML for Transaction analysis
func GenerateTransactionHTML(result types.TransactionAnalysisResult) string {
	html := `
    <h2 data-zh="复制事务" data-en="Replication Transaction">复制事务</h2>
    <div class="layout">
        <div class="left-panel">
            <div class="query-principle">
                <h3 data-zh="查询原理" data-en="Query Principle">查询原理</h3>
                <p data-zh="此分析解析 Ceph 日志中的复制事务操作。它识别复制事务的不同阶段：" data-en="This analysis parses replication transaction operations from the Ceph log. It identifies different stages of a replication transaction:">此分析解析 Ceph 日志中的复制事务操作。它识别复制事务的不同阶段：</p>
                <ul>
                    <li data-zh="<strong>开始事件</strong>：包含 "new_repop" 和 "rep_tid" 的日志行" data-en="<strong>Start events</strong>: Log lines containing "new_repop" with "rep_tid"">开始事件：包含 "new_repop" 和 "rep_tid" 的日志行</li>
                    <li data-zh="<strong>发出事件</strong>：包含 "issue_repop" 和 "rep_tid" 的日志行" data-en="<strong>Issue events</strong>: Log lines containing "issue_repop" with "rep_tid"">发出事件：包含 "issue_repop" 和 "rep_tid" 的日志行</li>
                    <li data-zh="<strong>回复事件</strong>：包含 "do_repop_reply" 和 "tid" 的日志行" data-en="<strong>Reply events</strong>: Log lines containing "do_repop_reply" with "tid"">回复事件：包含 "do_repop_reply" 和 "tid" 的日志行</li>
                    <li data-zh="<strong>完成事件</strong>：包含 "repop_all_committed" 和 "repop tid" 的日志行" data-en="<strong>Complete events</strong>: Log lines containing "repop_all_committed" with "repop tid"">完成事件：包含 "repop_all_committed" 和 "repop tid" 的日志行</li>
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
            <div class="pagination">
                <div class="pagination-info" id="transaction-table-pagination-info"></div>
                <div class="pagination-buttons" id="transaction-table-pagination-buttons"></div>
            </div>
        </div>
    </div>`

	return html
}
