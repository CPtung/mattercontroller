#!/bin/bash

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
	sed -n 's/.*PTY is \/dev\/pts\/\([0-9]*\).*/\1/p' /var/log/rcp.log