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

func leastInterval(tasks []byte, n int) int {
	freqMap := make(map[byte]int)
	for _, t := range tasks {
		freqMap[t]++
	}

	maxHeap := &TaskHeap{
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

type Task struct {
	name        byte
	freq        int
	availableAt int
}

type TaskHeap struct {
	items []Task
	less  func(Task, Task) bool
}

func (h *TaskHeap) Len() int           { return len(h.items) }
func (h *TaskHeap) Less(i, j int) bool { return h.less(h.items[i], h.items[j]) }
func (h *TaskHeap) Swap(i, j int)      { h.items[i], h.items[j] = h.items[j], h.items[i] }

func (h *TaskHeap) Push(x interface{}) {
	h.items = append(h.items, x.(Task))
}

func (h *TaskHeap) Pop() interface{} {
	item := h.items[len(h.items)-1]
	h.items = h.items[:len(h.items)-1]
	return item
}
