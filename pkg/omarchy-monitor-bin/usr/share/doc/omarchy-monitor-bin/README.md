<div align="center">

# 🖥️ Omarchy System Monitor (OSM)

[![Go Version](https://img.shields.io/badge/Go-1.25.4-00ADD8?style=for-the-badge&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Linux-orange?style=for-the-badge&logo=linux)](https://www.linux.org/)

**A beautiful, lightweight terminal-based system monitor built with Go**

*Real-time system metrics visualization with multiple themes and interactive controls*

[Features](#-features) • [Installation](#-installation) • [Usage](#-usage) • [Keybindings](#%EF%B8%8F-keybindings) • [Themes](#-themes)

</div>

---

## 📖 Overview

Omarchy System Monitor (OSM) is a powerful yet lightweight terminal user interface (TUI) application for monitoring system resources in real-time. Built with Go and leveraging the excellent `tview` and `tcell` libraries, OSM provides an intuitive, feature-rich alternative to traditional system monitoring tools like `htop` and `top`.

## ✨ Features

### 📊 **Comprehensive Monitoring**
- **Process Management** - View all running processes with CPU, memory usage, and user information
- **Disk Usage** - Monitor disk I/O statistics and partition usage
- **Network Statistics** - Track network interface traffic in real-time
- **Temperature Sensors** - Display system temperature readings
- **CPU Details** - Per-core CPU usage with historical graphs
- **Memory Analytics** - Detailed memory usage visualization

### 🎨 **Beautiful UI**
- **Multiple Themes** - Choose from Omarchy Default, Cyberpunk, and Retro CRT themes
- **Real-time Graphs** - Visual history tracking for CPU and memory
- **Interactive Tables** - Sortable columns for processes (CPU, Memory, PID, Name)
- **Clean Header** - System overview with hostname, load averages, and uptime

### ⚡ **Performance**
- **Lightweight** - Minimal resource footprint
- **Fast Updates** - 1-second refresh rate for responsive monitoring
- **Efficient** - Modular architecture with separated core logic and UI components

### 🔧 **Process Control**
- **Kill Processes** - Terminate processes directly from the monitor
- **Sorting Options** - Multiple sorting criteria for process list
- **Interactive Navigation** - Keyboard-driven interface for fast operation

## 🚀 Installation

### Prerequisites
- Go 1.22 or higher

### Building from Source

```bash
# Clone the repository
git clone https://github.com/manish12ys/osm.git
cd osm

# Download dependencies
go mod download

# Build the binary
go build -o omarchy-monitor

# Run the monitor
./omarchy-monitor
```

### Quick Install Script

```bash
# One-line install and run
git clone https://github.com/manish12ys/osm.git && cd osm && go build -o omarchy-monitor && ./omarchy-monitor
```

## 📚 Usage

### Starting the Monitor

```bash
./omarchy-monitor
```

### Quick Restart
Use the included restart script for rapid development iteration:

```bash
./restart.sh
```

This script will:
1. Kill any running instances
2. Wait for clean shutdown
3. Launch the new version

## ⌨️ Keybindings

| Key | Action |
|-----|--------|
| `q` | Quit the application |
| `1` | Switch to **Processes** view |
| `2` | Switch to **Disks** view |
| `3` | Switch to **Network** view |
| `4` | Switch to **Temperature** view |
| `5` | Switch to **CPU Detail** view |
| `6` | Switch to **Memory Detail** view |
| `t` | Cycle through themes |
| `c` | Sort processes by **CPU** usage |
| `m` | Sort processes by **Memory** usage |
| `p` | Sort processes by **PID** |
| `n` | Sort processes by **Name** |
| `k` | **Kill** selected process (Processes view only) |

## 🎨 Themes

OSM comes with three built-in themes. Press `t` to cycle through them:

### 1. **Omarchy Default**
Clean, professional theme with blue borders and traditional colors.

### 2. **Cyberpunk**
Neon pink and cyan color scheme for a futuristic look.

### 3. **Retro CRT**
Monochromatic green terminal aesthetic reminiscent of classic computing.


## 🔧 Technical Details

### Dependencies
- **[tview](https://github.com/rivo/tview)** `v0.42.0` - Terminal UI framework
- **[tcell](https://github.com/gdamore/tcell)** `v2.8.1` - Terminal handling
- **[gopsutil](https://github.com/shirou/gopsutil)** `v3.24.5` - Cross-platform system information

### Architecture
- **Modular Design** - Separation between core logic (`core/`) and UI (`ui/`)
- **Concurrent Updates** - Goroutine-based stats collection with 1-second intervals
- **Page-based Navigation** - Efficient view switching with tview Pages
- **Theme System** - Centralized color management for consistent styling

## 🤝 Contributing

Contributions are welcome! Here's how you can help:

1. **Fork the repository**
2. **Create a feature branch**
   ```bash
   git checkout -b feature/amazing-feature
   ```
3. **Make your changes**
   - Follow Go conventions and best practices
   - Test your changes thoroughly
   - Update documentation if needed
4. **Commit your changes**
   ```bash
   git commit -m 'Add amazing feature'
   ```
5. **Push to the branch**
   ```bash
   git push origin feature/amazing-feature
   ```
6. **Open a Pull Request**

### Ideas for Contributions
- 🎨 Add new themes
- 📊 Implement additional monitoring views (GPU, Battery, etc.)
- 🔍 Add search/filter functionality for processes
- 💾 Export statistics to file
- 🐳 Docker support
- 🪟 Windows and macOS compatibility

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [tview](https://github.com/rivo/tview) by Oliver Kuederle
- System stats powered by [gopsutil](https://github.com/shirou/gopsutil)
- Inspired by htop, btop, and other terminal system monitors

## 📬 Contact

**Project Link:** [https://github.com/manish12ys/osm](https://github.com/manish12ys/osm)

---

<div align="center">


*Star ⭐ this repository if you find it helpful!*

</div>
