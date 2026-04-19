package linked_list

import (
	. "gocode/types"
	"sync"
)

type LinkedHashMap struct {
	nodes map[int]*Node // num → linked list node
	head  *Node         // sentinel
	tail  *Node         // sentinel

	// sync.Mutex:   one goroutine at a time — reads block each other unnecessarily
	// sync.RWMutex: multiple concurrent readers allowed
	// writer gets exclusive access — blocks all readers and writers
	mu sync.RWMutex // read-write lock
}

// remove a node from the linked list — O(1)
func (lhm *LinkedHashMap) remove(node *Node) {
	node.Prev.Next = node.Next
	node.Next.Prev = node.Prev
}

// append to tail — O(1)
func (lhm *LinkedHashMap) append(node *Node) {
	node.Prev = lhm.tail.Prev
	node.Next = lhm.tail
	lhm.tail.Prev.Next = node
	lhm.tail.Prev = node
}

// peek first — O(1)
func (lhm *LinkedHashMap) first() *Node {
	if lhm.head.Next == lhm.tail {
		return nil // empty
	}
	return lhm.head.Next
}

/**
 * 1429. First Unique Number
 * You have a queue of integers, you need to retrieve the first unique integer in the queue.
 * Implement the FirstUnique class:
 *
 * FirstUnique(nums) — initializes the object with the numbers in the queue
 * showFirstUnique() — returns the value of the first unique integer, or -1 if none exists
 * add(value) — inserts value into the queue
 *
 * HashMap    → O(1) key lookup
 * LinkedList → maintains insertion order
 *
 */
type FirstUnique struct {
	freq map[int]int // num → count
	lhm  LinkedHashMap
}

func ConstructorLinkedHashmap(nums []int) FirstUnique {
	fu := FirstUnique{
		freq: make(map[int]int),
		lhm: LinkedHashMap{
			nodes: make(map[int]*Node),
			head:  &Node{},
			tail:  &Node{},
		},
	}
	fu.lhm.head.Next = fu.lhm.tail
	fu.lhm.tail.Prev = fu.lhm.head

	for _, num := range nums {
		fu.Add(num)
	}
	return fu
}

func (fu *FirstUnique) Add(num int) {
	fu.lhm.mu.Lock() // exclusive write lock
	defer fu.lhm.mu.Unlock()

	fu.freq[num]++
	if fu.freq[num] == 1 {
		// first occurrence — add to linked list
		node := &Node{Val: num}
		fu.lhm.nodes[num] = node
		fu.lhm.append(node)
	} else if fu.freq[num] == 2 {
		// became duplicate — remove from linked list
		fu.lhm.remove(fu.lhm.nodes[num])
	}
	// freq > 2: already removed, nothing to do
}

func (fu *FirstUnique) ShowFirstUnique() int {
	fu.lhm.mu.RLock() // shared read lock — multiple readers allowed
	defer fu.lhm.mu.RUnlock()

	if fu.lhm.head.Next == fu.lhm.tail {
		return -1
	}
	return fu.lhm.head.Next.Val
}
