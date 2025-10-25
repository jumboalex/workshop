#!/bin/bash
echo "Starting mock server..."
./satellite -mode=server -port=8091 -unreliability=0.1 -slowness=100ms > /dev/null 2>&1 &
SERVER_PID=$!

sleep 1

echo "Running client..."
./satellite -mode=client -url=http://localhost:8091

echo ""
echo "Stopping server..."
kill $SERVER_PID 2>/dev/null

echo "✅ Test complete!"
