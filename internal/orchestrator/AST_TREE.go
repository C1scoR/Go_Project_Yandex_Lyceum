package orchestrator

import (
	"fmt"
)

type Node struct {
	value  string
	right  *Node
	left   *Node
	next   *Node
	Status string
}

type Stack struct {
	head *Node
}

func NewNode(value string) *Node {
	return &Node{
		value: value,
		right: nil,
		left:  nil,
		next:  nil,
	}
}

// рекурсивно обходит дерево и печатает все его значения
func PrintOrder(node *Node) {
	if node == nil {
		return
	}
	PrintOrder(node.left)
	fmt.Print(node.value)
	PrintOrder(node.right)
}

func (st *Stack) Push(node *Node) {
	if st.head == nil {
		st.head = node
	}
	node.next = st.head
	st.head = node
}

func (st *Stack) Pop() *Node {
	if st.head != nil {
		popped_element := st.head
		st.head = st.head.next
		return popped_element
	}
	return nil
}

func TranslateToASTTree(Postfix_expression []string) *Node {
	//Postfix_expression = []string{"A", "B", "+", "C", "D", "+", "+", "E", "F", "+", "+"}
	var Stack Stack
	var x, y, tree *Node
	for _, element := range Postfix_expression {
		switch element {
		case "+", "-", "*", "/":
			tree = NewNode(element)
			tree.Status = StatusFree
			x = Stack.Pop()
			y = Stack.Pop()
			tree.left = y
			tree.right = x
			Stack.Push(tree)
		default:
			tree = NewNode(element)
			Stack.Push(tree)
		}
	}
	//PrintOrder(tree)
	return tree
}
