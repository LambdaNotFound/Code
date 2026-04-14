package containers

// hitWindow is the size of the sliding window in seconds.
const hitWindow = 300

/**
 * HitCounter — Circular-buffer hit counter (LeetCode 362)
 *
 * Counts hits in a sliding 300-second window using a fixed array of 300 slots.
 * Each slot represents one second bucket: slot i holds hits for any timestamp t
 * where t % 300 == i. A slot is valid only when its stored timestamp matches
 * the incoming timestamp; a mismatch means the slot is from a prior cycle and
 * is silently reset on the next write.
 *
 * Window invariant:
 *   A hit at time t is counted at query time now iff t > now-300 (exclusive).
 *   Hits exactly 300 seconds before now are NOT counted.
 *   This matches LeetCode 362's definition exactly.
 *
 * Boundary limitation:
 *   The slot index is timestamp%300, so timestamps 300 seconds apart share a slot
 *   (e.g. t=1 and t=301 both map to slot 1). When Hit(301) fires, it resets slot 1,
 *   overwriting t=1. With the exclusive boundary this is safe — t=1 at query now=301
 *   is already excluded (1 > 1 is false). With an inclusive boundary it would lose a
 *   valid hit. Use HitCounterQueue if you need inclusive semantics.
 *
 * Implementation notes:
 *   - Zero value of [300]hitSlot has timestamp=0, count=0. The zero timestamp
 *     is never falsely counted because count is 0 until a real Hit writes it.
 *   - NewHitCounter is provided for API consistency with the rest of this package.
 *   - Timestamps are assumed to be monotonically non-decreasing (problem guarantee).
 *
 * Complexity:
 *   Hit      O(1)
 *   GetHits  O(300) — fixed 300-iteration scan regardless of input
 *   Space    O(300) — fixed-size array, no heap allocation beyond the struct
 */

// hitSlot is a single bucket in the circular buffer.
type hitSlot struct {
	timestamp int
	count     int
}

// HitCounter counts hits within a sliding 300-second window using a circular buffer.
// The zero value is usable; use NewHitCounter for API consistency.
type HitCounter struct {
	slots [hitWindow]hitSlot
}

// NewHitCounter returns a ready-to-use HitCounter.
func NewHitCounter() *HitCounter {
	return &HitCounter{}
}

// Hit records a hit at the given timestamp.
// If the slot for this timestamp belongs to a prior cycle it is reset first.
// Time: O(1).
func (h *HitCounter) Hit(timestamp int) {
	i := timestamp % hitWindow
	if h.slots[i].timestamp == timestamp {
		h.slots[i].count++
	} else {
		h.slots[i] = hitSlot{timestamp: timestamp, count: 1}
	}
}

// GetHits returns the number of hits in (timestamp-300, timestamp] (exclusive lower bound).
// Hits exactly 300 seconds before timestamp are not counted.
// Time: O(300).
func (h *HitCounter) GetHits(timestamp int) int {
	total := 0
	for _, s := range h.slots {
		if s.timestamp > timestamp-hitWindow {
			total += s.count
		}
	}
	return total
}

/**
 * HitCounterQueue — FIFO-queue hit counter (LeetCode 362, queue variant)
 *
 * Counts hits in a sliding 300-second window using a FIFO queue where each
 * element is the raw timestamp of one hit. GetHits lazily evicts stale entries
 * (those with timestamp < now-300) from the front before returning the length.
 *
 * Window invariant:
 *   Same as HitCounter: t >= now-300 (inclusive lower bound).
 *   The drain condition is: front < now-300 (strict), so front == now-300 is kept.
 *
 * Implementation notes:
 *   - Backed by Queue[int] from queue.go (slice-based FIFO).
 *   - Multiple hits at the same timestamp each occupy a separate queue slot,
 *     giving exact counts at the cost of O(k) space for k same-second hits.
 *   - Each timestamp is enqueued once and dequeued at most once across all calls,
 *     so GetHits is O(1) amortized (O(k) worst-case for k stale entries to drain).
 *   - The monotonic timestamp guarantee means once the front is within the window,
 *     all subsequent entries are also within the window — no further scan needed.
 *
 * Complexity:
 *   Hit      O(1)
 *   GetHits  O(k) amortized, k = stale entries drained; O(1) if none
 *   Space    O(n), n = hits currently within the window
 */

// HitCounterQueue counts hits within a sliding 300-second window using a FIFO queue.
// The zero value is not usable; use NewHitCounterQueue.
type HitCounterQueue struct {
	q Queue[int]
}

// NewHitCounterQueue returns a ready-to-use HitCounterQueue.
func NewHitCounterQueue() *HitCounterQueue {
	return &HitCounterQueue{}
}

// Hit records a hit at the given timestamp.
// Time: O(1).
func (h *HitCounterQueue) Hit(timestamp int) {
	h.q.Enqueue(timestamp)
}

// GetHits returns the number of hits in [timestamp-300, timestamp] (inclusive).
// Stale entries (timestamp < now-300) are drained from the front lazily.
// Time: O(k) amortized, where k is the number of stale entries drained.
func (h *HitCounterQueue) GetHits(timestamp int) int {
	cutoff := timestamp - hitWindow // entries with t >= cutoff are counted
	for !h.q.IsEmpty() && h.q.Front() < cutoff {
		h.q.Dequeue()
	}
	return h.q.Size()
}
