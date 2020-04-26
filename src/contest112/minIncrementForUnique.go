package main

import (
	"fmt"
	"sort"
)

func minIncrementForUnique(A []int) int {
	sort.Ints(A)
	if len(A) == 0 {
		return 0
	}
	last := A[0]
	count := 1
	ret := 0
	res := 0
	for i, v := range A {
		if i == 0 {
			continue
		}
		if v == last { // duplicate
			count++
		} else { // new value
			ret = ret + count - 1 // new ret value
			res += ret
			med := last + 1
			for ret != 0 && med < v {
				ret--
				res += ret
				med++
			}
			last = v
			count = 1
		}
	}
	ret = ret + count - 1
	res += ret
	for ret != 0 {
		ret--
		res += ret
	}
	return res
}

func main() {
	fmt.Println("Hello World!")
	fmt.Printf("ex1: %d\n", minIncrementForUnique([]int{1, 2, 2}))
	fmt.Printf("ex2: %d\n", minIncrementForUnique([]int{1, 2, 2, 2, 3, 3}))
	fmt.Printf("ex3: %d\n", minIncrementForUnique([]int{3, 2, 1, 2, 1, 7}))
}
