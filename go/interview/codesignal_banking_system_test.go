package interview

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ── CreateAccount ─────────────────────────────────────────────────────────────

func Test_BankingSystem_CreateAccount(t *testing.T) {
	testCases := []struct {
		name      string
		accountId string
		want      bool
	}{
		{name: "new_account_succeeds", accountId: "acc1", want: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := NewBankingSystem()
			got := b.CreateAccount(1, tc.accountId)
			assert.Equal(t, tc.want, got)
		})
	}

	t.Run("duplicate_account_fails", func(t *testing.T) {
		b := NewBankingSystem()
		assert.True(t, b.CreateAccount(1, "acc1"))
		assert.False(t, b.CreateAccount(2, "acc1"))
	})
}

// ── Deposit ───────────────────────────────────────────────────────────────────

func Test_BankingSystem_Deposit(t *testing.T) {
	testCases := []struct {
		name      string
		accountId string
		amount    int
		setup     func(*BankingSystem)
		wantNil   bool
		wantBal   int
	}{
		{
			name:      "deposit_to_existing_account",
			accountId: "acc1",
			amount:    500,
			setup:     func(b *BankingSystem) { b.CreateAccount(1, "acc1") },
			wantBal:   500,
		},
		{
			name:      "deposit_to_nonexistent_account_returns_nil",
			accountId: "ghost",
			amount:    100,
			setup:     func(_ *BankingSystem) {},
			wantNil:   true,
		},
		{
			name:      "multiple_deposits_accumulate",
			accountId: "acc1",
			amount:    200,
			setup: func(b *BankingSystem) {
				b.CreateAccount(1, "acc1")
				b.Deposit(2, "acc1", 300)
			},
			wantBal: 500,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := NewBankingSystem()
			tc.setup(b)
			result := b.Deposit(10, tc.accountId, tc.amount)
			if tc.wantNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, tc.wantBal, *result)
			}
		})
	}
}

// ── Transfer ─────────────────────────────────────────────────────────────────

func Test_BankingSystem_Transfer(t *testing.T) {
	testCases := []struct {
		name    string
		setup   func(*BankingSystem)
		srcId   string
		dstId   string
		amount  int
		wantNil bool
		wantBal int // source balance after transfer
	}{
		{
			name: "basic_transfer_succeeds",
			setup: func(b *BankingSystem) {
				b.CreateAccount(1, "src")
				b.CreateAccount(1, "dst")
				b.Deposit(2, "src", 1000)
			},
			srcId: "src", dstId: "dst", amount: 400,
			wantBal: 600,
		},
		{
			name: "self_transfer_returns_nil",
			setup: func(b *BankingSystem) {
				b.CreateAccount(1, "acc")
				b.Deposit(2, "acc", 500)
			},
			srcId: "acc", dstId: "acc", amount: 100,
			wantNil: true,
		},
		{
			name: "insufficient_funds_returns_nil",
			setup: func(b *BankingSystem) {
				b.CreateAccount(1, "src")
				b.CreateAccount(1, "dst")
				b.Deposit(2, "src", 50)
			},
			srcId: "src", dstId: "dst", amount: 100,
			wantNil: true,
		},
		{
			name:    "nonexistent_source_returns_nil",
			setup:   func(b *BankingSystem) { b.CreateAccount(1, "dst") },
			srcId:   "ghost", dstId: "dst", amount: 10,
			wantNil: true,
		},
		{
			name:    "nonexistent_target_returns_nil",
			setup:   func(b *BankingSystem) { b.CreateAccount(1, "src"); b.Deposit(2, "src", 100) },
			srcId:   "src", dstId: "ghost", amount: 10,
			wantNil: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := NewBankingSystem()
			tc.setup(b)
			result := b.Transfer(10, tc.srcId, tc.dstId, tc.amount)
			if tc.wantNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, tc.wantBal, *result)
			}
		})
	}
}

// ── TopSpenders ───────────────────────────────────────────────────────────────

