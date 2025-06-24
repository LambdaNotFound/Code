package graph

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

func ladderLengthBidirectionalBFS(beginWord string, endWord string, wordList []string) int {
    wordSet := make(map[string]bool)
    for _, word := range wordList {
        wordSet[word] = true
    }
    if _, ok := wordSet[endWord]; !ok {
        return 0
    }

    queueFromStart, queueFromEnd := []string{beginWord}, []string{endWord}
    visitedFromStart, visitedFromEnd := make(map[string]bool), make(map[string]bool)
    visitedFromStart[beginWord] = true
    visitedFromEnd[endWord] = true
    
    level := 1
    for len(queueFromStart) > 0 && len(queueFromEnd) > 0 {
        if len(queueFromStart) > len(queueFromEnd) { // swap
            queueFromStart, queueFromEnd = queueFromEnd, queueFromStart
            visitedFromStart, visitedFromEnd = visitedFromEnd, visitedFromStart
        }

        size := len(queueFromStart)
        for i := 0; i < size; i++ {
            word := queueFromStart[0]
            queueFromStart = queueFromStart[1:]

            chArr := []rune(word)
            for j := 0; j < len(chArr); j++ {
                for ch := 'a'; ch <= 'z'; ch++ {
                    if ch == rune(chArr[j]) {
                        continue
                    }
                    chArr[j] = ch
                    nextWord := string(chArr)
                    if _, ok := wordSet[nextWord]; ok && !visitedFromStart[nextWord] {
                        if visitedFromEnd[nextWord] { // find a sequence to endWord
                            return level + 1
                        }

                        visitedFromStart[nextWord] = true
                        queueFromStart = append(queueFromStart, nextWord)
                    }
                }
                chArr = []rune(word) // reset to original word
            }
        } // finished a level
        level++
    }

    return 0
}
