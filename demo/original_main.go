package main

import (
	"flag"
	"fmt"
	"strings"
)

func main() {
	// Command line flags for file server functionality
	port := flag.String("port", ":8080", "Port to run the server on")
	help := flag.Bool("help", false, "Show help information")
	demo := flag.Bool("demo", false, "Run a quick demo of file reading functions")
	server := flag.Bool("server", false, "Start the HTTP file server")
	fileReader := flag.Bool("filereader", false, "Run original demo with file reading capabilities")

	flag.Parse()

	if *help {
		fmt.Println("Demo Application with File Reader HTTP Server")
		fmt.Println("Usage:")
		fmt.Println("  go run . -server -port=:8080     # Start HTTP server on port 8080")
		fmt.Println("  go run . -demo                   # Run file reading demo")
		fmt.Println("  go run . -filereader             # Run original demo + file reading")
		fmt.Println("  go run . -help                   # Show this help")
		fmt.Println("  go run .                         # Run original demo (default)")
		fmt.Println("\nHTTP Endpoints (when server mode is active):")
		fmt.Println("  GET /                            # API documentation")
		fmt.Println("  GET /read?file=filename          # Read file as string")
		fmt.Println("  GET /readlines?file=filename     # Read file line by line")
		fmt.Println("  GET /saferead?file=filename      # Safe file reading")
		fmt.Println("  GET /status                      # Server status")
		return
	}

	if *server {
		// Start the HTTP server
		fmt.Printf("Starting File Reader HTTP Server...\n")
		fmt.Printf("Open your browser and go to: http://localhost%s\n", *port)
		RunServer(*port)
		return
	}

	if *demo {
		// Run file reading demo
		runFileReaderDemo()
		return
	}

	if *fileReader {
		// Run original demo + file reading capabilities
		fmt.Println("=== Original Demo with File Reading ===")
		runOriginalDemo()
		fmt.Println("\n" + strings.Repeat("=", 50))
		runFileReaderDemo()
		return
	}

	// Default behavior: run original demo
	runOriginalDemo()
}

// runOriginalDemo runs the original demo functionality
func runOriginalDemo() {
	fmt.Println("Starting the demo application...")
	SayHello()

	// Call the Add function with different examples
	fmt.Println("\nTesting the Add function:")
	Add(5, 6)  // Should output: the addition is 11
	Add(-5, 6) // Should output: the addition is 1
}

// runFileReaderDemo demonstrates the file reading functions
func runFileReaderDemo() {
	fmt.Println("=== File Reader Demo ===")

	fileReader := NewFileReader()

	// Test files to read (trying some common files)
	testFiles := []string{
		"sample.txt", // Our test file
		"go.mod",
		"../README.md",     // Try parent directory
		"original_main.go", // This file itself
	}

	for _, filename := range testFiles {
		fmt.Printf("\n--- Testing file: %s ---\n", filename)

		// Test ReadFileAsString
		content, err := fileReader.ReadFileAsString(filename)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
		} else {
			// Show first 200 characters
			if len(content) > 200 {
				fmt.Printf("Content (first 200 chars): %s...\n", content[:200])
			} else {
				fmt.Printf("Content: %s\n", content)
			}
		}

		// Test ReadFileByLines
		lines, err := fileReader.ReadFileByLines(filename)
		if err != nil {
			fmt.Printf("Error reading lines: %v\n", err)
		} else {
			fmt.Printf("Total lines: %d\n", len(lines))
			if len(lines) > 0 {
				fmt.Printf("First line: %s\n", lines[0])
			}
		}
	}

	// Test error conditions
	fmt.Printf("\n--- Testing Error Conditions ---\n")

	// Test non-existent file
	_, err := fileReader.ReadFileAsString("nonexistent.txt")
	fmt.Printf("Non-existent file error: %v\n", err)

	// Test directory instead of file
	_, err = fileReader.ReadFileAsString(".")
	fmt.Printf("Directory error: %v\n", err)

	// Test SafeReadFile
	fmt.Printf("\n--- Testing SafeReadFile ---\n")
	content, err := fileReader.SafeReadFile("go.mod")
	if err != nil {
		fmt.Printf("SafeReadFile error: %v\n", err)
	} else {
		fmt.Printf("SafeReadFile success: %d characters read\n", len(content))
	}
}

// Example usage:
// go run .                         # Run original demo (default)
// go run . -server -port=:8080     # Start server on port 8080
// go run . -demo                   # Run file reading demonstration
// go run . -filereader             # Run both original demo and file reading demo
// go run . -help                   # Show help
