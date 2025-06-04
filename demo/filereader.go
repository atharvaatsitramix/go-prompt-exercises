package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// FileReader provides methods for reading files with comprehensive error handling
type FileReader struct{}

// NewFileReader creates a new FileReader instance
func NewFileReader() *FileReader {
	return &FileReader{}
}

// ReadFileAsString reads the entire file content as a string
// This function demonstrates various error handling scenarios when reading files
func (fr *FileReader) ReadFileAsString(filename string) (string, error) {
	// Check if file exists and get file info
	fileInfo, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("file does not exist: %s", filename)
		}
		return "", fmt.Errorf("cannot access file: %v", err)
	}

	// Check if it's actually a file (not a directory)
	if fileInfo.IsDir() {
		return "", fmt.Errorf("path is a directory, not a file: %s", filename)
	}

	// Check file permissions
	if fileInfo.Mode().Perm()&0400 == 0 {
		return "", fmt.Errorf("file is not readable: %s", filename)
	}

	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close() // Always close the file when done

	// Read the file content
	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file content: %v", err)
	}

	return string(content), nil
}

// ReadFileByLines reads a file line by line and returns a slice of strings
// This approach is memory-efficient for large files
func (fr *FileReader) ReadFileByLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	// Read line by line
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return lines, nil
}

// SafeReadFile provides additional safety checks and error handling
func (fr *FileReader) SafeReadFile(filename string) (string, error) {
	// Clean and validate the file path
	cleanPath := filepath.Clean(filename)

	// Get absolute path
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %v", err)
	}

	// Additional security check: ensure we're not reading outside allowed directories
	// This is a basic example - in production, you'd have more sophisticated path validation
	if filepath.Dir(absPath) == "/" {
		return "", fmt.Errorf("reading from root directory is not allowed")
	}

	return fr.ReadFileAsString(absPath)
}

/*
ERROR HANDLING GUIDE FOR FILE OPERATIONS:

1. FILE NOT FOUND ERRORS:
   - Use os.IsNotExist(err) to check if a file doesn't exist
   - Provide clear error messages indicating the missing file path
   - Consider creating default files or prompting user for alternative paths

2. PERMISSION ERRORS:
   - Check file permissions before attempting to read
   - Use os.IsPermission(err) to identify permission-related errors
   - Consider running with appropriate privileges or changing file permissions

3. DIRECTORY VS FILE ERRORS:
   - Always check if the path points to a file or directory using FileInfo.IsDir()
   - Handle cases where a directory is provided instead of a file

4. LARGE FILE HANDLING:
   - For large files, consider reading in chunks or line-by-line to avoid memory issues
   - Use bufio.Scanner for efficient line-by-line reading
   - Implement progress indicators for very large files

5. CONCURRENT ACCESS:
   - Be aware of other processes modifying the file while reading
   - Consider file locking mechanisms for critical applications
   - Handle cases where files might be deleted during reading

6. ENCODING ISSUES:
   - Be aware of file encoding (UTF-8, ASCII, etc.)
   - Consider using appropriate readers for different encodings
   - Handle BOM (Byte Order Mark) if present

7. NETWORK FILES:
   - Network-mounted files may have additional latency and failure modes
   - Implement retry mechanisms for network-related failures
   - Consider timeout settings for network file operations

8. MEMORY MANAGEMENT:
   - Always use defer file.Close() to ensure files are closed
   - Be mindful of memory usage when reading large files
   - Consider streaming approaches for processing large datasets

BEST PRACTICES:
- Always close files using defer
- Check errors at every step
- Provide meaningful error messages
- Log errors appropriately
- Consider using context for cancellation in long-running operations
- Validate input paths to prevent security issues
- Test error conditions thoroughly
*/
