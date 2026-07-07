package apidesign

/**
 * 146. LRU Cache
 *
 * 1. hashmap + doubly linked list
 *     dummy node for head & tail
 *     real tail is pointed by dummy tail, updated during insert() & remove()
 *
 * 2. container/list
 */
type DoublyLinkedNode struct {
	key, val   int
	prev, next *DoublyLinkedNode
}

type LRUCache struct {
	head, tail *DoublyLinkedNode // dummyHead & dummyTail
	cache      map[int]*DoublyLinkedNode
	capacity   int
}

func ConstructorLRUCache(capacity int) LRUCache {
	head, tail := &DoublyLinkedNode{}, &DoublyLinkedNode{}
	head.next = tail
	tail.prev = head
	return LRUCache{
		capacity: capacity,
		head:     head,
		tail:     tail,
		cache:    make(map[int]*DoublyLinkedNode),
	}
}

func (l *LRUCache) Get(key int) int {
	if node, found := l.cache[key]; found {
		l.remove(node)
		l.prepend(node)
		return node.val
	}
	return -1
}

func (l *LRUCache) Put(key int, value int) {
	if node, found := l.cache[key]; found {
		l.remove(node)
		l.prepend(node)
		node.val = value
	} else {
		newNode := &DoublyLinkedNode{key: key, val: value}
		l.cache[key] = newNode
		l.prepend(newNode)
		if len(l.cache) > l.capacity {
			tail := l.tail.prev // tail <- dummyTail
			l.remove(tail)
			delete(l.cache, tail.key)
		}
	}
}

func (l *LRUCache) remove(node *DoublyLinkedNode) {
	prev, next := node.prev, node.next
	prev.next, next.prev = next, prev
	node.prev, node.next = nil, nil
}

func (l *LRUCache) prepend(node *DoublyLinkedNode) {
	oldHead := l.head.next
	l.head.next = node // dummyHead -> node

	node.prev = l.head
	node.next = oldHead // node -> dummyTail
	oldHead.prev = node // node <- dummyTail
}
