# MLeak

## A memory leak detection tool built with Cobra.

[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

MLeak is a lightweight CLI tool designed to monitor and report memory usage in native applications, with built-in heuristics to detect potential memory leaks. It provides developers with insights into memory consumption patterns and growth rates.

## Features

### Core Monitoring
- Monitor memory usage of specified processes by PID
- Real-time RSS (Resident Set Size) tracking with configurable intervals
- Display real-time memory statistics with customizable intervals
- Support for multiple output formats (table, JSON, plain text)
- Cross-platform support (Linux, Windows)

### Leak Detection
- Trend analysis of memory usage over time
- Growth rate calculation to identify abnormal memory increases
- Heuristic-based leak risk classification (Low/Medium/High)
- Stability detection to differentiate between normal allocations and leaks

### Planned Features

#### Core Monitoring
- [ ] VSS (Virtual Set Size) and other memory metrics tracking
- [ ] Time-based leak thresholds (e.g., --threshold=5MB)
- [ ] macOS support
- [ ] CSV export functionality

#### Advanced Detection
- [ ] Allocation hooking (malloc/free interception)
- [ ] Stack trace capture at allocation time
- [ ] Debug symbol resolution
- [ ] String leak patterns detection
- [ ] Image buffer leak detection
- [ ] Historical data storage and visualization
- [ ] Memory snapshot comparison functionality

#### Alerting and Reporting
- [ ] Alert notifications for high leak risk
- [ ] Automated memory leak diagnosis with suggestions
- [ ] Export reports to PDF formats

#### Enterprise Features
- [ ] Kubernetes/container monitoring
- [ ] Grafana/Prometheus integration
- [ ] CI/CD plugin system
- [ ] Process filtering capabilities
- [ ] Automated remediation hints

## Installation

```bash
# From source (requires Go 1.24+)
git clone https://github.com/Feinq/mleak
cd mleak
go build -o mleak
sudo install mleak /usr/local/bin/
```

## Usage

Monitor a process:
```bash
mleak monitor <PID>
```

Options:
```
  -f, --format string     Output format (table, text, json) (default "table")
  -i, --interval duration Monitoring interval (default 10s)
```

Sample outputs:

Table format (default):
```
[MLeak] Monitoring PID 12345 (chrome) | Interval: 10s
┌────────────┬───────────┬───────────┐
│ Timestamp  │ RSS       │ Leak Risk │
├────────────┼───────────┼───────────┤
│ 15:04:05   │ 45.20 MB  │ Low       │
│ 15:04:15   │ 45.75 MB  │ Low       │ +0MB/10s
│ 15:04:25   │ 47.82 MB  │ Medium    │ ! +2MB/10s
└────────────┴───────────┴───────────┘
```

JSON format:
```json
{"timestamp":"2023-07-15T15:04:05Z","rss":47406080,"leak_risk":"Medium"}
```

Text format:
```
Timestamp: 2023-07-15T15:04:05Z
RSS: 45.20 MB
Leak Risk: Low
```

## Contributing
Contributions, issues, and feature requests are welcome. Feel free to clone the project and submit a pull request (PR) for any changes. Once your PR is accepted, it will be deployed.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.