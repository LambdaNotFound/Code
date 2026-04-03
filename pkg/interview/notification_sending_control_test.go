package interview

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ShouldSendNotification(t *testing.T) {
	type call struct {
		service  string
		day      int
		expected bool
	}
	testCases := []struct {
		name  string
		k     int
		calls []call
	}{
		{
			name: "first call always allowed",
			k:    3,
			calls: []call{
				{"email", 1, true},
			},
		},
		{
			name: "same day calls do not increment streak",
			k:    2,
			calls: []call{
				{"email", 1, true},
				{"email", 1, true},  // same day, streak stays 1
				{"email", 2, true},  // consecutive, streak becomes 2
				{"email", 2, true},  // same day at max streak, allowed (2 <= k=2)
				{"email", 3, false}, // would push streak to 3 > k=2
			},
		},
		{
			name: "blocked on k+1 consecutive day",
			k:    2,
			calls: []call{
				{"email", 1, true},
				{"email", 2, true},
				{"email", 3, false}, // streak would be 3 > k=2, blocked, state unchanged
			},
		},
		{
			name: "gap resets streak and allows again",
			k:    2,
			calls: []call{
				{"email", 1, true},
				{"email", 2, true},
				{"email", 3, false}, // blocked, lastDay stays at 2
				{"email", 5, true},  // gap from lastDay=2 → reset, streak=1
				{"email", 6, true},  // consecutive, streak=2
				{"email", 7, false}, // blocked again
			},
		},
		{
			name: "k=1 allows once then blocks",
			k:    1,
			calls: []call{
				{"sms", 1, true},
				{"sms", 2, false}, // consecutive, streak would be 2 > k=1
				{"sms", 4, true},  // gap → reset
			},
		},
		{
			name: "multiple services are tracked independently",
			k:    2,
			calls: []call{
				{"email", 1, true},
				{"push", 1, true},
				{"email", 2, true},
				{"push", 2, true},
				{"email", 3, false}, // email blocked
				{"push", 3, false},  // push blocked independently
				{"sms", 5, true},    // sms first call
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rl := NewRateLimiter(tc.k)
			for _, c := range tc.calls {
				result := rl.ShouldSendNotification(c.service, c.day)
				assert.Equal(t, c.expected, result, "service=%q day=%d", c.service, c.day)
			}
		})
	}
}

func Test_ShouldSendNotification2(t *testing.T) {
	type call struct {
		service  string
		day      int
		expected bool
	}
	testCases := []struct {
		name     string
		k        int
		cooldown map[string]int // per-service cooldown (0 = no cooldown)
		calls    []call
	}{
		{
			name:     "first call always allowed",
			k:        3,
			cooldown: map[string]int{},
			calls: []call{
				{"email", 1, true},
			},
		},
		{
			name:     "blocked on k+1 consecutive day without cooldown",
			k:        2,
			cooldown: map[string]int{},
			calls: []call{
				{"email", 1, true},
				{"email", 2, true},
				{"email", 3, false}, // blocked, lastBlockedDay=3
				{"email", 4, true},  // cooldown=0 expired (4 > 3+0), gap → reset
			},
		},
		{
			name:     "cooldown keeps service blocked after streak exceeded",
			k:        2,
			cooldown: map[string]int{"email": 3},
			calls: []call{
				{"email", 1, true},
				{"email", 2, true},
				{"email", 3, false}, // blocked, lastBlockedDay=3, cooldown=3
				{"email", 4, false}, // 4 <= 3+3=6, still cooling down
				{"email", 6, false}, // 6 <= 6, still cooling down
				{"email", 7, true},  // 7 > 6, cooldown expired, gap from lastDay=2 → reset
			},
		},
		{
			name:     "same day calls do not increment streak",
			k:        2,
			cooldown: map[string]int{},
			calls: []call{
				{"push", 1, true},
				{"push", 1, true},
				{"push", 2, true},
				{"push", 3, false},
			},
		},
		{
			name:     "cooldown only applies to blocked service",
			k:        1,
			cooldown: map[string]int{"email": 5},
			calls: []call{
				{"email", 1, true},
				{"email", 2, false}, // blocked, cooldown=5
				{"push", 2, true},   // push unaffected
				{"push", 3, false},  // push blocked by its own streak
				{"email", 7, false}, // email: 7 <= 2+5=7, still cooling
				{"email", 8, true},  // 8 > 7, cooldown expired, gap → reset
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rl := NewRateLimiter(tc.k)
			rl.cooldown = tc.cooldown
			rl.lastBlockedDay = make(map[string]int)
			for _, c := range tc.calls {
				result := rl.ShouldSendNotification2(c.service, c.day)
				assert.Equal(t, c.expected, result, "service=%q day=%d", c.service, c.day)
			}
		})
	}
}
