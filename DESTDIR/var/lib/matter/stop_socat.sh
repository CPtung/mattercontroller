#!/bin/bash

PID_FILE="/var/run/rcp.pid"

# check if pid file exists
if [ ! -f "$PID_FILE" ]; then
    echo "PID file $PID_FILE not found. Process may not be running."
    exit 1
fi

# get PID
PID=$(cat "$PID_FILE")

# check proccess still running
if ! kill -0 "$PID" 2>/dev/null; then
    echo "Process with PID $PID is not running."
    rm -f "$PID_FILE"
    exit 0
fi

kill "$PID