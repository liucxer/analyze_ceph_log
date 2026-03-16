package types

import (
	"time"
)

// AIOEvent represents an AIO start or finish event
type AIOEvent struct {
	Timestamp time.Time
	BlockAddr string
	RangeStr  string
	Length    int
	Duration  time.Duration
}

// RepopEvent represents an OSD repop event
type RepopEvent struct {
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
	OpID      string
}

// OSDOpEvent represents an OSD operation event
type OSDOpEvent struct {
	Timestamp         time.Time
	OpID              string
	PgID              string
	Object            string
	OpType            string
	RangeStr          string
	InBytes           int
	OutBytes          int
	Latency           float64 // in milliseconds
	FirstRepopReply   float64 // in milliseconds
	SecondRepopReply  float64 // in milliseconds
}

// Event represents a generic event with start/end times and duration
type Event struct {
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
	RangeStr  string
	Length    int
	BlockType string
	OpID      string
}

// AnalysisResult represents the result of an analysis
type AnalysisResult struct {
	Events      []Event
	TotalEvents int
	TotalDuration time.Duration
	MaxDuration  time.Duration
	MinDuration  time.Duration
	DurationCounts map[int]int
}

// TransactionEvent represents a transaction event with different stages
type TransactionEvent struct {
	TID              string
	StartTime        time.Time
	IssueTime        time.Time
	FirstReplyTime   time.Time
	SecondReplyTime  time.Time
	CompleteTime     time.Time
	TotalDuration    time.Duration
	IssueDuration    time.Duration
	FirstReplyDuration time.Duration
	SecondReplyDuration time.Duration
	OpID             string
	Object           string
	RangeStr         string
}

// TransactionAnalysisResult represents the result of transaction analysis
type TransactionAnalysisResult struct {
	Events            []TransactionEvent
	TotalTransactions int
	TotalDuration     time.Duration
	MaxDuration       time.Duration
	MinDuration       time.Duration
	DurationCounts    map[int]int
}

// OSDOpAnalysisResult represents the result of OSD operation analysis
type OSDOpAnalysisResult struct {
	Events       []OSDOpEvent
	TotalOps     int
	TotalLatency float64
	MaxLatency   float64
	TotalInBytes int
	TotalOutBytes int
	LatencyCounts map[string]int
}

// MetadataSyncEvent represents a metadata sync event
type MetadataSyncEvent struct {
	Timestamp time.Time
	Committed int
	Cleaned   int
	Duration  time.Duration
	FlushTime time.Duration
	KVCommitTime time.Duration
}

// MetadataSyncAnalysisResult represents the result of metadata sync analysis
type MetadataSyncAnalysisResult struct {
	Events      []MetadataSyncEvent
	TotalEvents int
	TotalDuration time.Duration
	MaxDuration  time.Duration
	MinDuration  time.Duration
	DurationCounts map[int]int
}
