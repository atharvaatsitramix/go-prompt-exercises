package service

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
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

	// Insert the new todo in sorted order (by ID)
	todos = ts.insertTodoSorted(todos, newTodo)

	if err := ts.saveTodos(todos); err != nil {
		return nil, fmt.Errorf("failed to save todo: %w", err)
	}

	return &newTodo, nil
}

// UpdateTodo updates an existing todo by ID using binary search
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

	// Use binary search to find the todo
	index := ts.binarySearchTodoByID(todos, id)
	if index == -1 {
		return nil, fmt.Errorf("todo with ID %d not found", id)
	}

	// Update the todo
	todos[index].Contents = strings.TrimSpace(req.Contents)
	updatedTodo := todos[index]

	// Save updated todos
	if err := ts.saveTodos(todos); err != nil {
		return nil, fmt.Errorf("failed to save todo: %w", err)
	}

	return &updatedTodo, nil
}

// DeleteTodo deletes a todo by ID using binary search
func (ts *TodoService) DeleteTodo(id int) error {
	// Load existing todos
	todos, err := ts.loadTodos()
	if err != nil {
		return fmt.Errorf("failed to load todos: %w", err)
	}

	// Use binary search to find the todo
	index := ts.binarySearchTodoByID(todos, id)
	if index == -1 {
		return fmt.Errorf("todo with ID %d not found", id)
	}

	// Remove the todo at the found index
	newTodos := make([]models.Todo, 0, len(todos)-1)
	newTodos = append(newTodos, todos[:index]...)
	newTodos = append(newTodos, todos[index+1:]...)

	// Save updated todos
	if err := ts.saveTodos(newTodos); err != nil {
		return fmt.Errorf("failed to save todos: %w", err)
	}

	return nil
}

// BinarySearchTodoByID performs binary search to find a todo by ID (exported for testing)
// Returns the index of the todo if found, -1 if not found
func (ts *TodoService) BinarySearchTodoByID(todos []models.Todo, targetID int) int {
	return ts.binarySearchTodoByID(todos, targetID)
}

// binarySearchTodoByID performs binary search to find a todo by ID
// Returns the index of the todo if found, -1 if not found
func (ts *TodoService) binarySearchTodoByID(todos []models.Todo, targetID int) int {
	left, right := 0, len(todos)-1

	for left <= right {
		mid := left + (right-left)/2
		midID := todos[mid].ID

		if midID == targetID {
			return mid
		} else if midID < targetID {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return -1 // Not found
}

// insertTodoSorted inserts a todo into the slice while maintaining sorted order by ID
func (ts *TodoService) insertTodoSorted(todos []models.Todo, newTodo models.Todo) []models.Todo {
	// If empty slice, just append
	if len(todos) == 0 {
		return append(todos, newTodo)
	}

	// Find the correct position to insert using binary search
	left, right := 0, len(todos)

	for left < right {
		mid := left + (right-left)/2
		if todos[mid].ID < newTodo.ID {
			left = mid + 1
		} else {
			right = mid
		}
	}

	// Insert at the found position
	result := make([]models.Todo, len(todos)+1)
	copy(result[:left], todos[:left])
	result[left] = newTodo
	copy(result[left+1:], todos[left:])

	return result
}

// validateTodoRequest validates the todo request
func (ts *TodoService) validateTodoRequest(req models.TodoRequest) error {
	if strings.TrimSpace(req.Contents) == "" {
		return fmt.Errorf("contents cannot be empty")
	}
	return nil
}

// loadTodos reads todos from JSON file and ensures they are sorted by ID
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

	// Ensure todos are sorted by ID for binary search to work
	sort.Slice(todos, func(i, j int) bool {
		return todos[i].ID < todos[j].ID
	})

	return todos, nil
}

// saveTodos writes todos to JSON file (todos should already be sorted)
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
