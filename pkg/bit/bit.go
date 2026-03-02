package bit

/*
*
* 371. Sum of Two Integers
* After understanding that a ^ b stands for the sum of each digit
* (ignoring carry over) and a & b stands for whether carry over occurs on each digit,
* you will easily see that you need to add carry over (a & b) << 1
* repeatedly until no carry over occurred
 */
func getSum(a int, b int) int {
	sum, carryover := a^b, a&b
	for carryover != 0 {
		carryover = carryover << 1
		sum, carryover = sum^carryover, sum&carryover
	}
	return sum
}
