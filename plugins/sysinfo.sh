#!/bin/bash

# Force color output for tools that support it
export CLICOLOR_FORCE=1
export TERM=xterm-256color

OUTPUT=""
if command -v neofetch &> /dev/null; then
    OUTPUT=$(neofetch)
elif command -v fastfetch &> /dev/null; then
    OUTPUT=$(fastfetch)
else
    OUTPUT="Neither fastfetch nor neofetch found. Please install one of them to see system info."
fi

# Use python3 to safely create JSON if available
if command -v python3 &> /dev/null; then
    python3 -c "import json, sys; print(json.dumps({'title': 'System Info', 'value': sys.stdin.read()}))" <<< "$OUTPUT"
else
    # Fallback: Simple escaping (may not handle all edge cases)
    # Escape backslashes and double quotes
    ESCAPED=$(echo "$OUTPUT" | sed 's/\\/\\\\/g' | sed 's/"/\\"/g' | awk '{printf "%s\\n", $0}')
    echo "{\"title\": \"System Info\", \"value\": \"$ESCAPED\"}"
fi
