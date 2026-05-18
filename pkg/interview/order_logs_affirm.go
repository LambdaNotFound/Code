package interview

/*
// Read line by line
scanner := bufio.NewScanner(os.Stdin)
for scanner.Scan() { line := scanner.Text(); ... }

// Read JSON
json.NewDecoder(os.Stdin).Decode(&result)

// Read CSV
records, _ := csv.NewReader(os.Stdin).ReadAll()

Part 1
Given two lists of log lines (representing transactions from different sources), identify established users. A user is established if they meet both conditions:

They have used at least 2 different loan types
They have taken loans on at least 2 different dates
Each log line follows this format:

YYYY-MM-DD,loanType,userId,amount
Examples:

"2025-09-03,Web,uuid1,100"
"2025-08-07,Mobile,uuid2,50"
"2025-09-02,Store,uuid4,200"
Part 2
Given an incoming transaction (userId, loanType, amount), compute its credit score. Only established users get scored; everyone else returns 0.

Scoring rules
Rule	Condition	Points
1	User has used the same loanType before	+50
2	Amount is within [minAmount, maxAmount] of all user transactions	+50
3	Amount is outside rule 2 bounds but within ±10% of those bounds	−10
4	Amount is outside rule 3 bounds but within ±40% of those bounds	−20
Rules 2, 3, and 4 are mutually exclusive tiers. Example with min=100, max=200:

[60 ......... 90 ... 100 ========= 200 ... 220 ......... 280]
     -20        -10         +50          -10         -20
*/
import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type UserStats struct {
	LoanTypes map[string]bool
	Dates     map[string]bool
	MinAmount float64
	MaxAmount float64
}

// CreditScorer aggregates user history and answers scoring queries.
type CreditScorer struct {
	stats       map[string]*UserStats
	established map[string]bool
}

func NewCreditScorer(logSources ...[]string) *CreditScorer {
	cs := &CreditScorer{
		stats:       make(map[string]*UserStats),
		established: make(map[string]bool),
	}

	// Phase 1: aggregate stats from all log sources
	for _, logs := range logSources {
		for _, line := range logs {
			cs.ingest(line)
		}
	}

	// Phase 2: determine established users
	for userID, s := range cs.stats {
		if len(s.LoanTypes) >= 2 && len(s.Dates) >= 2 {
			cs.established[userID] = true
		}
	}

	return cs
}

// ingest parses a single log line and updates the user's stats.
func (cs *CreditScorer) ingest(line string) {
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

	s, exists := cs.stats[userID]
	if !exists {
		// Sentinel init: MinAmount starts at +inf, MaxAmount at -inf
		// so the first min()/max() call correctly captures the real value.
		s = &UserStats{
			LoanTypes: make(map[string]bool),
			Dates:     make(map[string]bool),
			MinAmount: math.MaxFloat64,
			MaxAmount: -math.MaxFloat64,
		}
		cs.stats[userID] = s
	}

	s.LoanTypes[loanType] = true
	s.Dates[date] = true
	s.MinAmount = min(s.MinAmount, amount)
	s.MaxAmount = max(s.MaxAmount, amount)
}

// EstablishedUsers returns all established user IDs.
func (cs *CreditScorer) EstablishedUsers() []string {
	users := make([]string, 0, len(cs.established))
	for userID := range cs.established {
		users = append(users, userID)
	}
	return users
}

// Score returns the credit score for a new incoming transaction.
// Returns 0 for non-established users.
func (cs *CreditScorer) Score(userID, loanType string, amount float64) int {
	if !cs.established[userID] {
		return 0
	}

	s := cs.stats[userID]
	score := 0

	// Rule 1: same loan type used before
	if s.LoanTypes[loanType] {
		score += 50
	}

	// Rules 2-4: tiered amount scoring (mutually exclusive)
	lo, hi := s.MinAmount, s.MaxAmount
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

func main() {
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

	cs := NewCreditScorer(webLogs, storeLogs)

	fmt.Println("Established users:", cs.EstablishedUsers())
	// uuid1: 3 loanTypes, 3 dates → established
	// uuid2: 2 loanTypes, 2 dates → established
	// uuid4: 1 loanType, 1 date → NOT established

	// uuid1's history: min=100, max=200, types={Web, Mobile, Store}
	fmt.Println(cs.Score("uuid1", "Web", 150))    // +50 + 50 = 100
	fmt.Println(cs.Score("uuid1", "Web", 95))     // +50 - 10 = 40
	fmt.Println(cs.Score("uuid1", "Crypto", 280)) // 0 - 20 = -20
	fmt.Println(cs.Score("uuid1", "Web", 500))    // +50 + 0 = 50
	fmt.Println(cs.Score("uuid4", "Web", 100))    // 0 (not established)
}
