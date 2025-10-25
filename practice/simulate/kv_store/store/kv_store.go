package store

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// DeltaSnapshot stores only the keys that changed and their old values
type DeltaSnapshot struct {
	ChangedKeys map[string]*int // key -> old value (nil if key didn't exist)
	DeletedKeys map[string]int  // key -> old value (for keys that were deleted)
}

// KVStore represents a key-value store with delta-based snapshot capabilities
type KVStore struct {
	mu          sync.RWMutex     // Protects data and checkpoints
	data        map[string]int
	valueCount  map[int]int      // value -> count
	checkpoints []*DeltaSnapshot // Stack of delta snapshots
	tracking    map[string]*int  // Tracks original values since last checkpoint (nil = new key)
}

// NewKVStore creates a new KV store instance
func NewKVStore() *KVStore {
	return &KVStore{
		data:        make(map[string]int),
		valueCount:  make(map[int]int),
		checkpoints: []*DeltaSnapshot{},
		tracking:    make(map[string]*int),
	}
}

// Put sets a key-value pair in the store
func (kv *KVStore) Put(key string, value int) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	// Track original value if we have checkpoints and haven't tracked this key yet
	if len(kv.checkpoints) > 0 {
		if _, alreadyTracked := kv.tracking[key]; !alreadyTracked {
			if oldValue, existed := kv.data[key]; existed {
				oldValueCopy := oldValue
				kv.tracking[key] = &oldValueCopy
			} else {
				kv.tracking[key] = nil // Mark as new key
			}
		}
	}

	// If key exists, decrement the old value's count
	if oldValue, exists := kv.data[key]; exists {
		kv.valueCount[oldValue]--
		if kv.valueCount[oldValue] == 0 {
			delete(kv.valueCount, oldValue)
		}
	}

	// Set the new value and increment its count
	kv.data[key] = value
	kv.valueCount[value]++
}

// Get retrieves the value for a given key
func (kv *KVStore) Get(key string) (int, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	value, exists := kv.data[key]
	return value, exists
}

// Delete removes a key-value pair from the store
// Returns true if the key existed and was deleted, false otherwise
func (kv *KVStore) Delete(key string) bool {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	if oldValue, exists := kv.data[key]; exists {
		// Track original value if we have checkpoints and haven't tracked this key yet
		if len(kv.checkpoints) > 0 {
			if _, alreadyTracked := kv.tracking[key]; !alreadyTracked {
				oldValueCopy := oldValue
				kv.tracking[key] = &oldValueCopy
			}
		}

		delete(kv.data, key)
		kv.valueCount[oldValue]--
		if kv.valueCount[oldValue] == 0 {
			delete(kv.valueCount, oldValue)
		}
		return true
	}
	return false
}

// CountValue returns the number of keys that have the given value
// O(1) lookup using the valueCount map
func (kv *KVStore) CountValue(value int) int {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	return kv.valueCount[value]
}

// Checkpoint creates a delta snapshot storing the current tracking info
// After this, changes continue to be tracked for the next revert
func (kv *KVStore) Checkpoint() {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	// Create a delta snapshot with CURRENT tracking info
	// This represents changes since the PREVIOUS checkpoint (or start)
	delta := &DeltaSnapshot{
		ChangedKeys: make(map[string]*int),
		DeletedKeys: make(map[string]int),
	}

	// Copy current tracking to the snapshot
	for key, oldValue := range kv.tracking {
		// Check if key still exists in current data
		if _, exists := kv.data[key]; exists {
			delta.ChangedKeys[key] = oldValue
		} else {
			// Key was deleted
			if oldValue != nil {
				delta.DeletedKeys[key] = *oldValue
			}
		}
	}

	kv.checkpoints = append(kv.checkpoints, delta)

	// Clear tracking - next checkpoint should track from THIS point
	kv.tracking = make(map[string]*int)
}

