package math

/**
 * 1884. Egg Drop With 2 Eggs and N Floors
 *
 * 2 eggs, Egg 1 is your "probe" — once it breaks, it's gone Egg 2 must be used linearly from the last safe floor upward
 * — because you can't afford to break it
 *
 */
func twoEggDrop(n int) int {
	// find minimum t such that t*(t+1)/2 >= n
	t := 1
	for t*(t+1)/2 < n {
		t++
	}
	return t
}
