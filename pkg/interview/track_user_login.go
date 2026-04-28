package interview

import (
	"fmt"
	"sort"
	"time"
)

type LoginRecord struct {
	UserID    string
	LoginDate time.Time
}

type UserStats struct {
	TotalLogins   int
	LongestStreak int
	StreakPeriods int // periods with >= 2 consecutive days
}

func calculateStats(records []LoginRecord) map[string]UserStats {
	// Group login dates by user, deduplicate by date
	userDates := make(map[string]map[string]bool)
	for _, r := range records {
		dateKey := r.LoginDate.Format("2006-01-02")
		if userDates[r.UserID] == nil {
			userDates[r.UserID] = make(map[string]bool)
		}
		userDates[r.UserID][dateKey] = true
	}

	now := time.Now()
	cutoff := now.AddDate(0, 0, -15)
	result := make(map[string]UserStats)

	for userID, dateSet := range userDates {
		// Collect all dates as time.Time, sorted
		var dates []time.Time
		for dateStr := range dateSet {
			d, _ := time.Parse("2006-01-02", dateStr)
			dates = append(dates, d)
		}
		sort.Slice(dates, func(i, j int) bool {
			return dates[i].Before(dates[j])
		})

		// 1. Total logins in past 15 days
		totalLogins := 0
		for _, d := range dates {
			if d.After(cutoff) {
				totalLogins++
			}
		}

		// 2. Longest streak & 3. Count of streaks >= 2 days
		longestStreak := 0
		streakPeriods := 0
		currentStreak := 1

		for i := 1; i < len(dates); i++ {
			diff := int(dates[i].Sub(dates[i-1]).Hours() / 24)
			if diff == 1 {
				currentStreak++
			} else {
				if currentStreak > longestStreak {
					longestStreak = currentStreak
				}
				if currentStreak >= 2 {
					streakPeriods++
				}
				currentStreak = 1
			}
		}
		// flush the last streak
		if currentStreak > longestStreak {
			longestStreak = currentStreak
		}
		if currentStreak >= 2 {
			streakPeriods++
		}

		result[userID] = UserStats{
			TotalLogins:   totalLogins,
			LongestStreak: longestStreak,
			StreakPeriods: streakPeriods,
		}
	}

	return result
}

func testUserLogin() {
	now := time.Now()
	records := []LoginRecord{
		{"alice", now.AddDate(0, 0, -1)},
		{"alice", now.AddDate(0, 0, -2)},
		{"alice", now.AddDate(0, 0, -3)},
		{"alice", now.AddDate(0, 0, -10)},
		{"alice", now.AddDate(0, 0, -20)}, // outside 15 days
		{"bob", now.AddDate(0, 0, -1)},
		{"bob", now.AddDate(0, 0, -5)},
		{"bob", now.AddDate(0, 0, -6)},
		{"bob", now.AddDate(0, 0, -7)},
	}

	stats := calculateStats(records)
	for userID, s := range stats {
		fmt.Printf("User: %s | Last15Days: %d | LongestStreak: %d | StreakPeriods(>=2): %d\n",
			userID, s.TotalLogins, s.LongestStreak, s.StreakPeriods)
	}
}
