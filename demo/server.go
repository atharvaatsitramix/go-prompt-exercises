package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

// Server represents our HTTP server with file reading capabilities
type Server struct {
	fileReader *FileReader
	port       string
}

// Response structures for JSON responses
type FileResponse struct {
	Success bool     `json:"success"`
	Data    string   `json:"data,omitempty"`
	Error   string   `json:"error,omitempty"`
	Lines   []string `json:"lines,omitempty"`
}

type StatusResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

// NewServer creates a new server instance
func NewServer(port string) *Server {
	return &Server{
		fileReader: NewFileReader(),
		port:       port,
	}
}

// readFileHandler handles requests to read file content as a string
func (s *Server) readFileHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get filename from query parameter
	filename := r.URL.Query().Get("file")
	if filename == "" {
		response := FileResponse{
			Success: false,
			Error:   "Missing 'file' query parameter",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Clean the filename and add basic security checks
	filename = filepath.Clean(filename)

	// Prevent directory traversal attacks
	if strings.Contains(filename, "..") {
		response := FileResponse{
			Success: false,
			Error:   "Invalid file path: directory traversal not allowed",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Try to read the file
	content, err := s.fileReader.ReadFileAsString(filename)
	if err != nil {
		response := FileResponse{
			Success: false,
			Error:   err.Error(),
		}

		// Set appropriate HTTP status code based on error type
		if strings.Contains(err.Error(), "does not exist") {
			w.WriteHeader(http.StatusNotFound)
		} else if strings.Contains(err.Error(), "not readable") || strings.Contains(err.Error(), "permission") {
			w.WriteHeader(http.StatusForbidden)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	// Success response
	response := FileResponse{
		Success: true,
		Data:    content,
	}
	json.NewEncoder(w).Encode(response)
}

// readFileLinesHandler handles requests to read file content line by line
func (s *Server) readFileLinesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filename := r.URL.Query().Get("file")
	if filename == "" {
		response := FileResponse{
			Success: false,
			Error:   "Missing 'file' query parameter",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Security checks
	filename = filepath.Clean(filename)
	if strings.Contains(filename, "..") {
		response := FileResponse{
			Success: false,
			Error:   "Invalid file path: directory traversal not allowed",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	lines, err := s.fileReader.ReadFileByLines(filename)
	if err != nil {
		response := FileResponse{
			Success: false,
			Error:   err.Error(),
		}

		if strings.Contains(err.Error(), "does not exist") {
			w.WriteHeader(http.StatusNotFound)
		} else if strings.Contains(err.Error(), "permission") {
			w.WriteHeader(http.StatusForbidden)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	response := FileResponse{
		Success: true,
		Lines:   lines,
	}
	json.NewEncoder(w).Encode(response)
}

// safeReadHandler handles requests using the safe read function
func (s *Server) safeReadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filename := r.URL.Query().Get("file")
	if filename == "" {
		response := FileResponse{
			Success: false,
			Error:   "Missing 'file' query parameter",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	content, err := s.fileReader.SafeReadFile(filename)
	if err != nil {
		response := FileResponse{
			Success: false,
			Error:   err.Error(),
		}

		if strings.Contains(err.Error(), "does not exist") {
			w.WriteHeader(http.StatusNotFound)
		} else if strings.Contains(err.Error(), "permission") || strings.Contains(err.Error(), "not allowed") {
			w.WriteHeader(http.StatusForbidden)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	response := FileResponse{
		Success: true,
		Data:    content,
	}
	json.NewEncoder(w).Encode(response)
}

// statusHandler provides server status information
func (s *Server) statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	response := StatusResponse{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Message:   "File Reader HTTP Server is running",
	}
	json.NewEncoder(w).Encode(response)
}

// homeHandler provides API documentation
func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	html := `
<!DOCTYPE html>
<html>
<head>
    <title>File Reader API</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .endpoint { background: #f5f5f5; padding: 10px; margin: 10px 0; border-left: 4px solid #007cba; }
        .method { color: #007cba; font-weight: bold; }
        code { background: #f0f0f0; padding: 2px 4px; }
    </style>
</head>
<body>
    <h1>File Reader HTTP Server API</h1>
    <p>This server provides endpoints to read files with comprehensive error handling.</p>
    
    <h2>Available Endpoints:</h2>
    
    <div class="endpoint">
        <p><span class="method">GET</span> <code>/read?file=filename</code></p>
        <p>Reads the entire file content as a string</p>
    </div>
    
    <div class="endpoint">
        <p><span class="method">GET</span> <code>/readlines?file=filename</code></p>
        <p>Reads the file content line by line, returns an array of lines</p>
    </div>
    
    <div class="endpoint">
        <p><span class="method">GET</span> <code>/saferead?file=filename</code></p>
        <p>Reads file with additional security checks and path validation</p>
    </div>
    
    <div class="endpoint">
        <p><span class="method">GET</span> <code>/status</code></p>
        <p>Returns server status and health information</p>
    </div>
    
    <h2>Example Usage:</h2>
    <p><code>curl "http://localhost:8080/read?file=README.md"</code></p>
    <p><code>curl "http://localhost:8080/readlines?file=main.go"</code></p>
    
    <h2>Error Handling:</h2>
    <p>The API handles various error conditions including:</p>
    <ul>
        <li>File not found (404)</li>
        <li>Permission denied (403)</li>
        <li>Directory traversal attempts (400)</li>
        <li>Invalid file paths (400)</li>
        <li>Internal server errors (500)</li>
    </ul>
</body>
</html>`

	fmt.Fprint(w, html)
}

// loggingMiddleware logs all incoming requests
func (s *Server) loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		log.Printf("%s %s %s %v", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
	}
}

// SetupRoutes configures all the server routes
func (s *Server) SetupRoutes() {
	http.HandleFunc("/", s.loggingMiddleware(s.homeHandler))
	http.HandleFunc("/read", s.loggingMiddleware(s.readFileHandler))
	http.HandleFunc("/readlines", s.loggingMiddleware(s.readFileLinesHandler))
	http.HandleFunc("/saferead", s.loggingMiddleware(s.safeReadHandler))
	http.HandleFunc("/status", s.loggingMiddleware(s.statusHandler))
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.SetupRoutes()

	log.Printf("Starting File Reader HTTP Server on port %s", s.port)
	log.Printf("Available endpoints:")
	log.Printf("  GET  /           - API documentation")
	log.Printf("  GET  /read       - Read file as string")
	log.Printf("  GET  /readlines  - Read file line by line")
	log.Printf("  GET  /saferead   - Safe file reading with security checks")
	log.Printf("  GET  /status     - Server status")
	log.Printf("\nExample: curl \"http://localhost%s/read?file=README.md\"", s.port)

	return http.ListenAndServe(s.port, nil)
}

// RunServer is a convenience function to create and start the server
func RunServer(port string) {
	server := NewServer(port)

	if err := server.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
