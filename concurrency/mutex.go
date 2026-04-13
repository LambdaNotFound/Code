package concurrency

import (
	"fmt"
	"sync"
)

/**
 * sync.Mutex — mutual exclusion lock
 *
 * Use when multiple goroutines read and write shared state.
 * Only one goroutine holds the lock at a time; others block until it is released.
 *
 * Rules:
 *   - Always pair Lock() with Unlock() — use defer to guarantee release.
 *   - Do not copy a Mutex after first use (pass by pointer).
 *   - Do not lock a Mutex you already hold — Go mutexes are not reentrant.
 */

// -----------------------------------------------------------------------------
// Example 1: Basic Mutex — safe counter
//
// Without a mutex, concurrent increments on a plain int race and produce
// incorrect totals. Wrapping the read-modify-write in Lock/Unlock makes each
// increment atomic with respect to other goroutines.
// -----------------------------------------------------------------------------

type SafeCounter struct {
	mu    sync.Mutex
	count int
}

func (c *SafeCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

func ExampleSafeCounter() {
	counter := &SafeCounter{}
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}

	wg.Wait()
	fmt.Println(counter.Value()) // always 1000
}

// -----------------------------------------------------------------------------
// Example 2: RWMutex — optimised for read-heavy workloads
//
// sync.RWMutex distinguishes readers from writers:
//   - Multiple goroutines may hold a read lock simultaneously (RLock/RUnlock).
//   - A write lock is exclusive — it blocks all readers and other writers (Lock/Unlock).
//
// Prefer RWMutex when reads greatly outnumber writes (e.g. caches, registries).
// -----------------------------------------------------------------------------

type Registry struct {
	mu    sync.RWMutex
	items map[string]string
}

func NewRegistry() *Registry {
	return &Registry{items: make(map[string]string)}
}

func (r *Registry) Set(key, value string) {
	r.mu.Lock() // exclusive — blocks readers and other writers
	defer r.mu.Unlock()
	r.items[key] = value
}

func (r *Registry) Get(key string) (string, bool) {
	r.mu.RLock() // shared — concurrent readers are fine
	defer r.mu.RUnlock()
	v, ok := r.items[key]
	return v, ok
}

func ExampleRegistry() {
	reg := NewRegistry()
	var wg sync.WaitGroup

	// one writer
	wg.Add(1)
	go func() {
		defer wg.Done()
		reg.Set("lang", "Go")
	}()

	// many concurrent readers
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			reg.Get("lang")
		}()
	}

	wg.Wait()
	v, _ := reg.Get("lang")
	fmt.Println(v) // Go
}

// -----------------------------------------------------------------------------
// Example 3: sync.Once — one-time initialisation
//
// sync.Once guarantees that a function runs exactly once regardless of how many
// goroutines call Do concurrently. Useful for lazy singleton initialisation.
// -----------------------------------------------------------------------------

type Singleton struct {
	once  sync.Once
	value string
}

func (s *Singleton) Init(value string) {
	s.once.Do(func() {
		s.value = value
	})
}

func (s *Singleton) Value() string {
	return s.value
}

func ExampleSyncOnce() {
	s := &Singleton{}
	var wg sync.WaitGroup

	// Only the first call to Init takes effect; subsequent calls are no-ops.
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			s.Init(fmt.Sprintf("init-%d", n))
		}(i)
	}

	wg.Wait()
	fmt.Println("initialized:", s.Value() != "") // true
}

// -----------------------------------------------------------------------------
// Example 4: Mutex protecting a slice (append is not goroutine-safe)
//
// append may reallocate the underlying array; concurrent appends without
// synchronisation corrupt the slice header and lose writes.
// -----------------------------------------------------------------------------

type SafeLog struct {
	mu      sync.Mutex
	entries []string
}

func (l *SafeLog) Add(entry string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = append(l.entries, entry)
}

func (l *SafeLog) Snapshot() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	// return a copy so the caller cannot mutate internal state
	snapshot := make([]string, len(l.entries))
	copy(snapshot, l.entries)
	return snapshot
}

func ExampleSafeLog() {
	log := &SafeLog{}
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			log.Add(fmt.Sprintf("event-%d", n))
		}(i)
	}

	wg.Wait()
	fmt.Println("entries:", len(log.Snapshot())) // always 10
}
