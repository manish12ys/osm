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

### Arch Linux (Recommended)

```bash
# Clone the repository
git clone https://github.com/manish12ys/osm.git
cd osm

# Build and install Arch package
make install-pkg

# Or use the build script
./build-pkg.sh
sudo pacman -U omarchy-monitor-bin-*.pkg.tar.zst

# Run the monitor
osm
```

### From Source

**Prerequisites**: Go 1.22 or higher

```bash
# Clone the repository
git clone https://github.com/manish12ys/osm.git
cd osm

# Build the binary
make build

# Install to /usr/bin
sudo make install

# Run the monitor
osm
```

### Quick Development Build

```bash
# Build and run directly
make run

# Or build manually
go build -o omarchy-monitor
./omarchy-monitor
```

## 📚 Usage

### Running the Monitor

After installation:
```bash
osm
```

Or from source directory:
```bash
./omarchy-monitor
```

### Development Workflow

```bash
# Quick restart during development
./restart.sh

# Build and package
make package

# Clean build artifacts
make clean
```

The restart script will:
1. Kill any running instances
2. Wait for clean shutdown
3. Rebuild and launch the new version

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

### Building & Packaging

The project includes multiple build methods:

| Method | Command | Use Case |
|--------|---------|----------|
| **Make** | `make build` | Quick development builds |
| **Make Install** | `make install` | Install binary to `/usr/bin/osm` |
| **Arch Package** | `make package` | Create Arch Linux package |
| **Arch Install** | `make install-pkg` | Build and install as package |

Available PKGBUILD files:
- `PKGBUILD` - For stable GitHub releases
- `PKGBUILD.local` - For local development builds
- `PKGBUILD.bin` - For pre-compiled binaries

See `Makefile` for all available commands.

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
- 📊 Implement additional monitoring views (GPU, Battery, etc.)
- 🔍 Add search/filter functionality for processes
- 💾 Export statistics to file
- 📦 Package for other distributions (AUR, DEB, RPM)
- 🪟 Windows and macOS compatibility

## 🏗️ Building for Distribution

### Arch Linux Package

```bash
# Build package
./build-pkg.sh

# Install
sudo pacman -U omarchy-monitor-bin-*.pkg.tar.zst
```

### For AUR Submission

```bash
# Generate .SRCINFO
makepkg --printsrcinfo > .SRCINFO

# Follow AUR submission guidelines
```

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
