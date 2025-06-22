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

    type Item struct {
        Word string
        Step int
    }
    queue := list.New()
    queue.PushBack(Item{beginWord, 1})

    for queue.Len() > 0 {
        front := queue.Remove(queue.Front()).(Item)
        word, length := front.Word, front.Step

        if word == endWord {
            return length
        }

        wordArr := []rune(word)
        for i := 0; i < len(wordArr); i++ {
            original := wordArr[i]
            for ch := 'a'; ch <= 'z'; ch++ {
                if ch  == rune(wordArr[i]) {
                    continue
                }
                wordArr[i] = ch
                newWord := string(wordArr)
                if dict[newWord] {
                    queue.PushBack(Item{newWord, length + 1})
                    delete(dict, newWord) // remove otherwise loop in the BFS
                }
            }
            wordArr[i] = original
        }
    }
    return 0
}

func ladderLengthBiDirectionalBFS(beginWord string, endWord string, wordList []string) int {
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
        if len(queueFromStart) > len(queueFromEnd) {
            queueFromStart, queueFromEnd = queueFromEnd, queueFromStart
            visitedFromStart, visitedFromEnd = visitedFromEnd, visitedFromStart
        }
        size := len(queueFromStart)
        for i := 0; i < size; i++ {
            currentWord := queueFromStart[0]
            queueFromStart = queueFromStart[1:]
            for j := 0; j < len(currentWord); j++ {
                for k := 'a'; k <= 'z'; k++ {
                    if rune(currentWord[j]) == k {
                        continue
                    }
                    nextWord := currentWord[:j] + string(k) + currentWord[j+1:]
                    if _, ok := wordSet[nextWord]; ok && !visitedFromStart[nextWord] {
                        if visitedFromEnd[nextWord] {
                            return level + 1
                        }
                        visitedFromStart[nextWord] = true
                        queueFromStart = append(queueFromStart, nextWord)
                    }
                }
            }
        }
        level++
    }
    return 0
}
