package two_pointers

import "bytes"

/**
 * 838. Push Dominoes
 *
 * You are given a string dominoes representing the initial state where:
 *
 * dominoes[i] = 'L', if the ith domino has been pushed to the left,
 * dominoes[i] = 'R', if the ith domino has been pushed to the right, and
 * dominoes[i] = '.', if the ith domino has not been pushed.
 * Return a string representing the final state.
 *
 */
// 4 cases
// L...L → all L
// R...R → all R
// L...R → all . (forces push outward, nothing moves)
// R...L → RRR...LLL (forces push inward, meet in middle)
func pushDominoes(dominoes string) string {
	s := "L" + dominoes + "R"
	var res []byte

	prev := 0
	for curr := 1; curr < len(s); curr++ {
		if s[curr] == '.' {
			continue
		}

		span := curr - prev - 1
		if prev > 0 {
			res = append(res, s[prev])
		}

		switch {
		case s[prev] == s[curr]:
			res = append(res, bytes.Repeat([]byte{s[prev]}, span)...)
		case s[prev] == 'L' && s[curr] == 'R':
			res = append(res, bytes.Repeat([]byte{'.'}, span)...)
		default: // R...L
			res = append(res, bytes.Repeat([]byte{'R'}, span/2)...)
			res = append(res, bytes.Repeat([]byte{'.'}, span%2)...)
			res = append(res, bytes.Repeat([]byte{'L'}, span/2)...)
		}
		prev = curr
	}
	return string(res)
}

func pushDominoesVisualized(dominoes string) string {
	n := len(dominoes)
	lforce, rforce := make([]int, n), make([]int, n)

	force := 0
	for i := 0; i < n; i++ {
		if dominoes[i] == '.' {
			force = max(0, force-1)
		} else if dominoes[i] == 'R' {
			force = n
		} else {
			force = 0
		}
		rforce[i] = force
	}

	force = 0
	for i := n - 1; i >= 0; i-- {
		if dominoes[i] == '.' {
			force = max(0, force-1)
		} else if dominoes[i] == 'L' {
			force = n
		} else {
			force = 0
		}
		lforce[i] = force
	}

	result := make([]byte, n)
	for i := range dominoes {
		if dominoes[i] == 'R' || dominoes[i] == 'L' {
			result[i] = dominoes[i]
		} else if rforce[i] > lforce[i] {
			result[i] = 'R'
		} else if rforce[i] < lforce[i] {
			result[i] = 'L'
		} else {
			result[i] = '.'
		}
	}
	return string(result)
}

func pushDominoesBFS(dominoes string) string {
	n := len(dominoes)
	lforce, rforce := make([]int, n), make([]int, n)

	force := 0
	for i := 0; i < n; i++ {
		if dominoes[i] == '.' {
			force = max(0, force-1)
		} else if dominoes[i] == 'R' {
			force = n
		} else {
			force = 0
		}
		rforce[i] = force
	}

	force = 0
	for i := n - 1; i >= 0; i-- {
		if dominoes[i] == '.' {
			force = max(0, force-1)
		} else if dominoes[i] == 'L' {
			force = n
		} else {
			force = 0
		}
		lforce[i] = force
	}

	result := make([]byte, n)
	for i := range dominoes {
		if dominoes[i] == 'R' || dominoes[i] == 'L' {
			result[i] = dominoes[i]
		} else if rforce[i] > lforce[i] {
			result[i] = 'R'
		} else if rforce[i] < lforce[i] {
			result[i] = 'L'
		} else {
			result[i] = '.'
		}
	}
	return string(result)
}