// Revert restores the state from the last checkpoint
func (kv *KVStore) Revert() error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	if len(kv.checkpoints) == 0 {
		return fmt.Errorf("no checkpoints to revert to")
	}

	// Get the last checkpoint delta
	delta := kv.checkpoints[len(kv.checkpoints)-1]
	kv.checkpoints = kv.checkpoints[:len(kv.checkpoints)-1]

	// Undo ONLY the current tracking (changes since last checkpoint)
	for key, oldValue := range kv.tracking {
		if currentValue, exists := kv.data[key]; exists {
			// Key exists and was changed
			kv.valueCount[currentValue]--
			if kv.valueCount[currentValue] == 0 {
				delete(kv.valueCount, currentValue)
			}

			if oldValue == nil {
				// Key didn't exist before, delete it
				delete(kv.data, key)
			} else {
				// Restore old value
				kv.data[key] = *oldValue
				kv.valueCount[*oldValue]++
			}
		} else {
			// Key was deleted after checkpoint, restore it
			if oldValue != nil {
				kv.data[key] = *oldValue
				kv.valueCount[*oldValue]++
			}
		}
	}

	// Now restore the tracking from the checkpoint delta
	// This becomes our new tracking for potential next revert
	kv.tracking = make(map[string]*int)
	for key, oldValue := range delta.ChangedKeys {
		kv.tracking[key] = oldValue
	}
	for key, oldValue := range delta.DeletedKeys {
		oldValueCopy := oldValue
		kv.tracking[key] = &oldValueCopy
	}

	return nil
}

// SaveToDisk saves the current state to a file
func (kv *KVStore) SaveToDisk(filename string) error {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	state := struct {
		Data        map[string]int      `json:"data"`
		ValueCount  map[int]int         `json:"valueCount"`
		Checkpoints []*DeltaSnapshot    `json:"checkpoints"`
		Tracking    map[string]*int     `json:"tracking"`
	}{
		Data:        kv.data,
		ValueCount:  kv.valueCount,
		Checkpoints: kv.checkpoints,
		Tracking:    kv.tracking,
	}

	if err := encoder.Encode(state); err != nil {
		return fmt.Errorf("failed to encode state: %w", err)
	}

	return nil
}

// LoadFromDisk loads the state from a file
func (kv *KVStore) LoadFromDisk(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var state struct {
		Data        map[string]int      `json:"data"`
		ValueCount  map[int]int         `json:"valueCount"`
		Checkpoints []*DeltaSnapshot    `json:"checkpoints"`
		Tracking    map[string]*int     `json:"tracking"`
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&state); err != nil {
		return fmt.Errorf("failed to decode state: %w", err)
	}

	kv.mu.Lock()
	defer kv.mu.Unlock()

	kv.data = state.Data
	kv.valueCount = state.ValueCount
	kv.checkpoints = state.Checkpoints
	kv.tracking = state.Tracking

	return nil
}

// GetCheckpointCount returns the number of checkpoints
func (kv *KVStore) GetCheckpointCount() int {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	return len(kv.checkpoints)
}

// GetAllData returns a copy of all key-value pairs and value counts
func (kv *KVStore) GetAllData() (map[string]int, map[int]int) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	dataCopy := make(map[string]int, len(kv.data))
	for k, v := range kv.data {
		dataCopy[k] = v
	}

	countCopy := make(map[int]int, len(kv.valueCount))
	for v, count := range kv.valueCount {
		countCopy[v] = count
	}

	return dataCopy, countCopy
}

// Print displays the current state (for debugging)
func (kv *KVStore) Print() {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	fmt.Println("Current State:")
	if len(kv.data) == 0 {
		fmt.Println("  (empty)")
	} else {
		for k, v := range kv.data {
			fmt.Printf("  %s: %d\n", k, v)
		}
	}
	fmt.Printf("Checkpoints: %d\n", len(kv.checkpoints))
}
