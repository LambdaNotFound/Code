package apidesign

/**
 * 146. LRU Cache
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
    cache    map[int]*DoublyLinkedList
}

func Constructor(capacity int) LRUCache {
    head := &DoublyLinkedList{key: 0, val: 0}
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
    prev := node.prev
    next := node.next
    prev.next = next
    next.prev = prev
}

func (l *LRUCache) insert(node *DoublyLinkedList) {
    headNext := l.head.next
    l.head.next = node
    node.prev = l.head
    node.next = headNext
    headNext.prev = node
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
        l.remove(node)
        node.val = value
        l.insert(node)
    } else {
        newNode := &DoublyLinkedList{key: key, val: value}
        l.cache[key] = newNode
        l.insert(newNode)
        if len(l.cache) > l.capacity {
            lru := l.tail.prev
            l.remove(lru)
            delete(l.cache, lru.key)
        }
    }
}
