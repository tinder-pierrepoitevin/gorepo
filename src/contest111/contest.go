package main

import (
	"fmt"
	"strings"
)

/*
Example 1:

Input: [2,1]
Output: false
Example 2:

Input: [3,5,5]
Output: false
Example 3:

Input: [0,3,2,1]
Output: true
*/
func validMountainArray(A []int) bool {
	descending := false
	previous := 0
	if len(A) < 3 {
		return false
	}
	if A[0] >= A[1] {
		return false
	}
	for i, v := range A {
		if i != 0 {
			if v > previous {
				if descending {
					return false
				}
			} else if v == previous {
				return false
			} else { // v < previous
				if !descending {
					descending = true
				}
			}
		}
		previous = v
	}
	return descending
}

func testValidMountainArray() {
	fmt.Printf("validMountainArray([0,3,2,1]): %t\n", validMountainArray([]int{0, 3, 2, 1}))
	fmt.Printf("validMountainArray([0,1]): %t\n", validMountainArray([]int{0, 1}))
	fmt.Printf("validMountainArray([0,1,1]): %t\n", validMountainArray([]int{0, 1, 1}))
	fmt.Printf("validMountainArray([0,1,0]): %t\n", validMountainArray([]int{0, 1, 0}))
	fmt.Printf("validMountainArray([2,1,0]): %t\n", validMountainArray([]int{2, 1, 1}))
	fmt.Printf("validMountainArray([3,5,5]): %t\n", validMountainArray([]int{3, 5, 5}))
}

func minDeletionSize(A []string) int {
	count := 0
	if len(A) == 0 {
		return 0
	}
	strSize := len(A[0])
	for i := 0; i < strSize; i = i + 1 {
		previous := A[0][0]
		for ind, val := range A {
			// fmt.Printf("%s <> %s")
			if ind != 0 {
				if previous > val[i] {
					count = count + 1
					break
				}
			}
			previous = val[i]
		}
	}
	return count
}

/*
Example 1:

Input: ["cba","daf","ghi"]
Output: 1
Explanation:
After choosing D = {1}, each column ["c","d","g"] and ["a","f","i"] are in non-decreasing sorted order.
If we chose D = {}, then a column ["b","a","h"] would not be in non-decreasing sorted order.
Example 2:

Input: ["a","b"]
Output: 0
Explanation: D = {}
Example 3:

Input: ["zyx","wvu","tsr"]
Output: 3
Explanation: D = {0, 1, 2}
*/
func testMinDeletionSize() {
	fmt.Printf("ex1: %d\n", minDeletionSize([]string{"cba", "daf", "ghi"}))
	fmt.Printf("ex2: %d\n", minDeletionSize([]string{"a", "b"}))
	fmt.Printf("ex3: %d\n", minDeletionSize([]string{"zyx", "wvu", "tsr"}))
}

func diStringMatch(S string) []int {
	N := len(S)
	res := make([]int, N+1)
	res[0] = 0
	min := 0
	max := 0
	for i := 0; i < N; i++ {
		if S[i] == 'I' {
			max++
			res[i+1] = max
		} else {
			min--
			res[i+1] = min
		}
	}
	for i, v := range res {
		res[i] = v - min
	}
	return res
}

/*
Example 1:

Input: "IDID"
Output: [0,4,1,3,2]
Example 2:

Input: "III"
Output: [0,1,2,3]
Example 3:

Input: "DDI"
Output: [3,2,0,1]
*/
func testDiStringMatch() {
	testString("IDID")
	testString("III")
	testString("DDI")
}

func testString(S string) {
	fmt.Printf("diStringMatch(\"%s\") -> %v\n", S, diStringMatch(S))
}

/*
Given an array A of strings, find any smallest string that contains each string in A as a substring.

We may assume that no string in A is substring of another string in A.
*/
func shortestSuperstring(A []string) string {
	// Make a graph of overlapping and go through all the possible chains.
	N := len(A)
	graph := make([][]int, N)
	for i := 0; i < N; i++ {
		graph[i] = make([]int, N)
	}
	for i, v1 := range A {
		for j, v2 := range A {
			if i == j {
				graph[i][i] = -1
			} else {
				graph[i][j] = maxOverlap(v1, v2)
				graph[j][i] = maxOverlap(v2, v1)
			}
		}
	}
	fmt.Printf("%v\n", graph)
	return ""
}

func bestPermutation(graph *[][]int) {
	
}
func maxOverlap(first string, second string) int {
	F := len(first)
	S := len(second)
	res := 0
	for i := 0; i < F && i < S; i++ {
		v := strings.LastIndex(first, second[:i+1])
		if v != -1 && v+i == F-1 {
			res = i + 1
		}
		if v == -1 {
			break
		}
	}
	return res
}

func testMaxOverlap() {
	fmt.Printf("0 : %d\n", maxOverlap("", ""))
	fmt.Printf("1 : %d\n", maxOverlap("a", "a"))
	fmt.Printf("1 : %d\n", maxOverlap("ab", "ba"))
	fmt.Printf("0 : %d\n", maxOverlap("ab", "cd"))
	fmt.Printf("2 : %d\n", maxOverlap("abcdef", "efghij"))
	fmt.Printf("3 : %d\n", maxOverlap("baba", "abab"))
	fmt.Printf("4 : %d\n", maxOverlap("baba", "baba"))
	fmt.Printf("1 : %d\n", maxOverlap("abcdef", "f"))
	fmt.Printf("1 : %d\n", maxOverlap("f", "fedcba"))
	fmt.Printf("6 : %d\n", maxOverlap("abdcbdbac", "cbdbaccdb"))
}

/*
Example 1:

Input: ["alex","loves","leetcode"]
Output: "alexlovesleetcode"
Explanation: All permutations of "alex","loves","leetcode" would also be accepted.
Example 2:

Input: ["catg","ctaagt","gcta","ttca","atgcatc"]
Output: "gctaagttcatgcatc"
*/
func testShortestSuperstring() {
	fmt.Printf("ex1 -> %s\n", shortestSuperstring([]string{"alex", "loves", "leetcode"}))
	fmt.Printf("ex2 -> %s\n", shortestSuperstring([]string{"catg", "ctaagt", "gcta", "ttca", "atgcatc"}))
}

func main() {
	testMaxOverlap()
	testShortestSuperstring()
}
