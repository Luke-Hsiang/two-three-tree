package main

import (
	"fmt"
)

// Value is a helper class since int cannot be nil
type Value struct {
	Val int
}

// Node is a node in 2-3 tree
//
// A 2-node must have its value set in Val1, and with its Val2 set to nil
// and its subtrees (if exists) saved in Left / Mid, with its Right node set to nil
//
// Parent is just a helper node for insertion and deletion
type Node struct {
	Val1   *Value
	Val2   *Value
	Left   *Node
	Mid    *Node
	Right  *Node
	Parent *Node
}

func main() {
	root := &Node{}
	root = root.Insert(20)
	root = root.Insert(10)
	root = root.Insert(30)
	root = root.Insert(50)
	root = root.Insert(40)
	root = root.Insert(35)
	root = root.Insert(33)
	root = root.Insert(47)
	root.Print()
}

// Str print current node values
func (node *Node) Str() string {
	if node == nil {
		return ""
	}
	if node.Val1 == nil {
		return ""
	}

	if node.Val2 == nil {
		return fmt.Sprintf("(%v, nil)", node.Val1.Val)
	}

	return fmt.Sprintf("(%v, %v)", node.Val1.Val, node.Val2.Val)
}

// Print the tree by layer
func (node *Node) Print() {
	if node == nil {
		return
	}

	fmt.Printf("node: %v\n", node.Str())

	layer := 1

	curLayer := []*Node{node}

	for len(curLayer) > 0 {
		fmt.Printf("Layer %v:\n", layer)
		nextLayer := make([]*Node, 0)

		for _, n := range curLayer {
			fmt.Printf("%v", n.Str())
			if n.Parent != nil {
				fmt.Printf(", Parent: %v\n", n.Parent.Str())
			} else {
				fmt.Printf("\n")
			}

			if n.Left != nil {
				nextLayer = append(nextLayer, n.Left)
			}
			if n.Mid != nil {
				nextLayer = append(nextLayer, n.Mid)
			}
			if n.Right != nil {
				nextLayer = append(nextLayer, n.Right)
			}
		}
		curLayer = nextLayer

		layer++
	}
}

// Insert a value to a tree. Insertion will not succeed if node is not root of a tree
// return new root of the tree after insertion
func (node *Node) Insert(val int) *Node {
	if node.Parent != nil {
		return node
	}

	// if node has no Val, set Val1 to val.
	// this only happens when a tree is empty
	if node.Val1 == nil {
		node.Val1 = &Value{val}
	} else {
		node.innerInsert(val)
	}

	tmpRoot := node
	for tmpRoot.Parent != nil {
		tmpRoot = tmpRoot.Parent
	}
	return tmpRoot
}

func (node *Node) innerInsert(val int) {
	//insertion only happens on leaves. So if current node is not a leaf, insert the new value to the corresponding subtree
	if node.Left != nil {
		// current node is not leaf
		if node.Val1.Val > val {
			node.Left.innerInsert(val)
		} else if node.Val2 == nil || node.Val2.Val > val {
			node.Mid.innerInsert(val)
		} else {
			node.Right.innerInsert(val)
		}
	} else {
		// current node is leaf
		if node.Val2 == nil {
			// there is only one value in this leaf node,
			// just set value to the correspoding position
			if node.Val1.Val > val {
				node.Val2 = &Value{node.Val1.Val}
				node.Val1.Val = val
			} else {
				node.Val2 = &Value{val}
			}
		} else {
			// there are two values in this leaf node,
			// need to split current node, then merge it with parent
			var left int
			var mid int
			var right int
			if node.Val1.Val > val {
				left = val
				mid = node.Val1.Val
				right = node.Val2.Val
			} else if node.Val2.Val > val {
				left = node.Val1.Val
				mid = val
				right = node.Val2.Val
			} else {
				left = node.Val1.Val
				mid = node.Val2.Val
				right = val
			}

			node.Val1.Val = mid
			node.Val2 = nil
			node.Left = &Node{Val1: &Value{left}, Val2: nil, Parent: node}
			node.Mid = &Node{Val1: &Value{right}, Val2: nil, Parent: node}

			node.mergeNode()
		}
	}
}

func (node *Node) mergeNode() {
	if node.Parent == nil {
		return
	}

	if node.Parent.Val2 == nil {
		node.Left.Parent = node.Parent
		node.Mid.Parent = node.Parent

		if node.Parent.Val1.Val > node.Val1.Val {
			node.Left.Parent = node.Parent
			node.Parent.Val2 = &Value{node.Parent.Val1.Val}
			node.Parent.Val1 = node.Val1
			node.Parent.Right = node.Parent.Mid
			node.Parent.Left = node.Left
			node.Parent.Mid = node.Mid
		} else {
			node.Parent.Val2 = node.Val1
			node.Parent.Mid = node.Left
			node.Parent.Right = node.Mid
		}
	} else {
		if node.Parent.Val1.Val > node.Val1.Val {
			node.Parent.Left = node
			node.Parent.Mid = &Node{
				Val1:   node.Parent.Val2,
				Val2:   nil,
				Left:   node.Parent.Mid,
				Mid:    node.Parent.Right,
				Right:  nil,
				Parent: node.Parent,
			}
			node.Parent.Val2 = nil
			node.Parent.Right = nil
		} else if node.Parent.Val2.Val > node.Val1.Val {
			node.Parent.Left = &Node{
				Val1:   node.Parent.Val1,
				Val2:   nil,
				Left:   node.Parent.Left,
				Mid:    node.Left,
				Right:  nil,
				Parent: node.Parent,
			}
			node.Parent.Mid = &Node{
				Val1:   node.Parent.Val2,
				Val2:   nil,
				Left:   node.Mid,
				Mid:    node.Parent.Right,
				Right:  nil,
				Parent: node.Parent,
			}
			node.Parent.Val1 = node.Val1
			node.Parent.Val2 = nil
			node.Parent.Right = nil
		} else {
			node.Parent.Left = &Node{
				Val1:   node.Parent.Val1,
				Val2:   nil,
				Left:   node.Parent.Left,
				Mid:    node.Parent.Mid,
				Right:  nil,
				Parent: node.Parent,
			}
			node.Parent.Mid = node
			node.Parent.Val1 = node.Parent.Val2
			node.Parent.Val2 = nil
			node.Parent.Right = nil
		}
		//reassign parent
		if node.Parent.Left.Left != nil {
			node.Parent.Left.Left.Parent = node.Parent.Left
			node.Parent.Left.Mid.Parent = node.Parent.Left
		}
		if node.Parent.Mid.Left != nil {
			node.Parent.Mid.Left.Parent = node.Parent.Mid
			node.Parent.Mid.Mid.Parent = node.Parent.Mid
		}
		node.Parent.mergeNode()
	}
}
