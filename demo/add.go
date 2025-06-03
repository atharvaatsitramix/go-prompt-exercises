package main

import "fmt"

// Add takes two integer parameters and performs addition
// It calculates the sum of the two numbers and prints the result
// in the format "the addition is X" where X is the sum
// Parameters:
//   - a: first integer to add
//   - b: second integer to add
//
// Example usage:
//   - Add(5, 6) outputs "the addition is 11"
//   - Add(-5, 6) outputs "the addition is 1"
func Add(a, b int) {
	// Calculate the sum of the two input numbers
	sum := a + b

	// Print the result in the specified format
	fmt.Printf("the addition is %d\n", sum)
}
