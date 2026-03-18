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

	queue := []Task{}
	maxHeap := &TaskHeap{
		items: make([]Task, 0),
		less: func(a Task, b Task) bool {
			return a.freq > b.freq
		},
	}
	for k, v := range freqMap {
		task := Task{
			name:        k,
			freq:        v,
			availableAt: 0,
		}
		heap.Push(maxHeap, task)
	}

	time := 0
	for maxHeap.Len() != 0 || len(queue) != 0 {
		if maxHeap.Len() != 0 {
			task := heap.Pop(maxHeap).(Task)
			time++
			if task.freq-1 > 0 {
				task.freq--
				task.availableAt = time + n + 1
				queue = append(queue, task)
			}
		}

		if len(queue) != 0 {
			next := queue[0]
			if next.availableAt <= time {
				queue = queue[1:]
				heap.Push(maxHeap, next)
			}
		}

		if maxHeap.Len() == 0 && len(queue) != 0 {
			next := queue[0]
			time = next.availableAt
		}
	}

	return time
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

func (t *TaskHeap) Len() int {
	return len(t.items)
}

func (t *TaskHeap) Less(i int, j int) bool {
	return t.less(t.items[i], t.items[j])
}

func (t *TaskHeap) Swap(i int, j int) {
	t.items[i], t.items[j] = t.items[j], t.items[i]
}

func (t *TaskHeap) Push(task interface{}) {
	t.items = append(t.items, task.(Task))
}

func (t *TaskHeap) Pop() interface{} {
	item := t.items[len(t.items)-1]
	t.items = t.items[:len(t.items)-1]
	return item
}
