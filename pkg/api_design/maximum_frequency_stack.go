package apidesign

/**
 * 895. Maximum Frequency Stack
 *
 * Design a stack-like data structure to push elements to the stack and pop the most frequent element from the stack. Implement the FreqStack class:
 *
 * FreqStack() constructs an empty frequency stack.
 * void push(int val) pushes an integer val onto the top of the stack.
 * int pop() removes and returns the most frequent element in the stack. If there is a tie for the most frequent element, the element closest to the stack's top is removed and returned. LeetCode
 */

type FreqStack struct {
	freq    map[int]int   // val → its current frequency
	group   map[int][]int // frequency → stack of vals at that freq
	maxFreq int
}

func ConstructorFreqStack() FreqStack {
	return FreqStack{
		freq:  make(map[int]int),
		group: make(map[int][]int),
	}
}

func (f *FreqStack) Push(val int) {
	f.freq[val]++
	curFreq := f.freq[val]
	f.group[curFreq] = append(f.group[curFreq], val)
	f.maxFreq = max(f.maxFreq, curFreq)
}

func (f *FreqStack) Pop() int {
	stack := f.group[f.maxFreq]
	top := stack[len(stack)-1]

	stack = stack[:len(stack)-1]
	f.group[f.maxFreq] = stack
	if len(stack) == 0 {
		f.maxFreq--
	}

	f.freq[top]--
	return top
}
