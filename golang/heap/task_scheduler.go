package heap

/**
 * 621. Task Scheduler
 *
 * Input: tasks = ["A","A","A","B","B","B"], n = 2
 *
 * Output: 8
 *
 * Heap + Queue
 * Heap to track the next highest freq task
 * Queue to queue(FIFO) the tasks
 *
 */
import "container/heap"

type Task struct {
	name        byte
	freq        int
	availableAt int
}

func leastInterval(tasks []byte, n int) int {
	freq := make(map[byte]int)
	for _, t := range tasks {
		freq[t]++
	}

	maxHeap := &Heap[Task]{
		less: func(a, b Task) bool { return a.freq > b.freq },
	}
	for name, count := range freq {
		heap.Push(maxHeap, Task{name: name, freq: count})
	}

	cooldown := []Task{}
	elapsed := 0
	for len(cooldown) > 0 || maxHeap.Len() > 0 {
		for len(cooldown) > 0 && cooldown[0].availableAt <= elapsed {
			heap.Push(maxHeap, cooldown[0])
			cooldown = cooldown[1:]
		}

		if maxHeap.Len() > 0 {
			cur := heap.Pop(maxHeap).(Task)
			if cur.freq > 1 {
				cooldown = append(cooldown, Task{
					name:        cur.name,
					freq:        cur.freq - 1,
					availableAt: elapsed + n + 1,
				})
			}
		}

		elapsed++
	}

	return elapsed
}

func leastIntervalCalude(tasks []byte, n int) int {
	type Task struct {
		name        byte
		freq        int
		availableAt int
	}

	freqMap := make(map[byte]int)
	for _, t := range tasks {
		freqMap[t]++
	}

	maxHeap := &Heap[Task]{
		items: make([]Task, 0, len(freqMap)),
		less:  func(a, b Task) bool { return a.freq > b.freq },
	}
	for name, freq := range freqMap {
		heap.Push(maxHeap, Task{name: name, freq: freq})
	}

	cooldownQueue := make([]Task, 0)
	currentTime := 0

	for maxHeap.Len() > 0 || len(cooldownQueue) > 0 {
		// Step 1: move cooled-down tasks back to heap FIRST
		if len(cooldownQueue) > 0 && cooldownQueue[0].availableAt <= currentTime {
			heap.Push(maxHeap, cooldownQueue[0])
			cooldownQueue = cooldownQueue[1:]
		}

		// Step 2: schedule highest frequency available task
		if maxHeap.Len() > 0 {
			task := heap.Pop(maxHeap).(Task)
			currentTime++
			if task.freq-1 > 0 {
				task.freq--
				task.availableAt = currentTime + n
				cooldownQueue = append(cooldownQueue, task)
			}
		}

		// Step 3: no task available — jump to next available task
		if maxHeap.Len() == 0 && len(cooldownQueue) > 0 {
			currentTime = cooldownQueue[0].availableAt
		}
	}

	return currentTime
}
