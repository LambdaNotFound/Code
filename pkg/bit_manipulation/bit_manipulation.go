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
