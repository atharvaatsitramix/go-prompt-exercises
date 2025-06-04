package main

import (
	"flag"
	"fmt"
)

func runFileServer() {
	// Command line flags
	port := flag.String("port", ":8080", "Port to run the server on")
	help := flag.Bool("help", false, "Show help information")
	demo := flag.Bool("demo", false, "Run a quick demo of file reading functions")

	flag.Parse()

	if *help {
		fmt.Println("File Reader HTTP Server")
		fmt.Println("Usage:")
		fmt.Println("  go run filereader.go server.go fileserver.go -port=:8080          # Start HTTP server on port 8080")
		fmt.Println("  go run filereader.go server.go fileserver.go -demo                # Run file reading demo")
		fmt.Println("  go run filereader.go server.go fileserver.go -help                # Show this help")
		fmt.Println("\nHTTP Endpoints:")
		fmt.Println("  GET /                         # API documentation")
		fmt.Println("  GET /read?file=filename       # Read file as string")
		fmt.Println("  GET /readlines?file=filename  # Read file line by line")
		fmt.Println("  GET /saferead?file=filename   # Safe file reading")
		fmt.Println("  GET /status                   # Server status")
		return
	}

	if *demo {
		runDemo()
		return
	}

	// Start the HTTP server
	fmt.Printf("Starting File Reader HTTP Server...\n")
	fmt.Printf("Open your browser and go to: http://localhost%s\n", *port)
	RunServer(*port)
}

// runDemo demonstrates the file reading functions
func runDemo() {
	fmt.Println("=== File Reader Demo ===")

	fileReader := NewFileReader()

	// Test files to read (trying some common files)
	testFiles := []string{
		"README.md",
		"main.go",
		"go.mod",
		"../README.md", // Try parent directory
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

// To run the file server, uncomment the main function below and comment out the existing main in main.go
/*
func main() {
	runFileServer()
}
*/

// Example usage:
// go run filereader.go server.go fileserver.go -port=:8080     # Start server on port 8080
// go run filereader.go server.go fileserver.go -demo           # Run demonstration
// go run filereader.go server.go fileserver.go -help           # Show help
