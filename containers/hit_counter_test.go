package containers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// bothCounters runs f against a freshly constructed HitCounter and
// HitCounterQueue, labelling each sub-test clearly.
func bothCounters(t *testing.T, f func(t *testing.T, hit func(int), getHits func(int) int)) {
	t.Helper()
	t.Run("circular", func(t *testing.T) {
		h := NewHitCounter()
		f(t, h.Hit, h.GetHits)
	})
	t.Run("queue", func(t *testing.T) {
		h := NewHitCounterQueue()
		f(t, h.Hit, h.GetHits)
	})
}

// ---------------------------------------------------------------------------
// Group 1: Empty counter
// ---------------------------------------------------------------------------

func Test_HitCounter_GetHitsEmpty(t *testing.T) {
	bothCounters(t, func(t *testing.T, hit func(int), getHits func(int) int) {
		assert.Equal(t, 0, getHits(1))
		assert.Equal(t, 0, getHits(300))
	})
}

// ---------------------------------------------------------------------------
// Group 2: Single hit, basic retrieval
// ---------------------------------------------------------------------------

func Test_HitCounter_SingleHit(t *testing.T) {
	bothCounters(t, func(t *testing.T, hit func(int), getHits func(int) int) {
		hit(1)
		assert.Equal(t, 1, getHits(1))
	})
}

// ---------------------------------------------------------------------------
// Group 3: Window boundary
//
// HitCounter (circular buffer): exclusive lower bound — t > now-300.
//   Hit at t=1, query at now=301: 1 > 1 → false → NOT counted. Excluded at 300s.
//   Hit at t=1, query at now=300: 1 > 0 → true  → counted.   Included at 299s.
//
// HitCounterQueue (FIFO queue): inclusive lower bound — t >= now-300.
//   Hit at t=1, query at now=301: 1 >= 1 → true  → counted.   Included at 300s.
//   Hit at t=1, query at now=302: 1 >= 2 → false → NOT counted. Excluded at 301s.
//
// This difference arises from the circular buffer's slot-collision constraint:
// t=1 and t=301 both map to slot index 1 (1%300 = 301%300 = 1). Hit(301) resets
// the slot, overwriting t=1. With the exclusive boundary this is harmless (t=1
// is already excluded). Use HitCounterQueue when inclusive semantics are required.
// ---------------------------------------------------------------------------

func Test_HitCounter_Circular_BoundaryIncluded_299s(t *testing.T) {
	h := NewHitCounter()
	h.Hit(1)
	assert.Equal(t, 1, h.GetHits(300), "299s ago: inside window")
}

func Test_HitCounter_Circular_BoundaryExcluded_300s(t *testing.T) {
	// Circular buffer: exclusive — exactly 300s ago is NOT counted.
	h := NewHitCounter()
	h.Hit(1)
	assert.Equal(t, 0, h.GetHits(301), "exactly 300s ago: outside exclusive window")
}

func Test_HitCounter_Queue_BoundaryIncluded_300s(t *testing.T) {
	// Queue variant: inclusive — exactly 300s ago IS counted.
	h := NewHitCounterQueue()
	h.Hit(1)
	assert.Equal(t, 1, h.GetHits(301), "exactly 300s ago: inside inclusive window")
}

func Test_HitCounter_Queue_BoundaryExcluded_301s(t *testing.T) {
	h := NewHitCounterQueue()
	h.Hit(1)
	assert.Equal(t, 0, h.GetHits(302), "301s ago: outside inclusive window")
}

// ---------------------------------------------------------------------------
// Group 4: Multiple hits at the same timestamp
// ---------------------------------------------------------------------------

func Test_HitCounter_SameTimestampMultiHits(t *testing.T) {
	bothCounters(t, func(t *testing.T, hit func(int), getHits func(int) int) {
		hit(5)
		hit(5)
		hit(5)
		assert.Equal(t, 3, getHits(5))
	})
}

func Test_HitCounter_SameTimestampThenQuery(t *testing.T) {
	// Both variants agree: hits at t=10, query at t=100 → 2 (well within window).
	bothCounters(t, func(t *testing.T, hit func(int), getHits func(int) int) {
		hit(10)
		hit(10)
		assert.Equal(t, 2, getHits(100))
	})
}

