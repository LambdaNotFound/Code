package oodesign

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
 * Iterator pattern
 *
 * Provides sequential access to a collection's elements without exposing
 * how that collection stores them — the caller only ever calls HasNext/Next
 * and never needs to know it's a map, a slice, a tree, or anything else.
 *
 * This is the classic "external iterator" shape (caller pulls one element
 * at a time, holding the traversal position between calls). Idiomatic Go
 * usually reaches for `range` directly, or — since Go 1.23 — a range-over-func
 * `iter.Seq`/`iter.Seq2` "internal iterator" instead; see
 * containers/treemap_iter.go for that style applied to a red-black tree.
 * The external form here still earns its keep whenever the caller needs to
 * pause a traversal and resume it later, across separate function calls,
 * rather than staying inside a single `for range` loop body.
 */

// Iterator is what every concrete iterator implements.
type Iterator interface {
	HasNext() bool
	Next() int
}

// Collection is the aggregate interface: anything that can produce an
// Iterator over itself.
type Collection interface {
	CreateIterator() Iterator
}

// IntSet is a concrete collection backed by a map — deliberately an
// unordered structure, to make the point that the iterator (not the caller)
// owns how a traversal order gets derived from internal storage.
type IntSet struct {
	items map[int]struct{}
}

func NewIntSet(values ...int) *IntSet {
	items := make(map[int]struct{}, len(values))
	for _, v := range values {
		items[v] = struct{}{}
	}
	return &IntSet{items: items}
}

func (s *IntSet) CreateIterator() Iterator {
	keys := make([]int, 0, len(s.items))
	for k := range s.items {
		keys = append(keys, k)
	}
	sort.Ints(keys) // impose a deterministic order; the caller never sees the map
	return &intSetIterator{keys: keys}
}

type intSetIterator struct {
	keys []int
	pos  int
}

func (it *intSetIterator) HasNext() bool { return it.pos < len(it.keys) }

func (it *intSetIterator) Next() int {
	v := it.keys[it.pos]
	it.pos++
	return v
}

// IntList is a second concrete collection, backed by a slice instead of a
// map, to show that client code iterating via the Iterator interface never
// needs to change based on which storage the collection uses underneath.
type IntList struct {
	items []int
}

func NewIntList(values ...int) *IntList {
	return &IntList{items: values}
}

func (l *IntList) CreateIterator() Iterator {
	return &intListIterator{items: l.items}
}

type intListIterator struct {
	items []int
	pos   int
}

func (it *intListIterator) HasNext() bool { return it.pos < len(it.items) }

func (it *intListIterator) Next() int {
	v := it.items[it.pos]
	it.pos++
	return v
}

// collectAll drains an Iterator from any Collection — it works identically
// regardless of whether c is an IntSet or an IntList.
func collectAll(c Collection) []int {
	var out []int
	it := c.CreateIterator()
	for it.HasNext() {
		out = append(out, it.Next())
	}
	return out
}

func Test_IntSet_iteratesAllElementsInSortedOrder(t *testing.T) {
	set := NewIntSet(3, 1, 2)

	assert.Equal(t, []int{1, 2, 3}, collectAll(set))
}

func Test_IntList_iteratesInInsertionOrder(t *testing.T) {
	list := NewIntList(3, 1, 2)

	assert.Equal(t, []int{3, 1, 2}, collectAll(list))
}

func Test_Iterator_emptyCollectionHasNoNext(t *testing.T) {
	set := NewIntSet()

	it := set.CreateIterator()

	assert.False(t, it.HasNext())
}

func Test_Iterator_HasNextStaysFalseOnceDrained(t *testing.T) {
	it := NewIntList(1).CreateIterator()

	assert.True(t, it.HasNext())
	assert.Equal(t, 1, it.Next())

	assert.False(t, it.HasNext())
	assert.False(t, it.HasNext()) // calling it again doesn't flip it back
}

func Test_Collection_callerNeverNamesAConcreteCollection(t *testing.T) {
	collections := []Collection{
		NewIntSet(2, 1),
		NewIntList(5, 4),
	}

	var all [][]int
	for _, c := range collections {
		all = append(all, collectAll(c))
	}

	assert.Equal(t, [][]int{{1, 2}, {5, 4}}, all)
}
