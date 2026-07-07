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
	list    *list.List
	capcity int
	cache   map[int]*list.Element
}

func Constructor(capacity int) LRUCache {
	return LRUCache{
		list:    list.New(),
		capcity: capacity,
		cache:   make(map[int]*list.Element),
	}
}

type entry struct {
	key, val int
}

func (l *LRUCache) Get(key int) int {
	if e, found := l.cache[key]; found {
		l.list.Remove(e)
		newElem := l.list.PushFront(e.Value) // new element, new pointer
		l.cache[key] = newElem               // must! Remove + PushFront destroy ptr
		return e.Value.(entry).val
	} else {
		return -1
	}
}

func (l *LRUCache) Put(key int, value int) {
	if e, found := l.cache[key]; found {
		l.list.Remove(e)
		newElem := l.list.PushFront(entry{key: key, val: value})
		l.cache[key] = newElem
	} else {
		item := l.list.PushFront(entry{key: key, val: value})
		l.cache[key] = item

		if l.list.Len() > l.capcity {
			last := l.list.Back()
			l.list.Remove(last)
			delete(l.cache, last.Value.(entry).key)
		}
	}
}
