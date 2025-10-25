#!/bin/bash
# Test flag-based CLI

echo "=== Testing Flag-Based CLI ==="
echo ""

# Put some data
./bin/kv-cli put --key name --value Alice
./bin/kv-cli put --key age --value 30
./bin/kv-cli put --key status --value active
./bin/kv-cli put --key role --value admin

# Save to file
./bin/kv-cli save --file test_cli.json

echo ""
echo "=== Loading and querying ==="
# Load and query
./bin/kv-cli load --file test_cli.json
./bin/kv-cli get --key name
./bin/kv-cli count --value active
./bin/kv-cli list

echo ""
echo "=== Testing checkpoint/revert ==="
./bin/kv-cli load --file test_cli.json
./bin/kv-cli checkpoint
./bin/kv-cli put --key name --value Bob
./bin/kv-cli get --key name
./bin/kv-cli revert
./bin/kv-cli get --key name

# Cleanup
rm -f test_cli.json

echo ""
echo "=== All tests complete! ==="
