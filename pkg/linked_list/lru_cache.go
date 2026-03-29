package linked_list

import "container/list"

/**
 * 146. LRU Cache
 *
 * using container/list APIs
 *
 * LRU -> FIFO queue, PushFront(), Back(), Remove()
 *        FILO stack  PushBack(), Back(), Remove()
 */
type LRUCache struct {
	capcity int
	hashmap map[int]*list.Element
	list    *list.List
}

func Constructor(capacity int) LRUCache {
	return LRUCache{
		capcity: capacity,
		hashmap: make(map[int]*list.Element),
		list:    list.New(),
	}
}

type entry struct {
	key, val int
}

func (l *LRUCache) Get(key int) int {
	e, found := l.hashmap[key]
	if !found {
		return -1
	}

	l.list.Remove(e)
	newElem := l.list.PushFront(e.Value) // new element, new pointer
	l.hashmap[key] = newElem             // ← must update hashmap!
	return e.Value.(entry).val
}

func (l *LRUCache) Put(key int, value int) {
	if e, found := l.hashmap[key]; found {
		l.list.Remove(e)
		newElem := l.list.PushFront(entry{key: key, val: value})
		l.hashmap[key] = newElem
		return
	}

	if l.list.Len() == l.capcity {
		last := l.list.Back()
		l.list.Remove(last)
		delete(l.hashmap, last.Value.(entry).key)
	}

	item := l.list.PushFront(entry{key: key, val: value})
	l.hashmap[key] = item
}
