package apidesign

/**
 * 232. Implement Queue using Stacks
 */
type MyQueue struct {
	inbox  []int // push stack
	outbox []int // pop stack
}

func ConstructorMyQueue() MyQueue {
	return MyQueue{}
}

func (q *MyQueue) Push(x int) {
	q.inbox = append(q.inbox, x)
}

func (q *MyQueue) Pop() int {
	q.transfer()
	val := q.outbox[len(q.outbox)-1]
	q.outbox = q.outbox[:len(q.outbox)-1]
	return val
}

func (q *MyQueue) Peek() int {
	q.transfer()
	return q.outbox[len(q.outbox)-1]
}

func (q *MyQueue) Empty() bool {
	return len(q.inbox) == 0 && len(q.outbox) == 0
}

// transfer moves all elements from inbox to outbox
// only when outbox is empty — preserving FIFO order
func (q *MyQueue) transfer() {
	if len(q.outbox) == 0 {
		for len(q.inbox) > 0 {
			q.outbox = append(q.outbox, q.inbox[len(q.inbox)-1])
			q.inbox = q.inbox[:len(q.inbox)-1]
		}
	}
}
