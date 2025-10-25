package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"workshop/practice/simulate/kv_store/store"
)

var (
	kvStore     = store.NewKVStore()
	defaultFile = ".kv_store.json"
	autoLoad    = true
	autoSave    = true
)

// Command-line flags
var (
	key   = flag.String("key", "", "Key for get/put/delete operations")
	value = flag.String("value", "", "Value for put/count operations")
	file  = flag.String("file", "kv_store.json", "File path for save/load operations")
)

func main() {
	// Custom usage message
	flag.Usage = printUsage

	// Parse flags
	flag.Parse()

	// Auto-load from default file if it exists
	if autoLoad {
		if _, err := os.Stat(defaultFile); err == nil {
			kvStore.LoadFromDisk(defaultFile)
		}
	}

	// Get command (first non-flag argument)
	args := flag.Args()
	if len(args) < 1 {
		// No command provided - enter interactive mode
		runInteractiveMode()
		return
	}

	command := args[0]
	executeCommand(command)

	// Auto-save after each command (except load and help)
	if autoSave && command != "load" && command != "help" {
		kvStore.SaveToDisk(defaultFile)
	}
}

func runInteractiveMode() {
	fmt.Println("🗄️  Mini KV Store - Interactive Mode")
	fmt.Println("Type 'help' for available commands or 'exit' to quit")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("kv> ")

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())

		if input == "" {
			continue
		}

		if input == "exit" || input == "quit" {
			fmt.Println("Goodbye!")
			break
		}

		executeInteractiveCommand(input)

		// Auto-save after each command
		if autoSave {
			kvStore.SaveToDisk(defaultFile)
		}
	}
}

func executeInteractiveCommand(input string) {
	parts := parseInput(input)
	if len(parts) == 0 {
		return
	}

	command := parts[0]

	switch command {
	case "put":
		if len(parts) < 3 {
			fmt.Println("Usage: put <key> <value>")
			return
		}
		key := parts[1]
		valueStr := strings.Join(parts[2:], " ")
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			fmt.Printf("❌ Error: value must be an integer, got '%s'\n", valueStr)
			return
		}
		kvStore.Put(key, value)
		fmt.Printf("✅ Set '%s' = %d\n", key, value)

	case "get":
		if len(parts) < 2 {
			fmt.Println("Usage: get <key>")
			return
		}
		key := parts[1]
		val, exists := kvStore.Get(key)
		if exists {
			fmt.Printf("✅ '%s' = %d\n", key, val)
		} else {
			fmt.Printf("❌ Key '%s' not found\n", key)
		}

	case "delete", "del":
		if len(parts) < 2 {
			fmt.Println("Usage: delete <key>")
			return
		}
		key := parts[1]
		if deleted := kvStore.Delete(key); deleted {
			fmt.Printf("✅ Deleted key '%s'\n", key)
		} else {
			fmt.Printf("❌ Key '%s' not found\n", key)
		}

	case "count":
		if len(parts) < 2 {
			fmt.Println("Usage: count <value>")
			return
		}
		valueStr := strings.Join(parts[1:], " ")
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			fmt.Printf("❌ Error: value must be an integer, got '%s'\n", valueStr)
			return
		}
		count := kvStore.CountValue(value)
		fmt.Printf("✅ Value %d appears %d time(s)\n", value, count)

	case "checkpoint", "cp":
		kvStore.Checkpoint()
		count := kvStore.GetCheckpointCount()
		fmt.Printf("✅ Checkpoint created (total: %d)\n", count)

	case "revert", "rv":
		err := kvStore.Revert()
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
		} else {
			fmt.Println("✅ Reverted to last checkpoint")
		}

	case "save":
		filename := defaultFile
		if len(parts) > 1 {
			filename = parts[1]
		}
		err := kvStore.SaveToDisk(filename)
		if err != nil {
			fmt.Printf("❌ Error saving: %v\n", err)
		} else {
			fmt.Printf("✅ Saved to '%s'\n", filename)
		}

	case "load":
		filename := defaultFile
		if len(parts) > 1 {
			filename = parts[1]
		}
		err := kvStore.LoadFromDisk(filename)
		if err != nil {
			fmt.Printf("❌ Error loading: %v\n", err)
		} else {
			fmt.Printf("✅ Loaded from '%s'\n", filename)
		}

	case "list", "ls":
		printList()

	case "clear":
		kvStore = store.NewKVStore()
		fmt.Println("✅ Store cleared")

	case "help", "?":
		printInteractiveHelp()

	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Type 'help' for available commands")
	}
}

