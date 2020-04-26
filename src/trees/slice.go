package main

import (
	"fmt"
)

func main() {
	appendToSlice()
	removeFromSlice()
}

func initArray() {
	var A [3]int
	A[0] = 1
	A[1] = 2
	A[2] = 3
	fmt.Println(A)
}

func initArray2() {
	A := [...]int{1, 2, 3, 4}
	fmt.Println(A)
}

func initSlice() {
	A := []int{1, 2, 3, 4, 5}
	fmt.Println(A)
	B := []int{}
	fmt.Println(B)
}

func initSliceMake() {
	A := make([]int, 0)
	fmt.Println(A)
	B := make([]int, 3)
	fmt.Println(B)
}

func traverseSlice() {
	A := make([]int, 11)
	for i := range A {
		A[i] = i * i
	}
	fmt.Println(A)
}

func appendToSlice() {
	fmt.Println("append to slice")
	A := []int{1, 2, 3, 4, 5}
	B := append(A, 6)
	fmt.Println(A)
	fmt.Println(B)
}

func removeFromSlice() {
	fmt.Println("remove slice")
	A := []int{1, 2, 3, 4, 5}
	B := A[1:4]
	fmt.Println(A)
	fmt.Println(B)
}
