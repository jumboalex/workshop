package store

import (
	"fmt"
	"os"
	"sync"
	"testing"
)

// TestPutAndGet tests basic put and get operations
func TestPutAndGet(t *testing.T) {
	kv := NewKVStore()

	// Test putting and getting a value
	kv.Put("name", 100)
	value, exists := kv.Get("name")

	if !exists {
		t.Error("Expected key 'name' to exist")
	}

	if value != 100 {
		t.Errorf("Expected value 'Alice', got %d", value)
	}

	// Test getting a non-existent key
	_, exists = kv.Get("nonexistent")
	if exists {
		t.Error("Expected key 'nonexistent' to not exist")
	}
}

// TestPutOverwrite tests overwriting existing keys
func TestPutOverwrite(t *testing.T) {
	kv := NewKVStore()

	kv.Put("key", 10)
	kv.Put("key", 20)

	value, exists := kv.Get("key")
	if !exists {
		t.Error("Expected key to exist")
	}

	if value != 20 {
		t.Errorf("Expected value 'value2', got %d", value)
	}

	// Check that value count was properly updated
	count1 := kv.CountValue(10)
	count2 := kv.CountValue(20)

	if count1 != 0 {
		t.Errorf("Expected count for 'value1' to be 0, got %d", count1)
	}

	if count2 != 1 {
		t.Errorf("Expected count for 'value2' to be 1, got %d", count2)
	}
}

// TestCountValue tests value counting functionality
func TestCountValue(t *testing.T) {
	kv := NewKVStore()

	kv.Put("user1", 1)
	kv.Put("user2", 1)
	kv.Put("user3", 2)
	kv.Put("user4", 1)

	activeCount := kv.CountValue(1)
	inactiveCount := kv.CountValue(2)
	nonexistentCount := kv.CountValue(999)

	if activeCount != 3 {
		t.Errorf("Expected active count 3, got %d", activeCount)
	}

	if inactiveCount != 1 {
		t.Errorf("Expected inactive count 1, got %d", inactiveCount)
	}

	if nonexistentCount != 0 {
		t.Errorf("Expected nonexistent count 0, got %d", nonexistentCount)
	}

	// Update a value and check counts again
	kv.Put("user1", 2)

	activeCount = kv.CountValue(1)
	inactiveCount = kv.CountValue(2)

	if activeCount != 2 {
		t.Errorf("Expected active count 2 after update, got %d", activeCount)
	}

	if inactiveCount != 2 {
		t.Errorf("Expected inactive count 2 after update, got %d", inactiveCount)
	}
}

// TestCheckpointAndRevert tests checkpoint and revert functionality
func TestCheckpointAndRevert(t *testing.T) {
	kv := NewKVStore()

	// Set initial data
	kv.Put("key1", 10)
	kv.Put("key2", 20)

	// Create checkpoint
	kv.Checkpoint()

	if kv.GetCheckpointCount() != 1 {
		t.Errorf("Expected 1 checkpoint, got %d", kv.GetCheckpointCount())
	}

	// Modify data
	kv.Put("key1", 99)
	kv.Put("key3", 88)

	value, _ := kv.Get("key1")
	if value != 99 {
		t.Errorf("Expected 99, got %d", value)
	}

	// Revert to checkpoint
	err := kv.Revert()
	if err != nil {
		t.Errorf("Revert failed: %v", err)
	}

	// Check that data was restored
	value, exists := kv.Get("key1")
	if !exists || value != 10 {
		t.Errorf("Expected 10 after revert, got %d", value)
	}

	_, exists = kv.Get("key3")
	if exists {
		t.Error("Expected key3 to not exist after revert")
	}

	// Check value counts were restored
	count := kv.CountValue(10)
	if count != 1 {
		t.Errorf("Expected count 1 after revert, got %d", count)
	}
}

// TestMultipleCheckpoints tests multiple levels of checkpoints
func TestMultipleCheckpoints(t *testing.T) {
	kv := NewKVStore()

	kv.Put("state", 11)
	kv.Checkpoint()

	kv.Put("state", 22)
	kv.Checkpoint()

	kv.Put("state", 33)

	if kv.GetCheckpointCount() != 2 {
		t.Errorf("Expected 2 checkpoints, got %d", kv.GetCheckpointCount())
	}

	// Revert once
	kv.Revert()
	value, _ := kv.Get("state")
	if value != 22 {
		t.Errorf("Expected 22 after first revert, got %d", value)
	}

	// Revert again
	kv.Revert()
	value, _ = kv.Get("state")
	if value != 11 {
		t.Errorf("Expected 11 after second revert, got %d", value)
	}

	if kv.GetCheckpointCount() != 0 {
		t.Errorf("Expected 0 checkpoints after reverting all, got %d", kv.GetCheckpointCount())
	}
}

