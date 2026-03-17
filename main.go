package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/liucxer/analyze_ceph_log/pkg/analyzer"
	"github.com/liucxer/analyze_ceph_log/pkg/html"
	"github.com/liucxer/analyze_ceph_log/pkg/log"
	"github.com/liucxer/analyze_ceph_log/pkg/types"
)

func main() {
	// Check if a log file is provided
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <log_file> <analysis_type> [output.html]")
		fmt.Println("Analysis types:")
		fmt.Println("  aio - Analyze AIO operations")
		fmt.Println("  repop - Analyze OSD repop operations")
		fmt.Println("  op - Analyze OSD operations")
		fmt.Println("  transaction - Analyze transaction operations")
		fmt.Println("  metadata - Analyze metadata sync operations")
		fmt.Println("  client - Analyze client operations")
		fmt.Println("  dequeue - Analyze dequeue operations")
		fmt.Println("  all - Analyze all operation types")
		os.Exit(1)
	}

	logFile := os.Args[1]
	analysisType := os.Args[2]

	// Determine output file
	outputFile := "analysis.html"
	if len(os.Args) > 3 {
		outputFile = os.Args[3]
	}

	// Open the log file
	file, err := os.Open(logFile)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Analyze based on type
	switch analysisType {
	case "aio":
		analyzeAIO(file, outputFile)
	case "repop":
		analyzeRepop(file, outputFile)
	case "op":
		analyzeOSDOp(file, outputFile)
	case "transaction":
		analyzeTransaction(file, outputFile)
	case "metadata":
		analyzeMetadataSync(file, outputFile)
	case "client":
		analyzeClientOp(file, outputFile)
	case "dequeue":
		analyzeDequeue(file, outputFile)
	case "all":
		analyzeAll(file, outputFile)
	default:
		fmt.Printf("Unknown analysis type: %s\n", analysisType)
		os.Exit(1)
	}

	fmt.Printf("Analysis completed. Results saved to %s\n", outputFile)
}

// analyzeAIO analyzes AIO operations
func analyzeAIO(file *os.File, outputFile string) {
	// Reset file pointer
	file.Seek(0, 0)

	// Parse AIO events
	scanner := bufio.NewScanner(file)
	// Increase scanner buffer size for large log lines
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	events, err := log.ParseAIOEvents(scanner)
	if err != nil {
		fmt.Printf("Error parsing AIO events: %v\n", err)
		os.Exit(1)
	}

	// Analyze events
	result := analyzer.AnalyzeAIOEvents(events)

	// Generate HTML
	htmlContent := html.GenerateHTML(result, types.AnalysisResult{}, types.OSDOpAnalysisResult{}, types.TransactionAnalysisResult{}, types.MetadataSyncAnalysisResult{}, types.ClientOpAnalysisResult{}, types.DequeueAnalysisResult{}, types.AnalysisResult{})

	// Write to file
	err = os.WriteFile(outputFile, []byte(htmlContent), 0644)
	if err != nil {
		fmt.Printf("Error writing HTML file: %v\n", err)
		os.Exit(1)
	}
}

// analyzeRepop analyzes OSD repop operations
func analyzeRepop(file *os.File, outputFile string) {
	// Reset file pointer
	file.Seek(0, 0)

	// Parse Repop events
	scanner := bufio.NewScanner(file)
	// Increase scanner buffer size for large log lines
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	events, err := log.ParseRepopEvents(scanner)
	if err != nil {
		fmt.Printf("Error parsing Repop events: %v\n", err)
		os.Exit(1)
	}

	// Analyze events
	result := analyzer.AnalyzeRepopEvents(events)

	// Generate HTML
	htmlContent := html.GenerateHTML(types.AnalysisResult{}, result, types.OSDOpAnalysisResult{}, types.TransactionAnalysisResult{}, types.MetadataSyncAnalysisResult{}, types.ClientOpAnalysisResult{}, types.DequeueAnalysisResult{}, types.AnalysisResult{})

	// Write to file
	err = os.WriteFile(outputFile, []byte(htmlContent), 0644)
	if err != nil {
		fmt.Printf("Error writing HTML file: %v\n", err)
		os.Exit(1)
	}
}

