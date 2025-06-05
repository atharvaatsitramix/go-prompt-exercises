package models

// Todo represents a single todo item
type Todo struct {
	ID       int    `json:"id"`
	Contents string `json:"contents"`
}

// TodoRequest represents the input structure for creating/updating todos
type TodoRequest struct {
	ID       int    `json:"id,omitempty"`
	Contents string `json:"contents"`
}

// TodoResponse represents the API response structure
type TodoResponse struct {
	Message string `json:"message"`
	Data    []Todo `json:"data,omitempty"`
}
