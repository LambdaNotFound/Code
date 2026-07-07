package apidesign

/**
 * 155. Min Stack
 *
 * MinStack() initializes the stack object.
 *     void push(int val) pushes the element val onto the stack.
 *     void pop() removes the element on the top of the stack.
 *     int top() gets the top element of the stack.
 *     int getMin() retrieves the minimum element in the stack.
 *
 * You must implement a solution with O(1) time complexity for each function.
 */

type MinStack struct {
	data []int // elements stack
	mins []int // minimums stack
}

func ConstructorMinStack() MinStack {
	return MinStack{}
}

func (s *MinStack) Push(val int) {
	s.data = append(s.data, val)
	if len(s.mins) == 0 || val <= s.GetMin() {
		s.mins = append(s.mins, val)
	}
}

func (s *MinStack) Pop() {
	if s.Top() == s.GetMin() {
		s.mins = s.mins[:len(s.mins)-1]
	}
	s.data = s.data[:len(s.data)-1]
}

func (s *MinStack) Top() int {
	return s.data[len(s.data)-1]
}

func (s *MinStack) GetMin() int {
	return s.mins[len(s.mins)-1]
}
