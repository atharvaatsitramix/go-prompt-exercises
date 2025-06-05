package service

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"todo-backend/internal/models"
)

// TodoService handles todo business logic and file storage
type TodoService struct {
	filename string
	mutex    sync.RWMutex
}

// NewTodoService creates a new todo service
func NewTodoService(filename string) *TodoService {
	return &TodoService{
		filename: filename,
	}
}

// GetAllTodos retrieves all todos from storage
func (ts *TodoService) GetAllTodos() ([]models.Todo, error) {
	return ts.loadTodos()
}

// CreateTodo creates a new todo with auto-generated ID
func (ts *TodoService) CreateTodo(req models.TodoRequest) (*models.Todo, error) {
	// Validate request
	if err := ts.validateTodoRequest(req); err != nil {
		return nil, err
	}

	// Load existing todos
	todos, err := ts.loadTodos()
	if err != nil {
		return nil, fmt.Errorf("failed to load todos: %w", err)
	}

	// Create new todo
	newTodo := models.Todo{
		ID:       ts.getNextID(todos),
		Contents: strings.TrimSpace(req.Contents),
	}

	// Add to collection and save
	todos = append(todos, newTodo)
	if err := ts.saveTodos(todos); err != nil {
		return nil, fmt.Errorf("failed to save todo: %w", err)
	}

	return &newTodo, nil
}

// UpdateTodo updates an existing todo by ID
func (ts *TodoService) UpdateTodo(id int, req models.TodoRequest) (*models.Todo, error) {
	// Validate request
	if err := ts.validateTodoRequest(req); err != nil {
		return nil, err
	}

	// Load existing todos
	todos, err := ts.loadTodos()
	if err != nil {
		return nil, fmt.Errorf("failed to load todos: %w", err)
	}

	// Find and update todo
	found := false
	var updatedTodo models.Todo
	for i, todo := range todos {
		if todo.ID == id {
			todos[i].Contents = strings.TrimSpace(req.Contents)
			updatedTodo = todos[i]
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("todo with ID %d not found", id)
	}

	// Save updated todos
	if err := ts.saveTodos(todos); err != nil {
		return nil, fmt.Errorf("failed to save todo: %w", err)
	}

	return &updatedTodo, nil
}

// DeleteTodo deletes a todo by ID
func (ts *TodoService) DeleteTodo(id int) error {
	// Load existing todos
	todos, err := ts.loadTodos()
	if err != nil {
		return fmt.Errorf("failed to load todos: %w", err)
	}

	// Filter out the todo to delete
	found := false
	var newTodos []models.Todo
	for _, todo := range todos {
		if todo.ID != id {
			newTodos = append(newTodos, todo)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("todo with ID %d not found", id)
	}

	// Save updated todos
	if err := ts.saveTodos(newTodos); err != nil {
		return fmt.Errorf("failed to save todos: %w", err)
	}

	return nil
}

// validateTodoRequest validates the todo request
func (ts *TodoService) validateTodoRequest(req models.TodoRequest) error {
	if strings.TrimSpace(req.Contents) == "" {
		return fmt.Errorf("contents cannot be empty")
	}
	return nil
}

// loadTodos reads todos from JSON file
func (ts *TodoService) loadTodos() ([]models.Todo, error) {
	ts.mutex.RLock()
	defer ts.mutex.RUnlock()

	file, err := os.Open(ts.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.Todo{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var todos []models.Todo
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&todos)
	if err != nil {
		return nil, err
	}

	return todos, nil
}

// saveTodos writes todos to JSON file
func (ts *TodoService) saveTodos(todos []models.Todo) error {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	file, err := os.Create(ts.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(todos)
}

// getNextID returns the next available ID
func (ts *TodoService) getNextID(todos []models.Todo) int {
	maxID := 0
	for _, todo := range todos {
		if todo.ID > maxID {
			maxID = todo.ID
		}
	}
	return maxID + 1
}
