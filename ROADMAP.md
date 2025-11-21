# 🗺️ Omarchy System Monitor (OSM) Roadmap

This document outlines the future development plans for the Omarchy System Monitor.

## 🟢 Phase 1: Enhanced Functionality (Short-term)
*Focus on closing feature gaps with established monitors like htop/btop.*

*   **Advanced Process Management**:
	*   **Search & Filter**: ✅ Implemented (Press `/` to search)
	*   **Tree View**: ✅ Implemented (Display processes in a hierarchical tree view to show parent/child relationships.)
*   **Hardware Expansion**:
	*   **GPU Monitoring**: ✅ Implemented (Add support for NVIDIA/AMD GPU usage, VRAM, and temperature stats.)
	*   **Battery Stats**: ✅ Implemented (Visible in Header)
*   **Data Export**:
	*   ✅ Implemented (Allow users to export current system snapshots or log statistics to a file (CSV/JSON) for later analysis.)

## 🟡 Phase 2: Platform & Distribution (Medium-term)
*Focus on making OSM available to a wider audience.*

*   **Cross-Platform Compatibility**:
	*   **Windows Support**: ✅ Implemented (Adapted `gopsutil` calls and UI rendering to work correctly on Windows terminals (PowerShell/CMD).)
	*   **macOS Support**: ✅ Implemented (Ensure full compatibility with macOS, specifically for CPU and temperature sensors.)
*   **Packaging**:
	*   **Official AUR Package**: ✅ Implemented (Submit `osm` to the Arch User Repository.)
## 🔴 Phase 3: Customization & Advanced Features (Long-term)
*Focus on power users and flexibility.*

*   **Configuration System**:
	*   ✅ Implemented (Create a `~/.config/osm/config.yaml` file to persist user preferences (default theme, refresh rate, default view, custom keybindings).)
*   **Remote Monitoring**:
	*   Implement a client-server mode where OSM can connect to a remote server to monitor it via SSH or a lightweight agent.
*   **Container Monitoring**:
	*   ✅ Implemented (Add a dedicated view for Docker/Podman containers, showing resource usage per container.)
*   **Plugin System**:
	*   Allow users to write simple scripts (Lua or Shell) to display custom metrics in a user-defined panel.

## 🎨 Design & UX Improvements (Ongoing)
*   **New Themes**: ✅ Implemented (Added "Dracula", "Solarized", and "Nord" themes.)
*   **Responsive Layout**: ✅ Implemented (Better handling of very small or very large terminal window sizes.)
*   **Mouse Support**: ✅ Implemented (Enable clicking on column headers to sort or clicking tabs to switch views.)
