package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// MockServer simulates an unreliable satellite API
type MockServer struct {
	mu              sync.RWMutex
	sensors         map[int]*Sensor
	nextID          int
	unreliability   float64 // 0.0 to 1.0 - probability of failure
	slowness        time.Duration
	resourceLimited bool
}

// NewMockServer creates a new mock satellite server
func NewMockServer(unreliability float64, slowness time.Duration) *MockServer {
	return &MockServer{
		sensors:       make(map[int]*Sensor),
		nextID:        1,
		unreliability: unreliability,
		slowness:      slowness,
	}
}

// simulateUnreliability randomly fails or delays requests
func (s *MockServer) simulateUnreliability(w http.ResponseWriter) bool {
	// Random delay to simulate slow satellite connection
	if s.slowness > 0 {
		delay := time.Duration(rand.Float64() * float64(s.slowness))
		time.Sleep(delay)
	}

	// Random chance of complete failure
	if rand.Float64() < s.unreliability {
		// Randomly choose type of failure
		failureType := rand.Intn(3)
		switch failureType {
		case 0:
			// Connection timeout (no response)
			time.Sleep(10 * time.Second)
			return false
		case 1:
			// Server error
			http.Error(w, "Satellite temporarily unavailable", http.StatusServiceUnavailable)
			return false
		case 2:
			// Internal server error
			http.Error(w, "Internal satellite error", http.StatusInternalServerError)
			return false
		}
	}

	return true
}

// updateSensorStatus simulates sensor lifecycle
func (s *MockServer) updateSensorStatus(sensor *Sensor) {
	go func() {
		// Simulate INITIALIZING -> ACTIVE transition
		time.Sleep(time.Duration(5+rand.Intn(15)) * time.Second)

		s.mu.Lock()
		defer s.mu.Unlock()

		if sensor.Status == StatusInitializing {
			// 90% chance of becoming active, 10% chance of failure
			if rand.Float64() < 0.9 {
				sensor.Status = StatusActive
				// Start generating measurements
				measurement := 50.0 + rand.Float64()*100.0
				sensor.Measurement = &measurement
			} else {
				sensor.Status = StatusFailed
			}
		}
	}()

	// Continuously update measurements for active sensors
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			s.mu.Lock()
			if sensor.Status == StatusActive {
				// Update measurement with some drift
				newValue := *sensor.Measurement + (rand.Float64()-0.5)*10.0
				sensor.Measurement = &newValue
			}
			s.mu.Unlock()
		}
	}()
}

// handleGetSensorIDs handles GET /sensor-ids
func (s *MockServer) handleGetSensorIDs(w http.ResponseWriter, r *http.Request) {
	if !s.simulateUnreliability(w) {
		return
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	ids := make([]int, 0, len(s.sensors))
	for id := range s.sensors {
		ids = append(ids, id)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ids)
}

// handleCreateSensor handles POST /sensors
func (s *MockServer) handleCreateSensor(w http.ResponseWriter, r *http.Request) {
	if !s.simulateUnreliability(w) {
		return
	}

	var req SensorCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate frequency
	if req.Frequency < 0 {
		http.Error(w, "Invalid frequency: must be non-negative", http.StatusBadRequest)
		return
	}

	// Simulate resource limitation (10% chance)
	if s.resourceLimited && rand.Float64() < 0.1 {
		http.Error(w, "Not enough resources", http.StatusInternalServerError)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Create new sensor
	sensor := &Sensor{
		ID:          s.nextID,
		Frequency:   req.Frequency,
		Status:      StatusInitializing,
		Measurement: nil,
	}
	s.sensors[s.nextID] = sensor
	s.nextID++

	// Start background process to update sensor status
	s.updateSensorStatus(sensor)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sensor)
}

// handleGetSensor handles GET /sensors/<id>
func (s *MockServer) handleGetSensor(w http.ResponseWriter, r *http.Request) {
	if !s.simulateUnreliability(w) {
		return
	}

	// Extract sensor ID from path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid sensor ID", http.StatusBadRequest)
		return
	}

	s.mu.RLock()
	sensor, exists := s.sensors[id]
	s.mu.RUnlock()

	if !exists {
		http.Error(w, "Sensor not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sensor)
}

// Start starts the mock server
func (s *MockServer) Start(port int) {
	// Create a new ServeMux for this server instance
	mux := http.NewServeMux()

	mux.HandleFunc("/sensor-ids", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		s.handleGetSensorIDs(w, r)
	})

	mux.HandleFunc("/sensors/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		s.handleGetSensor(w, r)
	})

	mux.HandleFunc("/sensors", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		s.handleCreateSensor(w, r)
	})

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("🛰️  Mock satellite server starting on %s\n", addr)
	fmt.Printf("   Unreliability: %.0f%%\n", s.unreliability*100)
	fmt.Printf("   Max slowness: %v\n", s.slowness)
	fmt.Println()

	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
