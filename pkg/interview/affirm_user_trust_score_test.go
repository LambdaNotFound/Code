package interview

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupLoanManager() *LoanManager {
	webLogs := []string{
		"2025-09-03,Web,uuid1,100",
		"2025-09-05,Mobile,uuid1,200",
		"2025-08-07,Mobile,uuid2,50",
	}
	storeLogs := []string{
		"2025-09-02,Store,uuid1,150",
		"2025-09-04,Store,uuid4,200",
		"2025-09-06,Web,uuid2,80",
	}
	return NewLoanManager(webLogs, storeLogs)
}

func Test_LoanManager_EstablishedUsers(t *testing.T) {
	lm := setupLoanManager()

	// uuid1: 3 loan types (Web, Mobile, Store), 3 dates → established
	assert.True(t, lm.EstablishedUsers["uuid1"])
	// uuid2: 2 loan types (Mobile, Web), 2 dates → established
	assert.True(t, lm.EstablishedUsers["uuid2"])
	// uuid4: 1 loan type (Store), 1 date → NOT established
	assert.False(t, lm.EstablishedUsers["uuid4"])
}

func Test_LoanManager_Score(t *testing.T) {
	lm := setupLoanManager()
	// uuid1: min=100, max=200, types={Web, Mobile, Store}

	tests := []struct {
		userID   string
		loanType string
		amount   float64
		want     int
		desc     string
	}{
		{"uuid1", "Web", 150, 100, "known type + within [min,max]: +50+50"},
		{"uuid1", "Web", 95, 40, "known type + within ±10%: +50-10"},
		{"uuid1", "Crypto", 280, -20, "unknown type + within ±40%: 0-20"},
		{"uuid1", "Web", 500, 50, "known type + outside all tiers: +50+0"},
		{"uuid4", "Web", 100, 0, "non-established user always scores 0"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert.Equal(t, tt.want, lm.Score(tt.userID, tt.loanType, tt.amount))
		})
	}
}
