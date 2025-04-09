package monotone

import . "gocode/types"

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

func trap_slice(height []int) int {
    stack := []int{}
    res := 0
    for i := 0; i < len(height); i += 1 {
        right := height[i]
        for len(stack) != 0 && height[stack[len(stack)-1]] < right {
            bottom := height[stack[len(stack)-1]]
            stack = stack[:len(stack)-1]
            if len(stack) == 0 {
                break
            }
            left := height[stack[len(stack)-1]]
            length := i - 1 - stack[len(stack)-1]
            if left > right {
                res += (right - bottom) * length
            } else {
                res += (left - bottom) * length
            }
        }
        stack = append(stack, i)
    }
    return res
}

/**
 * 84. Largest Rectangle in Histogram <- Monotonic Ascending Stack
 *
 * Monotonic ascending stack storing the index
 * popping all the heights on stack
 * when hitting a right height < stack.top()
 *
 *            A
 *        X X A
 *      X X X A #
 *    X X X X A #
 *    -----------
 *          t h i
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

func largestRectangleArea_slice(heights []int) int {
    heights = append(heights, 0)

    stack := []int{}
    res := 0
    for i := 0; i < len(heights); i += 1 {
        right := heights[i]
        for len(stack) != 0 && heights[stack[len(stack)-1]] > right {
            height := heights[stack[len(stack)-1]]
            stack = stack[:len(stack)-1]

            width := 0
            if len(stack) == 0 {
                width = i
            } else {
                width = i - 1 - stack[len(stack)-1]
            }

            area := height * width
            if res < area {
                res = area
            }
        }

        stack = append(stack, i)
    }

    return res
}
