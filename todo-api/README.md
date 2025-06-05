# Todo API - Clean Architecture

A well-structured REST API backend for a todo list application built with Go. This project follows clean architecture principles with clear separation of concerns between business logic, HTTP handling, and data models.

## Architecture

The project follows a layered architecture pattern:

```
todo-api/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── models/
│   │   └── todo.go              # Data structures and DTOs
│   ├── service/
│   │   └── todo_service.go      # Business logic layer
│   └── handlers/
│       └── todo_handlers.go     # HTTP handlers layer
├── bin/                         # Compiled binaries
├── go.mod                       # Go module file
├── go.sum                       # Go dependencies
├── todos.json                   # Data storage (created automatically)
└── README.md                    # This documentation
```

### Layer Responsibilities

- **`cmd/server`**: Application entry point, server setup, and dependency injection
- **`internal/models`**: Data structures, request/response models
- **`internal/service`**: Business logic, validation, and data persistence
- **`internal/handlers`**: HTTP request handling, routing, and response formatting

## Features

- ✅ **Clean Architecture**: Separation of concerns with distinct layers
- ✅ **CRUD Operations**: Create, Read, Update, Delete todos
- ✅ **JSON Storage**: File-based persistence with thread-safe operations
- ✅ **Input Validation**: Comprehensive request validation
- ✅ **Error Handling**: Proper HTTP status codes and error messages
- ✅ **CORS Support**: Ready for frontend integration
- ✅ **Dependency Injection**: Loosely coupled components

## Quick Start

### 1. Build the application
```bash
go build -o bin/todo-server ./cmd/server
```

### 2. Run the server
```bash
./bin/todo-server
```

### 3. Alternative: Run directly with Go
```bash
go run ./cmd/server
```

The server will start on `http://localhost:8080`

## API Endpoints

### 1. Get All Todos
**GET** `/todos`

**Response:**
```json
{
  "message": "Todos retrieved successfully",
  "data": [
    {
      "id": 1,
      "contents": "Buy groceries"
    }
  ]
}
```

### 2. Create a New Todo
**POST** `/todos`

**Request:**
```json
{
  "contents": "Your todo content here"
}
```

**Response:**
```json
{
  "message": "Todo created successfully",
  "data": [
    {
      "id": 3,
      "contents": "Your todo content here"
    }
  ]
}
```

### 3. Update a Todo
**PUT** `/todos/{id}`

**Request:**
```json
{
  "contents": "Updated todo content"
}
```

**Response:**
```json
{
  "message": "Todo updated successfully",
  "data": [
    {
      "id": 1,
      "contents": "Updated todo content"
    }
  ]
}
```

### 4. Delete a Todo
**DELETE** `/todos/{id}`

**Response:**
```json
{
  "message": "Todo deleted successfully"
}
```

## Code Examples

### cURL Examples

```bash
# Get all todos
curl -X GET http://localhost:8080/todos

# Create a new todo
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"contents": "Learn Go programming"}'

# Update a todo
curl -X PUT http://localhost:8080/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"contents": "Learn advanced Go"}'

# Delete a todo
curl -X DELETE http://localhost:8080/todos/1
```

### Flutter Integration

```dart
import 'dart:convert';
import 'package:http/http.dart' as http;

class TodoService {
  static const String baseUrl = 'http://localhost:8080';

  // Get all todos
  static Future<List<Todo>> getAllTodos() async {
    final response = await http.get(Uri.parse('$baseUrl/todos'));
    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      return (data['data'] as List).map((item) => Todo.fromJson(item)).toList();
    }
    throw Exception('Failed to load todos');
  }

  // Create a new todo
  static Future<Todo> createTodo(String contents) async {
    final response = await http.post(
      Uri.parse('$baseUrl/todos'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({'contents': contents}),
    );
    if (response.statusCode == 201) {
      final data = jsonDecode(response.body);
      return Todo.fromJson(data['data'][0]);
    }
    throw Exception('Failed to create todo');
  }

  // Update a todo
  static Future<Todo> updateTodo(int id, String contents) async {
    final response = await http.put(
      Uri.parse('$baseUrl/todos/$id'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({'contents': contents}),
    );
    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      return Todo.fromJson(data['data'][0]);
    }
    throw Exception('Failed to update todo');
  }

  // Delete a todo
  static Future<void> deleteTodo(int id) async {
    final response = await http.delete(Uri.parse('$baseUrl/todos/$id'));
    if (response.statusCode != 200) {
      throw Exception('Failed to delete todo');
    }
  }
}
```

## Development

### Project Structure Explanation

1. **`cmd/server/main.go`**: Entry point that:
   - Initializes the service layer
   - Sets up HTTP handlers
   - Configures routing and middleware
   - Starts the server

2. **`internal/models/todo.go`**: Contains:
   - `Todo`: Core data structure
   - `TodoRequest`: Input validation model
   - `TodoResponse`: API response format

3. **`internal/service/todo_service.go`**: Business logic including:
   - CRUD operations
   - Data validation
   - File I/O operations
   - ID generation

4. **`internal/handlers/todo_handlers.go`**: HTTP layer that:
   - Handles HTTP requests/responses
   - Manages status codes
   - Formats JSON responses
   - Delegates to service layer

### Error Handling

The API returns appropriate HTTP status codes:

- **200 OK**: Successful operation
- **201 Created**: Todo created successfully
- **400 Bad Request**: Invalid request format or validation errors
- **404 Not Found**: Todo with specified ID doesn't exist
- **500 Internal Server Error**: Server-side errors

### Data Validation

- **Contents**: Cannot be empty or contain only whitespace
- **ID**: Must be a valid integer for update/delete operations
- **JSON**: Request body must be valid JSON format

## Dependencies

- **`github.com/gorilla/mux`**: HTTP router for handling different endpoints

## Benefits of This Architecture

1. **Separation of Concerns**: Each layer has a single responsibility
2. **Testability**: Easy to unit test individual components
3. **Maintainability**: Changes in one layer don't affect others
4. **Scalability**: Easy to add new features or modify existing ones
5. **Dependency Injection**: Loose coupling between components

## Running Tests

You can add tests for each layer:

```bash
# Test the service layer
go test ./internal/service

# Test the handlers layer
go test ./internal/handlers

# Test all packages
go test ./...
```

This clean architecture makes the codebase more maintainable, testable, and scalable for your Flutter todo application! 