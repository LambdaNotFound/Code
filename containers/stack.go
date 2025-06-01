package containers

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

/**
 * 232. Implement Queue using Stacks
 */

 type MyQueue struct {
    stack1 []int // Main stack for push operations
    stack2 []int // Temporary stack for reverse operations
}

func Constructor() MyQueue {
    return MyQueue{
        stack1: []int{},
        stack2: []int{},
    }
}

func (this *MyQueue) Push(x int) {
    this.stack1 = append(this.stack1, x) // Push directly to stack1
}

func (this *MyQueue) Pop() int {
    if len(this.stack2) == 0 {
        for len(this.stack1) > 0 {
            this.stack2 = append(this.stack2, this.stack1[len(this.stack1)-1])
            this.stack1 = this.stack1[:len(this.stack1)-1]
        }
    }
    val := this.stack2[len(this.stack2)-1]
    this.stack2 = this.stack2[:len(this.stack2)-1]
    return val
}

func (this *MyQueue) Peek() int {
    if len(this.stack2) == 0 {
        for len(this.stack1) > 0 {
            this.stack2 = append(this.stack2, this.stack1[len(this.stack1)-1])
            this.stack1 = this.stack1[:len(this.stack1)-1]
        }
    }
    return this.stack2[len(this.stack2)-1]
}

func (this *MyQueue) Empty() bool {
    return len(this.stack1) == 0 && len(this.stack2) == 0
}