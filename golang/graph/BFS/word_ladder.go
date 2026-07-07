package bfs

import (
	"container/list"
)

/**
 * 127. Word Ladder
 *
 * A transformation sequence from word beginWord to word endWord using a dictionary wordList is a sequence of words
 */
func ladderLength(beginWord string, endWord string, wordList []string) int {
	dict := make(map[string]bool)
	for _, word := range wordList {
		dict[word] = true
	}

	type Step struct {
		Word   string
		Length int
	}
	queue := list.New()
	queue.PushBack(Step{beginWord, 1})

	for queue.Len() > 0 {
		front := queue.Remove(queue.Front()).(Step)
		word, length := front.Word, front.Length
		if word == endWord {
			return length
		}

		chArr := []rune(word)
		for i := 0; i < len(chArr); i++ { // search for next word in dict
			for ch := 'a'; ch <= 'z'; ch++ {
				if ch == rune(chArr[i]) {
					continue
				}

				chArr[i] = ch
				nextWord := string(chArr)
				if dict[nextWord] {
					queue.PushBack(Step{nextWord, length + 1})
					delete(dict, nextWord) // remove word from dict, avoid looping circle in the graph
				}
			}
			chArr = []rune(word) // reset to original word
		}
	}

	return 0
}
