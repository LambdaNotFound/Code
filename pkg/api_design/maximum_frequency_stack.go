package apidesign

/**
 * 895. Maximum Frequency Stack
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