func Test_BankingSystem_TopSpenders(t *testing.T) {
	testCases := []struct {
		name  string
		setup func(*BankingSystem)
		n     int
		want  []string
	}{
		{
			name: "top_spenders_sorted_by_outgoing_desc",
			setup: func(b *BankingSystem) {
				b.CreateAccount(1, "alice")
				b.CreateAccount(1, "bob")
				b.CreateAccount(1, "carol")
				b.Deposit(2, "alice", 1000)
				b.Deposit(2, "bob", 1000)
				b.Deposit(2, "carol", 1000)
				b.Transfer(3, "alice", "carol", 300) // alice outgoing=300
				b.Transfer(3, "bob", "carol", 500)   // bob outgoing=500
			},
			n:    2,
			want: []string{"bob(500)", "alice(300)"},
		},
		{
			name: "alphabetical_tiebreak_on_equal_outgoing",
			setup: func(b *BankingSystem) {
				b.CreateAccount(1, "alice")
				b.CreateAccount(1, "bob")
				b.CreateAccount(1, "carol")
				b.Deposit(2, "alice", 500)
				b.Deposit(2, "bob", 500)
				b.Transfer(3, "alice", "carol", 200)
				b.Transfer(3, "bob", "carol", 200)
			},
			n:    2,
			want: []string{"alice(200)", "bob(200)"},
		},
		{
			name: "n_larger_than_account_count",
			setup: func(b *BankingSystem) {
				b.CreateAccount(1, "acc1")
				b.Deposit(2, "acc1", 100)
				b.CreateAccount(1, "acc2")
				b.Deposit(2, "acc2", 200)
				b.Transfer(3, "acc2", "acc1", 50)
			},
			n:    10,
			want: []string{"acc2(50)", "acc1(0)"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := NewBankingSystem()
			tc.setup(b)
			got := b.TopSpenders(20, tc.n)
			assert.Equal(t, tc.want, got)
		})
	}
}

// ── SchedulePayment / CancelPayment ──────────────────────────────────────────

func Test_BankingSystem_SchedulePayment(t *testing.T) {
	t.Run("scheduled_payment_executes_on_next_operation", func(t *testing.T) {
		b := NewBankingSystem()
		b.CreateAccount(1, "acc")
		b.Deposit(2, "acc", 1000)

		// Schedule payment of 200 at t=5 (delay=3 from t=2, executes at t=5)
		payId := b.SchedulePayment(2, "acc", 200, 3)
		assert.NotNil(t, payId)
		assert.Equal(t, "payment1", *payId)

		// Trigger at t=6 (past execution time) — balance should be 800
		bal := b.Deposit(6, "acc", 0)
		assert.NotNil(t, bal)
		assert.Equal(t, 800, *bal)
	})

	t.Run("schedule_to_nonexistent_account_returns_nil", func(t *testing.T) {
		b := NewBankingSystem()
		result := b.SchedulePayment(1, "ghost", 100, 5)
		assert.Nil(t, result)
	})

	t.Run("payment_skipped_if_insufficient_funds", func(t *testing.T) {
		b := NewBankingSystem()
		b.CreateAccount(1, "acc")
		b.Deposit(2, "acc", 50) // only 50
		b.SchedulePayment(2, "acc", 200, 3)

		// At t=10 insufficient funds — payment marked done but not deducted
		bal := b.Deposit(10, "acc", 0)
		assert.Equal(t, 50, *bal)
	})
}

func Test_BankingSystem_CancelPayment(t *testing.T) {
	t.Run("cancel_pending_payment_succeeds", func(t *testing.T) {
		b := NewBankingSystem()
		b.CreateAccount(1, "acc")
		b.Deposit(2, "acc", 1000)
		payId := b.SchedulePayment(2, "acc", 200, 10)

		ok := b.CancelPayment(3, "acc", *payId)
		assert.True(t, ok)

		// Balance unchanged after cancellation window
		bal := b.Deposit(20, "acc", 0)
		assert.Equal(t, 1000, *bal)
	})

	t.Run("cancel_executed_payment_fails", func(t *testing.T) {
		b := NewBankingSystem()
		b.CreateAccount(1, "acc")
		b.Deposit(2, "acc", 1000)
		payId := b.SchedulePayment(2, "acc", 200, 3) // executes at t=5

		// Advance past execution time
		b.Deposit(10, "acc", 0)

		ok := b.CancelPayment(11, "acc", *payId)
		assert.False(t, ok)
	})

	t.Run("cancel_second_payment_skips_first", func(t *testing.T) {
		b := NewBankingSystem()
		b.CreateAccount(1, "acc")
		b.Deposit(2, "acc", 1000)
		b.SchedulePayment(2, "acc", 100, 10) // payment1
		pay2 := b.SchedulePayment(2, "acc", 200, 10) // payment2

		// Cancel the second — loop must skip payment1 (p.id != paymentId → continue)
		ok := b.CancelPayment(3, "acc", *pay2)
		assert.True(t, ok)
	})

	t.Run("cancel_already_canceled_payment_fails", func(t *testing.T) {
		b := NewBankingSystem()
		b.CreateAccount(1, "acc")
		b.Deposit(2, "acc", 1000)
		payId := b.SchedulePayment(2, "acc", 200, 10)

		// Cancel once — succeeds
		assert.True(t, b.CancelPayment(3, "acc", *payId))
		// Cancel again — already canceled
		assert.False(t, b.CancelPayment(4, "acc", *payId))
	})

	t.Run("cancel_nonexistent_payment_fails", func(t *testing.T) {
		b := NewBankingSystem()
		b.CreateAccount(1, "acc")
		ok := b.CancelPayment(1, "acc", "payment99")
		assert.False(t, ok)
	})

	t.Run("cancel_wrong_account_fails", func(t *testing.T) {
		b := NewBankingSystem()
		b.CreateAccount(1, "acc1")
		b.CreateAccount(1, "acc2")
		b.Deposit(2, "acc1", 500)
		payId := b.SchedulePayment(2, "acc1", 100, 10)

		ok := b.CancelPayment(3, "acc2", *payId)
		assert.False(t, ok)
	})
}

