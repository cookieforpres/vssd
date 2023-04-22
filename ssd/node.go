package ssd

type Node struct {
	Name  string
	Start int
	End   int
}

func NewNode(name string) *Node {
	return &Node{
		Name:  name,
		Start: -1,
		End:   -1,
	}
}
