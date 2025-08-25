# Network Logging CLI Tool

## Overview

This is a **Go CLI tool** for logging network information. It allows you to:

1. Collect **local host information** (hostname, OS, IP addresses, username).
2. Scan the **local network** to gather basic info for other reachable devices (best-effort).
3. Save all collected data in a **single JSON file**.

---

## Requirements

* Go installed on your system (version 1.20+ recommended).
* Permission to ping devices on your local network.

---

## Usage

### 1. Run with default settings

```bash
go run main.go
```

* The JSON log will be saved in the `logs/` folder.
* File name format: `reports_<timestamp>.json` (e.g., `reports_2025-08-25_10-30-12.json`).

---

### 2. Run with custom filename

```bash
go run main.go --out custom_file_name.json
```

* The file will be created in the `logs/` folder with the specified name.

---

## Implementation Details

* The **network scanning** involves I/O-bound operations like pinging devices.

* To speed up the scan, the tool uses **concurrency with a worker pool**, ensuring:

  * Faster lookups across multiple devices
  * No race conditions during scanning or logging

* The tool also shows a **progress indicator** while scanning the network.

---

## Folder Structure

```
.
├── main.go            # Main CLI program
├── logs/              # Folder to store generated JSON reports
└── README.md          # Project documentation
```

---

## Notes

* Make sure the `logs/` folder exists; the program will not track log files in Git.
* All scan results are saved as JSON for easy parsing or integration with other tools.
