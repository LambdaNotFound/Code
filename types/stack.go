package types

type Stack[T any] []T

func (s *Stack[T]) IsEmpty() bool {
    return len(*s) == 0
}

func (s *Stack[T]) Push(item T) {
    *s = append(*s, item)
}

func (s *Stack[T]) Pop() (T, bool) {
    if s.IsEmpty() {
        var empty T
        return empty, false
    } else {
        idx := len(*s) - 1
        item := (*s)[idx]
        *s = (*s)[:idx]
        return item, true
    }
}

func (s *Stack[T]) Peek() (T, bool) {
    if s.IsEmpty() {
        var empty T
        return empty, false
    } else {
        index := len(*s) - 1
        return (*s)[index], true
    }
}

func (s *Stack[T]) Top() T {
    index := len(*s) - 1
    return (*s)[index]
}
