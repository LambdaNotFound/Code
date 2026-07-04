package interview

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ── Set / Get ─────────────────────────────────────────────────────────────────

func Test_InMemoryDB_SetGet(t *testing.T) {
	testCases := []struct {
		name    string
		setup   func(*InMemoryDB)
		key     string
		field   string
		wantNil bool
		wantVal string
	}{
		{
			name:    "get_nonexistent_key_returns_nil",
			setup:   func(_ *InMemoryDB) {},
			key:     "k1", field: "f1",
			wantNil: true,
		},
		{
			name: "get_nonexistent_field_returns_nil",
			setup: func(db *InMemoryDB) {
				db.Set("k1", "other", "v")
			},
			key: "k1", field: "missing",
			wantNil: true,
		},
		{
			name: "set_then_get_returns_value",
			setup: func(db *InMemoryDB) {
				db.Set("k1", "f1", "hello")
			},
			key: "k1", field: "f1",
			wantVal: "hello",
		},
		{
			name: "overwrite_existing_field",
			setup: func(db *InMemoryDB) {
				db.Set("k1", "f1", "old")
				db.Set("k1", "f1", "new")
			},
			key: "k1", field: "f1",
			wantVal: "new",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := NewInMemoryDB()
			tc.setup(db)
			result := db.Get(tc.key, tc.field)
			if tc.wantNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, tc.wantVal, *result)
			}
		})
	}
}

// ── SetAt / GetAt ─────────────────────────────────────────────────────────────

func Test_InMemoryDB_SetAtGetAt(t *testing.T) {
	t.Run("get_at_zero_treats_no_expiry", func(t *testing.T) {
		db := NewInMemoryDB()
		db.SetAt("k", "f", "v", 5)
		// timestamp=0 in GetAt means no expiry check (expireAt==0 means no TTL)
		result := db.GetAt("k", "f", 0)
		assert.NotNil(t, result)
		assert.Equal(t, "v", *result)
	})
}

// ── SetAtWithTTL / expiry ─────────────────────────────────────────────────────

func Test_InMemoryDB_TTL(t *testing.T) {
	testCases := []struct {
		name      string
		setTime   int
		ttl       int
		getTime   int
		wantNil   bool
		wantVal   string
	}{
		{
			name:    "value_accessible_before_expiry",
			setTime: 1, ttl: 10, getTime: 5,
			wantVal: "val",
		},
		{
			name:    "value_expired_at_expiry_time",
			setTime: 1, ttl: 10, getTime: 11, // expireAt = 11, getTime >= 11
			wantNil: true,
		},
		{
			name:    "value_accessible_at_one_before_expiry",
			setTime: 1, ttl: 10, getTime: 10, // expireAt = 11, getTime=10 < 11
			wantVal: "val",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := NewInMemoryDB()
			db.SetAtWithTTL("k", "f", "val", tc.setTime, tc.ttl)
			result := db.GetAt("k", "f", tc.getTime)
			if tc.wantNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, tc.wantVal, *result)
			}
		})
	}
}

// ── Delete / DeleteAt ─────────────────────────────────────────────────────────

func Test_InMemoryDB_Delete(t *testing.T) {
	testCases := []struct {
		name    string
		setup   func(*InMemoryDB)
		key     string
		field   string
		wantOk  bool
	}{
		{
			name:  "delete_existing_field_returns_true",
			setup: func(db *InMemoryDB) { db.Set("k", "f", "v") },
			key:   "k", field: "f", wantOk: true,
		},
		{
			name:  "delete_nonexistent_key_returns_false",
			setup: func(_ *InMemoryDB) {},
			key:   "k", field: "f", wantOk: false,
		},
		{
			name:  "delete_nonexistent_field_returns_false",
			setup: func(db *InMemoryDB) { db.Set("k", "other", "v") },
			key:   "k", field: "f", wantOk: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := NewInMemoryDB()
			tc.setup(db)
			ok := db.Delete(tc.key, tc.field)
			assert.Equal(t, tc.wantOk, ok)
			if tc.wantOk {
				// Field gone after deletion
				assert.Nil(t, db.Get(tc.key, tc.field))
			}
		})
	}
}

func Test_InMemoryDB_DeleteAt(t *testing.T) {
	t.Run("delete_expired_entry_returns_false", func(t *testing.T) {
		db := NewInMemoryDB()
		db.SetAtWithTTL("k", "f", "v", 1, 5) // expires at t=6
		ok := db.DeleteAt("k", "f", 10)       // t=10 > 6, already expired
		assert.False(t, ok)
	})

	t.Run("delete_live_entry_succeeds", func(t *testing.T) {
		db := NewInMemoryDB()
		db.SetAtWithTTL("k", "f", "v", 1, 10) // expires at t=11
		ok := db.DeleteAt("k", "f", 5)         // t=5 < 11, still live
		assert.True(t, ok)
	})
}

// ── Scan / ScanAt ─────────────────────────────────────────────────────────────

