# Ceph Log Analysis Tool

A comprehensive tool for analyzing Ceph OSD logs, including AIO operations, OSD repop operations, and OSD operations.

## Features

- **Multi-type Analysis**: Analyzes three types of operations:
  - AIO (Asynchronous I/O) operations
  - OSD repop operations
  - OSD operations
- **HTML Output**: Generates a well-structured HTML report with tabbed interface
- **Summary Statistics**: Displays comprehensive statistics for each operation type
- **Filtering Capabilities**: Allows filtering by time range and duration/latency
- **Duration Distribution**: Shows distribution of operations by duration/latency ranges
- **Legacy Scripts**: Includes legacy analysis scripts in the `old/` directory

## Installation

### Prerequisites
- Go 1.16 or later

### Setup
1. Clone the repository:
   ```bash
   git clone https://github.com/liucxer/analyze_ceph_log.git
   cd analyze_ceph_log
   ```

2. Build the project:
   ```bash
   go build -o analyze_ceph .
   ```

## Usage

### Command Syntax
```bash
# Using the built binary
./analyze_ceph <log_file> <analysis_type> [output.html]

# Using go run
go run analyze_ceph.go <log_file> <analysis_type> [output.html]
```

### Analysis Types
- `aio` - Analyze AIO operations
- `repop` - Analyze OSD repop operations
- `op` - Analyze OSD operations
- `all` - Analyze all operation types

### Examples

1. Analyze all operations and generate HTML report:
   ```bash
   go run analyze_ceph.go /path/to/ceph-osd.log all analysis.html
   ```

2. Analyze only AIO operations:
   ```bash
   go run analyze_ceph.go /path/to/ceph-osd.log aio aio_analysis.html
   ```

3. Analyze OSD operations:
   ```bash
   go run analyze_ceph.go /path/to/ceph-osd.log op osd_analysis.html
   ```

## Output

The tool generates an HTML report with:

- **Tabbed Interface**: Separate tabs for each operation type
- **Summary Section**: Key statistics at the top of each tab
- **Filter Form**: Time range and duration/latency filters for each table
- **Data Table**: Detailed list of operations with timestamps and durations
- **Duration Distribution**: Breakdown of operations by duration/latency ranges

## Project Structure

```
analyze_ceph_log/
├── analyze_ceph.go       # Main analysis tool
├── go.mod                # Go module file
├── .gitignore            # Git ignore file
├── README.md             # English README
├── README_zh.md          # Chinese README
└── old/                  # Legacy analysis scripts
    ├── analyze_aio.go
    ├── analyze_osd_op.go
    └── analyze_osd_repop.go
```

## How It Works

1. **Log Parsing**: Uses regular expressions to extract relevant information from Ceph logs
2. **Event Matching**: Matches start and end events to calculate durations
3. **Data Analysis**: Computes statistics and duration distributions
4. **HTML Generation**: Creates a well-structured HTML report with filtering capabilities

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open-source and available under the MIT License.
