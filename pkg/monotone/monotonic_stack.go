package monotone

/**
 * 42. Trapping Rain Water
 *
 * Monotonic descending stack, storing the index
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
 *
 *   X
 *   X X . . . #
 *   X X . . . #
 *   X X X X * #
 *   X X X X X #
 *   -----------
 *     t b     i
 */
func trap(height []int) int {
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