// analyzeOSDOp analyzes OSD operations
func analyzeOSDOp(file *os.File, outputFile string) {
	// Reset file pointer
	file.Seek(0, 0)

	// Parse OSD Op events
	scanner := bufio.NewScanner(file)
	// Increase scanner buffer size for large log lines
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	events, err := log.ParseOSDOpEvents(scanner)
	if err != nil {
		fmt.Printf("Error parsing OSD Op events: %v\n", err)
		os.Exit(1)
	}

	// Analyze events
	result := analyzer.AnalyzeOSDOpEvents(events)

	// Generate HTML
	htmlContent := html.GenerateHTML(types.AnalysisResult{}, types.AnalysisResult{}, result, types.TransactionAnalysisResult{}, types.MetadataSyncAnalysisResult{}, types.ClientOpAnalysisResult{}, types.DequeueAnalysisResult{}, types.AnalysisResult{})

	// Write to file
	err = os.WriteFile(outputFile, []byte(htmlContent), 0644)
	if err != nil {
		fmt.Printf("Error writing HTML file: %v\n", err)
		os.Exit(1)
	}
}

// analyzeTransaction analyzes transaction operations
func analyzeTransaction(file *os.File, outputFile string) {
	// Reset file pointer
	file.Seek(0, 0)

	// Parse Transaction events
	scanner := bufio.NewScanner(file)
	// Increase scanner buffer size for large log lines
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	events, err := log.ParseTransactionEvents(scanner)
	if err != nil {
		fmt.Printf("Error parsing Transaction events: %v\n", err)
		os.Exit(1)
	}

	// Analyze events
	result := analyzer.AnalyzeTransactionEvents(events)

	// Generate HTML
	htmlContent := html.GenerateHTML(types.AnalysisResult{}, types.AnalysisResult{}, types.OSDOpAnalysisResult{}, result, types.MetadataSyncAnalysisResult{}, types.ClientOpAnalysisResult{}, types.DequeueAnalysisResult{}, types.AnalysisResult{})

	// Write to file
	err = os.WriteFile(outputFile, []byte(htmlContent), 0644)
	if err != nil {
		fmt.Printf("Error writing HTML file: %v\n", err)
		os.Exit(1)
	}
}

// analyzeMetadataSync analyzes metadata sync operations
func analyzeMetadataSync(file *os.File, outputFile string) {
	// Reset file pointer
	file.Seek(0, 0)

	// Parse Metadata Sync events
	scanner := bufio.NewScanner(file)
	// Increase scanner buffer size for large log lines
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	events, err := log.ParseMetadataSyncEvents(scanner)
	if err != nil {
		fmt.Printf("Error parsing Metadata Sync events: %v\n", err)
		os.Exit(1)
	}

	// Analyze events
	result := analyzer.AnalyzeMetadataSyncEvents(events)

	// Generate HTML
	htmlContent := html.GenerateHTML(types.AnalysisResult{}, types.AnalysisResult{}, types.OSDOpAnalysisResult{}, types.TransactionAnalysisResult{}, result, types.ClientOpAnalysisResult{}, types.DequeueAnalysisResult{}, types.AnalysisResult{})

	// Write to file
	err = os.WriteFile(outputFile, []byte(htmlContent), 0644)
	if err != nil {
		fmt.Printf("Error writing HTML file: %v\n", err)
		os.Exit(1)
	}
}

// analyzeClientOp analyzes client operations
func analyzeClientOp(file *os.File, outputFile string) {
	// Reset file pointer
	file.Seek(0, 0)

	// Parse Client Op events
	scanner := bufio.NewScanner(file)
	// Increase scanner buffer size for large log lines
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	events, err := log.ParseClientOpEvents(scanner)
	if err != nil {
		fmt.Printf("Error parsing Client Op events: %v\n", err)
		os.Exit(1)
	}

	// Analyze events
	result := analyzer.AnalyzeClientOpEvents(events)

	// Generate HTML
	htmlContent := html.GenerateHTML(types.AnalysisResult{}, types.AnalysisResult{}, types.OSDOpAnalysisResult{}, types.TransactionAnalysisResult{}, types.MetadataSyncAnalysisResult{}, result, types.DequeueAnalysisResult{}, types.AnalysisResult{})

	// Write to file
	err = os.WriteFile(outputFile, []byte(htmlContent), 0644)
	if err != nil {
		fmt.Printf("Error writing HTML file: %v\n", err)
		os.Exit(1)
	}
}

// analyzeDequeue analyzes dequeue operations
func analyzeDequeue(file *os.File, outputFile string) {
	// Reset file pointer
	file.Seek(0, 0)

	// Parse Dequeue events
	scanner := bufio.NewScanner(file)
	// Increase scanner buffer size for large log lines
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	events, err := log.ParseDequeueEvents(scanner)
	if err != nil {
		fmt.Printf("Error parsing Dequeue events: %v\n", err)
		os.Exit(1)
	}

	// Analyze events
	result := analyzer.AnalyzeDequeueEvents(events)

	// Generate HTML
	htmlContent := html.GenerateHTML(types.AnalysisResult{}, types.AnalysisResult{}, types.OSDOpAnalysisResult{}, types.TransactionAnalysisResult{}, types.MetadataSyncAnalysisResult{}, types.ClientOpAnalysisResult{}, result, types.AnalysisResult{})

	// Write to file
	err = os.WriteFile(outputFile, []byte(htmlContent), 0644)
	if err != nil {
		fmt.Printf("Error writing HTML file: %v\n", err)
		os.Exit(1)
	}
}

