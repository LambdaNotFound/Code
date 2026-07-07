package interview

/**
Design a transaction logging system with two APIs:

Record a transaction with its amount and timestamp
Get the sum of all transactions from the last 60 minutes
Requirements:

Timestamps should have second-level precision
Both APIs will be called frequently and should have O(1) time complexity
*/

type Counter struct {
	queue  []int           // just timestamps, no need for a pair type
	counts map[int]float64 // count each timestamp bucket

	total float64
}

func NewCounter() Counter {
	return Counter{
		queue:  make([]int, 0),
		counts: make(map[int]float64),
	}
}

// Time: O(1)
// Space: O(W) W = number of distinct timestamps
func (c *Counter) putTransaction(amount float64, timestamp int) {
	c.counts[timestamp] += amount
	c.total += amount

	c.queue = append(c.queue, timestamp)
}

// Time: amortized O(1)
func (c *Counter) getTotalTransactionInLastOneHour(timestamp int) float64 {
	for len(c.queue) > 0 && c.queue[0] < timestamp-3600 {
		c.total -= c.counts[c.queue[0]] // look up live count
		delete(c.counts, c.queue[0])    // clean up map entry
		c.queue = c.queue[1:]
	}
	return c.total
}

/*
// For sparse traffic, the queue wins in practice;
// for dense or bursty traffic, the circular buffer's predictable latency matters more.
type slot struct {
	ts  int
	sum float64
}

type Counter struct {
	buf   [3600]slot
	total float64
}

func (c *Counter) putTransaction(amount float64, timestamp int) {
	s := &c.buf[timestamp%3600]
	if s.ts != timestamp {
		c.total -= s.sum // evict stale slot
		s.sum = 0
		s.ts = timestamp
	}
	s.sum += amount
	c.total += amount
}

func (c *Counter) getTotalTransactionInLastOneHour(timestamp int) float64 {
	total := 0.0
	for i := range c.buf {
		if c.buf[i].ts > timestamp-3600 {
			total += c.buf[i].sum
		}
	}
	return total
}
*/
