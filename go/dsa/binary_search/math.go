package binarysearch

/*
 * 50. Pow(x, n)
 */
func myPow(x float64, n int) float64 {
	if n < 0 {
		x = 1 / x
		n = -n
	}

	result := 1.0
	for n > 0 {
		if n%2 == 1 {
			result *= x // current bit is set — multiply in x
		}
		x *= x // square x for next bit
		n /= 2 // shift to next bit
	}

	return result
}

func myPowRecursive(x float64, n int) float64 {
	if n < 0 {
		x = 1 / x
		n = -n
	}
	return pow(x, n)
}

func pow(x float64, n int) float64 {
	// base case
	if n == 0 {
		return 1
	}

	// recursive case: compute half power once, reuse it
	half := pow(x, n/2)

	if n%2 == 0 {
		return half * half // x^n = (x^(n/2))^2
	} else {
		return half * half * x // x^n = (x^(n/2))^2 * x
	}
}
