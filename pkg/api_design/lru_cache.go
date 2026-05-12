package apidesign

/**
 * 146. LRU Cache
 *
 * 1. hashmap + doubly linked list
 *     dummy node for head & tail
 * 2. container/list
 */
type DoublyLinkedList struct {
	key  int
	val  int
	prev *DoublyLinkedList
	next *DoublyLinkedList
}

type LRUCache struct {
	capacity int
	head     *DoublyLinkedList
	tail     *DoublyLinkedList
	cache    map[int]*DoublyLinkedList // map cache value > double linked list node
}

func ConstructorLRUCache(capacity int) LRUCache {
	head := &DoublyLinkedList{key: 0, val: 0} // dummy node
	tail := &DoublyLinkedList{key: 0, val: 0}
	head.next = tail
	tail.prev = head
	return LRUCache{
		capacity: capacity,
		head:     head,
		tail:     tail,
		cache:    make(map[int]*DoublyLinkedList),
	}
}

func (l *LRUCache) remove(node *DoublyLinkedList) {
	prev, next := node.prev, node.next
	prev.next, next.prev = next, prev
}

func (l *LRUCache) insert(node *DoublyLinkedList) {
	oldHead := l.head.next
	l.head.next = node
	node.prev = l.head
	node.next = oldHead
	oldHead.prev = node
}

func (l *LRUCache) Get(key int) int {
	if node, ok := l.cache[key]; ok {
		l.remove(node)
		l.insert(node)
		return node.val
	}
	return -1
}

func (l *LRUCache) Put(key int, value int) {
	if node, ok := l.cache[key]; ok {
		node.val = value
		l.remove(node)
		l.insert(node)
	} else {
		newNode := &DoublyLinkedList{key: key, val: value}
		l.cache[key] = newNode
		l.insert(newNode)
		if len(l.cache) > l.capacity {
			tail := l.tail.prev
			l.remove(tail)
			delete(l.cache, tail.key)
		}
	}
}
