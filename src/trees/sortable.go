package main

import (
	"fmt"
	"sort"
)

func main() {
	numbers := []int{
		2,
		25,
		3,
		97,
		4,
		6,
		9,
		34,
		57,
	}
	sort.Sort(ByDivisor(numbers))
	fmt.Println(numbers)
}

/*
 * ByDivisor doc
 */
type ByDivisor []int

/*
 * ByDivisor doc
 */
func (sortable ByDivisor) Len() int {
	return len(sortable)
}

/*
 * ByDivisor doc
 */
func (sortable ByDivisor) Swap(i, j int) {
	x := sortable[i]
	sortable[i] = sortable[j]
	sortable[j] = x
}

/*
 * ByDivisor doc
 */
func (sortable ByDivisor) Less(i, j int) bool {
	if sortable[j] < sortable[i] {
		return !sortable.Less(j, i)
	}
	for x := 2; x <= sortable[i]; x++ {
		if sortable[i]%x == 0 {
			return true
		}
		if sortable[j]%x == 0 {
			return false
		}
	}
	// shouldn't be used
	return true
}
