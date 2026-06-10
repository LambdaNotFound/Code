package apidesign

import "math"

/**
 * 244. Shortest Word Distance II
 * Design a data structure that will be initialized with a string array, and then it should answer queries of the shortest distance between two different words.
 * Implement the WordDistance class:
 *
 * WordDistance(String[] wordsDict) initializes the object with the string array
 * int shortest(String word1, String word2) returns the shortest distance between the two words in the array
 */

type WordDistance struct {
	wordIndex map[string][]int
}

func NewWordDistance(wordDict []string) *WordDistance {
	index := make(map[string][]int)

	for idx, word := range wordDict {
		index[word] = append(index[word], idx)
	}
	return &WordDistance{
		wordIndex: index,
	}
}

func (w *WordDistance) shortest(word1, word2 string) int {
	minDist := math.MaxInt

	abs := func(a int) int {
		if a < 0 {
			return -a
		}
		return a
	}

	list1, list2 := w.wordIndex[word1], w.wordIndex[word2]
	ptr1, ptr2 := 0, 0
	for ptr1 < len(list1) && ptr2 < len(list2) {
		a, b := list1[ptr1], list2[ptr2]

		minDist = min(minDist, abs(a-b))
		if a > b {
			ptr2++
		} else {
			ptr1++
		}
	}

	return minDist
}
