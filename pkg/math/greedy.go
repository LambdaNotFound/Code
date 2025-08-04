package math

/**
 * 621. Task Scheduler
 *
 * but there's a constraint: there has to be a gap of at least n intervals between two tasks with the same label.
 */
func leastInterval(tasks []byte, n int) int {
    maxFreq := 0
    dict := make(map[byte]int)

    for i := 0; i < len(tasks); i++ {
        dict[tasks[i]]++
        maxFreq = max(maxFreq, dict[tasks[i]])
    }

    maxFreqCount := 0
    // if there are tasks with equal freq, then time increases
    for _, value := range dict {
        if value == maxFreq {
            maxFreqCount++
        }
    }

    partCount := maxFreq - 1
    partLength := n - (maxFreqCount - 1) // idol slots in each part
    emptySlots := partCount * partLength
    availableTasks := len(tasks) - maxFreq*maxFreqCount
    idles := max(0, emptySlots-availableTasks)

    return len(tasks) + idles
}
