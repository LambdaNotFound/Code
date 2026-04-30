package interview

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ── testUserLogin (demo function coverage) ────────────────────────────────────

func Test_testUserLogin(t *testing.T) {
	// testUserLogin() prints to stdout; calling it exercises its branches.
	testUserLogin()
}

// ── calculateStats ────────────────────────────────────────────────────────────
//
// calculateStats counts logins within the past 15 days from time.Now(),
// so all "recent" test records use time.Now().AddDate(0,0,-N) to stay live.
// Records placed at -20 days are outside the 15-day window.

func daysAgo(n int) time.Time {
	return time.Now().AddDate(0, 0, -n)
}

func Test_calculateStats(t *testing.T) {
	testCases := []struct {
		name    string
		records []LoginRecord
		userID  string
		want    UserStats
	}{
		{
			name:    "empty_records_returns_empty_map",
			records: []LoginRecord{},
			userID:  "", // no user to check
		},
		{
			name: "single_login_within_15_days",
			records: []LoginRecord{
				{"alice", daysAgo(1)},
			},
			userID: "alice",
			want: UserStats{
				TotalLogins:   1,
				LongestStreak: 1,
				StreakPeriods: 0, // single day, no streak >= 2
			},
		},
		{
			name: "login_outside_15_days_not_counted_in_total",
			records: []LoginRecord{
				{"alice", daysAgo(20)}, // outside window
			},
			userID: "alice",
			want: UserStats{
				TotalLogins:   0,
				LongestStreak: 1,
				StreakPeriods: 0,
			},
		},
		{
			name: "three_consecutive_days_streak",
			records: []LoginRecord{
				{"alice", daysAgo(3)},
				{"alice", daysAgo(2)},
				{"alice", daysAgo(1)},
			},
			userID: "alice",
			want: UserStats{
				TotalLogins:   3,
				LongestStreak: 3,
				StreakPeriods: 1,
			},
		},
		{
			name: "two_separate_streaks_counted_separately",
			records: []LoginRecord{
				// Streak 1: days 10, 9, 8 (length 3)
				{"bob", daysAgo(10)},
				{"bob", daysAgo(9)},
				{"bob", daysAgo(8)},
				// gap at day 7
				// Streak 2: days 3, 2 (length 2)
				{"bob", daysAgo(3)},
				{"bob", daysAgo(2)},
			},
			userID: "bob",
			want: UserStats{
				TotalLogins:   5,
				LongestStreak: 3,
				StreakPeriods: 2,
			},
		},
		{
			name: "duplicate_logins_on_same_day_deduplicated",
			records: []LoginRecord{
				{"carol", daysAgo(2)},
				{"carol", daysAgo(2)}, // same day, dedup'd
				{"carol", daysAgo(1)},
			},
			userID: "carol",
			want: UserStats{
				TotalLogins:   2,
				LongestStreak: 2,
				StreakPeriods: 1,
			},
		},
		{
			name: "non_consecutive_days_no_streak_periods",
			records: []LoginRecord{
				{"dave", daysAgo(10)},
				{"dave", daysAgo(7)},
				{"dave", daysAgo(4)},
			},
			userID: "dave",
			want: UserStats{
				TotalLogins:   3,
				LongestStreak: 1,
				StreakPeriods: 0,
			},
		},
		{
			name: "only_older_logins_total_zero_streak_still_counted",
			records: []LoginRecord{
				{"eve", daysAgo(16)}, // outside 15-day window
				{"eve", daysAgo(17)}, // outside 15-day window
			},
			userID: "eve",
			want: UserStats{
				TotalLogins:   0,
				LongestStreak: 2,
				StreakPeriods: 1,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := calculateStats(tc.records)

			if tc.userID == "" {
				assert.Empty(t, result)
				return
			}

			stats, ok := result[tc.userID]
			assert.True(t, ok, "user %s not found in result", tc.userID)
			assert.Equal(t, tc.want.TotalLogins, stats.TotalLogins, "TotalLogins mismatch")
			assert.Equal(t, tc.want.LongestStreak, stats.LongestStreak, "LongestStreak mismatch")
			assert.Equal(t, tc.want.StreakPeriods, stats.StreakPeriods, "StreakPeriods mismatch")
		})
	}
}

func Test_calculateStats_MultipleUsers(t *testing.T) {
	records := []LoginRecord{
		{"alice", daysAgo(2)},
		{"alice", daysAgo(1)},
		{"bob", daysAgo(5)},
	}

	result := calculateStats(records)

	assert.Len(t, result, 2)

	alice := result["alice"]
	assert.Equal(t, 2, alice.TotalLogins)
	assert.Equal(t, 2, alice.LongestStreak)
	assert.Equal(t, 1, alice.StreakPeriods)

	bob := result["bob"]
	assert.Equal(t, 1, bob.TotalLogins)
	assert.Equal(t, 1, bob.LongestStreak)
	assert.Equal(t, 0, bob.StreakPeriods)
}
