package monotone

import . "gocode/types"

/**
 * 42. Trapping Rain Water
 *
 * Monotonic descending stack storing the index
 * processing all the heights on stack less eq to right height
 *
 *    S = height * length
 *      = (shorter of (left, right) - bottom) * length
 *      = (min(height[left], height[i]) - bottom) * (i - 1 - left)
 *             stack.top()                stack.top(), pop()
 *
 *   X
 *   X X       #
 *   X X       #
 *   X X X X . #
 *   X X X X X #
 *   -----------
 *         t b i
 *                  stack stores left side index such that
 *   X                  height[stack.Top()] >= right
 *   X X . . . #    otherwise, keep popping and adding water
 *   X X . . . #
 *   X X X X * #    monotonic stack, elements < right will be popped
 *   X X X X X #
 *   -----------
 *     t b     i
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
            length := i - 1 - stack.Top()
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
 * 84. Largest Rectangle in Histogram
 *
 *
 */
func largestRectangleArea(heights []int) int {
    res := 0

    return res
}
