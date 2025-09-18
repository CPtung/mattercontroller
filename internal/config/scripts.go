package config

import (
	"fmt"
	"os"
)

var (
	startRCPScript = `#!/bin/bash

# start socat background process
socat -d -d EXEC:"ot-rcp 1",pty,raw,echo=0 pty,raw,echo=0 > /var/log/rcp.log 2>&1 &

# save PID
PID=$!

# print PID
echo $PID > /var/run/rcp.pid

if ! kill -0 $PID 2>/dev/null; then
    echo "Failed to start process"
    exit 1
fi

sleep 1 && \
	sed -n 's/.*PTY is \/dev\/pts\/\([0-9]*\).*/\1/p' /var/log/rcp.log`

	stopRCPScript = `#!/bin/bash

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

kill "$PID"`
)

func createScripts() error {

	if err := os.WriteFile(LibPath+"/start_rcp.sh", []byte(startRCPScript), 0755); err != nil {
		return fmt.Errorf("建立 start rcp script 失敗: %s\n", err.Error())
	}
	if err := os.WriteFile(LibPath+"/stop_rcp.sh", []byte(stopRCPScript), 0755); err != nil {
		return fmt.Errorf("建立 stop rcp script 失敗: %s\n", err.Error())
	}
	return nil
}
