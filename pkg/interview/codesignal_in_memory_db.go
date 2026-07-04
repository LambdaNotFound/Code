package interview

import (
	"sort"
	"strings"
)

/**
 * InMemoryDBV1 — level 1 baseline of the in-memory key/field/value store
 * (see InMemoryDB in in_memory_db2.go for the TTL + backup levels).
 *
 * Storage: two-level map, records[key][field] = value. No ordering is
 * maintained on write; Scan/ScanByPrefix build and sort their output
 * on read via the shared formatAndSort helper.
 *
 * Complexity: Set/Get/Delete are O(1) average (map access). Scan and
 * ScanByPrefix are O(f log f) where f is the number of fields under key,
 * dominated by sort.Strings.
 */
type InMemoryDBV1 struct {
	records map[string]map[string]string // key -> field -> value

}

func NewInMemoryDBV1() *InMemoryDBV1 {
	return &InMemoryDBV1{
		records: make(map[string]map[string]string),
	}
}

func (db *InMemoryDBV1) Set(key, field, value string) {
	if _, ok := db.records[key]; !ok {
		db.records[key] = make(map[string]string)
	}
	db.records[key][field] = value
}

func (db *InMemoryDBV1) Get(key, field string) *string {
	rec, ok := db.records[key]
	if !ok {
		return nil
	}
	val, ok := rec[field]
	if !ok {
		return nil
	}
	return &val
}

func (db *InMemoryDBV1) Delete(key, field string) bool {
	rec, ok := db.records[key]
	if !ok {
		return false
	}
	if _, ok := rec[field]; !ok {
		return false
	}
	delete(rec, field)
	return true
}

func (db *InMemoryDBV1) Scan(key string) []string {
	rec, ok := db.records[key]
	if !ok {
		return []string{}
	}
	return formatAndSort(rec, "")
}

func (db *InMemoryDBV1) ScanByPrefix(key, prefix string) []string {
	rec, ok := db.records[key]
	if !ok {
		return []string{}
	}
	return formatAndSort(rec, prefix)
}

// shared helper
func formatAndSort(rec map[string]string, prefix string) []string {
	var result []string
	for field, value := range rec {
		if strings.HasPrefix(field, prefix) {
			result = append(result, field+"("+value+")")
		}
	}
	sort.Strings(result) // lexicographic on field name, and since format is "field(value)", this works
	return result
}