// analyzeAll analyzes all operation types
func analyzeAll(file *os.File, outputFile string) {
	// Parse AIO events
	file.Seek(0, 0)
	scanner := bufio.NewScanner(file)
	// Increase scanner buffer size for large log lines
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	aioEvents, err := log.ParseAIOEvents(scanner)
	if err != nil {
		fmt.Printf("Error parsing AIO events: %v\n", err)
		os.Exit(1)
	}

	// Parse Repop events
	file.Seek(0, 0)
	scanner = bufio.NewScanner(file)
	// Increase scanner buffer size for large log lines
	buf = make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	repopEvents, err := log.ParseRepopEvents(scanner)
	if err != nil {
		fmt.Printf("Error parsing Repop events: %v\n", err)
		os.Exit(1)
	}

	// Parse OSD Op events
	file.Seek(0, 0)
	scanner = bufio.NewScanner(file)
	// Increase scanner buffer size for large log lines
	buf = make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	osdOpEvents, err := log.ParseOSDOpEvents(scanner)
	if err != nil {
		fmt.Printf("Error parsing OSD Op events: %v\n", err)
		os.Exit(1)
	}

	// Parse OSD Op V2 events (enqueue_op to sending reply on)
	file.Seek(0, 0)
	scanner = bufio.NewScanner(file)
	// Increase scanner buffer size for large log lines
	buf = make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	osdOpEventsV2, err := log.ParseOSDOpEventsV2(scanner)
	if err != nil {
		fmt.Printf("Error parsing OSD Op V2 events: %v\n", err)
		os.Exit(1)
	}

	// Parse Transaction events
	file.Seek(0, 0)
	scanner = bufio.NewScanner(file)
	// Increase scanner buffer size for large log lines
	buf = make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	transactionEvents, err := log.ParseTransactionEvents(scanner)
	if err != nil {
		fmt.Printf("Error parsing Transaction events: %v\n", err)
		os.Exit(1)
	}

	// Parse Metadata Sync events
	file.Seek(0, 0)
	scanner = bufio.NewScanner(file)
	// Increase scanner buffer size for large log lines
	buf = make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	metadataSyncEvents, err := log.ParseMetadataSyncEvents(scanner)
	if err != nil {
		fmt.Printf("Error parsing Metadata Sync events: %v\n", err)
		os.Exit(1)
	}

	// Parse Client Op events
	file.Seek(0, 0)
	scanner = bufio.NewScanner(file)
	// Increase scanner buffer size for large log lines
	buf = make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	clientOpEvents, err := log.ParseClientOpEvents(scanner)
	if err != nil {
		fmt.Printf("Error parsing Client Op events: %v\n", err)
		os.Exit(1)
	}

	// Parse Dequeue events
	file.Seek(0, 0)
	scanner = bufio.NewScanner(file)
	// Increase scanner buffer size for large log lines
	buf = make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	dequeueEvents, err := log.ParseDequeueEvents(scanner)
	if err != nil {
		fmt.Printf("Error parsing Dequeue events: %v\n", err)
		os.Exit(1)
	}

	// Analyze events
	aioResult := analyzer.AnalyzeAIOEvents(aioEvents)
	repopResult := analyzer.AnalyzeRepopEvents(repopEvents)
	osdOpResult := analyzer.AnalyzeOSDOpEvents(osdOpEvents)
	osdOpResultV2 := analyzer.AnalyzeRepopEvents(osdOpEventsV2) // Reuse the same analyzer function
	transactionResult := analyzer.AnalyzeTransactionEvents(transactionEvents)
	metadataSyncResult := analyzer.AnalyzeMetadataSyncEvents(metadataSyncEvents)
	clientResult := analyzer.AnalyzeClientOpEvents(clientOpEvents)
	dequeueResult := analyzer.AnalyzeDequeueEvents(dequeueEvents)

	// Generate HTML
	htmlContent := html.GenerateHTML(aioResult, repopResult, osdOpResult, transactionResult, metadataSyncResult, clientResult, dequeueResult, osdOpResultV2)

	// Write to file
	err = os.WriteFile(outputFile, []byte(htmlContent), 0644)
	if err != nil {
		fmt.Printf("Error writing HTML file: %v\n", err)
		os.Exit(1)
	}
}
