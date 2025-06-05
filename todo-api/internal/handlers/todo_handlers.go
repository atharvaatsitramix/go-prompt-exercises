package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"todo-backend/internal/models"
	"todo-backend/internal/service"

	"github.com/gorilla/mux"
)

// TodoHandler handles HTTP requests for todos
type TodoHandler struct {
	todoService *service.TodoService
}

// NewTodoHandler creates a new todo handler
func NewTodoHandler(todoService *service.TodoService) *TodoHandler {
	return &TodoHandler{
		todoService: todoService,
	}
}

// GetTodos handles GET /todos
func (th *TodoHandler) GetTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := th.todoService.GetAllTodos()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load todos: %v", err), http.StatusInternalServerError)
		return
	}

	response := models.TodoResponse{
		Message: "Todos retrieved successfully",
		Data:    todos,
	}

	th.writeJSONResponse(w, http.StatusOK, response)
}

// CreateTodo handles POST /todos
func (th *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req models.TodoRequest

	if err := th.readJSONRequest(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newTodo, err := th.todoService.CreateTodo(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := models.TodoResponse{
		Message: "Todo created successfully",
		Data:    []models.Todo{*newTodo},
	}

	th.writeJSONResponse(w, http.StatusCreated, response)
}

// UpdateTodo handles PUT /todos/{id}
func (th *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	var req models.TodoRequest
	if err := th.readJSONRequest(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedTodo, err := th.todoService.UpdateTodo(id, req)
	if err != nil {
		if err.Error() == fmt.Sprintf("todo with ID %d not found", id) {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := models.TodoResponse{
		Message: "Todo updated successfully",
		Data:    []models.Todo{*updatedTodo},
	}

	th.writeJSONResponse(w, http.StatusOK, response)
}

// DeleteTodo handles DELETE /todos/{id}
func (th *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	err = th.todoService.DeleteTodo(id)
	if err != nil {
		if err.Error() == fmt.Sprintf("todo with ID %d not found", id) {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to delete todo: %v", err), http.StatusInternalServerError)
		return
	}

	response := models.TodoResponse{
		Message: "Todo deleted successfully",
	}

	th.writeJSONResponse(w, http.StatusOK, response)
}

// readJSONRequest reads and unmarshals JSON request body
func (th *TodoHandler) readJSONRequest(r *http.Request, v interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body")
	}

	err = json.Unmarshal(body, v)
	if err != nil {
		return fmt.Errorf("invalid JSON format")
	}

	return nil
}

// writeJSONResponse writes a JSON response
func (th *TodoHandler) writeJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
