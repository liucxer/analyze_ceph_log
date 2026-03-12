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
	events, err := log.ParseAIOEvents(scanner)
	if err != nil {
		fmt.Printf("Error parsing AIO events: %v\n", err)
		os.Exit(1)
	}

	// Analyze events
	result := analyzer.AnalyzeAIOEvents(events)

	// Generate HTML
	htmlContent := html.GenerateHTML(result, types.AnalysisResult{}, types.OSDOpAnalysisResult{})

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
	events, err := log.ParseRepopEvents(scanner)
	if err != nil {
		fmt.Printf("Error parsing Repop events: %v\n", err)
		os.Exit(1)
	}

	// Analyze events
	result := analyzer.AnalyzeRepopEvents(events)

	// Generate HTML
	htmlContent := html.GenerateHTML(types.AnalysisResult{}, result, types.OSDOpAnalysisResult{})

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
	events, err := log.ParseOSDOpEvents(scanner)
	if err != nil {
		fmt.Printf("Error parsing OSD Op events: %v\n", err)
		os.Exit(1)
	}

	// Analyze events
	result := analyzer.AnalyzeOSDOpEvents(events)

	// Generate HTML
	htmlContent := html.GenerateHTML(types.AnalysisResult{}, types.AnalysisResult{}, result)

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
	aioEvents, err := log.ParseAIOEvents(scanner)
	if err != nil {
		fmt.Printf("Error parsing AIO events: %v\n", err)
		os.Exit(1)
	}

	// Parse Repop events
	file.Seek(0, 0)
	scanner = bufio.NewScanner(file)
	repopEvents, err := log.ParseRepopEvents(scanner)
	if err != nil {
		fmt.Printf("Error parsing Repop events: %v\n", err)
		os.Exit(1)
	}

	// Parse OSD Op events
	file.Seek(0, 0)
	scanner = bufio.NewScanner(file)
	osdOpEvents, err := log.ParseOSDOpEvents(scanner)
	if err != nil {
		fmt.Printf("Error parsing OSD Op events: %v\n", err)
		os.Exit(1)
	}

	// Analyze events
	aioResult := analyzer.AnalyzeAIOEvents(aioEvents)
	repopResult := analyzer.AnalyzeRepopEvents(repopEvents)
	osdOpResult := analyzer.AnalyzeOSDOpEvents(osdOpEvents)

	// Generate HTML
	htmlContent := html.GenerateHTML(aioResult, repopResult, osdOpResult)

	// Write to file
	err = os.WriteFile(outputFile, []byte(htmlContent), 0644)
	if err != nil {
		fmt.Printf("Error writing HTML file: %v\n", err)
		os.Exit(1)
	}
}
