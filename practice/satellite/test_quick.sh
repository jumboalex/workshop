#!/bin/bash
echo "Starting mock server in background..."
./satellite -mode=server -port=8090 -unreliability=0.1 -slowness=200ms > /dev/null 2>&1 &
SERVER_PID=$!

sleep 2

echo "Running client tests..."
./satellite -mode=client -url=http://localhost:8090

echo ""
echo "Killing server..."
kill $SERVER_PID 2>/dev/null

echo "Done!"