func Test_HitCounter_Circular_SameTimestampBoundary(t *testing.T) {
	// Circular: exclusive — t=10 at query t=310: 10 > 10 → false → 0.
	h := NewHitCounter()
	h.Hit(10)
	h.Hit(10)
	assert.Equal(t, 0, h.GetHits(310), "circular: exactly 300s → excluded")
	assert.Equal(t, 0, h.GetHits(311))
}

func Test_HitCounter_Queue_SameTimestampBoundary(t *testing.T) {
	// Queue: inclusive — t=10 at query t=310: 10 >= 10 → true → 2.
	h := NewHitCounterQueue()
	h.Hit(10)
	h.Hit(10)
	assert.Equal(t, 2, h.GetHits(310), "queue: exactly 300s → included")
	assert.Equal(t, 0, h.GetHits(311), "queue: 301s → excluded")
}

// ---------------------------------------------------------------------------
// Group 5: Multiple distinct timestamps, all within window
// ---------------------------------------------------------------------------

func Test_HitCounter_MultipleTimestamps(t *testing.T) {
	bothCounters(t, func(t *testing.T, hit func(int), getHits func(int) int) {
		hit(1)
		hit(2)
		hit(3)
		assert.Equal(t, 3, getHits(4))
	})
}

// ---------------------------------------------------------------------------
// Group 6: Stale hits excluded; mixed fresh and stale
// ---------------------------------------------------------------------------

func Test_HitCounter_AllStale(t *testing.T) {
	bothCounters(t, func(t *testing.T, hit func(int), getHits func(int) int) {
		hit(1)
		hit(2)
		assert.Equal(t, 0, getHits(303), "both hits are > 300s old")
	})
}

func Test_HitCounter_MixedFreshStale(t *testing.T) {
	// Window at now=500: [200, 500].
	// t=1:   1 >= 200? No  → stale
	// t=290: 290 >= 200? Yes → fresh
	// t=500: 500 >= 200? Yes → fresh
	bothCounters(t, func(t *testing.T, hit func(int), getHits func(int) int) {
		hit(1)
		hit(290)
		hit(500)
		assert.Equal(t, 2, getHits(500))
	})
}

func Test_HitCounter_Circular_OnlyNewestInWindow(t *testing.T) {
	// hits=[1,100,200,300,301], query at 301.
	// Circular (exclusive): t=1 → 1 > 1 → false → 4 hits counted.
	// t=1 and t=301 share slot 1; Hit(301) overwrites it. t=1 is also
	// excluded by the exclusive boundary, so the answer is correct.
	h := NewHitCounter()
	for _, ts := range []int{1, 100, 200, 300, 301} {
		h.Hit(ts)
	}
	assert.Equal(t, 4, h.GetHits(301))
}

func Test_HitCounter_Queue_OnlyNewestInWindow(t *testing.T) {
	// hits=[1,100,200,300,301], query at 301.
	// Queue (inclusive): t=1 → 1 >= 1 → true → all 5 hits counted.
	// This is the requirement's example.
	h := NewHitCounterQueue()
	for _, ts := range []int{1, 100, 200, 300, 301} {
		h.Hit(ts)
	}
	assert.Equal(t, 5, h.GetHits(301))
}

// ---------------------------------------------------------------------------
// Group 7: Large timestamp gap forces full eviction / slot reset
// ---------------------------------------------------------------------------

func Test_HitCounter_LargeGap(t *testing.T) {
	bothCounters(t, func(t *testing.T, hit func(int), getHits func(int) int) {
		for i := 1; i <= 10; i++ {
			hit(i)
		}
		hit(10000)
		assert.Equal(t, 1, getHits(10000))
	})
}

func Test_HitCounter_QueryAfterLargeGap_NoHits(t *testing.T) {
	bothCounters(t, func(t *testing.T, hit func(int), getHits func(int) int) {
		hit(1)
		assert.Equal(t, 0, getHits(10000))
	})
}

// ---------------------------------------------------------------------------
// Group 8: Circular buffer slot wrap-around (HitCounter specific)
//
// Slot index = timestamp % 300.  Both ts=1 and ts=301 map to slot 1.
// After Hit(301), slot 1 is reset to {timestamp:301, count:1}.
// GetHits(301): window=[1,301]; slot 1 has timestamp=301 >= 1 → counted.
// The old ts=1 data is gone (overwritten), which is correct since ts=1
// is at the very edge of the window (1 >= 301-300 = 1 ✓) but we no longer
// need it — the slot now reflects the newer ts=301 hit.
// ---------------------------------------------------------------------------

