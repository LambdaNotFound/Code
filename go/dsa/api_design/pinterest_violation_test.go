package apidesign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var violationRecords = []ViolationRecord{
	{PinID: 1, Policy: "spam", Date: "2024-01-01"},
	{PinID: 2, Policy: "spam", Date: "2024-01-01"},
	{PinID: 1, Policy: "spam", Date: "2024-01-02"}, // pin 1 again — should not double-count
	{PinID: 3, Policy: "nudity", Date: "2024-01-02"},
	{PinID: 4, Policy: "nudity", Date: "2024-01-03"},
	{PinID: 2, Policy: "nudity", Date: "2024-01-03"},
	{PinID: 5, Policy: "spam", Date: "2024-01-05"},
}

func Test_CountByPolicy(t *testing.T) {
	vt := NewViolationTracker(violationRecords)

	tests := []struct {
		name     string
		policy   string
		expected int
	}{
		{name: "spam", policy: "spam", expected: 3},   // pins 1, 2, 5
		{name: "nudity", policy: "nudity", expected: 3}, // pins 2, 3, 4
		{name: "unknown_policy", policy: "violence", expected: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, vt.CountByPolicy(tt.policy))
		})
	}
}

func Test_CountUniqueInRange(t *testing.T) {
	vt := NewViolationTracker(violationRecords)

	tests := []struct {
		name      string
		startDate string
		endDate   string
		expected  int
	}{
		{name: "single_day", startDate: "2024-01-01", endDate: "2024-01-01", expected: 2},        // pins 1, 2
		{name: "two_days", startDate: "2024-01-01", endDate: "2024-01-02", expected: 3},           // pins 1, 2, 3
		{name: "full_range", startDate: "2024-01-01", endDate: "2024-01-05", expected: 5},         // pins 1–5
		{name: "gap_skips_missing_date", startDate: "2024-01-04", endDate: "2024-01-04", expected: 0},
		{name: "before_any_record", startDate: "2023-12-01", endDate: "2023-12-31", expected: 0},
		{name: "end_only", startDate: "2024-01-05", endDate: "2024-01-05", expected: 1},           // pin 5
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, vt.CountUniqueInRange(tt.startDate, tt.endDate))
		})
	}
}

func Test_CountPerPolicyInRange(t *testing.T) {
	vt := NewViolationTracker(violationRecords)

	tests := []struct {
		name      string
		startDate string
		endDate   string
		expected  map[string]int
	}{
		{
			name:      "first_two_days",
			startDate: "2024-01-01",
			endDate:   "2024-01-02",
			expected:  map[string]int{"spam": 2, "nudity": 1}, // spam: pins 1,2; nudity: pin 3
		},
		{
			name:      "last_three_days",
			startDate: "2024-01-03",
			endDate:   "2024-01-05",
			expected:  map[string]int{"spam": 1, "nudity": 2}, // spam: pin 5; nudity: pins 2,4
		},
		{
			name:      "full_range",
			startDate: "2024-01-01",
			endDate:   "2024-01-05",
			expected:  map[string]int{"spam": 3, "nudity": 3},
		},
		{
			name:      "no_violations_in_range",
			startDate: "2024-01-04",
			endDate:   "2024-01-04",
			expected:  map[string]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, vt.CountPerPolicyInRange(tt.startDate, tt.endDate))
		})
	}
}
