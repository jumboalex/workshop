#!/bin/bash
# Test script for KV Store CLI

echo "Testing KV Store CLI with automated commands..."
echo ""

# Create a test input file
cat << 'EOF' | ./cli
put name Alice
put age 30
put city NYC
put status active
put role admin
get name
get age
count active
list
checkpoint
put name Bob
put newkey newvalue
list
revert
list
save test_store.json
clear
list
load test_store.json
list
exit
EOF

echo ""
echo "Test completed!"
echo ""
echo "Saved file contents:"
cat test_store.json