func Test_InMemoryDB_Scan(t *testing.T) {
	testCases := []struct {
		name  string
		setup func(*InMemoryDB)
		key   string
		want  []string
	}{
		{
			name:  "scan_nonexistent_key_returns_empty",
			setup: func(_ *InMemoryDB) {},
			key:   "k", want: []string{},
		},
		{
			name: "scan_returns_all_fields_sorted",
			setup: func(db *InMemoryDB) {
				db.Set("k", "b", "2")
				db.Set("k", "a", "1")
				db.Set("k", "c", "3")
			},
			key:  "k",
			want: []string{"a(1)", "b(2)", "c(3)"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := NewInMemoryDB()
			tc.setup(db)
			got := db.Scan(tc.key)
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_InMemoryDB_ScanAt(t *testing.T) {
	t.Run("scan_at_skips_expired_entries", func(t *testing.T) {
		db := NewInMemoryDB()
		db.SetAtWithTTL("k", "live", "yes", 1, 20) // expires at t=21
		db.SetAtWithTTL("k", "dead", "no", 1, 5)   // expires at t=6
		got := db.ScanAt("k", 10)
		assert.Equal(t, []string{"live(yes)"}, got)
	})
}

// ── ScanByPrefix / ScanByPrefixAt ────────────────────────────────────────────

func Test_InMemoryDB_ScanByPrefix(t *testing.T) {
	testCases := []struct {
		name   string
		setup  func(*InMemoryDB)
		key    string
		prefix string
		want   []string
	}{
		{
			name: "prefix_filters_matching_fields",
			setup: func(db *InMemoryDB) {
				db.Set("k", "alpha", "1")
				db.Set("k", "almond", "2")
				db.Set("k", "beta", "3")
			},
			key: "k", prefix: "al",
			want: []string{"almond(2)", "alpha(1)"}, // "almond" < "alpha" alphabetically ('m' < 'p')
		},
		{
			name: "empty_prefix_returns_all_sorted",
			setup: func(db *InMemoryDB) {
				db.Set("k", "z", "last")
				db.Set("k", "a", "first")
			},
			key: "k", prefix: "",
			want: []string{"a(first)", "z(last)"},
		},
		{
			name: "no_matching_prefix_returns_nil",
			setup: func(db *InMemoryDB) {
				db.Set("k", "foo", "bar")
			},
			key: "k", prefix: "xyz",
			want: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := NewInMemoryDB()
			tc.setup(db)
			got := db.ScanByPrefix(tc.key, tc.prefix)
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_InMemoryDB_ScanByPrefixAt(t *testing.T) {
	t.Run("prefix_scan_at_skips_expired", func(t *testing.T) {
		db := NewInMemoryDB()
		db.SetAtWithTTL("k", "alive_x", "yes", 1, 20)
		db.SetAtWithTTL("k", "dead_x", "no", 1, 3) // expires at t=4
		got := db.ScanByPrefixAt("k", "", 5)
		assert.Equal(t, []string{"alive_x(yes)"}, got)
	})
}

// ── Backup / Restore ──────────────────────────────────────────────────────────

func Test_InMemoryDB_Backup(t *testing.T) {
	testCases := []struct {
		name      string
		setup     func(*InMemoryDB)
		timestamp int
		wantCount int
	}{
		{
			name:      "backup_empty_db_returns_zero",
			setup:     func(_ *InMemoryDB) {},
			timestamp: 1, wantCount: 0,
		},
		{
			name: "backup_counts_live_keys",
			setup: func(db *InMemoryDB) {
				db.Set("k1", "f1", "v1")
				db.Set("k2", "f2", "v2")
			},
			timestamp: 1, wantCount: 2,
		},
		{
			name: "backup_excludes_expired_keys",
			setup: func(db *InMemoryDB) {
				db.SetAtWithTTL("live", "f", "v", 1, 10) // expires at t=11
				db.SetAtWithTTL("dead", "f", "v", 1, 2)  // expires at t=3
			},
			timestamp: 5, wantCount: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := NewInMemoryDB()
			tc.setup(db)
			count := db.Backup(tc.timestamp)
			assert.Equal(t, tc.wantCount, count)
		})
	}
}

func Test_InMemoryDB_Restore(t *testing.T) {
	t.Run("restore_brings_back_deleted_key", func(t *testing.T) {
		db := NewInMemoryDB()
		db.Set("k", "f", "original")
		db.Backup(10)    // snapshot at t=10 with k.f=original
		db.Delete("k", "f") // delete after backup
		db.Restore(20, 10)  // restore to t=10 state at t=20

		result := db.Get("k", "f")
		assert.NotNil(t, result)
		assert.Equal(t, "original", *result)
	})

	t.Run("restore_recalculates_ttl", func(t *testing.T) {
		db := NewInMemoryDB()
		db.SetAtWithTTL("k", "f", "val", 5, 10) // expireAt=15, TTL remaining at t=10: 5
		db.Backup(10)                             // backup at t=10 (remaining TTL = 5)

		// Restore at t=30: new expireAt = 30 + 5 = 35
		db.Restore(30, 10)

		// t=34 < 35: still live
		result := db.GetAt("k", "f", 34)
		assert.NotNil(t, result)
		assert.Equal(t, "val", *result)

		// t=35 >= 35: expired
		result = db.GetAt("k", "f", 35)
		assert.Nil(t, result)
	})

	t.Run("restore_overwrites_current_state", func(t *testing.T) {
		db := NewInMemoryDB()
		db.Set("k", "f", "before_backup")
		db.Backup(1)
		db.Set("k", "f", "after_backup")
		db.Set("new_k", "f", "extra")

		// Restore to t=1 state
		db.Restore(5, 1)

		result := db.Get("k", "f")
		assert.NotNil(t, result)
		assert.Equal(t, "before_backup", *result)

		// new_k was added after backup — should be gone
		assert.Nil(t, db.Get("new_k", "f"))
	})
}