// ── MergeAccounts ─────────────────────────────────────────────────────────────

func Test_BankingSystem_MergeAccounts(t *testing.T) {
	testCases := []struct {
		name    string
		setup   func(*BankingSystem)
		acc1    string
		acc2    string
		wantOk  bool
		wantBal int // acc1 balance after merge
	}{
		{
			name: "merge_two_accounts_combines_balance",
			setup: func(b *BankingSystem) {
				b.CreateAccount(1, "acc1")
				b.CreateAccount(1, "acc2")
				b.Deposit(2, "acc1", 300)
				b.Deposit(2, "acc2", 700)
			},
			acc1: "acc1", acc2: "acc2",
			wantOk: true, wantBal: 1000,
		},
		{
			name: "merge_same_account_fails",
			setup: func(b *BankingSystem) {
				b.CreateAccount(1, "acc1")
				b.Deposit(2, "acc1", 500)
			},
			acc1: "acc1", acc2: "acc1",
			wantOk: false,
		},
		{
			name: "merge_nonexistent_account_fails",
			setup: func(b *BankingSystem) {
				b.CreateAccount(1, "acc1")
			},
			acc1: "acc1", acc2: "ghost",
			wantOk: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := NewBankingSystem()
			tc.setup(b)
			ok := b.MergeAccounts(10, tc.acc1, tc.acc2)
			assert.Equal(t, tc.wantOk, ok)
			if tc.wantOk {
				bal := b.GetBalance(11, tc.acc1, 10)
				assert.NotNil(t, bal)
				assert.Equal(t, tc.wantBal, *bal)
			}
		})
	}
}

func Test_BankingSystem_MergeAccounts_PendingPayments(t *testing.T) {
	t.Run("pending_payment_on_merged_account_executes_from_acc1", func(t *testing.T) {
		b := NewBankingSystem()
		b.CreateAccount(1, "acc1")
		b.CreateAccount(1, "acc2")
		b.Deposit(2, "acc1", 1000)
		b.Deposit(2, "acc2", 500)

		// Schedule payment on acc2 that hasn't fired yet (delay=20 → fires at t=22)
		b.SchedulePayment(2, "acc2", 200, 20)

		// Merge acc2 into acc1 at t=5 (before payment fires)
		ok := b.MergeAccounts(5, "acc1", "acc2")
		assert.True(t, ok)

		// At t=25 the payment should execute against acc1 (deducting 200)
		// acc1 balance after merge = 1000 + 500 = 1500, then -200 = 1300
		bal := b.Deposit(25, "acc1", 0)
		assert.NotNil(t, bal)
		assert.Equal(t, 1300, *bal)
	})
}

// ── GetBalance ────────────────────────────────────────────────────────────────

func Test_BankingSystem_GetBalance(t *testing.T) {
	testCases := []struct {
		name    string
		setup   func(*BankingSystem)
		acc     string
		timeAt  int
		wantNil bool
		wantBal int
	}{
		{
			name: "balance_at_creation_time",
			setup: func(b *BankingSystem) {
				b.CreateAccount(5, "acc")
			},
			acc: "acc", timeAt: 5,
			wantBal: 0,
		},
		{
			name: "balance_before_creation_returns_nil",
			setup: func(b *BankingSystem) {
				b.CreateAccount(5, "acc")
			},
			acc: "acc", timeAt: 4,
			wantNil: true,
		},
		{
			name: "balance_after_deposit",
			setup: func(b *BankingSystem) {
				b.CreateAccount(1, "acc")
				b.Deposit(5, "acc", 400)
			},
			acc: "acc", timeAt: 5,
			wantBal: 400,
		},
		{
			name: "historical_balance_before_deposit",
			setup: func(b *BankingSystem) {
				b.CreateAccount(1, "acc")
				b.Deposit(5, "acc", 400)
				b.Deposit(10, "acc", 100)
			},
			acc: "acc", timeAt: 7, // between the two deposits
			wantBal: 400,
		},
		{
			name: "nonexistent_account_returns_nil",
			setup: func(_ *BankingSystem) {},
			acc:  "ghost", timeAt: 5,
			wantNil: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := NewBankingSystem()
			tc.setup(b)
			result := b.GetBalance(20, tc.acc, tc.timeAt)
			if tc.wantNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, tc.wantBal, *result)
			}
		})
	}
}
