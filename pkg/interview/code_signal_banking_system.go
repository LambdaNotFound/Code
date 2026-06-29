package interview

import (
	"sort"
	"strconv"
)

type Account struct {
	id      string
	balance int
	// level 2
	outgoing int
	// level 4
	history []BalanceSnapshot // ordered by timestamp
}

// level 3
type Payment struct {
	id        string
	accountId string
	amount    int
	executeAt int
	canceled  bool
	executed  bool
}

type BankingSystem struct {
	accounts map[string]*Account
	// level 3
	payments       []*Payment // ordered by creation
	paymentCounter int
}

// level 4
type BalanceSnapshot struct {
	timestamp int
	balance   int
}

func NewBankingSystem() *BankingSystem {
	return &BankingSystem{
		accounts: make(map[string]*Account),
	}
}

// level 3 Call this at the start of every operation
func (b *BankingSystem) processPending(timestamp int) {
	for _, p := range b.payments {
		if p.canceled || p.executed {
			continue
		}
		if p.executeAt > timestamp {
			continue
		}
		// Execute in creation order (slice is already ordered)
		acc, ok := b.accounts[p.accountId]
		if !ok || acc.balance < p.amount {
			p.executed = true // skip but mark done
			continue
		}
		acc.balance -= p.amount
		acc.outgoing += p.amount
		acc.recordSnapshot(p.executeAt) // level 4
		p.executed = true
	}
}

// level 4 helper — call after every balance change
func (a *Account) recordSnapshot(timestamp int) {
	a.history = append(a.history, BalanceSnapshot{timestamp: timestamp, balance: a.balance})
}

func (b *BankingSystem) CreateAccount(timestamp int, accountId string) bool {
	if _, exists := b.accounts[accountId]; exists {
		return false
	}
	acc := &Account{id: accountId, balance: 0}
	acc.recordSnapshot(timestamp) // level 4
	b.accounts[accountId] = acc
	return true
}

func (b *BankingSystem) Deposit(timestamp int, accountId string, amount int) *int {
	b.processPending(timestamp) // level 3

	acc, ok := b.accounts[accountId]
	if !ok {
		return nil
	}
	acc.balance += amount
	acc.recordSnapshot(timestamp) // level 4
	return &acc.balance
}

func (b *BankingSystem) Transfer(timestamp int, sourceId, targetId string, amount int) *int {
	b.processPending(timestamp) // level 3

	if sourceId == targetId {
		return nil
	}
	source, ok := b.accounts[sourceId]
	if !ok {
		return nil
	}
	target, ok := b.accounts[targetId]
	if !ok {
		return nil
	}
	if source.balance < amount {
		return nil
	}
	source.balance -= amount
	source.outgoing += amount        // level 2
	source.recordSnapshot(timestamp) // level 4

	target.balance += amount
	target.recordSnapshot(timestamp) // level 4
	return &source.balance
}

// level 2
func (b *BankingSystem) TopSpenders(timestamp int, n int) []string {
	b.processPending(timestamp) // level 3

	accounts := make([]*Account, 0, len(b.accounts))
	for _, acc := range b.accounts {
		accounts = append(accounts, acc)
	}

	sort.Slice(accounts, func(i, j int) bool {
		if accounts[i].outgoing != accounts[j].outgoing {
			return accounts[i].outgoing > accounts[j].outgoing // desc
		}
		return accounts[i].id < accounts[j].id // asc alphabetical tie-break
	})

	if n > len(accounts) {
		n = len(accounts)
	}
	result := make([]string, n)
	for i := 0; i < n; i++ {
		result[i] = accounts[i].id + "(" + strconv.Itoa(accounts[i].outgoing) + ")"
	}
	return result
}

// level 3
func (b *BankingSystem) SchedulePayment(timestamp int, accountId string, amount int, delay int) *string {
	b.processPending(timestamp)
	if _, ok := b.accounts[accountId]; !ok {
		return nil
	}
	b.paymentCounter++
	id := "payment" + strconv.Itoa(b.paymentCounter)
	p := &Payment{
		id:        id,
		accountId: accountId,
		amount:    amount,
		executeAt: timestamp + delay,
	}
	b.payments = append(b.payments, p)
	return &id
}

func (b *BankingSystem) CancelPayment(timestamp int, accountId string, paymentId string) bool {
	b.processPending(timestamp)
	for _, p := range b.payments {
		if p.id != paymentId {
			continue
		}
		if p.accountId != accountId {
			return false
		}
		if p.canceled || p.executed {
			return false
		}
		p.canceled = true
		return true
	}
	return false
}

// level 4
func (b *BankingSystem) MergeAccounts(timestamp int, accountId1, accountId2 string) bool {
	b.processPending(timestamp)
	if accountId1 == accountId2 {
		return false
	}
	acc1, ok1 := b.accounts[accountId1]
	acc2, ok2 := b.accounts[accountId2]
	if !ok1 || !ok2 {
		return false
	}

	// Merge balance
	acc1.balance += acc2.balance
	acc1.recordSnapshot(timestamp)

	// Merge outgoing
	acc1.outgoing += acc2.outgoing

	// Merge history — combine and sort by timestamp
	acc1.history = append(acc1.history, acc2.history...)
	sort.Slice(acc1.history, func(i, j int) bool {
		return acc1.history[i].timestamp < acc1.history[j].timestamp
	})

	// Move pending payments from acc2 to acc1
	for _, p := range b.payments {
		if !p.canceled && !p.executed && p.accountId == accountId2 {
			p.accountId = accountId1
		}
	}

	delete(b.accounts, accountId2)
	return true
}

func (b *BankingSystem) GetBalance(timestamp int, accountId string, timeAt int) *int {
	b.processPending(timestamp)
	acc, ok := b.accounts[accountId]
	if !ok {
		return nil
	}

	// Find latest snapshot at or before timeAt
	var result *int
	for _, snap := range acc.history {
		if snap.timestamp <= timeAt {
			val := snap.balance
			result = &val
		}
	}
	return result
}
