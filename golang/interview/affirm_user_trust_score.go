package interview

import (
	"bufio"
	"io"
	"math"
	"strconv"
	"strings"
	"sync"
)

// LoanLogEntry from a daily log file: <date, user id, order type, amount>.
type LoanLogEntry struct {
	LoanDate string
	UserID   string
	LoanType string
	Amount   float64
}

// UserStats for single user's loan stats
type UserStats struct {
	ByLoanType map[string][]*LoanLogEntry
	ByLoanDate map[string][]*LoanLogEntry
	MinAmount  float64
	MaxAmount  float64
}

type LoanManager struct {
	EstablishedUsers map[string]bool

	userStats map[string]*UserStats
	// follow up: stream input instead of log files
	mu sync.RWMutex
}

/**
 * Part 1
 * Given two lists of log lines (representing transactions from different sources),
 * identify established users. A user is established if they meet both conditions:
 *
 * They have used at least 2 different loan types
 * They have taken loans on at least 2 different dates
 * Each log line follows this format: YYYY-MM-DD,loanType,userId,amount
 *
 * Examples:
 *
 * "2025-09-03,Web,uuid1,100"
 * "2025-08-07,Mobile,uuid2,50"
 * "2025-09-02,Store,uuid4,200"
 *
 * Time: O(N) — one pass to parse all lines, O(1) per line (hash map ops); O(U) for the established-user scan, dominated by O(N).
 * Space: O(N) — every LoanLogEntry is stored twice (once in ByLoanType, once in ByLoanDate), plus O(U·(T+D)) for the map keys — all bounded by O(N).
 */
func NewLoanManager(logSources ...[]string) *LoanManager {
	lm := &LoanManager{
		EstablishedUsers: make(map[string]bool),

		userStats: make(map[string]*UserStats),
	}

	// Phase 1: aggregate stats from all log sources
	for _, logs := range logSources {
		for _, line := range logs {
			lm.Ingest(line)
		}
	}

	// Phase 2: calculate established users
	for userID, s := range lm.userStats {
		if len(s.ByLoanDate) >= 2 && len(s.ByLoanType) >= 2 {
			lm.EstablishedUsers[userID] = true
		}
	}

	return lm
}

// parses a single log line and updates the user's stats.
func (lm *LoanManager) Ingest(line string) {
	// Format: YYYY-MM-DD,loanType,userId,amount
	parts := strings.Split(line, ",")
	if len(parts) != 4 {
		return // malformed; skip
	}

	date, loanType, userID, amountStr := parts[0], parts[1], parts[2], parts[3]
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return
	}

	entry := &LoanLogEntry{
		LoanDate: date,
		UserID:   userID,
		LoanType: loanType,
		Amount:   amount,
	}

	stats, exists := lm.userStats[userID]
	if !exists {
		stats = &UserStats{
			ByLoanType: make(map[string][]*LoanLogEntry),
			ByLoanDate: make(map[string][]*LoanLogEntry),
			MinAmount:  math.MaxFloat64,
			MaxAmount:  -math.MaxFloat64,
		}
		lm.userStats[userID] = stats
	}

	stats.ByLoanType[loanType] = append(stats.ByLoanType[loanType], entry)
	stats.ByLoanDate[date] = append(stats.ByLoanDate[date], entry)
	stats.MinAmount = min(stats.MinAmount, amount)
	stats.MaxAmount = max(stats.MaxAmount, amount)
}

/**
 * Part 2
 * Given an incoming transaction (userId, loanType, amount), compute its credit score. Only established users get scored; everyone else returns 0.
 *
 * Scoring rules
 * Rule	Condition	Points
 * 1	User has used the same loanType before	+50
 * 2	Amount is within [minAmount, maxAmount] of all user transactions	+50
 * 3	Amount is outside rule 2 bounds but within ±10% of those bounds	−10
 * 4	Amount is outside rule 3 bounds but within ±40% of those bounds	−20
 * Rules 2, 3, and 4 are mutually exclusive tiers. Example with min=100, max=200:
 *
 * [60 ......... 90 ... 100 ========= 200 ... 220 ......... 280]
 * [    -20        -10         +50          -10         -20    ]
 *
 */

// Score returns the credit score for a new incoming transaction.
// Returns 0 for non-established users.
// Time: O(1) — two hash map lookups and a constant number of comparisons.
// Space: O(1) — no allocations.
func (lm *LoanManager) Score(userID, loanType string, amount float64) int {
	if !lm.EstablishedUsers[userID] {
		return 0
	}

	stat := lm.userStats[userID]
	score := 0

	// Rule 1: same loan type used before
	if _, found := stat.ByLoanType[loanType]; found {
		score += 50
	}

	// Rules 2-4: tiered amount scoring (mutually exclusive)
	lo, hi := stat.MinAmount, stat.MaxAmount
	switch {
	case amount >= lo && amount <= hi:
		score += 50
	case amount >= lo*0.9 && amount <= hi*1.1:
		score -= 10
	case amount >= lo*0.6 && amount <= hi*1.4:
		score -= 20
	}
	// else: outside all tiers → 0 contribution

	return score
}

// Parse reads log lines from r (e.g. os.Stdin) and returns a LoanManager
// seeded with all lines from that source.
//
// f, _ := os.Open("web.log")
// defer f.Close()
// lm := Parse(f)
func Parse(r io.Reader) *LoanManager {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return NewLoanManager(lines)
}
