package apidesign

import (
	"container/list"
	"math"
)

/**
 * 716. Max Stack (full)
 *
 */
type MaxStack struct {
	dll     *list.List
	treemap map[int][]*list.Element // value → stack of DLL nodes (handles duplicates)
}

func ConstructorMaxStack() MaxStack {
	return MaxStack{
		dll:     list.New(),
		treemap: make(map[int][]*list.Element),
	}
}

func (s *MaxStack) Push(val int) {
	elem := s.dll.PushBack(val)
	s.treemap[val] = append(s.treemap[val], elem)
}

func (s *MaxStack) Pop() int {
	elem := s.dll.Back()
	s.dll.Remove(elem)

	val := elem.Value.(int)
	nodes := s.treemap[val]
	s.treemap[val] = nodes[:len(nodes)-1]
	if len(s.treemap[val]) == 0 {
		delete(s.treemap, val)
	}
	return val
}

func (s *MaxStack) Top() int {
	return s.dll.Back().Value.(int)
}

func (s *MaxStack) PeekMax() int {
	return s.maxKey()
}

func (s *MaxStack) PopMax() int {
	maxVal := s.maxKey()
	nodes := s.treemap[maxVal]
	elem := nodes[len(nodes)-1] // most recently pushed node with this value

	s.dll.Remove(elem)
	s.treemap[maxVal] = nodes[:len(nodes)-1]
	if len(s.treemap[maxVal]) == 0 {
		delete(s.treemap, maxVal)
	}
	return maxVal
}

// maxKey returns the current maximum value in O(n) by scanning treemap keys.
// In a language with a built-in treemap this would be O(log n).
func (s *MaxStack) maxKey() int {
	maxVal := math.MinInt64
	for key := range s.treemap {
		if key > maxVal {
			maxVal = key
		}
	}
	return maxVal
}
