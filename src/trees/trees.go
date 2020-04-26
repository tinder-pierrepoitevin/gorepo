package main

import (
	"fmt"
)

/**
 * Definition for a binary tree node.
 */
type TreeNode struct {
	Left  *TreeNode
	Right *TreeNode
}

/**
 *  ValuedTreeNode is a TreeNode with a integer as composite data
 */
type ValuedTreeNode struct {
	Val   int
	Left  *ValuedTreeNode
	Right *ValuedTreeNode
}

func newLeaf() *TreeNode {
	return &TreeNode{nil, nil}
}

func newValuedLeaf(val int) *ValuedTreeNode {
	return &ValuedTreeNode{val, nil, nil}
}

func instantiateTree() *TreeNode {
	leftLeft := newLeaf()
	leftRight := newLeaf()
	left := &TreeNode{leftLeft, leftRight}
	right := newLeaf()
	root := &TreeNode{left, right}
	return root
}

func instantiateValuedTree() *ValuedTreeNode {
	leftLeft := newValuedLeaf(4)
	leftRight := newValuedLeaf(5)
	left := &ValuedTreeNode{2, leftLeft, leftRight}
	right := newValuedLeaf(3)
	root := &ValuedTreeNode{1, left, right}
	return root
}

func computeSize(node *TreeNode) int {
	if node == nil {
		return 0
	}
	leftSize := computeSize(node.Left)
	rightSize := computeSize(node.Right)
	return sizeCombination(leftSize, rightSize)
}

func sizeCombination(leftSize int, rightSize int) int {
	return leftSize + rightSize + 1
}

func computeHeight(node *TreeNode) int {
	if node == nil || (node.Left == nil && node.Right == nil) {
		return 0
	}
	leftHeight := computeHeight(node.Left)
	rightHeight := computeHeight(node.Right)
	return heightCombination(leftHeight, rightHeight)
}

func heightCombination(leftHeight int, rightHeight int) int {
	return max(leftHeight, rightHeight) + 1
}

func max(a int, b int) int {
	if a < b {
		return b
	}
	return a
}

func main() {
	root := instantiateTree()
	fmt.Printf("Tree size: %d\n", computeSize(root))
	fmt.Printf("Tree height: %d\n", computeHeight(root))
}
