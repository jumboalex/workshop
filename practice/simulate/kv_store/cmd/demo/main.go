package main

import (
	"fmt"
	"workshop/practice/simulate/kv_store/store"
)

func main() {
	fmt.Println("=== KV Store Demo ===\n")

	// Test 1: Basic Put/Get
	fmt.Println("Test 1: Basic Put/Get Operations")
	kv := store.NewKVStore()
	kv.Put("name", 100)
	kv.Put("age", 30)
	kv.Put("city", 10001)

	if val, exists := kv.Get("name"); exists {
		fmt.Printf("  name = %d\n", val)
	}
	if val, exists := kv.Get("age"); exists {
		fmt.Printf("  age = %d\n", val)
	}
	if val, exists := kv.Get("city"); exists {
		fmt.Printf("  city = %d\n", val)
	}
	fmt.Println()

	// Test 2: Value Count
	fmt.Println("Test 2: Value Count")
	kv.Put("user1", 1)
	kv.Put("user2", 1)
	kv.Put("user3", 0)
	kv.Put("user4", 1)

	activeCount := kv.CountValue(1)
	inactiveCount := kv.CountValue(0)
	fmt.Printf("  Active (1) count: %d\n", activeCount)
	fmt.Printf("  Inactive (0) count: %d\n", inactiveCount)
	fmt.Println()

	// Test 3: Checkpoint and Revert
	fmt.Println("Test 3: Checkpoint and Revert")
	kv.Put("name", 200)
	kv.Checkpoint()
	fmt.Println("Checkpoint created")

	kv.Put("name", 999)
	fmt.Printf("Modified: name = %d\n", func() int { v, _ := kv.Get("name"); return v }())

	kv.Revert()
	fmt.Printf("After revert: name = %d\n", func() int { v, _ := kv.Get("name"); return v }())
	fmt.Println()

	// Test 4: Multiple Checkpoints
	fmt.Println("Test 4: Multiple Checkpoints")
	kv.Put("state", 11)
	kv.Checkpoint()
	fmt.Println("Checkpoint 1 created (state=11)")

	kv.Put("state", 22)
	kv.Checkpoint()
	fmt.Println("Checkpoint 2 created (state=22)")

	kv.Put("state", 33)
	fmt.Printf("Current state value: %d\n", func() int { v, _ := kv.Get("state"); return v }())

	kv.Revert()
	fmt.Printf("After revert 1: %d\n", func() int { v, _ := kv.Get("state"); return v }())

	kv.Revert()
	fmt.Printf("After revert 2: %d\n", func() int { v, _ := kv.Get("state"); return v }())
	fmt.Println()

	// Test 5: Value count with updates
	fmt.Println("Test 5: Value Count with Updates")
	kv2 := store.NewKVStore()
	kv2.Put("k1", 5)
	kv2.Put("k2", 5)
	kv2.Put("k3", 6)
	fmt.Printf("  Value 5 count: %d\n", kv2.CountValue(5))
	fmt.Printf("  Value 6 count: %d\n", kv2.CountValue(6))
	fmt.Println()

	// Test 6: Save and Load
	fmt.Println("Test 6: Save and Load to Disk")
	kv3 := store.NewKVStore()
	kv3.Put("data", 777)
	kv3.Put("yes", 1)
	kv3.Checkpoint()
	kv3.Put("value", 50)

	err := kv3.SaveToDisk("test_demo.json")
	if err != nil {
		fmt.Printf("Error saving: %v\n", err)
		return
	}
	fmt.Println("  Saved to test_demo.json")

	kv4 := store.NewKVStore()
	err = kv4.LoadFromDisk("test_demo.json")
	if err != nil {
		fmt.Printf("Error loading: %v\n", err)
		return
	}
	fmt.Println("  Loaded from test_demo.json")

	if val, exists := kv4.Get("data"); exists {
		fmt.Printf("  Loaded data: %d\n", val)
	}
	fmt.Printf("  Checkpoint count after load: %d\n", kv4.GetCheckpointCount())
	fmt.Println()

	fmt.Println("=== All Tests Complete ===")
}
