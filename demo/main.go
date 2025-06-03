package main

import (
	"fmt"
)

func main() {
	fmt.Println("Starting the demo application...")
	SayHello()

	// Call the Add function with different examples
	fmt.Println("\nTesting the Add function:")
	Add(5, 6)  // Should output: the addition is 11
	Add(-5, 6) // Should output: the addition is 1
}
