package service

import (
	"os"
	"testing"
	"todo-backend/internal/models"
)

func TestBinarySearchTodoByID(t *testing.T) {
	service := NewTodoService("test_todos.json")

	// Test data - sorted by ID
	todos := []models.Todo{
		{ID: 1, Contents: "First todo"},
		{ID: 3, Contents: "Third todo"},
		{ID: 5, Contents: "Fifth todo"},
		{ID: 7, Contents: "Seventh todo"},
		{ID: 9, Contents: "Ninth todo"},
	}

	tests := []struct {
		targetID      int
		expectedIndex int
		description   string
	}{
		{1, 0, "Find first element"},
		{9, 4, "Find last element"},
		{5, 2, "Find middle element"},
		{3, 1, "Find element in first half"},
		{7, 3, "Find element in second half"},
		{2, -1, "Element not found (between existing)"},
		{0, -1, "Element not found (before first)"},
		{10, -1, "Element not found (after last)"},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			result := service.binarySearchTodoByID(todos, test.targetID)
			if result != test.expectedIndex {
				t.Errorf("Expected index %d for ID %d, got %d", test.expectedIndex, test.targetID, result)
			}
		})
	}
}

func TestInsertTodoSorted(t *testing.T) {
	service := NewTodoService("test_todos.json")

	// Start with sorted todos
	todos := []models.Todo{
		{ID: 1, Contents: "First todo"},
		{ID: 3, Contents: "Third todo"},
		{ID: 7, Contents: "Seventh todo"},
	}

	tests := []struct {
		newTodo     models.Todo
		expectedPos int
		description string
	}{
		{models.Todo{ID: 0, Contents: "Zero todo"}, 0, "Insert at beginning"},
		{models.Todo{ID: 2, Contents: "Second todo"}, 1, "Insert in middle (between 1 and 3)"},
		{models.Todo{ID: 5, Contents: "Fifth todo"}, 2, "Insert in middle (between 3 and 7)"},
		{models.Todo{ID: 10, Contents: "Tenth todo"}, 3, "Insert at end"},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			result := service.insertTodoSorted(todos, test.newTodo)

			// Check if the new todo is at the expected position
			if result[test.expectedPos].ID != test.newTodo.ID {
				t.Errorf("Expected todo with ID %d at position %d, got ID %d",
					test.newTodo.ID, test.expectedPos, result[test.expectedPos].ID)
			}

			// Check if the slice is still sorted
			for i := 1; i < len(result); i++ {
				if result[i-1].ID >= result[i].ID {
					t.Errorf("Slice is not sorted after insertion: ID %d >= ID %d at positions %d, %d",
						result[i-1].ID, result[i].ID, i-1, i)
				}
			}
		})
	}
}

func TestCRUDOperationsWithBinarySearch(t *testing.T) {
	// Use a temporary file for testing
	testFile := "test_crud_todos.json"
	defer os.Remove(testFile) // Clean up after test

	service := NewTodoService(testFile)

	// Test Create
	req1 := models.TodoRequest{Contents: "First todo"}
	todo1, err := service.CreateTodo(req1)
	if err != nil {
		t.Fatalf("Failed to create todo: %v", err)
	}
	if todo1.ID != 1 {
		t.Errorf("Expected ID 1, got %d", todo1.ID)
	}

	req2 := models.TodoRequest{Contents: "Second todo"}
	todo2, err := service.CreateTodo(req2)
	if err != nil {
		t.Fatalf("Failed to create second todo: %v", err)
	}
	if todo2.ID != 2 {
		t.Errorf("Expected ID 2, got %d", todo2.ID)
	}

	// Test Read (GetAllTodos)
	todos, err := service.GetAllTodos()
	if err != nil {
		t.Fatalf("Failed to get todos: %v", err)
	}
	if len(todos) != 2 {
		t.Errorf("Expected 2 todos, got %d", len(todos))
	}

	// Verify todos are sorted
	if todos[0].ID > todos[1].ID {
		t.Error("Todos are not sorted by ID")
	}

	// Test Update using binary search
	updateReq := models.TodoRequest{Contents: "Updated first todo"}
	updatedTodo, err := service.UpdateTodo(1, updateReq)
	if err != nil {
		t.Fatalf("Failed to update todo: %v", err)
	}
	if updatedTodo.Contents != "Updated first todo" {
		t.Errorf("Expected updated content, got: %s", updatedTodo.Contents)
	}

	// Test Update with non-existent ID
	_, err = service.UpdateTodo(999, updateReq)
	if err == nil {
		t.Error("Expected error when updating non-existent todo")
	}

	// Test Delete using binary search
	err = service.DeleteTodo(1)
	if err != nil {
		t.Fatalf("Failed to delete todo: %v", err)
	}

	// Verify deletion
	todos, err = service.GetAllTodos()
	if err != nil {
		t.Fatalf("Failed to get todos after deletion: %v", err)
	}
	if len(todos) != 1 {
		t.Errorf("Expected 1 todo after deletion, got %d", len(todos))
	}
	if todos[0].ID != 2 {
		t.Errorf("Expected remaining todo to have ID 2, got %d", todos[0].ID)
	}

	// Test Delete with non-existent ID
	err = service.DeleteTodo(999)
	if err == nil {
		t.Error("Expected error when deleting non-existent todo")
	}
}

func TestEmptyTodosOperations(t *testing.T) {
	service := NewTodoService("empty_test.json")
	defer os.Remove("empty_test.json")

	// Test binary search on empty slice
	result := service.binarySearchTodoByID([]models.Todo{}, 1)
	if result != -1 {
		t.Errorf("Expected -1 for empty slice, got %d", result)
	}

	// Test insert into empty slice
	newTodo := models.Todo{ID: 1, Contents: "First todo"}
	result_slice := service.insertTodoSorted([]models.Todo{}, newTodo)
	if len(result_slice) != 1 || result_slice[0].ID != 1 {
		t.Error("Failed to insert into empty slice")
	}
}
