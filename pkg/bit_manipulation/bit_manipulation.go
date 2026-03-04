package bit

/*
 *
 * 371. Sum of Two Integers
 *
 * a ^ b stands for the sum of each digit (w/o carry)
 * a & b stands for carry bit at each position
 *
 * add carry over (a & b) << 1 repeatedly into the higher bit
 * until no carry over occurred
 */
func getSum(a int, b int) int {
	sum, carry := a^b, a&b
	for carry != 0 {
		carry = carry << 1 // shift all bits to the left, to be added next
		sum, carry = sum^carry, sum&carry
	}
	return sum
}

/*
 * 191. Number of 1 Bits
 * Time Complexity:
 * If you treat the input as an arbitrary n and measure by its value,
 * the number of bits is about log₂(n), so you get O(log n) iterations.
 */
func hammingWeight(n int) int {
	count := 0
	for n != 0 {
		if n&1 == 1 {
			count += 1
		}
		n = n >> 1
	}

	return count
}

/*
 * 338. Counting Bits
 *
 * 0  --> 0    --> 0
 * 1  --> 1    --> 1
 * 2  --> 10   --> 1
 * 3  --> 11   --> 2
 * 4  --> 100  --> 1
 * 5  --> 101  --> 2
 * 6  --> 110  --> 2
 * 7  --> 111  --> 3
 * 8  --> 1000 --> 1
 * 9  --> 1001 --> 2
 * 10 --> 1010 --> 2
 * 11 --> 1011 --> 3
 * 12 --> 1100 --> 2
 * 13 --> 1101 --> 3
 */
func countBits(n int) []int {
	ans := make([]int, n+1)
	for i := 1; i <= n; i++ {
		ans[i] = ans[i>>1] + (i & 1)
	}
	return ans
}
