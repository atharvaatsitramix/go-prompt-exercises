package main

import (
	"fmt"
	"strings"
	"time"
	"todo-backend/internal/models"
	"todo-backend/internal/service"
)

// linearSearchTodoByID simulates the old linear search method
func linearSearchTodoByID(todos []models.Todo, targetID int) int {
	for i, todo := range todos {
		if todo.ID == targetID {
			return i
		}
	}
	return -1
}

func main() {
	fmt.Println("🚀 Todo API Performance Comparison: Linear Search vs Binary Search")
	fmt.Println(strings.Repeat("=", 70))

	// Create test data with various sizes
	testSizes := []int{100, 1000, 10000, 100000}

	for _, size := range testSizes {
		fmt.Printf("\n📊 Testing with %d todos:\n", size)

		// Generate test todos (sorted by ID)
		todos := make([]models.Todo, size)
		for i := 0; i < size; i++ {
			todos[i] = models.Todo{
				ID:       i + 1,
				Contents: fmt.Sprintf("Todo item %d", i+1),
			}
		}

		// Target ID to search for (worst case - near the end)
		targetID := size - 10

		// Test Linear Search
		start := time.Now()
		iterations := 1000
		for i := 0; i < iterations; i++ {
			linearSearchTodoByID(todos, targetID)
		}
		linearTime := time.Since(start)

		// Test Binary Search
		service := service.NewTodoService("temp.json")
		start = time.Now()
		for i := 0; i < iterations; i++ {
			service.BinarySearchTodoByID(todos, targetID)
		}
		binaryTime := time.Since(start)

		// Calculate improvement
		improvement := float64(linearTime) / float64(binaryTime)

		fmt.Printf("  Linear Search:  %v (average: %v per search)\n",
			linearTime, linearTime/time.Duration(iterations))
		fmt.Printf("  Binary Search:  %v (average: %v per search)\n",
			binaryTime, binaryTime/time.Duration(iterations))
		fmt.Printf("  🎯 Improvement: %.2fx faster with binary search\n", improvement)
	}

	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("📈 Performance Analysis:")
	fmt.Println("• Linear Search: O(n) - Performance degrades linearly with data size")
	fmt.Println("• Binary Search: O(log n) - Performance remains excellent even with large datasets")
	fmt.Println("• With 100,000 todos, binary search is significantly faster!")
	fmt.Println("• This optimization is especially important for production systems with large datasets")
}
