package main

import "fmt"

/**
 * Definition for a binary tree node.
 */
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func findBottomLeftValue(root *TreeNode) int {
	// Solution: do an in order traversale and remember the max level and the value encountered.
	res, _ := findBottomLeftValueRec(root, 0)
	return res
}

func findBottomLeftValueRec(root *TreeNode, height int) (int, int) {
	if root == nil {
		return 0, -1
	}
	leftR, leftH := findBottomLeftValueRec(root.Left, height+1)
	rightR, rightH := findBottomLeftValueRec(root.Right, height+1)
	if height > leftH {
		if height > rightH {
			return root.Val, height
		} else {
			return rightR, rightH
		}
	} else {
		if rightH > leftH {
			return rightR, rightH
		} else {
			// note that in the case heights are equal, we pick left
			return leftR, leftH
		}
	}
}

func main() {
	// 	        1
	// 	       / \
	//        2   3
	//       /   / \
	//      4   5   6
	// 	       /
	//        7
	tree7 := &TreeNode{7, nil, nil}
	tree4 := &TreeNode{4, nil, nil}
	tree5 := &TreeNode{5, tree7, nil}
	tree6 := &TreeNode{6, nil, nil}
	tree2 := &TreeNode{2, tree4, nil}
	tree3 := &TreeNode{3, tree5, tree6}
	tree1 := &TreeNode{1, tree2, tree3}
	res := findBottomLeftValue(tree1)
	fmt.Printf("res: %d\n", res)
}