// TestRevertWithoutCheckpoint tests reverting with no checkpoints
func TestRevertWithoutCheckpoint(t *testing.T) {
	kv := NewKVStore()

	err := kv.Revert()
	if err == nil {
		t.Error("Expected error when reverting without checkpoint")
	}
}

// TestSaveAndLoad tests disk persistence
func TestSaveAndLoad(t *testing.T) {
	kv1 := NewKVStore()

	kv1.Put("persistent", 777)
	kv1.Put("saved", 1)
	kv1.Checkpoint()
	kv1.Put("another", 50)

	filename := "/tmp/kv_store_test_unit.json"
	defer os.Remove(filename) // Cleanup

	// Save to disk
	err := kv1.SaveToDisk(filename)
	if err != nil {
		t.Fatalf("SaveToDisk failed: %v", err)
	}

	// Load into new store
	kv2 := NewKVStore()
	err = kv2.LoadFromDisk(filename)
	if err != nil {
		t.Fatalf("LoadFromDisk failed: %v", err)
	}

	// Verify data
	value, exists := kv2.Get("persistent")
	if !exists || value != 777 {
		t.Errorf("Expected 'data', got %d", value)
	}

	// Verify checkpoint was loaded
	if kv2.GetCheckpointCount() != 1 {
		t.Errorf("Expected 1 checkpoint after load, got %d", kv2.GetCheckpointCount())
	}

	// Verify value counts
	count := kv2.CountValue(777)
	if count != 1 {
		t.Errorf("Expected count 1 after load, got %d", count)
	}
}

// TestGetAllData tests the GetAllData method
func TestGetAllData(t *testing.T) {
	kv := NewKVStore()

	kv.Put("key1", 10)
	kv.Put("key2", 10)
	kv.Put("key3", 20)

	data, valueCounts := kv.GetAllData()

	if len(data) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(data))
	}

	if valueCounts[10] != 2 {
		t.Errorf("Expected count 2 for 'value1', got %d", valueCounts[10])
	}

	if valueCounts[20] != 1 {
		t.Errorf("Expected count 1 for 'value2', got %d", valueCounts[20])
	}

	// Verify it's a copy (mutation shouldn't affect store)
	data["key1"] = 99
	value, _ := kv.Get("key1")
	if value == 99 {
		t.Error("GetAllData should return a copy, not the original data")
	}
}

// TestConcurrentPut tests concurrent put operations
func TestConcurrentPut(t *testing.T) {
	kv := NewKVStore()
	var wg sync.WaitGroup

	numGoroutines := 10
	numOpsPerGoroutine := 100

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOpsPerGoroutine; j++ {
				key := fmt.Sprintf("key_%d_%d", id, j)
				value := id
				kv.Put(key, value)
			}
		}(i)
	}

	wg.Wait()

	data, _ := kv.GetAllData()
	expectedKeys := numGoroutines * numOpsPerGoroutine

	if len(data) != expectedKeys {
		t.Errorf("Expected %d keys, got %d", expectedKeys, len(data))
	}
}

// TestConcurrentGetAndPut tests concurrent reads and writes
func TestConcurrentGetAndPut(t *testing.T) {
	kv := NewKVStore()
	var wg sync.WaitGroup

	// Pre-populate with some data
	for i := 0; i < 100; i++ {
		kv.Put(fmt.Sprintf("key_%d", i), i)
	}

	// Concurrent reads
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				kv.Get(fmt.Sprintf("key_%d", j))
			}
		}(i)
	}

	// Concurrent writes
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				key := fmt.Sprintf("new_key_%d_%d", id, j)
				kv.Put(key, 50)
			}
		}(i)
	}

	wg.Wait()

	// Verify no race conditions
	data, _ := kv.GetAllData()
	if len(data) < 100 {
		t.Errorf("Expected at least 100 keys, got %d", len(data))
	}
}

// TestConcurrentCheckpoint tests concurrent checkpoint operations
func TestConcurrentCheckpoint(t *testing.T) {
	kv := NewKVStore()
	var wg sync.WaitGroup

	kv.Put("key", 50)

	// Multiple concurrent checkpoints
	numCheckpoints := 5
	wg.Add(numCheckpoints)
	for i := 0; i < numCheckpoints; i++ {
		go func() {
			defer wg.Done()
			kv.Checkpoint()
		}()
	}

	wg.Wait()

	count := kv.GetCheckpointCount()
	if count != numCheckpoints {
		t.Errorf("Expected %d checkpoints, got %d", numCheckpoints, count)
	}
}

