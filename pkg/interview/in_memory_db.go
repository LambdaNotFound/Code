package interview

/*
type InMemoryDB struct {
	records map[string]map[string]string // key -> field -> value

}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		records: make(map[string]map[string]string),
	}
}

func (db *InMemoryDB) Set(key, field, value string) {
	if _, ok := db.records[key]; !ok {
		db.records[key] = make(map[string]string)
	}
	db.records[key][field] = value
}

func (db *InMemoryDB) Get(key, field string) *string {
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

func (db *InMemoryDB) Delete(key, field string) bool {
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

func (db *InMemoryDB) Scan(key string) []string {
	rec, ok := db.records[key]
	if !ok {
		return []string{}
	}
	return formatAndSort(rec, "")
}

func (db *InMemoryDB) ScanByPrefix(key, prefix string) []string {
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

*/
