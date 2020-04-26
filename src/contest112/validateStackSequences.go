package main

import (
	"fmt"
)

func validateStackSequences(pushed []int, popped []int) bool {
	stack := make([]int, 1001)
	if len(pushed) == 0 {
		return true
	}
	poppedIndex := 0
	end := 0
	for _, v := range pushed {
		stack[end] = v
		end++
		for end > 0 && stack[end-1] == popped[poppedIndex] {
			end--
			poppedIndex++
		}
	}
	return end == 0
}

func main() {
	fmt.Println("Hello World!")
	fmt.Printf("ex1: %t\n", validateStackSequences([]int{1, 2, 2}, []int{1, 2, 2}))
	fmt.Printf("ex2: %t\n", validateStackSequences([]int{1, 2, 3, 4, 5}, []int{4, 3, 5, 1}))

}
