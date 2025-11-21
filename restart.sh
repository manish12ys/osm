#!/bin/bash
# Quick restart script for the Omarchy System Monitor

# Kill any running instances
pkill -f omarchy-monitor 2>/dev/null

# Small delay to ensure clean shutdown
sleep 0.5

# Run the new version
./omarchy-monitor
