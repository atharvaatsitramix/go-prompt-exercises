package main

import (
	"fmt"
	"log"
	"net/http"

	"todo-backend/internal/handlers"
	"todo-backend/internal/service"

	"github.com/gorilla/mux"
)

// corsMiddleware adds CORS headers
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Initialize service layer
	todoService := service.NewTodoService("todos.json")

	// Initialize handler layer
	todoHandler := handlers.NewTodoHandler(todoService)

	// Setup router
	r := mux.NewRouter()

	// Apply CORS middleware
	r.Use(corsMiddleware)

	// API routes
	r.HandleFunc("/todos", todoHandler.GetTodos).Methods("GET")
	r.HandleFunc("/todos", todoHandler.CreateTodo).Methods("POST")
	r.HandleFunc("/todos/{id}", todoHandler.UpdateTodo).Methods("PUT")
	r.HandleFunc("/todos/{id}", todoHandler.DeleteTodo).Methods("DELETE")

	// Handle preflight requests
	r.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {}).Methods("OPTIONS")
	r.HandleFunc("/todos/{id}", func(w http.ResponseWriter, r *http.Request) {}).Methods("OPTIONS")

	port := ":8080"
	fmt.Printf("Todo Backend Server starting on port %s\n", port)
	fmt.Println("Available endpoints:")
	fmt.Println("  GET    /todos       - Get all todos")
	fmt.Println("  POST   /todos       - Create a new todo")
	fmt.Println("  PUT    /todos/{id}  - Update a todo")
	fmt.Println("  DELETE /todos/{id}  - Delete a todo")

	log.Fatal(http.ListenAndServe(port, r))
}
