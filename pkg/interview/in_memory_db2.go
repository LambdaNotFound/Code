package interview

import (
	"sort"
	"strings"
)

type FieldEntry struct {
	value    string
	expireAt int // 0 means no expiry
}

// level 4
type BackupEntry struct {
	value      string
	expireAt   int // 0 = no expiry, relative to original timeline
	backupTime int // when this backup was taken
}

type Backup struct {
	timestamp int
	records   map[string]map[string]*BackupEntry
}

type InMemoryDB struct {
	records map[string]map[string]*FieldEntry
	// level 4
	backups []Backup // ordered by timestamp ascending
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		records: make(map[string]map[string]*FieldEntry),
	}
}

// helper — returns nil if entry is expired
func (db *InMemoryDB) getEntry(key, field string, timestamp int) *FieldEntry {
	rec, ok := db.records[key]
	if !ok {
		return nil
	}
	e, ok := rec[field]
	if !ok {
		return nil
	}
	if e.expireAt != 0 && timestamp >= e.expireAt {
		return nil
	}
	return e
}

// ==================== LEVEL 1 (updated) ====================

func (db *InMemoryDB) Set(key, field, value string) {
	db.SetAt(key, field, value, 0)
}

func (db *InMemoryDB) SetAt(key, field, value string, timestamp int) {
	if _, ok := db.records[key]; !ok {
		db.records[key] = make(map[string]*FieldEntry)
	}
	db.records[key][field] = &FieldEntry{value: value}
}

func (db *InMemoryDB) SetAtWithTTL(key, field, value string, timestamp, ttl int) {
	if _, ok := db.records[key]; !ok {
		db.records[key] = make(map[string]*FieldEntry)
	}
	db.records[key][field] = &FieldEntry{value: value, expireAt: timestamp + ttl}
}

func (db *InMemoryDB) Get(key, field string) *string {
	return db.GetAt(key, field, 0)
}

func (db *InMemoryDB) GetAt(key, field string, timestamp int) *string {
	e := db.getEntry(key, field, timestamp)
	if e == nil {
		return nil
	}
	return &e.value
}

func (db *InMemoryDB) Delete(key, field string) bool {
	return db.DeleteAt(key, field, 0)
}

func (db *InMemoryDB) DeleteAt(key, field string, timestamp int) bool {
	if db.getEntry(key, field, timestamp) == nil {
		return false
	}
	delete(db.records[key], field)
	return true
}

// ==================== LEVEL 2 (updated) ====================

func (db *InMemoryDB) Scan(key string) []string {
	return db.ScanAt(key, 0)
}

func (db *InMemoryDB) ScanAt(key string, timestamp int) []string {
	return db.scanByPrefixAt(key, "", timestamp)
}

func (db *InMemoryDB) ScanByPrefix(key, prefix string) []string {
	return db.ScanByPrefixAt(key, prefix, 0)
}

func (db *InMemoryDB) ScanByPrefixAt(key, prefix string, timestamp int) []string {
	return db.scanByPrefixAt(key, prefix, timestamp)
}

func (db *InMemoryDB) scanByPrefixAt(key, prefix string, timestamp int) []string {
	rec, ok := db.records[key]
	if !ok {
		return []string{}
	}
	var result []string
	for field, e := range rec {
		// skip expired
		if e.expireAt != 0 && timestamp >= e.expireAt {
			continue
		}
		if strings.HasPrefix(field, prefix) {
			result = append(result, field+"("+e.value+")")
		}
	}
	sort.Strings(result)
	return result
}

// level 4
func (db *InMemoryDB) Backup(timestamp int) int {
	snapshot := make(map[string]map[string]*BackupEntry)
	count := 0

	for key, rec := range db.records {
		fields := make(map[string]*BackupEntry)
		hasLive := false
		for field, e := range rec {
			// skip already expired
			if e.expireAt != 0 && timestamp >= e.expireAt {
				continue
			}
			fields[field] = &BackupEntry{
				value:      e.value,
				expireAt:   e.expireAt, // 0 or absolute timestamp
				backupTime: timestamp,
			}
			hasLive = true
		}
		if hasLive {
			snapshot[key] = fields
			count++
		}
	}

	db.backups = append(db.backups, Backup{timestamp: timestamp, records: snapshot})
	return count
}

// remainingTTL = expireAt - backupTimestamp
// newExpireAt  = restoreTimestamp + remainingTTL
func (db *InMemoryDB) Restore(timestamp, timestampToRestore int) {
	// Find latest backup at or before timestampToRestore
	var target *Backup
	for i := len(db.backups) - 1; i >= 0; i-- {
		if db.backups[i].timestamp <= timestampToRestore {
			target = &db.backups[i]
			break
		}
	}

	// Rebuild records from backup, recalculating expireAt
	db.records = make(map[string]map[string]*FieldEntry)
	for key, rec := range target.records {
		db.records[key] = make(map[string]*FieldEntry)
		for field, e := range rec {
			newEntry := &FieldEntry{value: e.value}
			if e.expireAt != 0 {
				// remaining TTL at backup time
				remainingTTL := e.expireAt - target.timestamp
				newEntry.expireAt = timestamp + remainingTTL
			}
			db.records[key][field] = newEntry
		}
	}
}