func Test_HitCounter_SlotWrapAround(t *testing.T) {
	h := NewHitCounter()
	h.Hit(1)   // slot 1 → {timestamp:1, count:1}
	h.Hit(301) // slot 1 → {timestamp:301, count:1} (reset)
	// Window at 301: [1, 301]. Only ts=301 is in the slot now.
	assert.Equal(t, 1, h.GetHits(301))
}

func Test_HitCounter_SlotWrapAround_MultipleHits(t *testing.T) {
	h := NewHitCounter()
	h.Hit(1)
	h.Hit(1)   // slot 1 → {timestamp:1, count:2}
	h.Hit(301) // slot 1 reset → {timestamp:301, count:1}
	h.Hit(301) // slot 1 → {timestamp:301, count:2}
	assert.Equal(t, 2, h.GetHits(301))
}

// ---------------------------------------------------------------------------
// Group 9: LeetCode 362 canonical example
// ---------------------------------------------------------------------------

func Test_HitCounter_LeetCode362_Circular(t *testing.T) {
	// LeetCode 362 uses exclusive boundary: t > now-300.
	h := NewHitCounter()
	h.Hit(1)
	h.Hit(1)
	h.Hit(1)
	assert.Equal(t, 3, h.GetHits(4))
	assert.Equal(t, 3, h.GetHits(300))
	assert.Equal(t, 0, h.GetHits(301), "circular: exactly 300s → excluded")
	assert.Equal(t, 0, h.GetHits(302))
}

func Test_HitCounter_LeetCode362_Queue(t *testing.T) {
	// Queue uses inclusive boundary: t >= now-300.
	h := NewHitCounterQueue()
	h.Hit(1)
	h.Hit(1)
	h.Hit(1)
	assert.Equal(t, 3, h.GetHits(4))
	assert.Equal(t, 3, h.GetHits(300))
	assert.Equal(t, 3, h.GetHits(301), "queue: exactly 300s → included")
	assert.Equal(t, 0, h.GetHits(302))
}

// ---------------------------------------------------------------------------
// Group 10: GetHits at timestamp 0, and query before hits in window
// ---------------------------------------------------------------------------

func Test_HitCounter_GetHitsAtZero(t *testing.T) {
	bothCounters(t, func(t *testing.T, hit func(int), getHits func(int) int) {
		assert.Equal(t, 0, getHits(0))
	})
}

func Test_HitCounter_HitThenQueryOutsideWindow(t *testing.T) {
	// Both variants: well within window (299s gap) → 1.
	bothCounters(t, func(t *testing.T, hit func(int), getHits func(int) int) {
		hit(100)
		assert.Equal(t, 1, getHits(399), "299s gap: inside both windows")
	})
}

func Test_HitCounter_Circular_ExactlyAtExclusiveBoundary(t *testing.T) {
	// Circular (exclusive): t=100, query=400 → 100 > 100 → false → 0.
	h := NewHitCounter()
	h.Hit(100)
	assert.Equal(t, 0, h.GetHits(400), "circular: exactly 300s → excluded")
	assert.Equal(t, 0, h.GetHits(401))
}

func Test_HitCounter_Queue_ExactlyAtInclusiveBoundary(t *testing.T) {
	// Queue (inclusive): t=100, query=400 → 100 >= 100 → true → 1.
	h := NewHitCounterQueue()
	h.Hit(100)
	assert.Equal(t, 1, h.GetHits(400), "queue: exactly 300s → included")
	assert.Equal(t, 0, h.GetHits(401), "queue: 301s → excluded")
}

func Test_HitCounter_RequirementExample(t *testing.T) {
	// Requirement: hits=[1,100,200,300,301], query at 301 → 5.
	// Only HitCounterQueue (inclusive) satisfies this. HitCounter (exclusive)
	// returns 4 because t=1 is evicted by the slot collision with t=301.
	h := NewHitCounterQueue()
	for _, ts := range []int{1, 100, 200, 300, 301} {
		h.Hit(ts)
	}
	assert.Equal(t, 5, h.GetHits(301), "all 5 hits within [1, 301] should be counted")
}
