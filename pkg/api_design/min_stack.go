package apidesign

/**
 * 155. Min Stack
 *
 * MinStack() initializes the stack object.
 * void push(int val) pushes the element val onto the stack.
 * void pop() removes the element on the top of the stack.
 * int top() gets the top element of the stack.
 * int getMin() retrieves the minimum element in the stack.
 * You must implement a solution with O(1) time complexity for each function.
 */

type MinStack struct {
    stack    []int
    minStack []int
}

func ConstructorMinStack() MinStack {
    return MinStack{
        stack:    []int{},
        minStack: []int{},
    }
}

func (m *MinStack) Push(val int) {
    m.stack = append(m.stack, val)

    len := len(m.minStack)
    if len == 0 || m.minStack[len-1] > val {
        m.minStack = append(m.minStack, val)
    } else {
        m.minStack = append(m.minStack, m.minStack[len-1])
    }
}

func (m *MinStack) Pop() {
    len := len(m.stack)

    m.stack = m.stack[:len-1]
    m.minStack = m.minStack[:len-1]
}

func (m *MinStack) Top() int {
    len := len(m.stack)

    return m.stack[len-1]
}

func (m *MinStack) GetMin() int {
    len := len(m.minStack)

    return m.minStack[len-1]
}