func parseInput(input string) []string {
	// Simple parsing that handles quoted strings
	var parts []string
	var current strings.Builder
	inQuotes := false

	for i := 0; i < len(input); i++ {
		char := input[i]

		switch char {
		case '"':
			inQuotes = !inQuotes
		case ' ':
			if inQuotes {
				current.WriteByte(char)
			} else if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
		default:
			current.WriteByte(char)
		}
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}

func executeCommand(command string) {
	switch command {
	case "put":
		if *key == "" || *value == "" {
			fmt.Println("Error: put requires --key and --value flags")
			fmt.Println("Usage: kv-cli put --key <key> --value <value>")
			os.Exit(1)
		}
		intValue, err := strconv.Atoi(*value)
		if err != nil {
			fmt.Printf("❌ Error: value must be an integer, got '%s'\n", *value)
			os.Exit(1)
		}
		kvStore.Put(*key, intValue)
		fmt.Printf("✅ Set '%s' = %d\n", *key, intValue)

	case "get":
		if *key == "" {
			fmt.Println("Error: get requires --key flag")
			fmt.Println("Usage: kv-cli get --key <key>")
			os.Exit(1)
		}
		val, exists := kvStore.Get(*key)
		if exists {
			fmt.Printf("✅ '%s' = %d\n", *key, val)
		} else {
			fmt.Printf("❌ Key '%s' not found\n", *key)
			os.Exit(1)
		}

	case "delete", "del":
		if *key == "" {
			fmt.Println("Error: delete requires --key flag")
			fmt.Println("Usage: kv-cli delete --key <key>")
			os.Exit(1)
		}
		if deleted := kvStore.Delete(*key); deleted {
			fmt.Printf("✅ Deleted key '%s'\n", *key)
		} else {
			fmt.Printf("❌ Key '%s' not found\n", *key)
			os.Exit(1)
		}

	case "count":
		if *value == "" {
			fmt.Println("Error: count requires --value flag")
			fmt.Println("Usage: kv-cli count --value <value>")
			os.Exit(1)
		}
		intValue, err := strconv.Atoi(*value)
		if err != nil {
			fmt.Printf("❌ Error: value must be an integer, got '%s'\n", *value)
			os.Exit(1)
		}
		count := kvStore.CountValue(intValue)
		fmt.Printf("✅ Value %d appears %d time(s)\n", intValue, count)

	case "checkpoint", "cp":
		kvStore.Checkpoint()
		count := kvStore.GetCheckpointCount()
		fmt.Printf("✅ Checkpoint created (total: %d)\n", count)

	case "revert", "rv":
		err := kvStore.Revert()
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Reverted to last checkpoint")

	case "save":
		err := kvStore.SaveToDisk(*file)
		if err != nil {
			fmt.Printf("❌ Error saving: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ Saved to '%s'\n", *file)

	case "load":
		err := kvStore.LoadFromDisk(*file)
		if err != nil {
			fmt.Printf("❌ Error loading: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ Loaded from '%s'\n", *file)

	case "list", "ls":
		printList()

	case "clear":
		kvStore = store.NewKVStore()
		fmt.Println("✅ Store cleared")

	case "help", "?", "-h", "--help":
		printHelp()

	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Run 'kv-cli help' for usage information")
		os.Exit(1)
	}
}

func printList() {
	data, valueCounts := kvStore.GetAllData()

	if len(data) == 0 {
		fmt.Println("Store is empty")
		return
	}

	fmt.Println("\nCurrent Key-Value Pairs:")
	fmt.Println("------------------------")
	for k, v := range data {
		fmt.Printf("  %s = %d\n", k, v)
	}

	fmt.Println("\nValue Counts:")
	fmt.Println("-------------")
	for v, count := range valueCounts {
		fmt.Printf("  %d → %d\n", v, count)
	}

	checkpointCount := kvStore.GetCheckpointCount()
	fmt.Printf("\nTotal keys: %d\n", len(data))
	fmt.Printf("Checkpoints: %d\n", checkpointCount)
}

func printUsage() {
	fmt.Println("Mini KV Store - Command Line Interface")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  kv-cli                    # Interactive mode")
	fmt.Println("  kv-cli [flags] <command>  # Flag-based mode")
	fmt.Println()
	fmt.Println("FLAGS:")
	fmt.Println("  -key <key>       Key for get/put/delete operations")
	fmt.Println("  -value <value>   Value for put/count operations")
	fmt.Println("  -file <path>     File path for save/load (default: kv_store.json)")
	fmt.Println()
	fmt.Println("COMMANDS:")
	fmt.Println("  put              Store a key-value pair (requires -key and -value)")
	fmt.Println("  get              Retrieve value for a key (requires -key)")
	fmt.Println("  delete, del      Delete a key-value pair (requires -key)")
	fmt.Println("  count            Count keys with the given value (requires -value)")
	fmt.Println("  checkpoint, cp   Create a snapshot of current state")
	fmt.Println("  revert, rv       Revert to last checkpoint")
	fmt.Println("  save             Save to disk (optional -file)")
	fmt.Println("  load             Load from disk (optional -file)")
	fmt.Println("  list, ls         Show all key-value pairs")
	fmt.Println("  clear            Clear the entire store")
	fmt.Println("  help, ?, -h      Show this help message")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  kv-cli -key name -value Alice put")
	fmt.Println("  kv-cli -key message -value \"hello world\" put")
	fmt.Println("  kv-cli -key name get")
	fmt.Println("  kv-cli -key name delete")
	fmt.Println("  kv-cli -value active count")
	fmt.Println("  kv-cli checkpoint")
	fmt.Println("  kv-cli revert")
	fmt.Println("  kv-cli -file backup.json save")
	fmt.Println("  kv-cli -file backup.json load")
	fmt.Println("  kv-cli list")
	fmt.Println()
	fmt.Println("NOTES:")
	fmt.Println("  • Run without arguments to enter interactive mode")
	fmt.Println("  • In flag-based mode, flags must come BEFORE the command")
	fmt.Println("  • Data automatically persists to .kv_store.json")
	fmt.Println("  • Use quotes for values with spaces")
}

func printHelp() {
	printUsage()
}

func printInteractiveHelp() {
	fmt.Println("\nAvailable Commands:")
	fmt.Println("-------------------")
	fmt.Println("  put <key> <value>    Store a key-value pair")
	fmt.Println("  get <key>            Retrieve value for a key")
	fmt.Println("  delete <key>         Delete a key-value pair")
	fmt.Println("  count <value>        Count keys with the given value")
	fmt.Println("  checkpoint           Create a snapshot of current state")
	fmt.Println("  revert               Revert to last checkpoint")
	fmt.Println("  save [file]          Save to disk (default: .kv_store.json)")
	fmt.Println("  load [file]          Load from disk (default: .kv_store.json)")
	fmt.Println("  list                 Show all key-value pairs")
	fmt.Println("  clear                Clear the entire store")
	fmt.Println("  help                 Show this help message")
	fmt.Println("  exit, quit           Exit interactive mode")
	fmt.Println("\nExamples:")
	fmt.Println("  put name Alice")
	fmt.Println("  put message \"hello world\"")
	fmt.Println("  get name")
	fmt.Println("  count active")
	fmt.Println("  save backup.json")
	fmt.Println()
}
