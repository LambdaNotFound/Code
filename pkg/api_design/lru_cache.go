package apidesign

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

func (this *LRUCache) remove(node *DoublyLinkedList) {
    prev := node.prev
    next := node.next
    prev.next = next
    next.prev = prev
}

func (this *LRUCache) insert(node *DoublyLinkedList) {
    headNext := this.head.next
    this.head.next = node
    node.prev = this.head
    node.next = headNext
    headNext.prev = node
}

func (this *LRUCache) Get(key int) int {
    if node, ok := this.cache[key]; ok {
        this.remove(node)
        this.insert(node)
        return node.val
    }
    return -1
}

func (this *LRUCache) Put(key int, value int) {
    if node, ok := this.cache[key]; ok {
        this.remove(node)
        node.val = value
        this.insert(node)
    } else {
        newNode := &DoublyLinkedList{key: key, val: value}
        this.cache[key] = newNode
        this.insert(newNode)
        if len(this.cache) > this.capacity {
            lru := this.tail.prev
            this.remove(lru)
            delete(this.cache, lru.key)
        }
    }
}
