package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

const defaultURL = "http://localhost:8080"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	// Define flags for all commands
	createCmd := flag.NewFlagSet("create-sensor", flag.ExitOnError)
	createFrequency := createCmd.Int("frequency", 0, "Sensor frequency (required)")

	getCmd := flag.NewFlagSet("get-sensor", flag.ExitOnError)
	getID := getCmd.Int("id", 0, "Sensor ID (required)")

	listCmd := flag.NewFlagSet("list-sensors", flag.ExitOnError)

	serverCmd := flag.NewFlagSet("server", flag.ExitOnError)
	serverPort := serverCmd.Int("port", 8080, "Server port")
	serverUnreliability := serverCmd.Float64("unreliability", 0.2, "Server unreliability (0.0-1.0)")
	serverSlowness := serverCmd.Duration("slowness", 500*time.Millisecond, "Max server delay")

	demoCmd := flag.NewFlagSet("demo", flag.ExitOnError)
	demoPort := demoCmd.Int("port", 8080, "Server port")
	demoUnreliability := demoCmd.Float64("unreliability", 0.2, "Server unreliability (0.0-1.0)")
	demoSlowness := demoCmd.Duration("slowness", 500*time.Millisecond, "Max server delay")

	switch command {
	case "create-sensor":
		createCmd.Parse(os.Args[2:])
		cmdCreateSensor(*createFrequency)

	case "get-sensor":
		getCmd.Parse(os.Args[2:])
		cmdGetSensor(*getID)

	case "list-sensors":
		listCmd.Parse(os.Args[2:])
		cmdListSensors()

	case "server":
		serverCmd.Parse(os.Args[2:])
		runServer(*serverPort, *serverUnreliability, *serverSlowness)

	case "demo":
		demoCmd.Parse(os.Args[2:])
		runDemo(*demoPort, *demoUnreliability, *demoSlowness)

	case "help", "--help", "-h":
		printUsage()

	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func cmdCreateSensor(frequency int) {
	if frequency == 0 {
		fmt.Println("Error: --frequency is required")
		os.Exit(1)
	}

	client := NewSatelliteInterface(defaultURL)
	sensor, err := client.CreateSensor(frequency)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Sensor created successfully!\n")
	fmt.Printf("  ID: %d\n", sensor.ID)
	fmt.Printf("  Frequency: %d\n", sensor.Frequency)
	fmt.Printf("  Status: %s\n", sensor.Status)
	if sensor.Measurement != nil {
		fmt.Printf("  Measurement: %.3f\n", *sensor.Measurement)
	}
}

func cmdGetSensor(id int) {
	if id == 0 {
		fmt.Println("Error: --id is required")
		os.Exit(1)
	}

	client := NewSatelliteInterface(defaultURL)
	sensor, err := client.GetSensor(id)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Sensor details:\n")
	fmt.Printf("  ID: %d\n", sensor.ID)
	fmt.Printf("  Frequency: %d\n", sensor.Frequency)
	fmt.Printf("  Status: %s\n", sensor.Status)
	if sensor.Measurement != nil {
		fmt.Printf("  Measurement: %.3f\n", *sensor.Measurement)
	} else {
		fmt.Printf("  Measurement: null\n")
	}
}

func cmdListSensors() {
	client := NewSatelliteInterface(defaultURL)
	ids, err := client.GetSensorIDs()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if len(ids) == 0 {
		fmt.Println("No sensors found")
	} else {
		fmt.Printf("Found %d sensor(s):\n", len(ids))
		for _, id := range ids {
			fmt.Printf("  - %d\n", id)
		}
	}
}

func runServer(port int, unreliability float64, slowness time.Duration) {
	fmt.Println("=== Satellite Mock Server ===")
	server := NewMockServer(unreliability, slowness)
	server.Start(port)
}

func runClient(baseURL string) {
	fmt.Println("=== Satellite API Client ===")
	fmt.Printf("Connecting to: %s\n", baseURL)

	interfaceClient := NewSatelliteInterface(baseURL)

	fmt.Println("\n" + "==================================================")
	fmt.Println("🚀 Satellite Interface Demo")
	fmt.Println("==================================================")

	// 1. Attempt to create a sensor
	newSensor, err := interfaceClient.CreateSensor(42)
	if err != nil {
		handleError("CreateSensor", err)
	} else {
		fmt.Printf("\n✅ Sensor created successfully: %+v\n", newSensor)
	}

	// Exit if sensor creation failed
	if err != nil {
		return
	}

	// 2. Attempt to get all sensor IDs
	ids, err := interfaceClient.GetSensorIDs()
	if err != nil {
		handleError("GetSensorIDs", err)
	} else {
		fmt.Printf("\n✅ All sensor IDs retrieved: %v\n", ids)
	}

	// 3. Get sensor details (will automatically wait for ACTIVE status)
	fmt.Printf("\n⏳ Retrieving sensor %d (will wait for ACTIVE status)...\n", newSensor.ID)
	sensor, err := interfaceClient.GetSensor(newSensor.ID)
	if err != nil {
		handleError("GetSensor", err)
	} else {
		fmt.Printf("\n✅ Sensor is ACTIVE!\n")
		fmt.Printf("   Details: %+v\n", sensor)
		if sensor.Measurement != nil {
			fmt.Printf("   Measurement: %.3f\n", *sensor.Measurement)
		}
	}
}

func runDemo(port int, unreliability float64, slowness time.Duration) {
	fmt.Println("=== Integrated Demo: Unreliable Satellite Simulation ===")

	// Start mock server in background
	fmt.Printf("Starting mock satellite server on port %d...\n", port)
	fmt.Printf("  Unreliability: %.0f%%\n", unreliability*100)
	fmt.Printf("  Max slowness: %v\n\n", slowness)

	server := NewMockServer(unreliability, slowness)
	go server.Start(port)

	// Wait for server to start
	time.Sleep(500 * time.Millisecond)

	// Run client against local server
	baseURL := fmt.Sprintf("http://localhost:%d", port)
	runClient(baseURL)

	fmt.Println("\n=== Demo Complete ===")
	fmt.Println("Server is still running. Press Ctrl+C to exit.")

	// Keep running
	select {}
}

func printUsage() {
	fmt.Println("Satellite Sensor API Client")
	fmt.Println("\nUsage:")
	fmt.Println("  satellite <command> [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  create-sensor    Create a new sensor")
	fmt.Println("  get-sensor       Get sensor details by ID")
	fmt.Println("  list-sensors     List all sensor IDs")
	fmt.Println("  server           Start mock satellite server")
	fmt.Println("  demo             Run integrated demo")
	fmt.Println("  help             Show this help message")
	fmt.Println("\nExamples:")
	fmt.Println("  ./satellite create-sensor --frequency=42")
	fmt.Println("  ./satellite get-sensor --id=1")
	fmt.Println("  ./satellite list-sensors")
	fmt.Println("  ./satellite server --port=8080 --unreliability=0.2")
	fmt.Println("  ./satellite demo")
	fmt.Printf("\nNote: Client commands connect to %s by default\n", defaultURL)
}
