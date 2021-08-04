package boltdbtoipldhamt

import "github.com/emirpasic/gods/stacks/linkedliststack"

type (
	TreeList []*TreeNode

	TreeNode struct {
		Level    int         `json:"level"`
		Data     interface{} `json:"data"`
		Parent   *TreeNode   `json:"parent"`
		Children TreeList    `json:"children"`
	}
)

var visitor = NewVisitor()
var ListStack = linkedliststack.New()

func NewTreeNode(data interface{}, parent *TreeNode) *TreeNode {
	return &TreeNode{Data: data, Parent: parent, Level: 0}
}

func NewTreeNodeWithLevel(data interface{}, parent *TreeNode, level int) *TreeNode {
	return &TreeNode{Data: data, Parent: parent, Level: level}
}

func (t *TreeNode) AddChild(data interface{}) *TreeNode {
	newNode := NewTreeNodeWithLevel(data, t, t.Level+1)
	t.Children = append(t.Children, newNode)
	return newNode
}

func TransverseTreeNode(treeNode *TreeNode, visitFunc OnVisitFunc) {
	stack := linkedliststack.New()
	stack.Push(treeNode)
	for !stack.Empty() {
		n, ok := stack.Pop()
		if !ok {
			continue
		}
		ListStack.Push(n)

		if visitor.Visited(n) {
			visitor.Visit(n, visitFunc)
			for _, w := range n.(*TreeNode).Children {
				if visitor.Visited(w) {
					stack.Push(w)
				}
			}
		}
	}
}
