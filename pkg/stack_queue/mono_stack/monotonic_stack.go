package monostack

import . "gocode/containers"

/**
 * Monotonic stack template:
 *
 *    for loop over items {
 *        while !stack.empty() && stack.top() <> item{
 *            stack.pop()
 *            processing
 *        }
 *
 *        stack.push(item)
 *    }
 */

/**
 * 42. Trapping Rain Water <- Monotonic Descending Stack
 *
 * Monotonic descending stack storing the index
 * popping all the heights on stack
 * when hitting a right height > stack.top()
 *
 *    S = height * length
 *      = (shorter of (left, right) - bottom) * length
 *      = (min(height[left], height[i]) - bottom) * (i - 1 - left)
 *             stack.top()                stack.top(), pop()
 *
 *    X
 *    X X       #
 *    X X       #
 *    X X X X . #
 *    X X X X X #
 *    -----------
 *          t b i
 *                   stack stores left side index such that
 *    X                  height[stack.Top()] >= right
 *    X X . . . #    otherwise, keep popping and adding water
 *    X X . . . #
 *    X X X X * #    monotonic stack, elements < right will be popped
 *    X X X X X #
 *    -----------
 *      t b     i
 */
func trap(height []int) int {
	stack := Stack[int]{}
	res := 0
	for i := 0; i < len(height); i += 1 {
		right := height[i]
		for !stack.IsEmpty() && height[stack.Top()] < right {
			bottom := height[stack.Top()]
			stack.Pop()
			if stack.IsEmpty() {
				break
			}
			left := height[stack.Top()]
			width := i - 1 - stack.Top()
			if left > right {
				res += (right - bottom) * width
			} else {
				res += (left - bottom) * width
			}
		}
		stack = append(stack, i)
	}
	return res
}

func trapSlice(height []int) int {
	res := 0
	st := make([]int, 0) // monotonic decreasing stack of indices

	for i := 0; i < len(height); i++ {
		for len(st) > 0 && height[st[len(st)-1]] < height[i] {
			bottomIdx := st[len(st)-1]
			st = st[:len(st)-1] // pop bottom

			if len(st) == 0 {
				break // no left wall — can't trap water
			}

			leftIdx := st[len(st)-1]

			waterHeight := min(height[leftIdx], height[i]) - height[bottomIdx]
			waterWidth := i - leftIdx - 1
			res += waterHeight * waterWidth
		}
		st = append(st, i)
	}

	return res
}

/**
 * 84. Largest Rectangle in Histogram <- Monotonic Ascending Stack
 *
 * Given an array of integers heights representing the histogram's bar height
 * where the width of each bar is 1, return the area of the largest rectangle in the histogram.
 *
 * Monotonic ascending stack storing the index
 * popping all the heights on stack
 * when hitting a right height < stack.top()
 *
 *            A                        X
 *        X X A                        X
 *      X X X A #                X     X
 *    X X X X A #              X X     X x         monotonic stack, popping out all items > current
 *    -----------              -----------         between i, j, the values are larger than height[j]
 *          t h i                i     j cur
 *
 *            X
 *        X B B
 *      X X B B #     by def, when popping height h,
 *    X X X B B #     all the bars to the right, are >= h
 *    -----------
 *        t h   i
 *
 *            X
 *        C C C
 *      X C C C #
 *    X X C C C #
 *    -----------
 *      t h     i
 */
func largestRectangleAreaSlice(heights []int) int {
	st := make([]int, 0)
	res := 0
	heights = append(heights, 0) // sentinel to flush stack

	for i := 0; i < len(heights); i++ {
		for len(st) > 0 && heights[st[len(st)-1]] > heights[i] {
			height := heights[st[len(st)-1]]
			st = st[:len(st)-1] // pop first

			width := i       // empty stack → extends to index 0
			if len(st) > 0 { // st top is the left boundary
				width = i - st[len(st)-1] - 1 // i-1 - (st.top+1) + 1
			}
			res = max(res, height*width)
		}
		st = append(st, i)
	}

	return res
}

func largestRectangleArea(heights []int) int {
	heights = append(heights, 0)

	stack := Stack[int]{}
	res := 0
	for i := 0; i < len(heights); i += 1 {
		right := heights[i]
		for !stack.IsEmpty() && heights[stack.Top()] > right {
			height := heights[stack.Top()]
			stack.Pop()

			width := 0
			if stack.IsEmpty() {
				width = i
			} else {
				width = i - 1 - stack.Top()
			}

			area := height * width
			if res < area {
				res = area
			}
		}

		stack.Push(i)
	}

	return res
}

/**
 * 739. Daily Temperatures
 *
 * days to wait until next *warmer* day
 * Input: temperatures = [73,74,75,71,69,72,76,73]
 * Output: [1,1,4,2,1,1,0,0]
 */
func dailyTemperatures(temperatures []int) []int {
	res := make([]int, len(temperatures))
	st := []int{}
	for i := 0; i < len(temperatures); i++ {
		for len(st) > 0 && temperatures[st[len(st)-1]] < temperatures[i] {
			j := st[len(st)-1]
			st = st[:len(st)-1]
			res[j] = i - j
		}
		st = append(st, i)
	}
	return res
}

func dailyTemperaturesRightToLeft(temperatures []int) []int {
	res := make([]int, len(temperatures))
	st := []int{}
	for i := len(temperatures) - 1; i >= 0; i-- {
		for len(st) > 0 && temperatures[st[len(st)-1]] <= temperatures[i] {
			st = st[:len(st)-1]
		}
		if len(st) == 0 {
			res[i] = 0
		} else {
			res[i] = st[len(st)-1] - i
		}
		st = append(st, i)
	}
	return res
}

/**
 * 316. Remove Duplicate Letters
 *
 * smallest in lexicographical order:
 *     A string a is lexicographically smaller than a string b if in the first position where a and b differ,
 *     string a has a letter that appears earlier in the alphabet than the corresponding letter in b.
 *     If the first min(a.length, b.length) characters do not differ, then the shorter string is the lexicographically smaller one.
 *
 * Input: s = "cbacdcbc"
 * Output: "acdb"
 *
 */
func removeDuplicateLetters(s string) string {
	lastIndex := map[byte]int{}
	for i := range s {
		lastIndex[s[i]] = i
	}

	stack := []byte{}
	inStack := map[byte]bool{}
	for i := range s {
		if inStack[s[i]] {
			continue
		}

		for len(stack) > 0 && stack[len(stack)-1] > s[i] && lastIndex[stack[len(stack)-1]] > i {
			inStack[stack[len(stack)-1]] = false
			stack = stack[:len(stack)-1]
		}

		stack = append(stack, s[i])
		inStack[s[i]] = true

	}

	return string(stack)
}
