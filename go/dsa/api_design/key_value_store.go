package apidesign

/*
 * BEGIN        → push a new layer onto the stack
 * COMMIT       → merge top layer into the one below (or into base if stack is empty)
 * ROLLBACK     → discard top layer
 * GET          → search layers top to bottom, return first hit
 * SET/DELETE   → write only to the top layer
 */

type KVStore struct {
	base   map[string]string
	layers []map[string]string // transaction stack
}

func NewKVStore() *KVStore {
	return &KVStore{base: make(map[string]string)}
}

func (kv *KVStore) active() map[string]string {
	if len(kv.layers) == 0 {
		return kv.base
	}
	return kv.layers[len(kv.layers)-1]
}

func (kv *KVStore) Begin() {
	kv.layers = append(kv.layers, make(map[string]string))
}

func (kv *KVStore) Rollback() {
	if len(kv.layers) == 0 {
		panic("no active transaction")
	}
	kv.layers = kv.layers[:len(kv.layers)-1]
}

func (kv *KVStore) Commit() {
	if len(kv.layers) == 0 {
		panic("no active transaction")
	}
	top := kv.layers[len(kv.layers)-1]
	kv.layers = kv.layers[:len(kv.layers)-1]

	target := kv.base
	if len(kv.layers) > 0 {
		target = kv.layers[len(kv.layers)-1]
	}
	for k, v := range top {
		target[k] = v
	}
}

func (kv *KVStore) Set(key, val string) {
	kv.active()[key] = val
}

func (kv *KVStore) Get(key string) (string, bool) {
	// search top to bottom
	for i := len(kv.layers) - 1; i >= 0; i-- {
		if v, ok := kv.layers[i][key]; ok {
			return v, true
		}
	}
	v, ok := kv.base[key]
	return v, ok
}

func (kv *KVStore) Delete(key string) {
	// sentinel empty string marks deletion
	kv.active()[key] = ""
}
