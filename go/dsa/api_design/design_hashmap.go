package apidesign

/**
 * 705. Design HashSet
 *
 */
type MyHashSet struct {
	buckets    [][]int
	numBuckets int
}

func NewMyHashSet() MyHashSet {
	n := 1009
	return MyHashSet{
		buckets:    make([][]int, n),
		numBuckets: n,
	}
}

func (s *MyHashSet) hash(key int) int {
	return key % s.numBuckets
}

func (s *MyHashSet) Add(key int) {
	if s.Contains(key) {
		return // no duplicates in a set
	}
	idx := s.hash(key)
	s.buckets[idx] = append(s.buckets[idx], key)
}

func (s *MyHashSet) Remove(key int) {
	idx := s.hash(key)
	for i, k := range s.buckets[idx] {
		if k == key {
			s.buckets[idx] = append(s.buckets[idx][:i], s.buckets[idx][i+1:]...)
			return
		}
	}
}

func (s *MyHashSet) Contains(key int) bool {
	for _, k := range s.buckets[s.hash(key)] {
		if k == key {
			return true
		}
	}
	return false
}

/**
 * 706. Design HashMap
 */
type mapEntry struct {
	key   int
	value int
}

type MyHashMap struct {
	buckets    [][]mapEntry
	numBuckets int
}

func NewMyHashMap() MyHashMap {
	n := 1009
	return MyHashMap{
		buckets:    make([][]mapEntry, n),
		numBuckets: n,
	}
}

func (m *MyHashMap) hash(key int) int {
	return key % m.numBuckets
}

func (m *MyHashMap) Put(key int, value int) {
	idx := m.hash(key)
	// If key exists, update in place — don't append a duplicate.
	for i, e := range m.buckets[idx] {
		if e.key == key {
			m.buckets[idx][i].value = value
			return
		}
	}
	m.buckets[idx] = append(m.buckets[idx], mapEntry{key, value})
}

func (m *MyHashMap) Get(key int) int {
	for _, e := range m.buckets[m.hash(key)] {
		if e.key == key {
			return e.value
		}
	}
	return -1
}

func (m *MyHashMap) Remove(key int) {
	idx := m.hash(key)
	for i, e := range m.buckets[idx] {
		if e.key == key {
			m.buckets[idx] = append(m.buckets[idx][:i], m.buckets[idx][i+1:]...)
			return
		}
	}
}
