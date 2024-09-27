package main

import (
	"fmt"
	"math"
	"crypto/rand"
	"os"
	"os/exec"
	"sync"
	"encoding/json"
//	"bytes"
//	"encoding/gob"
)

func blockingIO() ([]byte, error) {
	data := make([]byte, 100)
	_, err := rand.Read(data)  // This fills the data with random bytes
	if err != nil {
		return nil, err
	}
	return data, nil
}

// cpuBound performs a CPU-intensive task
func cpuBound() int {
	sum := 0
	for i := 0; i < int(math.Pow(10, 7)); i++ {
		sum += i * i
	}
	return sum
}

// runCpuBoundInSubprocess runs the cpuBound function in a separate subprocess with isolated environment variables
func runCpuBoundInSubprocess() error {
	// Create a nested map (equivalent to a complex JSON structure)
	jsonData := map[string]interface{}{
		"task": "example",
		"details": map[string]interface{}{
			"user":  "Giulio",
			"items": []string{"item1", "item2", "item3"},
		},
	}

	// Serialize the data to JSON
	data, err := json.Marshal(jsonData)
	if err != nil {
		return fmt.Errorf("Serialization error: %v", err)
	}


	// Register the type
	/*gob.Register(map[string]interface{}{})

	// Serialize the data with gob
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("Error encoding:", err)
	}
	data := buffer.Bytes()
	*/

	// Write the JSON data to a temporary file
	tempFile, err := os.CreateTemp("", "data.json")
	if err != nil {
		return fmt.Errorf("Error creating temp file: %v", err)

	}
	defer os.Remove(tempFile.Name()) // Clean up

	if err := os.WriteFile(tempFile.Name(), data, 0644); err != nil {
		return fmt.Errorf("Error writing temp file: %v", err)
	}

	cmd := exec.Command(os.Args[0], tempFile.Name()) // Run the same program and pass the JSON data filename as an argument
	cmd.Env = append(os.Environ(), "CUSTOM_ENV_VAR=custom_value") // Set custom environment variables

	output, err := cmd.CombinedOutput() // Capture output from the subprocess
	if err != nil {
		return fmt.Errorf("subprocess error: %v, output: %s", err, string(output))
	}
	fmt.Println("Subprocess output:", string(output))
	return nil
}

func main() {
	var wg sync.WaitGroup

	// 1. Run blockingIO in a goroutine (default thread pool equivalent)
	wg.Add(1)
	go func() {
		defer wg.Done()
		result, err := blockingIO()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Default thread pool:", result)
	}()

	// 2. Run blockingIO in a custom thread pool (simulated using goroutines and sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()
		result, err := blockingIO()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Custom thread pool:", result)
	}()

	// 3. Run cpuBound in a separate subprocess with isolated environment
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := runCpuBoundInSubprocess()
		if err != nil {
			fmt.Println("Error:", err)
		}
	}()

	// Wait for all goroutines to finish
	wg.Wait()
	// Get the value of the environment variable
	customEnvVar := os.Getenv("CUSTOM_ENV_VAR")

	// Check if the environment variable is set
	if customEnvVar == "" {
		fmt.Println("In Main => CUSTOM_ENV_VAR is not set")
	} else {
		fmt.Println("In Main => CUSTOM_ENV_VAR:", customEnvVar)
	}
}

func readJSONFromFile(filename string) (map[string]interface{}, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Error reading file: %v", err)
	}
	var result map[string]interface{}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("Error unmarshalling JSON: %v", err)
	}

	/*gob.Register(map[string]interface{}{})
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)

	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("Error decoding:", err)
	}*/

	return result, nil
}

// cpu function is executed when the subprocess runs the "cpu" argument
func cpu() {
	result := cpuBound()
	fmt.Printf("CPU-bound result: %d\n", result)

	// Get the value of the environment variable
	customEnvVar := os.Getenv("CUSTOM_ENV_VAR")

	// Check if the environment variable is set
	if customEnvVar == "" {
		fmt.Println("In Subprocess => CUSTOM_ENV_VAR is not set")
	} else {
		fmt.Println("In Subprocess => CUSTOM_ENV_VAR:", customEnvVar)
	}

	// Deserialize the JSON filename passed as an argument
	filename := os.Args[1]
	receivedData, err := readJSONFromFile(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Process the data
	fmt.Println("Received data in child process:", receivedData)

	fmt.Println("User: ", receivedData["details"].(map[string]interface{})["user"].(string))

	// Perform some operation and print a result
	fmt.Println("Processing complete")

}

// Entry point - checks if we are in subprocess mode
func init() {
	if len(os.Args) > 1  { //&& os.Args[1] == "cpu" {
		cpu() // Execute CPU-bound task in subprocess
		os.Exit(0) // Ensure the subprocess exits after completing its work
	}
}
