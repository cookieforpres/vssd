package ssd

import (
	"errors"
	"sync"
)

type SSD struct {
	Name        string
	Size        int
	CurrentSize int
	Data        []byte
	Nodes       []*Node
	*sync.Mutex
}

var ErrNoSpace = errors.New("there is no space left on this ssd to write the data")
var ErrNoNode = errors.New("there is no node by that name")

func New(name string, size int) *SSD {
	return &SSD{
		Name:        name,
		Size:        size,
		CurrentSize: 0,
		Data:        make([]byte, size),
		Nodes:       make([]*Node, 0),
		Mutex:       &sync.Mutex{},
	}
}

func (s *SSD) AddNode(node *Node) {
	if s.TryLock() {
		defer s.Unlock()
	}

	s.Nodes = append(s.Nodes, node)
}

func (s *SSD) RemoveNode(node *Node) {
	if s.TryLock() {
		defer s.Unlock()
	}

	for i, n := range s.Nodes {
		if n == node {
			s.Nodes = append(s.Nodes[:i], s.Nodes[i+1:]...)
			for j := node.Start; j < node.End; j++ {
				s.Data[j] = 0
			}
			break
		}
	}
}

func (s *SSD) Write(node *Node, data []byte) error {
	for _, n := range s.Nodes {
		if n.Name == node.Name {
			s.RemoveNode(n)
			break
		}
	}

	if s.TryLock() {
		defer s.Unlock()
	}

	for i := 0; i < len(s.Data); i++ {
		start := i
		for j := 0; j < len(data); j++ {
			if (i + j) >= len(s.Data) {
				if s.Data[i+j] != 0 {
					i += j
					break
				}
			}

			if j == len(data)-1 {
				for k := 0; k < len(data); k++ {
					s.Data[i+k] = data[k]
				}

				node.Start = start
				node.End = i + len(data)

				s.CurrentSize += len(data)
				s.AddNode(node)

				return nil
			}
		}
	}

	return ErrNoSpace
}

func (s *SSD) Read(name string) ([]byte, error) {
	if s.TryLock() {
		defer s.Unlock()
	}

	for _, n := range s.Nodes {
		if n.Name == name {
			return s.Data[n.Start:n.End], nil
		}
	}

	return nil, ErrNoNode
}

func (s *SSD) Delete(name string) error {
	if s.TryLock() {
		defer s.Unlock()
	}

	for _, n := range s.Nodes {
		if n.Name == name {
			for i := n.Start; i < n.End; i++ {
				s.Data[i] = 0
			}

			s.RemoveNode(n)
			s.CurrentSize -= n.End - n.Start

			return nil
		}
	}

	return ErrNoNode
}

func (s *SSD) HasSpace(size int) bool {
	if s.TryLock() {
		defer s.Unlock()
	}
	return s.CurrentSize+size <= s.Size
}