// TestEmptyStore tests operations on empty store
func TestEmptyStore(t *testing.T) {
	kv := NewKVStore()

	// Get from empty store
	_, exists := kv.Get("key")
	if exists {
		t.Error("Expected key to not exist in empty store")
	}

	// Count in empty store
	count := kv.CountValue(50)
	if count != 0 {
		t.Errorf("Expected count 0 in empty store, got %d", count)
	}

	// GetAllData on empty store
	data, valueCounts := kv.GetAllData()
	if len(data) != 0 || len(valueCounts) != 0 {
		t.Error("Expected empty maps from empty store")
	}

	// Checkpoint count on empty store
	if kv.GetCheckpointCount() != 0 {
		t.Errorf("Expected 0 checkpoints in new store, got %d", kv.GetCheckpointCount())
	}
}

// TestCheckpointPreservesValueCounts tests that checkpoints preserve value counts correctly
func TestCheckpointPreservesValueCounts(t *testing.T) {
	kv := NewKVStore()

	kv.Put("k1", 5)
	kv.Put("k2", 5)
	kv.Put("k3", 6)

	kv.Checkpoint()

	// Change values
	kv.Put("k1", 6)

	// Counts should have changed
	if kv.CountValue(5) != 1 {
		t.Errorf("Expected foo count 1 before revert, got %d", kv.CountValue(5))
	}

	if kv.CountValue(6) != 2 {
		t.Errorf("Expected bar count 2 before revert, got %d", kv.CountValue(6))
	}

	// Revert
	kv.Revert()

	// Counts should be restored to checkpoint values
	if kv.CountValue(5) != 2 {
		t.Errorf("Expected foo count 2 after revert, got %d", kv.CountValue(5))
	}

	if kv.CountValue(6) != 1 {
		t.Errorf("Expected bar count 1 after revert, got %d", kv.CountValue(6))
	}
}

// TestLoadInvalidFile tests loading from non-existent or invalid file
func TestLoadInvalidFile(t *testing.T) {
	kv := NewKVStore()

	// Test non-existent file
	err := kv.LoadFromDisk("/nonexistent/path/file.json")
	if err == nil {
		t.Error("Expected error when loading non-existent file")
	}

	// Test invalid JSON
	invalidFile := "/tmp/invalid_kv_test.json"
	os.WriteFile(invalidFile, []byte("invalid json"), 0644)
	defer os.Remove(invalidFile)

	err = kv.LoadFromDisk(invalidFile)
	if err == nil {
		t.Error("Expected error when loading invalid JSON")
	}
}

// BenchmarkPut benchmarks put operations
func BenchmarkPut(b *testing.B) {
	kv := NewKVStore()
	for i := 0; i < b.N; i++ {
		kv.Put(fmt.Sprintf("key_%d", i), i)
	}
}

// BenchmarkGet benchmarks get operations
func BenchmarkGet(b *testing.B) {
	kv := NewKVStore()
	// Pre-populate
	for i := 0; i < 1000; i++ {
		kv.Put(fmt.Sprintf("key_%d", i), i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kv.Get(fmt.Sprintf("key_%d", i%1000))
	}
}

// BenchmarkCountValue benchmarks value counting
func BenchmarkCountValue(b *testing.B) {
	kv := NewKVStore()
	// Pre-populate
	for i := 0; i < 1000; i++ {
		kv.Put(fmt.Sprintf("key_%d", i), 42)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kv.CountValue(42)
	}
}

// BenchmarkCheckpoint benchmarks checkpoint creation
func BenchmarkCheckpoint(b *testing.B) {
	kv := NewKVStore()
	// Pre-populate
	for i := 0; i < 100; i++ {
		kv.Put(fmt.Sprintf("key_%d", i), i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kv.Checkpoint()
	}
}

// BenchmarkRevert benchmarks revert operations
func BenchmarkRevert(b *testing.B) {
	kv := NewKVStore()
	// Pre-populate and checkpoint
	for i := 0; i < 100; i++ {
		kv.Put(fmt.Sprintf("key_%d", i), i)
	}

	for i := 0; i < b.N; i++ {
		kv.Checkpoint()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kv.Revert()
		if i < b.N-1 {
			kv.Checkpoint() // Re-checkpoint for next iteration
		}
	}
}
